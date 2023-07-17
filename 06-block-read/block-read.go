package main

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/evmos/ethermint/encoding"

	"github.com/haqq-network/haqq/app"
)

func main() {
	// Connect to the node
	node, err := client.NewClientFromNode("https://rpc.tm.testedge2.haqq.network:443")
	if err != nil {
		panic(err)
	}

	// Retrieve block
	height := int64(4600600)
	resBlock, err := node.Block(context.Background(), &height)
	if err != nil {
		panic(err)
	}

	if len(resBlock.Block.Data.Txs) == 0 {
		panic("No transactions in block")
	}

	// Decode the transactions and print result
	encodingConfig := encoding.MakeConfig(app.ModuleBasics)
	for i, tx := range resBlock.Block.Data.Txs {
		decodedTx, err := encodingConfig.TxConfig.TxDecoder()(tx)
		if err != nil {
			panic(err)
		}

		// Get unsigned TX in JSON
		encoded, err := encodingConfig.TxConfig.TxJSONEncoder()(decodedTx)
		if err != nil {
			panic(err)
		}

		fmt.Printf("-------- TX %d --------\n", i)
		fmt.Println(string(encoded))
	}
}
