package main

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/evmos/ethermint/encoding"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"github.com/haqq-network/haqq/app"
)

const (
	denom = "aISLM"
	exp   = 10e17
)

func main() {
	// Set up a sample transaction parameters
	encodingConfig := encoding.MakeConfig(app.ModuleBasics)
	chainId := "haqq_54211-3"
	gasLimit := uint64(500000)
	memo := "Sample transaction"

	// Get Private Key
	// ### PUBLIC TEST KEY ###
	// haqq1tjdjfavsy956d25hvhs3p0nw9a7pfghqm0up92
	// steel clarify purity foil garage cry permit razor enforce fetch pulp picnic flee bulk more conduct kiwi enrich winner gossip hotel alpha snake space
	// Private key in HEX — 1BAC9A1B44A2115843860BA5CF3A730314C76E20E7756EE4CDE932E61BA1C17F
	key, err := crypto.HexToECDSA("1BAC9A1B44A2115843860BA5CF3A730314C76E20E7756EE4CDE932E61BA1C17F")
	if err != nil {
		panic(err)
	}
	privkey := &ethsecp256k1.PrivKey{
		Key: crypto.FromECDSA(key),
	}

	// Get sender and recipient Addresses
	sender := privkey.PubKey().Address()
	senderAcc := sdk.AccAddress(sender.Bytes())
	recipient := sdk.MustAccAddressFromBech32("haqq1rflt3kdkhefxzhkt94cwknmdfvpps39cpr55m0")

	// Prepare amount to Send via bank module
	amount := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(1*exp)))
	// or
	amount, err = sdk.ParseCoinsNormalized("1000000000000000000" + denom)
	if err != nil {
		panic(err.Error())
	}

	// Setup fees
	fees, err := sdk.ParseCoinsNormalized("20" + denom)
	if err != nil {
		panic(err)
	}

	// Create a MsgSend transaction
	msg := banktypes.NewMsgSend(senderAcc, recipient, amount)
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

	// Retrieve account number and sequence from the blockchain
	accountNumber, sequence, err := getAccount(senderAcc, encodingConfig)
	if err != nil {
		panic(err)
	}

	// Sign Transaction
	signerData := authsigning.SignerData{
		ChainID:       chainId,
		AccountNumber: accountNumber,
		Sequence:      sequence,
		PubKey:        privkey.PubKey(),
		Address:       senderAcc.String(),
	}

	signMode := encodingConfig.TxConfig.SignModeHandler().DefaultMode()
	// For SIGN_MODE_DIRECT, calling SetSignatures calls setSignerInfos on
	// TxBuilder under the hood, and SignerInfos is needed to generated the
	// sign bytes. This is the reason for setting SetSignatures here, with a
	// nil signature.
	//
	// Note: this line is not needed for SIGN_MODE_LEGACY_AMINO, but putting it
	// also doesn't affect its generated sign bytes, so for code's simplicity
	// sake, we put it here.
	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   privkey.PubKey(),
		Data:     &sigData,
		Sequence: sequence,
	}

	var sigs []signing.SignatureV2
	sigs = []signing.SignatureV2{sig}

	if err := tx.SetSignatures(sigs...); err != nil {
		panic(err)
	}

	// Generate the bytes to be signed.
	bytesToSign, err := encodingConfig.TxConfig.SignModeHandler().GetSignBytes(signMode, signerData, tx.GetTx())
	if err != nil {
		panic(err)
	}

	// Sign those bytes
	sigBytes, err := privkey.Sign(bytesToSign)
	if err != nil {
		panic(err)
	}

	// Construct the SignatureV2 struct
	sigData = signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: sigBytes,
	}
	sig = signing.SignatureV2{
		PubKey:   privkey.PubKey(),
		Data:     &sigData,
		Sequence: sequence,
	}

	err = tx.SetSignatures(sig)
	if err != nil {
		panic(err)
	}

	// Get unsigned TX in JSON
	json, err := encodingConfig.TxConfig.TxJSONEncoder()(tx.GetTx())
	if err != nil {
		panic(err)
	}

	// Print the serialized transaction
	fmt.Println("Signed tx for " + chainId)
	fmt.Println(string(json))
}

func getAccount(addr sdk.AccAddress, encodingConfig params.EncodingConfig) (uint64, uint64, error) {
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
		return 0, 0, err
	}

	in := &authtypes.QueryAccountRequest{
		Address: addr.String(),
	}
	out := new(authtypes.QueryAccountResponse)

	err = grpcConn.Invoke(context.Background(), "/cosmos.auth.v1beta1.Query/Account", in, out, grpc.Header(&header))
	if err != nil {
		return 0, 0, err
	}

	var acc authtypes.AccountI
	if err := encodingConfig.InterfaceRegistry.UnpackAny(out.Account, &acc); err != nil {
		return 0, 0, err
	}

	return acc.GetAccountNumber(), acc.GetSequence(), nil
}
