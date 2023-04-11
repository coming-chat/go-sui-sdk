package client

import (
	"context"
	"testing"

	"github.com/coming-chat/go-sui/sui_types"

	"github.com/coming-chat/go-sui/account"
	"github.com/coming-chat/go-sui/types"
	"github.com/stretchr/testify/require"
)

func shouldSimulate() bool {
	if whoami == "gg" {
		return true
	}
	return false
}

func shouldExecute() bool {
	if whoami == "gg" {
		return true
	}
	return false
}

const (
	M1Coin1 = "0x0501ebf5518912e380e8b3b68f93548418fb1bce59ed025f68ad6d236f012f92"
	M1Coin2 = "0x0d19d099213c23af5a6562034ce2772555f6945920913e18809369d738042b91"
)

func TestClient_TransferObject(t *testing.T) {
	cli := ChainClient(t)
	signer := M1Address(t)
	recipient := signer
	coins, err := cli.GetSuiCoinsOwnedByAddress(context.TODO(), *signer)
	require.NoError(t, err)
	coin, err := coins.PickCoinNoLess(20000)
	require.NoError(t, err)

	txnBytes, err := cli.TransferObject(
		context.Background(), *signer, *recipient, coin.CoinObjectId, nil,
		types.NewSafeSuiBigInt(uint64(100000000)),
	)
	require.Nil(t, err)

	t.Log(simulateCheck(t, cli, txnBytes))
}

func TestClient_TransferSui(t *testing.T) {
	cli := ChainClient(t)
	signer := M1Address(t)
	recipient := signer
	coins, err := cli.GetSuiCoinsOwnedByAddress(context.TODO(), *signer)
	require.NoError(t, err)
	coin, err := coins.PickCoinNoLess(100000)
	require.NoError(t, err)

	txnBytes, err := cli.TransferSui(
		context.Background(), *signer, *recipient, coin.CoinObjectId, types.NewSafeSuiBigInt(uint64(100000)),
		types.NewSafeSuiBigInt(uint64(100000000)),
	)
	require.Nil(t, err)

	t.Log(simulateCheck(t, cli, txnBytes))
}

func TestClient_PayAllSui(t *testing.T) {
	cli := ChainClient(t)
	signer := M1Address(t)
	recipient := signer
	coins, err := cli.GetSuiCoinsOwnedByAddress(context.TODO(), *signer)
	require.NoError(t, err)
	coin, err := coins.PickCoinNoLess(1000)
	require.NoError(t, err)
	coin1, err := coins.PickCoinNoLess(50000)
	require.NoError(t, err)

	txnBytes, err := cli.PayAllSui(
		context.Background(), *signer, *recipient, []types.ObjectId{coin.CoinObjectId, coin1.CoinObjectId},
		types.NewSafeSuiBigInt(uint64(100000000)),
	)
	require.Nil(t, err)

	t.Log(simulateCheck(t, cli, txnBytes))
}

func TestClient_Pay(t *testing.T) {
	cli := ChainClient(t)
	signer := M1Address(t)
	recipient := Address
	coins, err := cli.GetSuiCoinsOwnedByAddress(context.TODO(), *signer)
	require.NoError(t, err)
	coin, err := coins.PickCoinNoLess(10000)
	require.NoError(t, err)

	txnBytes, err := cli.Pay(
		context.Background(),
		*signer,
		[]types.ObjectId{coin.CoinObjectId},
		[]types.Address{*recipient},
		[]types.SafeSuiBigInt[uint64]{
			types.NewSafeSuiBigInt(uint64(10000)),
		},
		nil,
		types.NewSafeSuiBigInt(uint64(100000000)),
	)
	require.Nil(t, err)

	t.Log(simulateCheck(t, cli, txnBytes))
}

func TestClient_PaySui(t *testing.T) {
	cli := ChainClient(t)
	signer := M1Address(t)
	recipient := Address
	coins, err := cli.GetSuiCoinsOwnedByAddress(context.TODO(), *signer)
	require.NoError(t, err)
	coin, err := coins.PickCoinNoLess(10000)
	require.NoError(t, err)

	txnBytes, err := cli.PaySui(
		context.Background(),
		*signer,
		[]types.ObjectId{coin.CoinObjectId},
		[]types.Address{*recipient},
		[]types.SafeSuiBigInt[uint64]{
			types.NewSafeSuiBigInt(uint64(1000)),
		},
		types.NewSafeSuiBigInt(uint64(100000000)),
	)
	require.Nil(t, err)

	t.Log(simulateCheck(t, cli, txnBytes))
}

func TestClient_SplitCoin(t *testing.T) {
	cli := ChainClient(t)
	signer := M1Address(t)
	coins, err := cli.GetSuiCoinsOwnedByAddress(context.TODO(), *Address)
	require.NoError(t, err)
	coin, err := coins.PickCoinNoLess(1e7)
	require.NoError(t, err)
	splitCoins := []types.SafeSuiBigInt[uint64]{types.NewSafeSuiBigInt(uint64(1e7))} // 1SUI

	txnBytes, err := cli.SplitCoin(
		context.Background(), *signer, coin.CoinObjectId, splitCoins, nil,
		types.NewSafeSuiBigInt(uint64(100000000)),
	)
	require.Nil(t, err)

	simulateCheck(t, cli, txnBytes)
}

func TestClient_SplitCoinEqual(t *testing.T) {
	cli := ChainClient(t)
	signer := M1Address(t)
	coins, err := cli.GetSuiCoinsOwnedByAddress(context.TODO(), *Address)
	require.NoError(t, err)
	coin, err := coins.PickCoinNoLess(100000)
	require.Nil(t, err)

	txnBytes, err := cli.SplitCoinEqual(
		context.Background(), *signer, coin.CoinObjectId, types.NewSafeSuiBigInt(uint64(2)),
		nil, types.NewSafeSuiBigInt(uint64(100000000)),
	)
	require.Nil(t, err)

	t.Log(simulateCheck(t, cli, txnBytes))
}

func TestClient_MergeCoins(t *testing.T) {
	cli := ChainClient(t)
	signer := Address
	coins, err := cli.GetSuiCoinsOwnedByAddress(context.TODO(), *Address)
	require.NoError(t, err)
	coin1, err := coins.PickCoinNoLess(1000)
	require.NoError(t, err)
	coin2, err := coins.PickCoinNoLess(20000)
	require.NoError(t, err)

	txnBytes, err := cli.MergeCoins(
		context.Background(), *signer, coin1.CoinObjectId, coin2.CoinObjectId, nil,
		types.NewSafeSuiBigInt(uint64(100000000)),
	)
	require.Nil(t, err)

	t.Log(simulateCheck(t, cli, txnBytes))
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

// params & return like `func simulateAndSendTxn`
// @param acc Never use
// @return types.DryRunTransactionBlockResponse
func simulateCheck(
	t *testing.T,
	cli *Client,
	txn *types.TransactionBytes,
) *types.DryRunTransactionBlockResponse {
	simulate, err := cli.DryRunTransaction(context.Background(), txn)
	require.Nil(t, err)
	require.Equal(t, simulate.Effects.Data.V1.Status.Error, "")
	require.True(t, simulate.Effects.Data.IsSuccess())
	return simulate
}

func simulateAndSendTxn(
	t *testing.T,
	cli *Client,
	txn *types.TransactionBytes,
	acc *account.Account,
) *types.SuiTransactionBlockResponse {
	if shouldExecute() || shouldSimulate() {
		simulate, err := cli.DryRunTransaction(context.Background(), txn)
		require.Nil(t, err)
		require.True(t, simulate.Effects.Data.IsSuccess())
	}
	if shouldExecute() {
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
	return nil
}
