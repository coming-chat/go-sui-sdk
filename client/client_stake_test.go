package client

import (
	"context"

	"testing"

	"github.com/coming-chat/go-sui/types"
	"github.com/stretchr/testify/require"
)

func TestGetDelegatedStakes(t *testing.T) {
	cli := TestnetClient(t)

	acc := M1Account(t)
	addr, _ := types.NewAddressFromHex(acc.Address)

	res, err := cli.GetDelegatedStakes(context.Background(), *addr)
	require.Nil(t, err)
	t.Log(res)
}

func TestGetValidators(t *testing.T) {
	cli := TestnetClient(t)

	res, err := cli.GetValidators(context.Background())
	require.Nil(t, err)
	// t.Log(res)
	for _, validator := range res {
		if string(validator.Name) == "Chainode Tech" {
			t.Log(validator)
		}
	}
}

func TestGetSuiSystemState(t *testing.T) {
	cli := TestnetClient(t)

	res, err := cli.GetSuiSystemState(context.Background())
	require.Nil(t, err)

	for _, v := range res.Validators.ActiveValidators {
		t.Logf("%v, %v\n", string(v.Metadata.Name), v.CalculateAPY(res.Epoch))
	}
}

func TestRequestAddDelegation(t *testing.T) {
	cli := TestnetClient(t)

	acc := M1Account(t)
	addr, _ := types.NewAddressFromHex(acc.Address)

	coin1 := "0x6619dfd5b9d955f7734ef1e680a222f40861302b" // 0.02
	coin2 := "0x7b2bb7b953a9529aa8b6ca267d8d20aa5113d92e" // 0.02 // 20000000
	c1, _ := types.NewHexData(coin1)
	c2, _ := types.NewHexData(coin2)
	coins := []types.ObjectId{*c1, *c2}
	amount := uint64(20000000 * 2)

	validatorAddress := "0xc397cc2c83165c78a850a0295fb08538de9566f1"
	validator, _ := types.NewAddressFromHex(validatorAddress)

	gasId := "0x3a1afb6e982cec8b9d5971f0e70fd53ae3e4cd4f"
	gas, _ := types.NewHexData(gasId)

	txn, err := cli.RequestAddDelegation(context.Background(), *addr, coins, amount, *validator, *gas, 20000)
	require.Nil(t, err)

	// simulate
	simulate, err := cli.DryRunTransaction(context.Background(), txn)
	require.Nil(t, err)
	t.Log(simulate)

	require.Equal(t, simulate.Status.Status, types.TransactionStatusSuccess)

	signedTxn := txn.SignSerializedSigWith(acc.PrivateKey)
	resp, err := cli.ExecuteTransactionSerializedSig(context.Background(), *signedTxn, types.TxnRequestTypeWaitForLocalExecution)
	require.Nil(t, err)
	t.Log(resp)
}

func TestRequestSwitchDelegation(t *testing.T) {
	cli := TestnetClient(t)

	acc := M1Account(t)
	addr, _ := types.NewAddressFromHex(acc.Address)

	// queryed actived delegated stake
	// actived status = DelegatedStake.DelegationStatus (ActiveDelegationStatus)
	delegationId := "0x7c6d82f689221d2763df7f4768be692bb684cb6f" // status.Id.Id
	delegation, _ := types.NewHexData(delegationId)
	stakedSuiId := "0x2b474dc1df5d995d2e45af2e16ac5f9fff57cf72" // status.StakedSuiId
	stakedID, _ := types.NewHexData(stakedSuiId)

	newValidator := "0xfdc2a0fe740d34e5424d1a06c0c4ac106b7096f2"
	newAddr, _ := types.NewAddressFromHex(newValidator)

	gasId := "0x3a1afb6e982cec8b9d5971f0e70fd53ae3e4cd4f"
	gas, _ := types.NewHexData(gasId)

	txn, err := cli.RequestSwitchDelegation(context.Background(), *addr, *delegation, *stakedID, *newAddr, *gas, 20000)
	require.Nil(t, err)

	// simulate
	simulate, err := cli.DryRunTransaction(context.Background(), txn)
	require.Nil(t, err)
	require.Equal(t, simulate.Status.Status, types.TransactionStatusSuccess)

	signedTxn := txn.SignSerializedSigWith(acc.PrivateKey)
	resp, err := cli.ExecuteTransactionSerializedSig(context.Background(), *signedTxn, types.TxnRequestTypeWaitForLocalExecution)
	require.Nil(t, err)
	t.Log(resp)
}

func TestXxx(t *testing.T) {
	cli := TestnetClient(t)

	acc := M1Account(t)
	addr, _ := types.NewAddressFromHex(acc.Address)

	// queryed actived delegated stake
	// actived status = DelegatedStake.DelegationStatus (ActiveDelegationStatus)
	delegationId := "0x012f0a2c86035d5f31d375de877a1c20220260ef" // status.Id.Id
	delegation, _ := types.NewHexData(delegationId)
	stakedSuiId := "0x8ad8686c34668aca5d3efa041ca1a10986d332a6" // status.StakedSuiId
	stakedID, _ := types.NewHexData(stakedSuiId)

	gasId := "0x3a1afb6e982cec8b9d5971f0e70fd53ae3e4cd4f"
	gas, _ := types.NewHexData(gasId)

	txn, err := cli.RequestWithdrawDelegation(context.Background(), *addr, *delegation, *stakedID, *gas, 20000)
	require.Nil(t, err)

	// simulate
	simulate, err := cli.DryRunTransaction(context.Background(), txn)
	require.Nil(t, err)
	require.Equal(t, simulate.Status.Status, types.TransactionStatusSuccess)

	signedTxn := txn.SignSerializedSigWith(acc.PrivateKey)
	resp, err := cli.ExecuteTransactionSerializedSig(context.Background(), *signedTxn, types.TxnRequestTypeWaitForLocalExecution)
	require.Nil(t, err)
	t.Log(resp)
}
