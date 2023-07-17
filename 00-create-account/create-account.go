package main

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/tendermint/tendermint/libs/bytes"

	cmdcfg "github.com/haqq-network/haqq/cmd/config"
)

func main() {
	// set the address prefixes
	config := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(config)
	cmdcfg.SetBip44CoinType(config)
	config.Seal()

	// Generate a new private key
	priv, err := ethsecp256k1.GenerateKey()
	if err != nil {
		panic(err)
	}

	pk, err := priv.ToECDSA()
	if err != nil {
		panic(err)
	}

	// Formats key for output
	privB := ethcrypto.FromECDSA(pk)
	keyS := strings.ToUpper(hexutil.Encode(privB)[2:])

	// Get address from private key
	address := priv.PubKey().Address()

	// Print addresses
	fmt.Println("New private key generated!")
	fmt.Println("PK: " + keyS)
	fmt.Println("Address bytes:", address)
	fmt.Printf("Address (hex): %s\n", bytes.HexBytes(address).String())
	fmt.Printf("Address (EIP-55): %s\n", common.BytesToAddress(address))
	fmt.Printf("Bech32 Acc: %s\n", sdk.AccAddress(address))
}
