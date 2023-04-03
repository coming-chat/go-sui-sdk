package client

import (
	"context"
	"testing"

	"github.com/coming-chat/go-sui/types"
	"github.com/stretchr/testify/require"
)

func TestClient_GetLatestSuiSystemState(t *testing.T) {
	cli := TestnetClient(t)
	state, err := cli.GetLatestSuiSystemState(context.Background())
	require.Nil(t, err)
	t.Logf("system state = %v", state)

	for _, v := range state.ActiveValidators {
		t.Logf("%v, %v\n", v.Name, v.CalculateAPY(state.Epoch))
	}
}

func TestGetDelegatedStakes(t *testing.T) {
	cli := DevnetClient(t)

	stakes, err := cli.GetStakes(context.Background(), *M1Address(t))
	require.Nil(t, err)

	for _, validator := range stakes {
		for _, stake := range validator.Stakes {
			t.Logf("earned amount %10v at %v", *stake.EstimatedReward, validator.ValidatorAddress)
		}
	}
}

func TestGetStakesByIds(t *testing.T) {
	cli := DevnetClient(t)

	id1, _ := types.NewHexData("0x4ad2f0a918a241d6a19573212aeb56947bb9255a14e921a7ec78b262536826f0")
	stakes, err := cli.GetStakesByIds(context.Background(), []types.ObjectId{*id1})
	require.Nil(t, err)

	for _, validator := range stakes {
		for _, stake := range validator.Stakes {
			t.Logf("earned amount %10v at %v", *stake.EstimatedReward, validator.ValidatorAddress)
		}
	}
}

func TestRequestAddDelegation(t *testing.T) {
	if true {
		coins := []string{
			"0x0501ebf5518912e380e8b3b68f93548418fb1bce59ed025f68ad6d236f012f92",
			"0x0d19d099213c23af5a6562034ce2772555f6945920913e18809369d738042b91",
		}
		amount := uint64(1000000000) // 1 SUI
		validatorAddress := "0x8ce890590fed55c37d44a043e781ad94254b413ee079a53fb5c037f7a6311304"
		// gasId := "0x11ce8b45348f6db3f46a8a54a5d06ab91d8381bbc3cb67d66bef8c7ce2b5a7c5"

		requestAddDelegation(t, coins, amount, validatorAddress)
		// ✅ https://explorer.sui.io/transaction/EcH9dK1wSLzgv15CNNHr2KvDEpARUW5mSeug3LHeFqGB?network=testnet
	}
}

func requestAddDelegation(t *testing.T, coinIds []string, amount uint64, validatorAddress string) {
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

	gasId := "0x11ce8b45348f6db3f46a8a54a5d06ab91d8381bbc3cb67d66bef8c7ce2b5a7c5"
	gas, _ := types.NewHexData(gasId)
	txn, err := cli.RequestAddStake(context.Background(), *addr, coins, amount, *validator, gas, 20000)
	require.Nil(t, err)

	resp := simulateCheck(t, cli, txn, acc)
	t.Log(resp)
}

func TestRequestWithdrawDelegation(t *testing.T) {
	if true {
		stakedSuiId := "0x4ad2f0a918a241d6a19573212aeb56947bb9255a14e921a7ec78b262536826f0"
		requestWithdrawDelegation(t, stakedSuiId, "")
		// ✅ https://explorer.sui.io/transaction/DVpYBQ8Djg7jtHhcAwy8hgeWCjRpTaTYkNzpNigco5HP?network=testnet
	}
}

func requestWithdrawDelegation(t *testing.T, stakedId, gasId string) {
	cli := DevnetClient(t)
	acc := M1Account(t)
	addr, _ := types.NewAddressFromHex(acc.Address)

	stakedID, err := types.NewHexData(stakedId) // status.StakedSuiId
	require.Nil(t, err)

	// gas, err := types.NewHexData(gasId)
	// require.Nil(t, err)

	txn, err := cli.RequestWithdrawStake(context.Background(), *addr, *stakedID, nil, 20000)
	require.Nil(t, err)

	resp := simulateCheck(t, cli, txn, acc)
	t.Log(resp)
}
