package grok

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/civilware/Gnomon/structures"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"
	"github.com/docopt/docopt-go"
	"github.com/sirupsen/logrus"
)

const GROKSCID = "80c093dc0def477ea962164bbf86432ccde656bfe4d91c9413dfa80c858f8ff1"

var logger = structures.Logger.WithFields(logrus.Fields{})

var command_line string = `Grokker
Grokker app, powered by Gnomon and dReams.

Usage:
  Grokker [options]
  Grokker -h | --help

Options:
  -h --help                      Show this screen.
  --daemon=<127.0.0.1:10102>     Set daemon rpc address to connect.
  --wallet=<127.0.0.1:10103>     Set wallet rpc address to connect.
  --login=<user:pass>     	 Wallet rpc user:pass for auth.
  --scid=<scid>	         	 Set SCID for Grokker to watch.
  --fastsync=<true>	         Gnomon option,  true/false value to define loading at chain height on start up.
  --num-parallel-blocks=<5>      Gnomon option,  defines the number of parallel blocks to index.`

// Service to automate owner actions for Grokked game
func RunGrokker() {
	n := runtime.NumCPU()
	runtime.GOMAXPROCS(n)

	version := rpc.DREAMSv
	arguments, err := docopt.ParseArgs(command_line, nil, version)

	if err != nil {
		logger.Fatalf("Error while parsing arguments: %s\n", err)
	}

	fastsync := true
	if arguments["--fastsync"] != nil {
		if arguments["--fastsync"].(string) == "false" {
			fastsync = false
		}
	}

	parallel := 5
	if arguments["--num-parallel-blocks"] != nil {
		s := arguments["--num-parallel-blocks"].(string)
		switch s {
		case "1":
			parallel = 1
		case "2":
			parallel = 2
		case "3":
			parallel = 3
		case "4":
			parallel = 4
		case "5":
			parallel = 5
		default:
			parallel = 5
		}
	}

	// Set default rpc params
	rpc.Daemon.Rpc = "127.0.0.1:10102"
	rpc.Wallet.Rpc = "127.0.0.1:10103"

	if arguments["--daemon"] != nil {
		if arguments["--daemon"].(string) != "" {
			rpc.Daemon.Rpc = arguments["--daemon"].(string)
		}
	}

	if arguments["--wallet"] != nil {
		if arguments["--wallet"].(string) != "" {
			rpc.Wallet.Rpc = arguments["--wallet"].(string)
		}
	}

	if arguments["--login"] != nil {
		if arguments["--login"].(string) != "" {
			rpc.Wallet.UserPass = arguments["--login"].(string)
		}
	}

	var scid string
	if arguments["--scid"] != nil {
		if arguments["--scid"].(string) != "" {
			scid = arguments["--scid"].(string)
		}
	}

	if scid == "" {
		logger.Fatalln("[Grokker] No --scid given")
	}

	menu.InitLogrusLog(logrus.InfoLevel)

	menu.Gnomes.Trim = true
	menu.Gnomes.Fast = fastsync
	menu.Gnomes.Para = parallel

	logger.Println("[Grokker]", version, runtime.GOOS, runtime.GOARCH)

	// Check for daemon connection
	rpc.Ping()
	if !rpc.Daemon.Connect {
		logger.Fatalf("[Grokker] Daemon %s not connected\n", rpc.Daemon.Rpc)
	}

	// Check for wallet connection
	rpc.GetAddress("Grokker")
	if !rpc.Wallet.Connect {
		os.Exit(1)
	}

	// Handle ctrl-c close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println()
		menu.Gnomes.Stop("Grokker")
		rpc.Wallet.Connected(false)
		menu.CloseAppSignal(true)
		logger.Println("[Grokker] Closing...")
	}()

	// Set up Gnomon search filters for Duel
	filter := []string{rpc.GetSCCode(GROKSCID)}

	// Set up SCID rating map
	menu.Control.Contract_rating = make(map[string]uint64)

	// Start Gnomon with search filters
	go menu.StartGnomon("Grokker", "boltdb", filter, 0, 0, nil)

	// Routine for checking daemon, wallet connection and Gnomon sync
	go func() {
		for !menu.ClosingApps() && !menu.Gnomes.IsInitialized() {
			time.Sleep(time.Second)
		}

		logger.Println("[Grokker] Starting when Gnomon is synced")
		for !menu.ClosingApps() && menu.Gnomes.IsRunning() && rpc.IsReady() {
			rpc.Ping()
			rpc.EchoWallet("Grokker")
			menu.Gnomes.IndexContains()
			if menu.Gnomes.Indexer.LastIndexedHeight >= menu.Gnomes.Indexer.ChainHeight-3 && menu.Gnomes.HasIndex(1) {
				menu.Gnomes.Synced(true)
			} else {
				menu.Gnomes.Synced(false)
				if !menu.Gnomes.Start && menu.Gnomes.IsInitialized() {
					diff := menu.Gnomes.Indexer.ChainHeight - menu.Gnomes.Indexer.LastIndexedHeight
					if diff > 3 {
						logger.Printf("[Grokker] Gnomon has %d blocks to go\n", diff)
					}
				}
			}
			time.Sleep(3 * time.Second)
		}
		menu.CloseAppSignal(true)
	}()

	// Wait for Gnomon to sync
	for !menu.ClosingApps() && !menu.Gnomes.IsSynced() {
		time.Sleep(time.Second)
	}

	time.Sleep(time.Second)

	grok := uint64(99)
	clock := uint64(9999673725)
	var valid, buffer, firstCase, secondCase bool
	scids := menu.Gnomes.GetAllOwnersAndSCIDs()
	for sc := range scids {
		if sc == scid {
			valid = true
			break
		}
	}

	if !valid {
		logger.Warningf("[Grokker] %s not a valid Grokked SCID", scid)
		menu.Gnomes.Stop("Grokker")
		rpc.Wallet.Connected(false)
		menu.CloseAppSignal(true)
		logger.Println("[Grokker] Closing...")
	}

	// Start Grokker
	for !menu.ClosingApps() {
		time.Sleep(3 * time.Second)
		if _, u := menu.Gnomes.GetSCIDValuesByKey(GROKSCID, "start"); u != nil {
			switch u[0] {
			case 0:
				if !firstCase {
					logger.Println("[Grokker] Waiting to set the game...")
					firstCase = true
					secondCase = false
				}
			case 1:
				if !secondCase {
					logger.Println("[Grokker] Waiting for player to join the game...")
					firstCase = false
					secondCase = true
				}
			case 2:
				firstCase = false
				secondCase = false

				if buffer {
					buffer = false
					continue
				}

				if _, in := menu.Gnomes.GetSCIDValuesByKey(GROKSCID, "in"); in != nil {
					if _, u := menu.Gnomes.GetSCIDValuesByKey(GROKSCID, "grok"); u != nil {
						if u[0] != grok {
							grok = u[0]
							clock = uint64(9999673725)
							logger.Printf("[Grokker] Grok is %d", grok)
						}
						if addr, _ := menu.Gnomes.GetSCIDValuesByKey(GROKSCID, u[0]); addr != nil {
							switch in[0] {
							case 1:
								if tx := Win(u[0]); tx != "" {
									rpc.ConfirmTx(tx, "Grokker", 90)
									time.Sleep(time.Second)
									buffer = true
								}
							default:
								var tf uint64
								now := uint64(time.Now().Unix())
								_, last := menu.Gnomes.GetSCIDValuesByKey(GROKSCID, "last")
								_, dur := menu.Gnomes.GetSCIDValuesByKey(GROKSCID, "duration")
								if last != nil && dur != nil {
									tf = last[0] + dur[0]
									if now > tf+10 {
										logger.Printf("[Grokker] Grokking %d, %d minutes past", u[0], (now-tf)/60)
										if tx := Grokked(); tx != "" {
											rpc.ConfirmTx(tx, "Grokker", 90)
											time.Sleep(time.Second)
											buffer = true
										}
									} else {
										new := (tf - now) / 60
										if new < clock {
											clock = new
											logger.Printf("[Grokker] Not overdue yet, %d minutes left\n", clock)
										}
									}
								}
							}
						}
					}

				} else {
					logger.Println("[Grokker] Can't read in value")
				}
			default:

			}
		}
	}

	logger.Println("[Grokker] Closed")
}