# Grokked
Goal of the game, don't be Grok.

Written in Go and using [Fyne Toolkit](https://fyne.io/), Gokked is a proof of attention game built on Dero's private L1. Powered by [Gnomon](https://github.com/civilware/Gnomon) and [dReams](https://github.com/dReam-dApps/dReams), the goal of the game is simple. Don't get caught being the Grok. Players can join existing Grokked games, or deploy their own Grokked SC and run games with their preferred preferences. All Grokked SCs are tied into dReams ratings system.

The game starts with a player being randomly selected as the Grok. If that player doesn't prove they are paying attention by interacting with a SC (passing the Grok) within a certain time frame, they are removed from the game and a new player becomes the Grok. Each time a player is removed or passes the Grok, the time frame shrinks meaning players have to pay closer attention if they don't want to be removed from the game. SC's tally wins on chain for a global leader board. The owner of the SC facilitates removing players from the game, if the owner is not paying attention any player can Grok the owner meaning all players get a win on the board and a share of the pot.

![goMod](https://img.shields.io/github/go-mod/go-version/SixofClubsss/Grokked.svg)![goReport](https://goreportcard.com/badge/github.com/SixofClubsss/Grokked)[![goDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/SixofClubsss/Grokked)

Grokked dApp is available for download from [dReams](https://dreamdapps.io)

![windowsOS](https://raw.githubusercontent.com/SixofClubsss/dreamdappsite/main/assets/os-windows-green.svg)![macOS](https://raw.githubusercontent.com/SixofClubsss/dreamdappsite/main/assets/os-macOS-green.svg)![linuxOS](https://raw.githubusercontent.com/SixofClubsss/dreamdappsite/main/assets/os-linux-green.svg)

### Owners
Service to automate owner actions
- Install latest [Go version](https://go.dev/doc/install)
- Install [Fyne](https://developer.fyne.io/started/) dependencies
- Clone repo and build using:
```
git clone https://github.com/SixofClubsss/Grokked.git
cd Grokked/cmd/Grokker
go build .
```
Start flags
```
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
  --num-parallel-blocks=<5>      Gnomon option,  defines the number of parallel blocks to index.
```
If using default ports, Grokker service can be started with
```
./Grokker --login=user:pass --scid=scid
```

### Donations
- **Dero Address**: dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn

![DeroDonations](https://raw.githubusercontent.com/SixofClubsss/dreamdappsite/main/assets/DeroDonations.jpg)

---

#### Licensing

Grokked is free and open source.   
The source code is published under the [MIT](https://github.com/SixofClubsss/Grokked/blob/main/LICENSE) License.   
Copyright © 2023 SixofClubs  