package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/evmos/ethermint/encoding"
	"github.com/haqq-network/haqq/app"

	"github.com/cosmos/cosmos-sdk/client"
)

func main() {
	txHash := "B0131AF60B1A16480C8951BA9C3D187F36B42C8FBC60044D7B6E1BFD1104AA70"

	// Connect to the node
	node, err := client.NewClientFromNode("https://rpc.tm.testedge2.haqq.network:443")
	if err != nil {
		panic(err)
	}

	hash, err := hex.DecodeString(txHash)
	if err != nil {
		panic(err)
	}

	// TODO: this may not always need to be proven
	// https://github.com/cosmos/cosmos-sdk/issues/6807
	resTx, err := node.Tx(context.Background(), hash, true)
	if err != nil {
		panic(err)
	}

	// Decode the transaction
	encodingConfig := encoding.MakeConfig(app.ModuleBasics)
	decodedTx, err := encodingConfig.TxConfig.TxDecoder()(resTx.Tx)

	// Get unsigned TX in JSON
	encoded, err := encodingConfig.TxConfig.TxJSONEncoder()(decodedTx)
	if err != nil {
		panic(err)
	}

	ecodedTxResponse, err := json.Marshal(resTx)
	if err != nil {
		panic(err)
	}

	// Print the serialized transaction
	fmt.Println("Response:")
	fmt.Println(string(ecodedTxResponse))
	fmt.Println("TX JSON:")
	fmt.Println(string(encoded))
}
