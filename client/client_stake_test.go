package client

import (
	"context"

	"testing"

	"github.com/coming-chat/go-sui/account"
	"github.com/coming-chat/go-sui/types"
	"github.com/stretchr/testify/require"
)

func TestGetDelegatedStakes(t *testing.T) {
	cli := DevnetClient(t)

	acc := M1Account(t)
	addr, _ := types.NewAddressFromHex(acc.Address)

	stakes, err := cli.GetDelegatedStakes(context.Background(), *addr)
	require.Nil(t, err)

	systemState, err := cli.GetSuiSystemState(context.Background())
	require.Nil(t, err)

	for _, stake := range stakes {
		earn, validator := stake.CalculateEarnAmount(systemState.Validators.ActiveValidators)
		t.Logf("earned amount %10v at %v", earn, string(validator.Metadata.Name))
	}
}

func TestGetValidators(t *testing.T) {
	cli := DevnetClient(t)

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
	cli := DevnetClient(t)

	res, err := cli.GetSuiSystemState(context.Background())
	require.Nil(t, err)

	for _, v := range res.Validators.ActiveValidators {
		t.Logf("%v, %v\n", string(v.Metadata.Name), v.CalculateAPY(res.Epoch))
	}
}

func TestRequestAddDelegation(t *testing.T) {
	if false {
		coins := []string{
			"0x6619dfd5b9d955f7734ef1e680a222f40861302b", // 20000000 // 0.02
			"0x7b2bb7b953a9529aa8b6ca267d8d20aa5113d92e", // 20000000
		}
		amount := uint64(20000000 * 2) // all use
		validatorAddress := "0xc397cc2c83165c78a850a0295fb08538de9566f1"
		gasId := "0x3a1afb6e982cec8b9d5971f0e70fd53ae3e4cd4f"
		requestAddDelegation(t, coins, amount, validatorAddress, gasId)
		// ✅ https://explorer.sui.io/transaction/EcH9dK1wSLzgv15CNNHr2KvDEpARUW5mSeug3LHeFqGB?network=testnet
	}

	if false {
		coins := []string{
			"0x82ba777672c570959a4a25b076bf2b80dfefd367", // 20000000 // 0.02
			"0x838b3d559569a132f0f86d13883826651105ecc7", // 20000000
		}
		amount := uint64(30000000) // only use 0.03
		validatorAddress := "0xc397cc2c83165c78a850a0295fb08538de9566f1"
		gasId := "0x3a1afb6e982cec8b9d5971f0e70fd53ae3e4cd4f"
		requestAddDelegation(t, coins, amount, validatorAddress, gasId)
		// ✅ https://explorer.sui.io/transaction/F3VXATTbZU2VtrJv7V95FRkAQ4KAwBqr3mZHmcPRs3fC?network=testnet
	}

	if false {
		onlyOneCoin := "0x1179965d09c9d8b547e4c746c200022173336569" // 0.01
		coins := []string{
			onlyOneCoin,
		}
		amount := uint64(5000000) //0.005
		validatorAddress := "0xfb2bbe688390e83170aff94b96d317397940ad33"
		gasId := onlyOneCoin
		requestAddDelegation(t, coins, amount, validatorAddress, gasId)
		// Error: "Mutable object 0x1179....36569 cannot appear in more than one single transactions in a batch"
	}
}

func TestRequestSwitchDelegation(t *testing.T) {
	if false {
		delegationId := "0x7c6d82f689221d2763df7f4768be692bb684cb6f" // status.Id.Id
		stakedSuiId := "0x2b474dc1df5d995d2e45af2e16ac5f9fff57cf72"  // status.StakedSuiId
		newValidator := "0xfdc2a0fe740d34e5424d1a06c0c4ac106b7096f2"
		gasId := "0x3a1afb6e982cec8b9d5971f0e70fd53ae3e4cd4f"
		requestSwitchDelegation(t, delegationId, stakedSuiId, newValidator, gasId)
		// ✅ https://explorer.sui.io/transaction/AV5fTh8RwV4nMUunJkWuN5C7HgfosR8MHHAaKWTcT768?network=testnet
	}
}

func TestRequestWithdrawDelegation(t *testing.T) {
	if false {
		delegationId := "0x012f0a2c86035d5f31d375de877a1c20220260ef" // status.Id.Id
		stakedSuiId := "0x8ad8686c34668aca5d3efa041ca1a10986d332a6"  // status.StakedSuiId
		gasId := "0x3a1afb6e982cec8b9d5971f0e70fd53ae3e4cd4f"
		requestWithdrawDelegation(t, delegationId, stakedSuiId, gasId)
		// ✅ https://explorer.sui.io/transaction/DVpYBQ8Djg7jtHhcAwy8hgeWCjRpTaTYkNzpNigco5HP?network=testnet
	}
}

func requestAddDelegation(t *testing.T, coinIds []string, amount uint64, validatorAddress string, gasId string) {
	cli := DevnetClient(t)
	acc := M1Account(t)
	addr, _ := types.NewAddressFromHex(acc.Address)

	var coins = []types.ObjectId{}
	for _, id := range coinIds {
		coin, err := types.NewHexData(id)
		require.Nil(t, err)
		coins = append(coins, *coin)
	}

	validator, err := types.NewAddressFromHex(validatorAddress)
	require.Nil(t, err)

	gas, err := types.NewHexData(gasId)
	require.Nil(t, err)

	txn, err := cli.RequestAddDelegation(context.Background(), *addr, coins, amount, *validator, *gas, 20000)
	require.Nil(t, err)

	resp := simulateAndSendTxn(t, cli, txn, acc)
	t.Log(resp)
}

func requestSwitchDelegation(t *testing.T, delegationId, stakedId, validator, gasId string) {
	cli := DevnetClient(t)
	acc := M1Account(t)
	addr, _ := types.NewAddressFromHex(acc.Address)

	// queryed actived delegated stake
	// actived status = DelegatedStake.DelegationStatus (ActiveDelegationStatus)
	delegation, err := types.NewHexData(delegationId) // status.Id.Id
	require.Nil(t, err)
	stakedID, err := types.NewHexData(stakedId) // status.StakedSuiId
	require.Nil(t, err)

	newValidator, err := types.NewAddressFromHex(validator)
	require.Nil(t, err)

	gas, err := types.NewHexData(gasId)
	require.Nil(t, err)

	txn, err := cli.RequestSwitchDelegation(context.Background(), *addr, *delegation, *stakedID, *newValidator, *gas, 20000)
	require.Nil(t, err)

	resp := simulateAndSendTxn(t, cli, txn, acc)
	t.Log(resp)
}

func requestWithdrawDelegation(t *testing.T, delegationId, stakedId, gasId string) {
	cli := DevnetClient(t)
	acc := M1Account(t)
	addr, _ := types.NewAddressFromHex(acc.Address)

	// queryed actived delegated stake
	// actived status = DelegatedStake.DelegationStatus (ActiveDelegationStatus)
	delegation, err := types.NewHexData(delegationId) // status.Id.Id
	require.Nil(t, err)
	stakedID, err := types.NewHexData(stakedId) // status.StakedSuiId
	require.Nil(t, err)

	gas, err := types.NewHexData(gasId)
	require.Nil(t, err)

	txn, err := cli.RequestWithdrawDelegation(context.Background(), *addr, *delegation, *stakedID, *gas, 20000)
	require.Nil(t, err)

	resp := simulateAndSendTxn(t, cli, txn, acc)
	t.Log(resp)
}

func simulateAndSendTxn(t *testing.T, cli *Client, txn *types.TransactionBytes, acc *account.Account) *types.ExecuteTransactionResponse {
	// simulate
	simulate, err := cli.DryRunTransaction(context.Background(), txn)
	require.Nil(t, err)
	t.Log(simulate)

	require.Equal(t, simulate.Status.Status, types.ExecutionStatusSuccess)

	signedTxn := txn.SignSerializedSigWith(acc.PrivateKey)
	resp, err := cli.ExecuteTransactionSerializedSig(context.Background(), *signedTxn, types.TxnRequestTypeWaitForLocalExecution)
	require.Nil(t, err)
	return resp
}
