package client

import (
	"context"
	"math/big"
	"testing"

	"github.com/coming-chat/go-sui/sui_types"
	"github.com/coming-chat/go-sui/types"
	"github.com/stretchr/testify/require"
)

const (
	ComingChatValidatorAddress = "0x520289e77c838bae8501ae92b151b99a54407288fdd20dee6e5416bfe943eb7a"
)

func TestClient_GetLatestSuiSystemState(t *testing.T) {
	cli := MainnetClient(t)
	state, err := cli.GetLatestSuiSystemState(context.Background())
	require.Nil(t, err)
	t.Logf("system state = %v", state)
}

func TestClient_GetValidatorsApy(t *testing.T) {
	cli := ChainClient(t)
	apys, err := cli.GetValidatorsApy(context.Background())
	require.Nil(t, err)
	t.Logf("current epoch %v", apys.Epoch)
	apyMap := apys.ApyMap()
	for idx := 0; idx < 10; idx++ {
		key := apys.Apys[idx].Address
		t.Logf("%v apy = %v", key, apyMap[key])
	}
}

func TestGetDelegatedStakes(t *testing.T) {
	cli := ChainClient(t)

	address, err := sui_types.NewAddressFromHex("0xd77955e670f42c1bc5e94b9e68e5fe9bdbed9134d784f2a14dfe5fc1b24b5d9f")
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

func TestGetStakesByIds(t *testing.T) {
	cli := TestnetClient(t)
	owner, err := sui_types.NewAddressFromHex("0xd77955e670f42c1bc5e94b9e68e5fe9bdbed9134d784f2a14dfe5fc1b24b5d9f")
	stakes, err := cli.GetStakes(context.Background(), *owner)
	require.Nil(t, err)
	require.GreaterOrEqual(t, len(stakes), 1)

	stake1 := stakes[0].Stakes[0].Data
	stakeId := stake1.StakedSuiId
	stakesFromId, err := cli.GetStakesByIds(context.Background(), []suiObjectID{stakeId})
	require.Nil(t, err)
	require.GreaterOrEqual(t, len(stakesFromId), 1)

	queriedStake := stakesFromId[0].Stakes[0].Data
	require.Equal(t, stake1, queriedStake)
	t.Log(stakesFromId)
}

func TestRequestAddDelegation(t *testing.T) {
	cli := TestnetClient(t)
	signer := Address

	coins, err := cli.GetCoins(context.Background(), *signer, nil, nil, 10)
	require.Nil(t, err)

	amount := SUI(1).Uint64()
	gasBudget := SUI(0.01).Uint64()
	pickedCoins, err := types.PickupCoins(coins, *big.NewInt(0).SetUint64(amount), 0, 0, 0)
	require.Nil(t, err)

	validatorAddress := ComingChatValidatorAddress
	validator, err := sui_types.NewAddressFromHex(validatorAddress)
	require.Nil(t, err)

	txBytes, err := BCS_RequestAddStake(*signer,
		pickedCoins.CoinRefs(),
		types.NewSafeSuiBigInt(amount),
		*validator,
		1000,
		gasBudget,
	)
	require.Nil(t, err)

	simulateCheck(t, cli, txBytes, true)
}

func TestRequestWithdrawDelegation(t *testing.T) {
	cli := TestnetClient(t)

	signer, err := sui_types.NewAddressFromHex("0xd77955e670f42c1bc5e94b9e68e5fe9bdbed9134d784f2a14dfe5fc1b24b5d9f")
	require.Nil(t, err)
	stakes, err := cli.GetStakes(context.Background(), *signer)
	require.Nil(t, err)
	require.True(t, len(stakes) > 0)
	require.True(t, len(stakes[0].Stakes) > 0)

	stakeId := stakes[0].Stakes[0].Data.StakedSuiId

	gasBudget := SUI(0.02).Decimal()
	txn, err := cli.RequestWithdrawStake(context.Background(), *signer, stakeId, nil, gasBudget)
	require.Nil(t, err)

	simulateCheck(t, cli, txn.TxBytes, true)
}
