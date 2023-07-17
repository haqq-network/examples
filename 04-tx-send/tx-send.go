package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/evmos/ethermint/encoding"

	"github.com/haqq-network/haqq/app"
)

//go:embed tx.json
var txBytes []byte

func main() {
	encodingConfig := encoding.MakeConfig(app.ModuleBasics)

	// Decode transaction from JSON
	tx, err := encodingConfig.TxConfig.TxJSONDecoder()(txBytes)
	if err != nil {
		panic(err)
	}

	// Encode Tx into Protobuf
	broadcastBytes, err := encodingConfig.TxConfig.TxEncoder()(tx)
	if err != nil {
		panic(err)
	}

	// Connect to the node
	node, err := client.NewClientFromNode("https://rpc.tm.testedge2.haqq.network:443")
	if err != nil {
		panic(err)
	}

	// Broadcast transaction
	res, err := node.BroadcastTxSync(context.Background(), broadcastBytes)
	if errRes := client.CheckTendermintError(err, broadcastBytes); errRes != nil {
		panic(errRes)
	}

	txResp := sdk.NewResponseFormatBroadcastTx(res)

	// Get unsigned TX in JSON
	encoded, err := json.Marshal(txResp)
	if err != nil {
		panic(err)
	}

	// Print the serialized transaction
	fmt.Println("TX Broadcasted with the Hash: " + txResp.TxHash)
	fmt.Println("Full JSON:")
	fmt.Println(string(encoded))
}
