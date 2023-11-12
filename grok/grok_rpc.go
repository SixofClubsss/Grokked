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
		logger.Errorln("[Set]", err)
		return
	}

	logger.Printf("[Set] Set TX: %s\n", txid)
	rpc.AddLog("Set TX: " + txid.TXID)

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
		logger.Errorln("[Join]", err)
		return
	}

	logger.Printf("[Join] Join TX: %s\n", txid)
	rpc.AddLog("Join TX: " + txid.TXID)

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
		logger.Errorln("[Start]", err)
		return
	}

	logger.Printf("[Start] Start TX: %s\n", txid)
	rpc.AddLog("Start TX: " + txid.TXID)

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
		logger.Errorln("[Pass]", err)
		return
	}

	logger.Printf("[Pass] Pass TX: %s\n", txid)
	rpc.AddLog("Pass TX: " + txid.TXID)

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
		logger.Errorln("[Grokked]", err)
		return
	}

	logger.Printf("[Grokked] TX: %s\n", txid)
	rpc.AddLog("Grokked TX: " + txid.TXID)

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
		logger.Errorln("[Win]", err)
		return
	}

	logger.Printf("[Win] TX: %s\n", txid)
	rpc.AddLog("Win TX: " + txid.TXID)

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
		logger.Errorln("[Refund]", err)
		return
	}

	logger.Printf("[Refund] Refund TX: %s\n", txid)
	rpc.AddLog("Refund TX: " + txid.TXID)

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
		logger.Errorln("[UploadContract]", err)
		return
	}

	logger.Println("[UploadContract] TXID:", txid)
	rpc.AddLog("Upload Grokked SC: " + txid.TXID)

	return txid.TXID
}
