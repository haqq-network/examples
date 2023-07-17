package main

import (
	"context"
	"crypto/tls"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	cmdcfg "github.com/haqq-network/haqq/cmd/config"
)

func main() {
	// set the address prefixes
	config := sdk.GetConfig()
	cmdcfg.SetBech32Prefixes(config)
	cmdcfg.SetBip44CoinType(config)
	config.Seal()

	addr := sdk.MustAccAddressFromBech32("haqq1rflt3kdkhefxzhkt94cwknmdfvpps39cpr55m0")

	var header metadata.MD
	var dialOpts []grpc.DialOption
	// Setup GRPC options
	dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
		MinVersion: tls.VersionTLS12,
	})))
	// In case of insecure GRPC connection, use the following line instead
	//dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	grpcConn, err := grpc.Dial("grpc.cosmos.testedge2.haqq.network:443", dialOpts...)
	if err != nil {
		panic(err)
	}

	in := &banktypes.QueryAllBalancesRequest{
		Address: addr.String(),
	}
	out := new(banktypes.QueryAllBalancesResponse)

	err = grpcConn.Invoke(context.Background(), "/cosmos.bank.v1beta1.Query/AllBalances", in, out, grpc.Header(&header))
	if err != nil {
		panic(err)
	}

	if len(out.Balances) == 0 {
		panic("no balance")
	}

	for _, balance := range out.Balances {
		fmt.Println(balance.String())
	}
}
