package grok

import (
	"github.com/dReam-dApps/dReams/rpc"
	dero "github.com/deroproject/derohe/rpc"
)

// Owner sets entry amount and pass duration
func Set(amt, dur uint64) (tx string) {
	rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()

	args := dero.Arguments{
		dero.Argument{Name: "entrypoint", DataType: "S", Value: "Set"},
		dero.Argument{Name: "amt", DataType: "U", Value: amt},
		dero.Argument{Name: "dur", DataType: "U", Value: dur},
	}

	t1 := dero.Transfer{
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        amt / 2,
	}

	t := []dero.Transfer{t1}
	txid := dero.Transfer_Result{}
	fee := rpc.GasEstimate(GROKSCID, "[Set]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     GROKSCID,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
		rpc.PrintError("[Grokked] Set: %s", err)
		return
	}

	rpc.PrintLog("[Grokked] Set TX: %s", txid)

	return txid.TXID
}

// Players join a set game
func Join(amt uint64) (tx string) {
	rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()

	args := dero.Arguments{dero.Argument{Name: "entrypoint", DataType: "S", Value: "Join"}}

	t1 := dero.Transfer{
		Destination: "dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn",
		Amount:      0,
		Burn:        amt,
	}

	t := []dero.Transfer{t1}
	txid := dero.Transfer_Result{}
	fee := rpc.GasEstimate(GROKSCID, "[Join]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     GROKSCID,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
		rpc.PrintError("[Grokked] Join: %s", err)
		return
	}

	rpc.PrintLog("[Grokked] Join TX: %s", txid)

	return txid.TXID
}

// Owner starts the game, must have 3+ players
func Start() (tx string) {
	rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()

	args := dero.Arguments{dero.Argument{Name: "entrypoint", DataType: "S", Value: "Start"}}

	t := []dero.Transfer{}
	txid := dero.Transfer_Result{}
	fee := rpc.GasEstimate(GROKSCID, "[Start]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     GROKSCID,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
		rpc.PrintError("[Grokked] Start: %s", err)
		return
	}

	rpc.PrintLog("[Grokked] Start TX: %s", txid)

	return txid.TXID
}

// Pass the Grok to another player
func Pass() (tx string) {
	rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()

	args := dero.Arguments{dero.Argument{Name: "entrypoint", DataType: "S", Value: "Pass"}}

	t := []dero.Transfer{}
	txid := dero.Transfer_Result{}
	fee := rpc.GasEstimate(GROKSCID, "[Pass]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     GROKSCID,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
		rpc.PrintError("[Grokked] Pass: %s", err)
		return
	}

	rpc.PrintLog("[Grokked] Pass TX: %s", txid)

	return txid.TXID
}

// Grok player for not paying attention
func Grokked() (tx string) {
	rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()

	args := dero.Arguments{dero.Argument{Name: "entrypoint", DataType: "S", Value: "Grokked"}}

	t := []dero.Transfer{}
	txid := dero.Transfer_Result{}
	fee := rpc.GasEstimate(GROKSCID, "[Grokked]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     GROKSCID,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
		rpc.PrintError("[Grokked] Grok: %s", err)
		return
	}

	rpc.PrintLog("[Grokked] Grok TX: %s", txid)

	return txid.TXID
}

// Win scenario when one player left
func Win(a uint64) (tx string) {
	rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()

	args := dero.Arguments{
		dero.Argument{Name: "entrypoint", DataType: "S", Value: "Win"},
		dero.Argument{Name: "a", DataType: "U", Value: a},
	}

	t := []dero.Transfer{}
	txid := dero.Transfer_Result{}
	fee := rpc.GasEstimate(GROKSCID, "[Win]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     GROKSCID,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
		rpc.PrintError("[Grokked] Win: %s", err)
		return
	}

	rpc.PrintLog("[Grokked] Win TX: %s", txid)

	return txid.TXID
}

// Grok the SC owner and payout all players
func Refund(p uint64) (tx string) {
	rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Rpc, rpc.Wallet.UserPass)
	defer cancel()

	args := dero.Arguments{
		dero.Argument{Name: "entrypoint", DataType: "S", Value: "Refund"},
		dero.Argument{Name: "p", DataType: "U", Value: p},
	}

	t := []dero.Transfer{}
	txid := dero.Transfer_Result{}
	fee := rpc.GasEstimate(GROKSCID, "[Refund]", args, t, rpc.LowLimitFee)
	params := &dero.Transfer_Params{
		Transfers: t,
		SC_ID:     GROKSCID,
		SC_RPC:    args,
		Ringsize:  2,
		Fees:      fee,
	}

	if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
		rpc.PrintError("[Grokked] Refund: %s", err)
		return
	}

	rpc.PrintLog("[Grokked] Refund TX: %s", txid)

	return txid.TXID
}

// Upload a new Grokked SC
func UploadContract(owner bool) (tx string) {
	rpcClientW, ctx, cancel := rpc.SetWalletClient(rpc.Wallet.Address, rpc.Wallet.UserPass)
	defer cancel()

	code := rpc.GetSCCode(GROKSCID)
	if code == "" {
		logger.Errorln("[UploadContract] Error getting Grokked SC")
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

	if err := rpcClientW.CallFor(ctx, &txid, "transfer", params); err != nil {
		rpc.PrintError("[Grokked] Upload: %s", err)
		return
	}

	rpc.PrintLog("[Grokked] Upload TX: %s", txid)

	return txid.TXID
}
