package main

import (
	_ "embed"
	"fmt"

	"github.com/evmos/ethermint/encoding"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/bytes"

	"github.com/haqq-network/haqq/app"
)

//go:embed tx.json
var txBytes []byte

func main() {
	encodingConfig := encoding.MakeConfig(app.ModuleBasics)

	fmt.Println("Embedded JSON:")
	fmt.Println(string(txBytes))
	fmt.Println("Hash should be: B0131AF60B1A16480C8951BA9C3D187F36B42C8FBC60044D7B6E1BFD1104AA70")

	tx, err := encodingConfig.TxConfig.TxJSONDecoder()(txBytes)
	if err != nil {
		panic(err)
	}

	broadcastBytes, err := encodingConfig.TxConfig.TxEncoder()(tx)
	if err != nil {
		panic(err)
	}

	txHash := tmhash.Sum(broadcastBytes)
	fmt.Println("Actual hash: " + bytes.HexBytes(txHash).String())
}
