package grok

import (
	"bytes"
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
	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/bundle"
	"github.com/dReam-dApps/dReams/dwidget"
	"github.com/dReam-dApps/dReams/gnomes"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/dReam-dApps/dReams/rpc"
)

func LayoutAll(d *dreams.AppObject) fyne.CanvasObject {
	var scid string
	var max *fyne.Container
	var sc_list *widget.List
	var contracts []gnomes.SC
	var disableFunc func()
	var synced, isOwner bool

	// List sync progress
	progress := widget.NewProgressBar()
	progress.TextFormatter = func() string {
		return ""
	}

	// SC view options
	sc_select := widget.NewSelect([]string{"All", "Joined", "Owned"}, nil)
	sc_select.SetSelected("Owned")
	sc_select.OnChanged = func(s string) {
		go func() {
			reset := max.Objects[0]
			max.Objects[0] = syncScreen(progress)

			sc_list.UnselectAll()
			contracts = []gnomes.SC{}
			switch s {
			case "Owned":
				contracts, _ = createGrokkedList(true, progress)
			case "Joined":
				var joined, all []gnomes.SC
				all, _ = createGrokkedList(false, progress)
				for _, sc := range all {
					for i := uint64(0); i < 31; i++ {
						if addr, _, _ := gnomon.GetLiveSCIDValuesByKey(sc.ID, i); addr != nil {
							if addr[0] == rpc.Wallet.Address {
								joined = append(joined, sc)
							}
						}
					}
				}
				contracts = joined

			default:
				contracts, _ = createGrokkedList(false, progress)
			}
			sc_list.Refresh()
			max.Objects[0] = reset
		}()
	}

	// Refresh viewed sc list
	sc_refresh := widget.NewButtonWithIcon("", dreams.FyneIcon("viewRefresh"), func() { sc_select.OnChanged(sc_select.Selected) })
	sc_refresh.Importance = widget.LowImportance

	// Main display label and image indicator
	label := widget.NewLabel("Select SCID")
	label.Alignment = fyne.TextAlignCenter
	ind := canvas.NewImageFromImage(nil)
	ind.SetMinSize(fyne.NewSize(400, 320))

	// Set game objects
	set_amt := dwidget.NewAmountEntry("", 0.1, 5)
	set_amt.SetPlaceHolder("DERO:")
	set_amt.AllowFloat = true

	set_dur := dwidget.NewAmountEntry("", 1, 0)
	set_dur.SetPlaceHolder("Minutes:")
	set_dur.AllowFloat = false
	set_dur.Validator = func(s string) error {
		dur, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return fmt.Errorf("requires int")
		}

		if dur < 5 {
			return fmt.Errorf("duration to short")
		} else if dur > 4320 {
			return fmt.Errorf("duration to long")
		}

		return nil
	}

	set_dep := dwidget.NewAmountEntry("", 0.1, 5)
	set_dep.SetPlaceHolder("DERO:")
	set_dep.AllowFloat = true

	set_spacer := canvas.NewRectangle(color.Transparent)
	set_spacer.SetMinSize(fyne.NewSize(200, 0))

	// Confirming txs
	var confirming bool
	confirmed := make(chan bool)

	// Set game button
	set_button := widget.NewButton("Set", nil)
	set_button.Importance = widget.HighImportance

	set_box := container.NewVBox(
		container.NewAdaptiveGrid(2, container.NewVBox(widget.NewLabel("Entry Fee"), set_amt), container.NewVBox(widget.NewLabel("Deposit"), set_dep)),
		set_spacer,
		container.NewBorder(widget.NewLabel("Pass Duration"), nil, nil, set_button, set_dur))
	set_box.Hide()

	// Install SC buttons
	unlock_button := widget.NewButton("Unlock SC", nil)
	unlock_button.Importance = widget.HighImportance
	unlock_button.OnTapped = func() {
		dialog.NewConfirm("Unlock SC", "Unlock Grokked SC\n\nTo help support the project a 1.00000 DERO donation is attached to this action", func(b bool) {
			if b {
				go func() {
					tx := UploadContract(isOwner)
					label.SetText("Confirming Install TX...")
					confirming = true
					disableFunc()
					go menu.ShowTxDialog("Upload SC", "", tx, 3*time.Second, d.Window)
					if rpc.ConfirmTx(tx, "Grokked", 45) {
						isOwner = true
					}
					time.Sleep(time.Second)

					confirmed <- true
				}()
			}
		}, d.Window).Show()
	}
	unlock_button.Hide()

	new_button := widget.NewButton("New SC", nil)
	new_button.Importance = widget.HighImportance
	new_button.OnTapped = func() {
		dialog.NewConfirm("New SC", "Install a new Grokked SC, gas fee is 0.07000 DERO", func(b bool) {
			if b {
				go func() {
					tx := UploadContract(isOwner)
					label.SetText("Confirming Install TX...")
					confirming = true
					disableFunc()
					go menu.ShowTxDialog("New SC", "", tx, 3*time.Second, d.Window)
					rpc.ConfirmTx(tx, "Grokked", 45)
					time.Sleep(time.Second)

					confirmed <- true
				}()
			}
		}, d.Window).Show()
	}
	new_button.Hide()

	// Update existing SCID to latest version
	update_button := widget.NewButton("Update SC", nil)
	update_button.Importance = widget.HighImportance
	update_button.Hide()
	update_button.OnTapped = func() {
		if gnomon.IsReady() {
			if _, v := gnomon.GetSCIDValuesByKey(scid, "v"); v != nil {
				if v[0] == scVersion {
					dialog.NewInformation("Update", fmt.Sprintf("SCID: %s\n\nThis SCID is already at latest version (v%d)", scid, scVersion), d.Window).Show()
					return
				}

				go func() {
					tx := UpdateGrokked(scid, v[0], scVersion, d)
					label.SetText("Confirming Update TX...")
					confirming = true
					disableFunc()
					go menu.ShowTxDialog("Update", "", tx, 3*time.Second, d.Window)
					rpc.ConfirmTx(tx, "Grokked", 45)
					time.Sleep(time.Second)

					confirmed <- true
				}()
			} else {
				dialog.NewInformation("Update", "Error getting SCID version", d.Window).Show()
			}
		}
	}

	// Set game button
	set_button.OnTapped = func() {
		if dur, err := strconv.ParseInt(set_dur.Text, 10, 64); err == nil {
			if dur < 5 {
				dialog.NewInformation("Short Duration", "Duration has to be longer than 5 minutes", d.Window).Show()
				return
			} else if dur > 4320 {
				dialog.NewInformation("Long Duration", "Duration has to be longer than 3 days (4320 minutes)", d.Window).Show()
				return
			}
		}

		dialog.NewConfirm("Set", fmt.Sprintf("Set game for %.5f DERO and %s minutes pass time\n\nYou will deposit %.5f DERO into SC", rpc.Float64Type(set_amt.Text), set_dur.Text, rpc.Float64Type(set_dep.Text)), func(b bool) {
			if b {
				go func() {
					seconds := rpc.Uint64Type(set_dur.Text) * 60
					tx := Set(scid, rpc.ToAtomic(set_amt.Text, 5), rpc.ToAtomic(set_dep.Text, 5), seconds)
					label.SetText("Confirming Set TX...")
					confirming = true
					disableFunc()
					go menu.ShowTxDialog("Set", "", tx, 3*time.Second, d.Window)
					rpc.ConfirmTx(tx, "Grokked", 45)
					time.Sleep(time.Second)

					confirmed <- true
				}()
			}
		}, d.Window).Show()
	}

	cancel_button := widget.NewButton("Cancel", nil)
	cancel_button.Importance = widget.HighImportance
	cancel_button.Hide()
	cancel_button.OnTapped = func() {
		dialog.NewConfirm("Cancel", "Cancel the round on this SCID?", func(b bool) {
			if b {
				go func() {
					tx := Cancel(scid)
					label.SetText("Confirming Cancel TX...")
					confirming = true
					disableFunc()
					go menu.ShowTxDialog("Cancel", "", tx, 3*time.Second, d.Window)
					rpc.ConfirmTx(tx, "Grokked", 45)
					time.Sleep(time.Second)

					confirmed <- true
				}()
			}
		}, d.Window).Show()
	}

	// Start game button
	start_button := widget.NewButton("Start", nil)
	start_button.Importance = widget.HighImportance
	start_button.Hide()
	start_button.OnTapped = func() {
		dialog.NewConfirm("Start", "Start this game?", func(b bool) {
			if b {
				go func() {
					tx := Start(scid)
					label.SetText("Confirming Start TX...")
					confirming = true
					disableFunc()
					go menu.ShowTxDialog("Start", "", tx, 3*time.Second, d.Window)
					rpc.ConfirmTx(tx, "Grokked", 45)
					time.Sleep(time.Second)

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
	pass_button.Importance = widget.HighImportance
	pass_button.OnTapped = func() {
		dialog.NewConfirm("Pass", "Pass Grok to new player", func(b bool) {
			if b {
				go func() {
					tx := Pass(scid)
					label.SetText("Confirming Pass TX...")
					confirming = true
					disableFunc()
					go menu.ShowTxDialog("Pass", "", tx, 3*time.Second, d.Window)
					rpc.ConfirmTx(tx, "Grokked", 45)
					time.Sleep(time.Second)

					confirmed <- true
				}()
			}
		}, d.Window).Show()
	}
	pass_button.Hide()

	// Join game button
	join_button := widget.NewButton("Join", nil)
	join_button.Importance = widget.HighImportance
	join_button.OnTapped = func() {
		if gnomon.IsReady() {
			if _, amt := gnomon.GetSCIDValuesByKey(scid, "amount"); amt != nil {
				dialog.NewConfirm("Join Game", fmt.Sprintf("Entry is %s DERO", rpc.FromAtomic(amt[0], 5)), func(b bool) {
					if b {
						go func() {
							tx := Join(scid, amt[0])
							confirming = true
							label.SetText("Confirming Join TX...")
							disableFunc()
							go menu.ShowTxDialog("Join", "", tx, 3*time.Second, d.Window)
							rpc.ConfirmTx(tx, "Grokked", 45)
							time.Sleep(time.Second)

							confirmed <- true
						}()
					}
				}, d.Window).Show()

				return
			}

			logger.Errorln("[Set] No amount given for TX")
		}
	}
	join_button.Hide()

	// Grok owner button, payout all players still in game
	grok_owner_button := widget.NewButton("Grok", nil)
	grok_owner_button.Importance = widget.HighImportance
	grok_owner_button.OnTapped = func() {
		num := uint64(99)
		if gnomon.IsReady() {
			// If waiting for payout use grok for Refund
			if _, in := gnomon.GetSCIDValuesByKey(scid, "in"); in != nil && in[0] == 1 {
				if _, winner := gnomon.GetSCIDValuesByKey(scid, "grok"); winner != nil {
					num = winner[0]
				}
			} else {
				for i := uint64(0); i < 31; i++ {
					if addr, _, _ := gnomon.GetLiveSCIDValuesByKey(scid, i); addr != nil {
						if addr[0] == rpc.Wallet.Address {
							num = i
							break
						}
					}
				}
			}
		}

		if num != 99 {
			dialog.NewConfirm("Grokked", "Grok owner?", func(b bool) {
				if b {
					go func() {
						tx := Refund(scid, num)
						label.SetText("Confirming Grokked TX...")
						confirming = true
						disableFunc()
						go menu.ShowTxDialog("Grokked", "", tx, 3*time.Second, d.Window)
						rpc.ConfirmTx(tx, "Grokked", 45)
						time.Sleep(time.Second)

						confirmed <- true
					}()
				}
			}, d.Window).Show()
		} else {
			dialog.NewInformation("Can't Grok", "You are not in this game", d.Window).Show()
		}
	}
	grok_owner_button.Hide()

	// Owner button to Grok player
	grok_button := widget.NewButton("Grok", nil)
	grok_button.Importance = widget.HighImportance
	grok_button.OnTapped = func() {
		dialog.NewConfirm("Grokked", "Grok Player?", func(b bool) {
			if b {
				go func() {
					tx := Grokked(scid)
					label.SetText("Confirming Grokked TX...")
					confirming = true
					disableFunc()
					go menu.ShowTxDialog("Grokked", "", tx, 3*time.Second, d.Window)
					rpc.ConfirmTx(tx, "Grokked", 45)
					time.Sleep(time.Second)

					confirmed <- true
				}()
			}
		}, d.Window).Show()
	}
	grok_button.Hide()

	// Payout when last player standing
	pay_button := widget.NewButton("Pay", nil)
	pay_button.Importance = widget.HighImportance
	pay_button.OnTapped = func() {
		if gnomon.IsReady() {
			if _, in := gnomon.GetSCIDValuesByKey(scid, "in"); in != nil {
				if _, u := gnomon.GetSCIDValuesByKey(scid, "grok"); u != nil {
					if addr, _ := gnomon.GetSCIDValuesByKey(scid, u[0]); addr != nil {
						dialog.NewConfirm("Pay winner", fmt.Sprintf("Pay %s", addr[0]), func(b bool) {
							if b {
								go func() {
									switch in[0] {
									case 1:
										tx := Win(scid, u[0])
										label.SetText("Confirming Payout TX...")
										confirming = true
										disableFunc()
										go menu.ShowTxDialog("Pay", "", tx, 3*time.Second, d.Window)
										rpc.ConfirmTx(tx, "Grokked", 45)
										time.Sleep(time.Second)
									default:
										dialog.NewInformation("To Many Players", "There are still more players to be Grokked", d.Window).Show()
									}

									confirmed <- true
								}()
							}
						}, d.Window).Show()
					}
				}
			}
		}
	}
	pay_button.Hide()

	// Rate scid
	rate_button := widget.NewButton("Rate SC", func() {
		if len(scid) != 64 {
			dialog.NewInformation("Rate", fmt.Sprintf("Invalid SCID %s", scid), d.Window).Show()
			return
		}

		if !gnomes.CheckOwner(scid) {
			menu.RateConfirm(scid, d)
		} else {
			dialog.NewInformation("Can't rate", "You are the owner of this SCID", d.Window).Show()
			logger.Warnln("[Grokked] Can't rate, you own this contract")
		}
	})

	disableFunc = func() {
		start_button.Hide()
		join_button.Hide()
		start_button.Hide()
		set_box.Hide()
		pass_button.Hide()
		pay_button.Hide()
		unlock_button.Hide()
		new_button.Hide()
		grok_button.Hide()
		grok_owner_button.Hide()
		rate_button.Hide()
		update_button.Hide()
		cancel_button.Hide()
	}

	// Main process for Grokked
	go func() {
		for {
			select {
			case <-d.Receive():
				if !rpc.Wallet.IsConnected() || !rpc.Daemon.IsConnected() {
					disableFunc()
					scid = ""
					synced = false
					isOwner = false
					sc_list.UnselectAll()
					contracts = []gnomes.SC{}
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
						unlock_button.Hide()
						new_button.Show()
					} else {
						new_button.Hide()
						unlock_button.Show()
					}
				}

				// Grok initial sync
				if !synced && gnomes.Scan(d.IsConfiguring()) {
					logger.Println("[Grokked] Syncing")
					contracts, isOwner = createGrokkedList(true, progress)
					synced = true
					rate_button.Show()
				}

				if len(scid) != 64 {
					ind.Resource = resourceGrokJpg
					ind.Refresh()
					label.SetText("Select a SCID...")
					d.WorkDone()
					continue
				}

				if gnomon.IsReady() {
					var playing bool
					players = []string{""}
					players_select.ClearSelected()
					// Find players in this round
					for i := uint64(0); i < 31; i++ {
						if addr, _, _ := gnomon.GetLiveSCIDValuesByKey(scid, i); addr != nil {
							players = append(players, fmt.Sprintf("(%d) %s", i, addr[0]))
							if addr[0] == rpc.Wallet.Address {
								playing = true
							}
						}
					}

					sort.Strings(players)
					players_select.Options = players

					// Find owner of SC
					var owned bool
					if owner, _ := gnomon.GetSCIDValuesByKey(scid, "owner"); owner != nil {
						if owner[0] == rpc.Wallet.Address {
							owned = true
						}
					}

					if _, u := gnomon.GetSCIDValuesByKey(scid, "start"); u != nil {
						switch u[0] {
						case 0:
							// Waiting for owner to set game
							grok_button.Hide()
							grok_owner_button.Hide()
							cancel_button.Hide()
							join_button.Hide()
							pass_button.Hide()
							pay_button.Hide()
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
							pay_button.Hide()
							set_box.Hide()
							var players uint64
							_, in := gnomon.GetSCIDValuesByKey(scid, "in")
							if in != nil {
								players = in[0]
							}

							if players == 0 {
								cancel_button.Show()
							} else {
								cancel_button.Hide()
							}

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
								if players < 3 {
									label.SetText(fmt.Sprintf("Waiting for %d more player(s)...", 3-players))
								} else {
									if !confirming {
										label.SetText(fmt.Sprintf("%d players in, waiting for owner to start the game...", players))
									}
								}
							}

							if owned {
								if !confirming {
									if players > 2 {
										start_button.Show()
									}
								}
							} else {
								start_button.Hide()
							}

							// Grok owner and refund if game hasn't started for over a week
							if _, last := gnomon.GetSCIDValuesByKey(scid, "last"); last != nil {
								now := uint64(time.Now().Unix())
								if now > last[0]+604800 {
									if !confirming {
										grok_owner_button.Show()
									} else {
										grok_owner_button.Hide()
										label.SetText("Confirming Grokked TX...")
									}
								} else {
									grok_owner_button.Hide()
								}
							} else {
								grok_owner_button.Hide()
							}
						case 2:
							join_button.Hide()
							set_box.Hide()
							start_button.Hide()
							cancel_button.Hide()
							if _, in := gnomon.GetSCIDValuesByKey(scid, "in"); in != nil {
								if _, u := gnomon.GetSCIDValuesByKey(scid, "grok"); u != nil {
									if addr, _ := gnomon.GetSCIDValuesByKey(scid, u[0]); addr != nil {
										// If only one player left, win situation
										if in[0] == 1 {
											grok_button.Hide()
											var amt uint64
											_, pot := gnomon.GetSCIDValuesByKey(scid, "pot")
											if pot != nil {
												amt = pot[0] / 2
											}

											if addr[0] != rpc.Wallet.Address {
												label.SetText(fmt.Sprintf("Try harder Grok, Winner is %d, (%s DERO)\n(%s)", u[0], rpc.FromAtomic(amt, 5), addr[0]))
												ind.Resource = resourceGrokJpg
												ind.Refresh()
											} else {
												label.SetText(fmt.Sprintf("Well done Jedi, You are the winner of this round, (%s DERO)\n(%s)", rpc.FromAtomic(amt, 5), addr[0]))
												ind.Resource = resourceJediJpg
												ind.Refresh()
											}

											// Grok the owner
											_, last := gnomon.GetSCIDValuesByKey(scid, "last")
											_, dur := gnomon.GetSCIDValuesByKey(scid, "duration")
											if last != nil && dur != nil {
												now := uint64(time.Now().Unix())
												if now > last[0]+dur[0]+600 {
													if !confirming {
														grok_owner_button.Show()
													} else {
														grok_owner_button.Hide()
														label.SetText("Confirming Grok TX...")
													}
												} else {
													grok_owner_button.Hide()
												}
											}

											// Payout
											if owned {
												if !confirming {
													pay_button.Show()
												} else {
													pay_button.Hide()
													label.SetText("Confirming Pay TX...")
												}
											} else {
												pay_button.Hide()
											}
										} else {
											// Find Grok time frame
											var tf uint64
											var overdue bool
											left := "?"
											now := uint64(time.Now().Unix())
											_, last := gnomon.GetSCIDValuesByKey(scid, "last")
											_, dur := gnomon.GetSCIDValuesByKey(scid, "duration")
											if last != nil && dur != nil {
												tf = last[0] + dur[0]
												if now < tf {
													grok_button.Hide()
													left = fmt.Sprintf("%d minutes left", (tf-now)/60)
												} else if tf != 0 {
													overdue = true
													left = fmt.Sprintf("%d minutes past", (now-tf)/60)
													if owned {
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

											var owner_overdue string
											if now < tf+600 {
												owner_overdue = fmt.Sprintf("Can Grok owner in %d minutes", ((tf+600)-now)/60)
											} else {
												owner_overdue = "You can Grok the owner"
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
														label.SetText(fmt.Sprintf("Waiting to be Grokked\n\nYou weren't paying enough attention %d, (%s)\n\n%s", u[0], left, owner_overdue))
													}
												} else {
													label.SetText("Confirming Pass TX...")
												}
											} else {
												if playing {
													ind.Resource = resourceJediJpg
													ind.Refresh()

													label.SetText(fmt.Sprintf("Well done Jedi, Grok is %d, (%s)\n\n%s", u[0], left, owner_overdue))
												} else {
													ind.Resource = resourceGrokJpg
													ind.Refresh()
													label.SetText("You've been Grokked, pay more attention next time...")
													grok_owner_button.Hide()
												}

												pass_button.Hide()
											}

											// Grok the owner
											if now > tf+600 {
												if !confirming {
													grok_owner_button.Show()
												} else {
													grok_owner_button.Hide()
													label.SetText("Confirming Grok TX...")
												}
											} else {
												grok_owner_button.Hide()
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

	sc_button := widget.NewButton("Select SCID", nil)

	sc_entry := widget.NewEntry()
	sc_entry.Disable()

	sc_spacer := canvas.NewRectangle(color.Transparent)
	sc_spacer.SetMinSize(fyne.NewSize(400, 0))

	grok_cont := container.NewBorder(
		container.NewBorder(
			container.NewCenter(container.NewVBox(canvas.NewLine(bundle.TextColor), dwidget.NewCanvasText("Grokked", 18, fyne.TextAlignCenter), canvas.NewLine(bundle.TextColor))),
			container.NewCenter(container.NewVBox(sc_spacer, sc_entry, container.NewCenter(label))),
			nil,
			nil,
			container.NewCenter(ind)),

		container.NewCenter(sc_button),
		nil,
		nil,
		container.NewVBox(
			container.NewVBox(container.NewCenter(players_select)),
			layout.NewSpacer(),
			container.NewVBox(
				container.NewCenter(set_box),
				container.NewCenter(cancel_button),
				container.NewCenter(start_button),
				container.NewCenter(
					container.NewHBox(
						container.NewCenter(grok_button),
						container.NewCenter(grok_owner_button))),
				container.NewCenter(join_button),
				container.NewCenter(pass_button),
				container.NewCenter(pay_button),
				layout.NewSpacer(),
			),
			layout.NewSpacer(),
			layout.NewSpacer()))

	max = container.NewStack(grok_cont)

	sc_list = widget.NewList(
		func() int {
			return len(contracts)
		},
		func() fyne.CanvasObject {
			spacer := canvas.NewImageFromImage(nil)
			spacer.SetMinSize(fyne.NewSize(60, 60))
			return container.NewStack(
				container.NewBorder(
					nil,
					nil,
					container.NewStack(container.NewCenter(spacer), canvas.NewImageFromImage(nil)),
					nil,
					container.NewBorder(
						nil,
						widget.NewLabel(""),
						container.NewStack(canvas.NewImageFromImage(nil), canvas.NewImageFromImage(nil)),
						nil,
						container.NewHBox(canvas.NewImageFromImage(nil), widget.NewLabel("")))))
		},
		func(id widget.ListItemID, c fyne.CanvasObject) {
			if len(contracts) < 1 {
				return
			}

			sc := contracts[id]

			if c.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).Text != sc.ID {
				c.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(sc.ID)
				c.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("%s   (v%d)   %s", sc.Header.Name, sc.Version, sc.Header.Description))

				badge := canvas.NewImageFromResource(menu.DisplayRating(sc.Rating))
				badge.SetMinSize(fyne.NewSize(35, 35))
				c.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = badge

				if sc.Header.IconURL != "" {
					if img, err := dreams.DownloadBytes(sc.Header.IconURL); err == nil {
						if canv := canvas.NewImageFromReader(bytes.NewReader(img), sc.Header.Name); canv != nil {
							canv.SetMinSize(fyne.NewSize(66, 66))
							c.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = container.NewCenter(canv)
						}
					} else {
						logger.Errorln("[Grokked] img", err)
					}

				} else {
					unknown := canvas.NewImageFromResource(bundle.ResourceGamesCirclePng)
					unknown.SetMinSize(fyne.NewSize(66, 66))
					c.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*fyne.Container).Objects[0] = container.NewCenter(unknown)
				}

				frame := canvas.NewImageFromResource(bundle.ResourceFramePng)
				frame.SetMinSize(fyne.NewSize(70, 70))
				c.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*fyne.Container).Objects[1] = frame
				c.Refresh()
			}
		})

	sc_list.OnSelected = func(id int) {
		scid = contracts[id].ID
		sc_entry.SetText(scid)
		if contracts[id].Version < scVersion {
			update_button.Show()
		} else {
			update_button.Hide()
		}
	}

	back_button := widget.NewButton("Back", func() {
		sc_entry.SetText(scid)
		max.Objects[0] = grok_cont
	})

	sc_button.OnTapped = func() {
		sc_entry.SetText(scid)
		sc_list.UnselectAll()

		max.Objects[0] = container.NewBorder(
			container.NewVBox(
				container.NewCenter(container.NewVBox(canvas.NewLine(bundle.TextColor), dwidget.NewCanvasText("Grokked SCs", 18, fyne.TextAlignCenter), canvas.NewLine(bundle.TextColor))),
				container.NewCenter(container.NewBorder(nil, nil, nil, sc_refresh, sc_select))),
			container.NewCenter(back_button),
			nil,
			nil,
			container.NewAdaptiveGrid(3,
				layout.NewSpacer(),
				container.NewBorder(
					nil,
					container.NewCenter(
						dwidget.NewSpacer(100, 180),
						container.NewVBox(
							container.NewCenter(rate_button),
							layout.NewSpacer(),
							container.NewCenter(unlock_button),
							layout.NewSpacer(),
							container.NewCenter(new_button),
							container.NewCenter(layout.NewSpacer()),
							layout.NewSpacer(),
							container.NewCenter(update_button))),
					nil,
					nil,
					sc_list),
				layout.NewSpacer()))
	}

	return max
}

// Create list of Grokked SCIDs from index
func createGrokkedList(owned bool, progress *widget.ProgressBar) (options []gnomes.SC, owner bool) {
	if gnomon.IsReady() && gnomon.HasChecked() {
		scids := gnomon.GetAllOwnersAndSCIDs()
		_, check := gnomon.GetSCIDValuesByKey(GROKSCID, "v")
		if check == nil {
			return
		}

		scVersion = check[0]

		if progress != nil {
			progress.Max = float64(len(scids))
			progress.SetValue(0)
		}

		for scid := range scids {
			if !gnomon.IsReady() {
				break
			}

			if _, start := gnomon.GetSCIDValuesByKey(scid, "start"); start != nil {
				if _, version := gnomon.GetSCIDValuesByKey(scid, "v"); version != nil {
					var new gnomes.SC
					headers := gnomes.GetSCHeaders(scid)
					new.ID = scid
					new.Version = version[0]
					if headers.Name != "" {
						new.Header.Name = headers.Name
						if headers.Description != "" {
							new.Header.Description = headers.Description
						}
						if headers.IconURL != "" {
							new.Header.IconURL = headers.IconURL
						}
					}

					var hidden bool
					_, restrict := gnomon.GetSCIDValuesByKey(rpc.RatingSCID, "restrict")
					_, rating := gnomon.GetSCIDValuesByKey(rpc.RatingSCID, scid)

					if restrict != nil && rating != nil {
						new.Rating = rating[0]
						if rating[0] <= restrict[0] {
							hidden = true
						}
					}

					if owned {
						if o, _ := gnomon.GetSCIDValuesByKey(scid, "owner"); o != nil {
							if o[0] == rpc.Wallet.Address {
								owner = true
								options = append(options, new)
								continue
							}
						}
					} else {
						if _, join := gnomon.GetSCIDValuesByKey(scid, "joined"); join != nil {
							if _, last := gnomon.GetSCIDValuesByKey(scid, "last"); last != nil {
								if version[0] == check[0] && !hidden {
									options = append(options, new)
									continue
								}
							}
						}
					}
				}
			}

			if progress != nil {
				progress.SetValue(progress.Value + 1)
			}
		}

		// Sort by ratings then name
		sort.Slice(options, func(i, j int) bool {
			if options[i].Rating > options[j].Rating {
				return true
			}

			if options[i].Rating == options[j].Rating && options[i].Header.Name > options[j].Header.Name {
				return true
			}

			return false
		})
	}

	return
}

// Splash screen for when scids list syncing
func syncScreen(bar *widget.ProgressBar) *fyne.Container {
	text := canvas.NewText("Syncing...", color.White)
	text.Alignment = fyne.TextAlignCenter
	text.TextSize = 21

	img := canvas.NewImageFromResource(bundle.ResourceGamesCirclePng)
	img.SetMinSize(fyne.NewSize(150, 150))

	max := container.NewStack(container.NewCenter(img, text), bar)

	return container.NewAdaptiveGrid(3, layout.NewSpacer(), max, layout.NewSpacer())
}
