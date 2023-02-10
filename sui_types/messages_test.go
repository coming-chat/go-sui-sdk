package sui_types

import (
	"bytes"
	"context"
	"encoding/hex"
	"github.com/coming-chat/go-sui/client"
	"github.com/coming-chat/go-sui/types"
	"github.com/fardream/go-bcs/bcs"
	"github.com/stretchr/testify/require"
	"math/big"
	"os"
	"testing"
)

var (
	Address, _ = types.NewAddressFromHex("0xb08ae4d187ca0057baa1666fe43fb9d7f3693a9a")
	M1Mnemonic = os.Getenv("WalletSdkTestM1")
)

func Test_BCSEncodeTransactionData(t *testing.T) {
	chain, err := client.Dial(types.DevNetRpcUrl)
	require.NoError(t, err)
	coins, err := chain.GetSuiCoinsOwnedByAddress(context.TODO(), *Address)
	require.NoError(t, err)
	transferCoins, coin, err := coins.PickSUICoinsWithGas(big.NewInt(1000), 1000, types.PickByOrder)
	require.NoError(t, err)
	var (
		coinRef = []*types.ObjectRef{coin.Reference}
		coinId  = []types.ObjectId{coin.Reference.ObjectId}
	)
	for _, v := range transferCoins {
		coinRef = append(coinRef, v.Reference)
		coinId = append(coinId, v.Reference.ObjectId)
	}
	tx := TransactionData{
		Kind: TransactionKind{
			Single: &SingleTransactionKind{
				PayAllSui: &PayAllSui{
					Coins:     coinRef,
					Recipient: *Address,
				},
			},
		},
		Sender:     *Address,
		GasPayment: *coin.Reference,
		GasPrice:   uint64(1),
		GasBudget:  uint64(1000),
	}
	encodeTx, err := bcs.Marshal(tx)

	currentTxEncode, err := chain.PayAllSui(context.TODO(), *Address, *Address, coinId, 1000)
	require.NoError(t, err)
	t.Logf("%x", encodeTx)
	t.Logf("%x", currentTxEncode.TxBytes.Data())
	if !bytes.Equal(encodeTx, currentTxEncode.TxBytes.Data()) {
		t.Fatal("encode failed")
	}
}

func TestBCS_EncodeMoveCall(t *testing.T) {
	chain, err := client.Dial(types.TestnetRpcUrl)
	require.NoError(t, err)
	coins, err := chain.GetSuiCoinsOwnedByAddress(context.TODO(), *Address)
	require.NoError(t, err)
	coin, err := coins.PickCoinNoLess(2000)
	require.NoError(t, err)
	packageId, err := types.NewHexData("0x20145a5630cd30756acc0b22f500410066f748e5")
	require.NoError(t, err)
	packageRead, err := chain.GetObject(context.TODO(), *packageId)
	require.NoError(t, err)
	globalProfile, err := types.NewHexData("0x6d36f1275c1ec0f459fa3919c9baadaada39b068")
	require.NoError(t, err)
	globalProfileObject, err := chain.GetObject(context.TODO(), *globalProfile)
	require.NoError(t, err)
	profile := "{\"name\":\"test\",\"bio\":\"Hello\",\"avatar\":\"\"}"
	signature, err := hex.DecodeString("d485020c6ac369e6f2b28be2dcca24ebfd827c53893b6462e9e65cf16dba3cedf004e8740b8c8c3579a4391269b9e103bcfc39627c6af729abb7675bc8004301")
	require.NoError(t, err)
	args := MoveCallArg{
		&SharedObject{
			Id:                   *globalProfile,
			InitialSharedVersion: globalProfileObject.Details.Owner.Shared.InitialSharedVersion,
		},
		[]byte(profile),
		signature,
	}
	callArgs, err := args.GetMoveCallArgs()
	require.NoError(t, err)
	state, err := chain.GetSuiSystemState(context.TODO())
	require.NoError(t, err)
	tx := TransactionData{
		Kind: TransactionKind{
			Single: &SingleTransactionKind{
				Call: &MoveCall{
					Package:       *packageRead.Details.Reference,
					Module:        "profile",
					Function:      "register",
					TypeArguments: []*TypeTag{},
					Arguments:     callArgs,
				},
			},
		},
		Sender:     *Address,
		GasPayment: *coin.Reference,
		GasPrice:   state.ReferenceGasPrice,
		GasBudget:  uint64(2000),
	}
	encodeTx, err := bcs.Marshal(tx)
	require.NoError(t, err)
	currentTxEncode, err := chain.MoveCall(context.TODO(), *Address, *packageId, "profile", "register", []string{}, []any{globalProfile, profile, "0xd485020c6ac369e6f2b28be2dcca24ebfd827c53893b6462e9e65cf16dba3cedf004e8740b8c8c3579a4391269b9e103bcfc39627c6af729abb7675bc8004301"}, &coin.Reference.ObjectId, 2000)
	require.NoError(t, err)
	t.Logf("%x", encodeTx)
	t.Logf("%x", currentTxEncode.TxBytes.Data())
	if !bytes.Equal(encodeTx, currentTxEncode.TxBytes.Data()) {
		t.Fatal("encode failed")
	}
}
