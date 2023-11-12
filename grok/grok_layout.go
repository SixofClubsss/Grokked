package grok

import (
	"fmt"
	"image/color"
	"sort"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/civilware/Gnomon/structures"
	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/dwidget"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"
)

func LayoutAllItems(d *dreams.AppObject) fyne.CanvasObject {
	var scid string
	var disableFunc func()
	var synced, isOwner bool

	// SCID select and filter options
	contract_select := widget.NewSelect([]string{GROKSCID}, func(s string) {
		scid = s
	})
	contract_select.PlaceHolder = "Grokked SCIDs:"

	show_opt := widget.NewRadioGroup([]string{"All", "Joined", "Owned"}, nil)
	show_opt.Horizontal = true
	show_opt.Required = true
	show_opt.OnChanged = func(s string) {
		go func() {
			contract_select.Options = []string{}
			switch s {
			case "Owned":
				contract_select.Options, _ = createGrokkedList(true)
			case "Joined":
				var joined, list []string
				list, _ = createGrokkedList(false)
				for _, sc := range list {
					for i := uint64(0); i < 31; i++ {
						var v []*structures.SCIDVariable
						if addr, _, _ := menu.Gnomes.Indexer.GetSCIDValuesByKey(v, sc, i, menu.Gnomes.Indexer.ChainHeight); addr != nil {
							if addr[0] == rpc.Wallet.Address {
								joined = append(joined, sc)
							}
						}
					}
				}
				contract_select.Options = joined

			default:
				contract_select.Options, _ = createGrokkedList(false)
			}
		}()
	}

	// Main display label and image indicator
	label := widget.NewLabel("Select SCID")
	label.Alignment = fyne.TextAlignCenter
	ind := canvas.NewImageFromImage(nil)
	ind.SetMinSize(fyne.NewSize(400, 320))

	// Set game objects
	set_amt := dwidget.NewDeroEntry("", 0.1, 5)
	set_amt.SetPlaceHolder("Dero:")
	set_amt.AllowFloat = true
	set_dur := dwidget.NewDeroEntry("", 1, 0)
	set_dur.SetPlaceHolder("Minutes")
	set_dur.AllowFloat = false

	set_spacer := canvas.NewRectangle(color.Transparent)
	set_spacer.SetMinSize(fyne.NewSize(200, 0))

	var confirming bool
	confirmed := make(chan bool)
	set_button := widget.NewButton("Set", nil)

	set_box := container.NewBorder(nil, set_spacer, nil, set_button, container.NewAdaptiveGrid(2, set_amt, set_dur))
	set_box.Hide()

	// Install SC buttons
	unlock_button := widget.NewButton("Unlock SC", nil)
	unlock_button.OnTapped = func() {
		dialog.NewConfirm("Unlock SC", "Unlock Grokked SC\n\nTo help support the project a 1.00000 DERO donation is attached to this action", func(b bool) {
			if b {
				go func() {
					if tx := UploadContract(isOwner); tx != "" {
						label.SetText("Confirming Install TX...")
						confirming = true
						disableFunc()
						rpc.ConfirmTx(tx, "Grokked", 90)
						time.Sleep(time.Second)
					}
					confirmed <- true
				}()
			}
		}, d.Window).Show()
	}
	unlock_button.Hide()

	new_button := widget.NewButton("New SC", nil)
	new_button.OnTapped = func() {
		dialog.NewConfirm("New SC", "Install a new Grokked SC", func(b bool) {
			if b {
				go func() {
					if tx := UploadContract(isOwner); tx != "" {
						label.SetText("Confirming Install TX...")
						confirming = true
						disableFunc()
						rpc.ConfirmTx(tx, "Grokked", 90)
						time.Sleep(time.Second)
					}
					confirmed <- true
				}()
			}
		}, d.Window).Show()
	}
	new_button.Hide()

	// Set game button
	set_button.OnTapped = func() {
		if dur, err := strconv.ParseInt(set_dur.Text, 10, 64); err == nil {
			if dur < 5 {
				dialog.NewInformation("Short Duration", "Duration has to be longer than 5 minutes", d.Window).Show()
				return
			}
		}

		var dep float64
		amt, err := strconv.ParseFloat(set_amt.Text, 64)
		if err == nil {
			dep = amt / 2
		}
		dialog.NewConfirm("Set", fmt.Sprintf("Set game for %.5f DERO and %s minutes pass time\n\nYou will deposit %.5f DERO into SC", amt, set_dur.Text, dep), func(b bool) {
			if b {
				go func() {
					seconds := rpc.Uint64Type(set_dur.Text) * 60
					if tx := Set(rpc.ToAtomic(set_amt.Text, 5), seconds); tx != "" {
						label.SetText("Confirming Set TX...")
						confirming = true
						disableFunc()
						rpc.ConfirmTx(tx, "Grokked", 90)
						time.Sleep(time.Second)
					}
					confirmed <- true
				}()
			}
		}, d.Window).Show()
	}

	// Start game button
	start_button := widget.NewButton("Start", nil)
	start_button.Hide()
	start_button.OnTapped = func() {
		dialog.NewConfirm("Start", "Start this game?", func(b bool) {
			if b {
				go func() {
					if tx := Start(); tx != "" {
						label.SetText("Confirming Start TX...")
						confirming = true
						disableFunc()
						rpc.ConfirmTx(tx, "Grokked", 90)
						time.Sleep(time.Second)
					}
					confirmed <- true
				}()
			}
		}, d.Window).Show()
	}

	// View current players in game
	var players []string
	players_select := widget.NewSelect(players, nil)
	players_select.PlaceHolder = "View Players"

	// Pass Grok button
	pass_button := widget.NewButton("Pass", nil)
	pass_button.OnTapped = func() {
		dialog.NewConfirm("Pass", "Pass Grok to new player", func(b bool) {
			if b {
				go func() {
					if tx := Pass(); tx != "" {
						label.SetText("Confirming Pass TX...")
						confirming = true
						disableFunc()
						rpc.ConfirmTx(tx, "Grokked", 90)
						time.Sleep(time.Second)
					}
					confirmed <- true
				}()
			}
		}, d.Window).Show()
	}
	pass_button.Hide()

	// Join game button
	join_button := widget.NewButton("Join", nil)
	join_button.OnTapped = func() {
		if _, amt := menu.Gnomes.GetSCIDValuesByKey(scid, "amount"); amt != nil {
			dialog.NewConfirm("Join Game", fmt.Sprintf("Entry is %s DERO", rpc.FromAtomic(amt[0], 5)), func(b bool) {
				if b {
					go func() {
						if tx := Join(amt[0]); tx != "" {
							confirming = true
							label.SetText("Confirming Join TX...")
							disableFunc()
							rpc.ConfirmTx(tx, "Grokked", 90)
							time.Sleep(time.Second)
						}
						confirmed <- true
					}()
				}
			}, d.Window).Show()

			return
		}

		logger.Errorln("[Set] No amount given for TX")
	}
	join_button.Hide()

	// Grok owner button, payout all players still in game
	grok_button := widget.NewButton("Grok", nil)
	grok_button.OnTapped = func() {
		num := uint64(99)
		for i := uint64(0); i < 31; i++ {
			var v []*structures.SCIDVariable
			if addr, _, _ := menu.Gnomes.Indexer.GetSCIDValuesByKey(v, scid, i, menu.Gnomes.Indexer.ChainHeight); addr != nil {
				if addr[0] == rpc.Wallet.Address {
					num = i
					break
				}
			}
		}

		if num != 99 {
			dialog.NewConfirm("Grokked", "Grok owner?", func(b bool) {
				if b {
					go func() {
						if tx := Refund(num); tx != "" {
							label.SetText("Confirming Grokked TX...")
							confirming = true
							disableFunc()
							rpc.ConfirmTx(tx, "Grokked", 90)
							time.Sleep(time.Second)
						}
						confirmed <- true
					}()
				}
			}, d.Window).Show()
		}

	}
	grok_button.Hide()

	disableFunc = func() {
		start_button.Hide()
		join_button.Hide()
		start_button.Hide()
		set_box.Hide()
		pass_button.Hide()
		unlock_button.Hide()
		new_button.Hide()
		grok_button.Hide()
	}

	// Main process for Grokked
	go func() {
		for {
			select {
			case <-d.Receive():
				if !rpc.Wallet.IsConnected() || !rpc.Daemon.IsConnected() {
					disableFunc()
					synced = false
					isOwner = false
					contract_select.ClearSelected()
					contract_select.Options = []string{}
					ind.Resource = resourceGrokJpg
					ind.Refresh()
					label.SetText("Connect your wallet and daemon...")
					d.WorkDone()
					continue
				}

				// Continue if confirming tx
				if confirming {
					disableFunc()
					select {
					case <-confirmed:
						confirming = false
						d.WorkDone()
						continue
					default:
						d.WorkDone()
						continue
					}
				} else {
					if isOwner {
						new_button.Show()
						unlock_button.Hide()
					} else {
						unlock_button.Show()
						new_button.Hide()
					}
				}

				// Grok initial sync
				if !synced && menu.GnomonScan(d.IsConfiguring()) {
					logger.Println("[Grokked] Syncing")
					contract_select.Options, isOwner = createGrokkedList(true)
					show_opt.SetSelected("Owned")
					synced = true

				}

				if len(scid) != 64 {
					ind.Resource = resourceGrokJpg
					ind.Refresh()
					label.SetText("Select a SCID...")
					d.WorkDone()
					continue
				}

				if menu.Gnomes.IsReady() {
					var playing bool
					players = []string{""}
					players_select.ClearSelected()
					// Find players in this round
					for i := uint64(0); i < 31; i++ {
						var v []*structures.SCIDVariable
						if addr, _, _ := menu.Gnomes.Indexer.GetSCIDValuesByKey(v, scid, i, menu.Gnomes.Indexer.ChainHeight); addr != nil {
							players = append(players, addr[0])
							if addr[0] == rpc.Wallet.Address {
								playing = true
							}
						}
					}

					sort.Strings(players)
					players_select.Options = players

					// Find owner of SC
					var owned bool
					if owner, _ := menu.Gnomes.GetSCIDValuesByKey(scid, "owner"); owner != nil {
						if owner[0] == rpc.Wallet.Address {
							owned = true
						}
					}

					if _, u := menu.Gnomes.GetSCIDValuesByKey(scid, "start"); u != nil {
						switch u[0] {
						case 0:
							// Waiting for owner to set game
							grok_button.Hide()
							join_button.Hide()
							pass_button.Hide()
							if owned {
								if !confirming {
									set_box.Show()
								}
							} else {
								set_box.Hide()
							}
							ind.Resource = resourceGrokkedJpg
							ind.Refresh()
							label.SetText("Waiting for owner to set the game...")
						case 1:
							// Game is set, waiting for players to join or owner to start game
							grok_button.Hide()
							pass_button.Hide()
							if !playing {
								if !confirming {
									join_button.Show()
									label.SetText("Join the game...")
								}
								ind.Resource = resourceGrokJpg
								ind.Refresh()
							} else {
								join_button.Hide()
								ind.Resource = resourceJediJpg
								ind.Refresh()
								label.SetText("Waiting for owner to start the game...")
							}

							if owned {
								if !confirming {
									start_button.Show()
								}
								set_box.Hide()
							} else {
								start_button.Hide()
								set_box.Hide()
							}

							// Grok owner and refund if game hasn't started for 48hrs
							if _, last := menu.Gnomes.GetSCIDValuesByKey(scid, "last"); last != nil {
								now := uint64(time.Now().Unix())
								if now > last[0]+173400 {
									if !confirming {
										grok_button.Show()
									} else {
										grok_button.Hide()
										label.SetText("Confirming Grokked TX...")
									}
								} else {
									grok_button.Hide()
								}
							}
						case 2:
							join_button.Hide()
							set_box.Hide()
							start_button.Hide()
							if _, in := menu.Gnomes.GetSCIDValuesByKey(scid, "in"); in != nil {
								if _, u := menu.Gnomes.GetSCIDValuesByKey(scid, "grok"); u != nil {
									if addr, _ := menu.Gnomes.GetSCIDValuesByKey(scid, u[0]); addr != nil {
										// If only one player left, win situation
										if in[0] == 1 {
											var amt uint64
											_, pot := menu.Gnomes.GetSCIDValuesByKey(scid, "pot")
											if pot != nil {
												amt = pot[0] / 2
											}

											if addr[0] != rpc.Wallet.Address {
												label.SetText(fmt.Sprintf("Try harder Grok, Winner is %d, (%s DERO)\n(%s)\n\n", u[0], rpc.FromAtomic(amt, 5), addr[0]))
												ind.Resource = resourceGrokJpg
												ind.Refresh()
											} else {
												label.SetText(fmt.Sprintf("Well done Jedi, You are the winner of this round, (%s DERO)\n(%s)\n\n", rpc.FromAtomic(amt, 5), addr[0]))
												ind.Resource = resourceJediJpg
												ind.Refresh()
											}
											pass_button.Hide()
										} else {
											// Find Grok time frame
											var tf uint64
											var overdue bool
											left := "?"
											now := uint64(time.Now().Unix())
											_, last := menu.Gnomes.GetSCIDValuesByKey(scid, "last")
											_, dur := menu.Gnomes.GetSCIDValuesByKey(scid, "duration")
											if last != nil && dur != nil {
												tf = last[0] + dur[0]
												if now < tf {
													left = fmt.Sprintf("%d minutes left", (tf-now)/60)
												} else if tf != 0 {
													overdue = true
													left = fmt.Sprintf("%d minutes past", (now-tf)/60)
												}
											}

											// Grok image and text
											if addr[0] == rpc.Wallet.Address {
												ind.Resource = resourceGrokJpg
												ind.Refresh()
												if !confirming {
													if !overdue {
														pass_button.Show()
														label.SetText(fmt.Sprintf("You are the Grok, %d, better pass soon (%s)\n\n", u[0], left))
													} else {
														pass_button.Hide()
														label.SetText(fmt.Sprintf("Waiting to be Grokked\n\nYou weren't paying enough attention %d, (%s)\n\n", u[0], left))
													}
												} else {
													label.SetText("Confirming Pass TX...")
												}
											} else {
												if playing {
													ind.Resource = resourceJediJpg
													ind.Refresh()
													var owner_overdue string
													if now < tf+600 {
														owner_overdue = fmt.Sprintf("Can Grok owner in %d minutes", ((tf+600)-now)/60)
													} else {
														owner_overdue = "You can Grok the owner"
													}
													label.SetText(fmt.Sprintf("Well done Jedi, Grok is %d, (%s)\n\n%s", u[0], left, owner_overdue))
												} else {
													ind.Resource = resourceGrokJpg
													ind.Refresh()
													label.SetText("You've been Grokked, pay more attention next time...")
													grok_button.Hide()
												}

												pass_button.Hide()
											}

											// Grok the owner
											if now > tf+600 {
												if !confirming {
													grok_button.Show()
												} else {
													grok_button.Hide()
													label.SetText("Confirming Grok TX...")
												}
											} else {
												grok_button.Hide()
											}
										}
									}
								}

							} else {
								label.SetText("Start the game...")
							}
						default:
							// Nothing
						}
					}
				}
				d.WorkDone()

			case <-d.CloseDapp():
				logger.Println("[Grokked] Done")
				return
			}
		}
	}()

	return container.NewBorder(
		container.NewBorder(container.NewCenter(container.NewVBox(show_opt, contract_select)), label, nil, nil, container.NewCenter(ind)),
		container.NewVBox(
			container.NewCenter(unlock_button),
			container.NewCenter(new_button)),
		nil,
		nil,
		container.NewVBox(
			layout.NewSpacer(),
			container.NewCenter(players_select),
			container.NewCenter(set_box),
			container.NewCenter(start_button),
			container.NewCenter(grok_button),
			container.NewCenter(join_button),
			container.NewCenter(pass_button),
			layout.NewSpacer()))
}

// Create list of Grokked SCIDs from index
func createGrokkedList(owned bool) (options []string, owner bool) {
	if menu.Gnomes.IsReady() {
		scids := menu.Gnomes.GetAllOwnersAndSCIDs()
		for scid := range scids {
			if !menu.Gnomes.IsReady() {
				break
			}

			if _, join := menu.Gnomes.GetSCIDValuesByKey(scid, "joined"); join != nil {
				if _, start := menu.Gnomes.GetSCIDValuesByKey(scid, "start"); start != nil {
					if _, last := menu.Gnomes.GetSCIDValuesByKey(scid, "last"); last != nil {
						_, version := menu.Gnomes.GetSCIDValuesByKey(scid, "v")
						if version != nil {
							v := version[0]

							if owned {
								o, _ := menu.Gnomes.GetSCIDValuesByKey(scid, "owner")
								if o != nil {
									if o[0] == rpc.Wallet.Address {
										owner = true
										options = append(options, scid)
										continue
									}
								}
							} else {
								if v > 1 {
									options = append(options, scid)
									continue
								}
							}
						}
					}
				}
			}
			sort.Strings(options)
		}
	}

	return
}
