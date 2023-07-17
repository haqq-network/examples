package main

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/evmos/ethermint/encoding"

	"github.com/haqq-network/haqq/app"
)

const (
	denom = "aISLM"
	exp   = 10e17
)

func main() {
	// Set up a sample transaction parameters
	encodingConfig := encoding.MakeConfig(app.ModuleBasics)
	chainId := "haqq_11235-1"
	gasLimit := uint64(500000)
	sender := sdk.MustAccAddressFromBech32("haqq1ru779uy8pwqc3f6q4wugx5p0su3zqgm9wrc54h")
	recipient := sdk.MustAccAddressFromBech32("haqq1rflt3kdkhefxzhkt94cwknmdfvpps39cpr55m0")
	memo := "Sample transaction"

	// Prepare amount to Send via bank module
	amount := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(1*exp)))
	// or
	amount, err := sdk.ParseCoinsNormalized("1000000000000000000" + denom)
	if err != nil {
		panic(err.Error())
	}

	// Setup fees
	fees, err := sdk.ParseCoinsNormalized("20" + denom)
	if err != nil {
		panic(err)
	}

	// Create a MsgSend transaction
	msg := banktypes.NewMsgSend(sender, recipient, amount)
	if err := msg.ValidateBasic(); err != nil {
		panic(err)
	}

	// Build Transaction
	tx := encodingConfig.TxConfig.NewTxBuilder()
	if err := tx.SetMsgs(msg); err != nil {
		panic(err)
	}
	tx.SetMemo(memo)
	tx.SetFeeAmount(fees)
	tx.SetGasLimit(gasLimit)

	// Get unsigned TX in JSON
	json, err := encodingConfig.TxConfig.TxJSONEncoder()(tx.GetTx())
	if err != nil {
		panic(err)
	}

	// Print the serialized transaction
	fmt.Println("unsigned tx for " + chainId)
	fmt.Println(string(json))
}
