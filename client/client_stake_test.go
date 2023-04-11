package client

import (
	"context"
	"testing"

	"github.com/coming-chat/go-sui/types"
	"github.com/stretchr/testify/require"
)

func TestClient_GetLatestSuiSystemState(t *testing.T) {
	cli := ChainClient(t)
	state, err := cli.GetLatestSuiSystemState(context.Background())
	require.Nil(t, err)
	t.Logf("system state = %v", state)

	for _, v := range state.ActiveValidators {
		t.Logf("%v, %v\n", v.Name, v.CalculateAPY(state.Epoch.Uint64()))
	}
}

func TestClient_GetAndCalculateRollingAverageApys(t *testing.T) {
	cli := ChainClient(t)
	apys, err := cli.GetAndCalculateRollingAverageApys(context.Background(), 98)
	require.Nil(t, err)
	for address, apy := range apys {
		t.Logf("%v apy = %v", address, apy)
	}
}

func TestGetDelegatedStakes(t *testing.T) {
	cli := ChainClient(t)

	address, err := types.NewAddressFromHex("0xd77955e670f42c1bc5e94b9e68e5fe9bdbed9134d784f2a14dfe5fc1b24b5d9f")
	require.Nil(t, err)
	stakes, err := cli.GetStakes(context.Background(), *address)
	require.Nil(t, err)

	for _, validator := range stakes {
		for _, stake := range validator.Stakes {
			if stake.Data.StakeStatus.Data.Active != nil {
				t.Logf(
					"earned amount %10v at %v",
					stake.Data.StakeStatus.Data.Active.EstimatedReward.Uint64(),
					validator.ValidatorAddress,
				)
			}
		}
	}
}

//func TestGetStakesByIds(t *testing.T) {
//	cli := ChainClient(t)
//
//	id1, _ := types.NewHexData("0x0e32ab08fe29b830ca2c04266297fe121128bf77d380ebec3256a4e1734144aa")
//	stakes, err := cli.GetStakesByIds(context.Background(), []types.ObjectId{*id1})
//	require.Nil(t, err)
//
//	for _, validator := range stakes {
//		for _, stake := range validator.Stakes {
//			t.Logf("earned amount %10v at %v", *stake.EstimatedReward, validator.ValidatorAddress)
//		}
//	}
//}

func TestRequestAddDelegation(t *testing.T) {
	if true {
		coins := []string{
			"0x0153883d60e0df7052b12bc04454dd2eec1c3723ee12145ca73522c6a3917523",
			"0x21d6e05e77325cbca6bf73410763b216c5614a1184c0efc414de68ebb80b842b",
		}
		amount := SUI(1).Decimal()
		validatorAddress := "0x8ce890590fed55c37d44a043e781ad94254b413ee079a53fb5c037f7a6311304"
		// gasId := "0x11ce8b45348f6db3f46a8a54a5d06ab91d8381bbc3cb67d66bef8c7ce2b5a7c5"

		requestAddDelegation(t, coins, amount, validatorAddress)
	}
}

func requestAddDelegation(t *testing.T, coinIds []string, amount types.SuiBigInt, validatorAddress string) {
	cli := ChainClient(t)
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

	gasBudget := SUI(0.02).Decimal()
	txn, err := cli.RequestAddStake(context.Background(), *addr, coins, amount, *validator, nil, gasBudget)
	require.Nil(t, err)

	resp := simulateCheck(t, cli, txn)
	t.Log(resp)
}

func TestRequestWithdrawDelegation(t *testing.T) {
	if true {
		stakedSuiId := "0x0e32ab08fe29b830ca2c04266297fe121128bf77d380ebec3256a4e1734144aa"
		requestWithdrawDelegation(t, stakedSuiId, "")
	}
}

func requestWithdrawDelegation(t *testing.T, stakedId, gasId string) {
	cli := ChainClient(t)
	acc := M1Account(t)
	addr, _ := types.NewAddressFromHex(acc.Address)

	stakedID, err := types.NewHexData(stakedId) // status.StakedSuiId
	require.Nil(t, err)

	gasBudget := SUI(0.02).Decimal()
	txn, err := cli.RequestWithdrawStake(context.Background(), *addr, *stakedID, nil, gasBudget)
	require.Nil(t, err)

	resp := simulateCheck(t, cli, txn)
	t.Log(resp)
}
