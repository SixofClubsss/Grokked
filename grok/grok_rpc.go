package grok

import (
	"fmt"

	"fyne.io/fyne/v2/dialog"
	dreams "github.com/dReam-dApps/dReams"
	"github.com/dReam-dApps/dReams/rpc"
	dero "github.com/deroproject/derohe/rpc"
)

// Owner sets entry amount and pass duration
func Set(scid string, amt, dep, dur uint64) (tx string) {
	args := dero.Arguments{
		dero.Argument{Name: "entrypoint", DataType: "S", Value: "Set"},
		dero.Argument{Name: "amt", DataType: "U", Value: amt},
		dero.Argument{Name: "dur", DataType: "U", Value: dur},
	}

	t1 := dero.Transfer{
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        dep,
	}

	t := []dero.Transfer{t1}
	txid := dero.Transfer_Result{}
	fee := rpc.GasEstimate(scid, "[Grokked]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     scid,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpc.Wallet.CallFor(&txid, "transfer", params); err != nil {
		rpc.PrintError("[Grokked] Set: %s", err)
		return
	}

	rpc.PrintLog("[Grokked] Set TX: %s", txid)

	return txid.TXID
}

// Cancel a game if no players have joined
func Cancel(scid string) (tx string) {
	args := dero.Arguments{dero.Argument{Name: "entrypoint", DataType: "S", Value: "Cancel"}}

	t := []dero.Transfer{}
	txid := dero.Transfer_Result{}
	fee := rpc.GasEstimate(scid, "[Grokked]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     scid,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpc.Wallet.CallFor(&txid, "transfer", params); err != nil {
		rpc.PrintError("[Grokked] Cancel: %s", err)
		return
	}

	rpc.PrintLog("[Grokked] Cancel TX: %s", txid)

	return txid.TXID
}

// Players join a set game
func Join(scid string, amt uint64) (tx string) {
	args := dero.Arguments{dero.Argument{Name: "entrypoint", DataType: "S", Value: "Join"}}

	t1 := dero.Transfer{
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        amt,
	}

	t := []dero.Transfer{t1}
	txid := dero.Transfer_Result{}
	fee := rpc.GasEstimate(scid, "[Grokked]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     scid,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpc.Wallet.CallFor(&txid, "transfer", params); err != nil {
		rpc.PrintError("[Grokked] Join: %s", err)
		return
	}

	rpc.PrintLog("[Grokked] Join TX: %s", txid)

	return txid.TXID
}

// Owner starts the game, must have 3+ players
func Start(scid string) (tx string) {
	args := dero.Arguments{dero.Argument{Name: "entrypoint", DataType: "S", Value: "Start"}}

	t := []dero.Transfer{}
	txid := dero.Transfer_Result{}
	fee := rpc.GasEstimate(scid, "[Grokked]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     scid,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpc.Wallet.CallFor(&txid, "transfer", params); err != nil {
		rpc.PrintError("[Grokked] Start: %s", err)
		return
	}

	rpc.PrintLog("[Grokked] Start TX: %s", txid)

	return txid.TXID
}

// Pass the Grok to another player
func Pass(scid string) (tx string) {
	args := dero.Arguments{dero.Argument{Name: "entrypoint", DataType: "S", Value: "Pass"}}

	t := []dero.Transfer{}
	txid := dero.Transfer_Result{}
	fee := rpc.GasEstimate(scid, "[Grokked]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     scid,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpc.Wallet.CallFor(&txid, "transfer", params); err != nil {
		rpc.PrintError("[Grokked] Pass: %s", err)
		return
	}

	rpc.PrintLog("[Grokked] Pass TX: %s", txid)

	return txid.TXID
}

// Grok player for not paying attention
func Grokked(scid string) (tx string) {
	args := dero.Arguments{dero.Argument{Name: "entrypoint", DataType: "S", Value: "Grokked"}}

	t := []dero.Transfer{}
	txid := dero.Transfer_Result{}
	fee := rpc.GasEstimate(scid, "[Grokked]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     scid,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpc.Wallet.CallFor(&txid, "transfer", params); err != nil {
		rpc.PrintError("[Grokked] Grok: %s", err)
		return
	}

	rpc.PrintLog("[Grokked] Grok TX: %s", txid)

	return txid.TXID
}

// Win scenario when one player left
func Win(scid string, a uint64) (tx string) {
	args := dero.Arguments{
		dero.Argument{Name: "entrypoint", DataType: "S", Value: "Win"},
		dero.Argument{Name: "a", DataType: "U", Value: a},
	}

	t := []dero.Transfer{}
	txid := dero.Transfer_Result{}
	fee := rpc.GasEstimate(scid, "[Grokked]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     scid,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpc.Wallet.CallFor(&txid, "transfer", params); err != nil {
		rpc.PrintError("[Grokked] Win: %s", err)
		return
	}

	rpc.PrintLog("[Grokked] Win TX: %s", txid)

	return txid.TXID
}

// Grok the SC owner and payout all players
func Refund(scid string, p uint64) (tx string) {
	args := dero.Arguments{
		dero.Argument{Name: "entrypoint", DataType: "S", Value: "Refund"},
		dero.Argument{Name: "p", DataType: "U", Value: p},
	}

	t := []dero.Transfer{}
	txid := dero.Transfer_Result{}
	fee := rpc.GasEstimate(scid, "[Grokked]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     scid,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpc.Wallet.CallFor(&txid, "transfer", params); err != nil {
		rpc.PrintError("[Grokked] Refund: %s", err)
		return
	}

	rpc.PrintLog("[Grokked] Refund TX: %s", txid)

	return txid.TXID
}

// Upload a new Grokked SC
func UploadContract(owner bool) (tx string) {
	code := rpc.GetSCCode(GROKSCID)
	if code == "" {
		rpc.PrintError("[Grokked] Upload: error getting Grokked SC")
		return
	}

	fee := rpc.UnlockFee / 3
	if owner {
		fee = 0
	}

	txid := dero.Transfer_Result{}
	t1 := dero.Transfer{
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      fee,
	}

	params := &dero.Transfer_Params{
		Transfers: []dero.Transfer{t1},
		SC_Code:   code,
		SC_Value:  0,
		SC_RPC:    dero.Arguments{},
		Ringsize:  2,
	}

	if err := rpc.Wallet.CallFor(&txid, "transfer", params); err != nil {
		rpc.PrintError("[Grokked] Upload: %s", err)
		return
	}

	rpc.PrintLog("[Grokked] Upload TX: %s", txid)

	return txid.TXID
}

// Update SC to the latest Grokked version
func UpdateGrokked(scid string, version, update uint64, d *dreams.AppObject) (tx string) {
	code := rpc.GetSCCode(GROKSCID)
	if code == "" {
		rpc.PrintError("[Grokked] Update: error getting Grokked SC")
		return
	}

	args := dero.Arguments{
		dero.Argument{Name: "entrypoint", DataType: "S", Value: "UpdateCode"},
		dero.Argument{Name: "code", DataType: "S", Value: code},
	}

	t := []dero.Transfer{}
	txid := dero.Transfer_Result{}
	fee := rpc.GasEstimate(scid, "[Grokked]", args, t, rpc.HighLimitFee*2)
	dialog.NewConfirm("Update SC", fmt.Sprintf("SCID: %s\n\nUpdate from (v%d), to latest version (v%d)? Gas fee is %s DERO", scid, version, update, rpc.FromAtomic(fee, 5)), func(b bool) {
		if b {
			params := &dero.Transfer_Params{
				Transfers: t,
				SC_ID:     scid,
				SC_RPC:    args,
				Ringsize:  2,
				Fees:      fee,
			}

			if err := rpc.Wallet.CallFor(&txid, "transfer", params); err != nil {
				rpc.PrintError("[Grokked] Update: %s", err)
				return
			}

			rpc.PrintLog("[Grokked] Update TX: %s", txid)
			tx = txid.TXID
		}
	}, d.Window).Show()

	return
}
