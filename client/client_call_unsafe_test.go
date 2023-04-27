package client

import (
	"context"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/coming-chat/go-sui/sui_types"

	"github.com/coming-chat/go-sui/account"
	"github.com/coming-chat/go-sui/types"
	"github.com/stretchr/testify/require"
)

func TestClient_TransferObject(t *testing.T) {
	cli := ChainClient(t)
	signer := M1Address(t)
	recipient := signer
	coins, err := cli.GetCoins(context.Background(), *signer, nil, nil, 10)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(coins.Data), 2)
	coin := coins.Data[0]

	txnBytes, err := cli.TransferObject(context.Background(), *signer, *recipient,
		coin.CoinObjectId, nil, types.NewSafeSuiBigInt(SUI(0.01).Uint64()),
	)
	require.Nil(t, err)

	simulateCheck(t, cli, txnBytes, true)
}

func TestClient_TransferSui(t *testing.T) {
	cli := ChainClient(t)
	signer := M1Address(t)
	recipient := signer
	coins, err := cli.GetCoins(context.Background(), *signer, nil, nil, 10)
	require.NoError(t, err)

	amount := SUI(0.0001).Uint64()
	gasBudget := SUI(0.01).Uint64()
	pickedCoins, err := types.PickupCoins(coins, *big.NewInt(0).SetUint64(amount + gasBudget), 1, false)
	require.Nil(t, err)

	txnBytes, err := cli.TransferSui(
		context.Background(), *signer, *recipient,
		pickedCoins.Coins[0].CoinObjectId,
		types.NewSafeSuiBigInt(amount),
		types.NewSafeSuiBigInt(gasBudget),
	)
	require.Nil(t, err)

	simulateCheck(t, cli, txnBytes, true)
}

func TestClient_PayAllSui(t *testing.T) {
	cli := ChainClient(t)
	signer := M1Address(t)
	recipient := signer
	coins, err := cli.GetCoins(context.Background(), *signer, nil, nil, 10)
	require.NoError(t, err)

	amount := SUI(0.001).Uint64()
	gasBudget := SUI(0.01).Uint64()
	pickedCoins, err := types.PickupCoins(coins, *big.NewInt(0).SetUint64(amount + gasBudget), 0, false)
	require.Nil(t, err)

	txnBytes, err := cli.PayAllSui(
		context.Background(), *signer, *recipient,
		pickedCoins.CoinIds(),
		types.NewSafeSuiBigInt(gasBudget),
	)
	require.Nil(t, err)

	simulateCheck(t, cli, txnBytes, true)
}

func TestClient_Pay(t *testing.T) {
	cli := ChainClient(t)
	signer := M1Address(t)
	recipient := Address
	coins, err := cli.GetCoins(context.Background(), *signer, nil, nil, 10)
	require.NoError(t, err)

	amount := SUI(0.001).Uint64()
	gasBudget := SUI(0.01).Uint64()
	pickedCoins, err := types.PickupCoins(coins, *big.NewInt(0).SetUint64(amount + gasBudget), 0, true)
	require.NoError(t, err)

	txnBytes, err := cli.Pay(
		context.Background(), *signer,
		pickedCoins.CoinIds(),
		[]types.Address{*recipient},
		[]types.SafeSuiBigInt[uint64]{
			types.NewSafeSuiBigInt(amount),
		},
		nil,
		types.NewSafeSuiBigInt(gasBudget),
	)
	require.Nil(t, err)

	simulateCheck(t, cli, txnBytes, true)
}

func TestClient_PaySui(t *testing.T) {
	cli := ChainClient(t)
	signer := M1Address(t)
	recipient := Address
	coins, err := cli.GetCoins(context.Background(), *signer, nil, nil, 10)
	require.NoError(t, err)

	amount := SUI(0.001).Uint64()
	gasBudget := SUI(0.01).Uint64()
	pickedCoins, err := types.PickupCoins(coins, *big.NewInt(0).SetUint64(amount + gasBudget), 0, false)
	require.NoError(t, err)

	txnBytes, err := cli.PaySui(
		context.Background(), *signer,
		pickedCoins.CoinIds(),
		[]types.Address{*recipient},
		[]types.SafeSuiBigInt[uint64]{
			types.NewSafeSuiBigInt(amount),
		},
		types.NewSafeSuiBigInt(gasBudget),
	)
	require.Nil(t, err)

	simulateCheck(t, cli, txnBytes, true)
}

func TestClient_SplitCoin(t *testing.T) {
	cli := ChainClient(t)
	signer := M1Address(t)
	coins, err := cli.GetCoins(context.Background(), *signer, nil, nil, 10)
	require.NoError(t, err)

	amount := SUI(0.01).Uint64()
	gasBudget := SUI(0.01).Uint64()
	pickedCoins, err := types.PickupCoins(coins, *big.NewInt(0).SetUint64(amount), 1, true)
	require.NoError(t, err)
	splitCoins := []types.SafeSuiBigInt[uint64]{types.NewSafeSuiBigInt(amount / 2)}

	txnBytes, err := cli.SplitCoin(
		context.Background(), *signer,
		pickedCoins.Coins[0].CoinObjectId,
		splitCoins,
		nil, types.NewSafeSuiBigInt(gasBudget),
	)
	require.Nil(t, err)

	simulateCheck(t, cli, txnBytes, false)
}

func TestClient_SplitCoinEqual(t *testing.T) {
	cli := ChainClient(t)
	signer := M1Address(t)
	coins, err := cli.GetCoins(context.Background(), *signer, nil, nil, 10)
	require.NoError(t, err)

	amount := SUI(0.01).Uint64()
	gasBudget := SUI(0.01).Uint64()
	pickedCoins, err := types.PickupCoins(coins, *big.NewInt(0).SetUint64(amount), 1, true)
	require.NoError(t, err)

	txnBytes, err := cli.SplitCoinEqual(
		context.Background(), *signer,
		pickedCoins.Coins[0].CoinObjectId,
		types.NewSafeSuiBigInt(uint64(2)),
		nil, types.NewSafeSuiBigInt(gasBudget),
	)
	require.Nil(t, err)

	simulateCheck(t, cli, txnBytes, true)
}

func TestClient_MergeCoins(t *testing.T) {
	cli := ChainClient(t)
	signer := Address
	coins, err := cli.GetCoins(context.Background(), *signer, nil, nil, 10)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(coins.Data), 3)

	coin1 := coins.Data[0]
	coin2 := coins.Data[1]
	coin3 := coins.Data[2] // gas coin

	txnBytes, err := cli.MergeCoins(
		context.Background(), *signer,
		coin1.CoinObjectId, coin2.CoinObjectId,
		&coin3.CoinObjectId, coin3.Balance,
	)
	require.Nil(t, err)

	simulateCheck(t, cli, txnBytes, true)
}

func TestClient_Publish(t *testing.T) {
	t.Log("TestClient_Publish TODO")
	// cli := ChainClient(t)

	// txnBytes, err := cli.Publish(context.Background(), *signer, *coin1, *coin2, nil, 10000)
	// require.Nil(t, err)
	// simulateCheck(t, cli, txnBytes, M1Account(t))
}

func TestClient_MoveCall(t *testing.T) {
	t.Log("TestClient_MoveCall TODO")
	// cli := ChainClient(t)

	// txnBytes, err := cli.MoveCall(context.Background(), *signer, *coin1, *coin2, nil, 10000)
	// require.Nil(t, err)
	// simulateCheck(t, cli, txnBytes, M1Account(t))
}

func TestClient_BatchTransaction(t *testing.T) {
	t.Log("TestClient_BatchTransaction TODO")
	// cli := ChainClient(t)

	// txnBytes, err := cli.BatchTransaction(context.Background(), *signer, *coin1, *coin2, nil, 10000)
	// require.Nil(t, err)
	// simulateCheck(t, cli, txnBytes, M1Account(t))
}

// @return types.DryRunTransactionBlockResponse
func simulateCheck(
	t *testing.T,
	cli *Client,
	txn *types.TransactionBytes,
	showJson bool,
) *types.DryRunTransactionBlockResponse {
	simulate, err := cli.DryRunTransaction(context.Background(), txn)
	require.Nil(t, err)
	require.Equal(t, simulate.Effects.Data.V1.Status.Error, "")
	require.True(t, simulate.Effects.Data.IsSuccess())
	if showJson {
		data, err := json.Marshal(simulate)
		require.Nil(t, err)
		t.Log(string(data))
	}
	return simulate
}

func executeTxn(
	t *testing.T,
	cli *Client,
	txn *types.TransactionBytes,
	acc *account.Account,
) *types.SuiTransactionBlockResponse {
	// First of all, make sure that there are no problems with simulated trading.
	simulate, err := cli.DryRunTransaction(context.Background(), txn)
	require.Nil(t, err)
	require.True(t, simulate.Effects.Data.IsSuccess())

	// sign and send
	signature, err := acc.SignSecureWithoutEncode(txn.TxBytes, sui_types.DefaultIntent())
	require.NoError(t, err)
	options := types.SuiTransactionBlockResponseOptions{
		ShowEffects: true,
	}
	resp, err := cli.ExecuteTransactionBlock(
		context.TODO(), txn.TxBytes, []any{signature}, &options,
		types.TxnRequestTypeWaitForLocalExecution,
	)
	require.NoError(t, err)
	t.Log(resp)
	return resp
}
