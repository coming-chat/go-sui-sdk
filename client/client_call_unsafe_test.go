package client

import (
	"context"
	"testing"

	"github.com/coming-chat/go-sui/sui_types"

	"github.com/coming-chat/go-sui/account"
	"github.com/coming-chat/go-sui/types"
	"github.com/shopspring/decimal"
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
	cli := DevnetClient(t)
	signer := M1Address(t)
	recipient := signer
	objId, err := types.NewHexData(M1Coin1)
	require.NoError(t, err)

	txnBytes, err := cli.TransferObject(context.Background(), *signer, *recipient, *objId, nil, 10000)
	require.Nil(t, err)

	simulateCheck(t, cli, txnBytes, M1Account(t))
}

func TestClient_TransferSui(t *testing.T) {
	cli := DevnetClient(t)
	signer := M1Address(t)
	recipient := signer
	objId, err := types.NewHexData(M1Coin1)
	require.NoError(t, err)

	txnBytes, err := cli.TransferSui(context.Background(), *signer, *recipient, *objId, 100000, 10000)
	require.Nil(t, err)

	simulateCheck(t, cli, txnBytes, M1Account(t))
}

func TestClient_PayAllSui(t *testing.T) {
	cli := DevnetClient(t)
	signer := M1Address(t)
	recipient := signer
	objId, err := types.NewHexData(M1Coin1)
	require.NoError(t, err)
	coin2, err := types.NewHexData(M1Coin2)
	require.NoError(t, err)

	txnBytes, err := cli.PayAllSui(context.Background(), *signer, *recipient, []types.ObjectId{*objId, *coin2}, 10000)
	require.Nil(t, err)

	simulateCheck(t, cli, txnBytes, M1Account(t))
}

func TestClient_Pay(t *testing.T) {
	cli := DevnetClient(t)
	signer := M1Address(t)
	recipient := Address
	objId, err := types.NewHexData(M1Coin1)
	require.NoError(t, err)

	amount := decimal.NewFromInt(10000)
	txnBytes, err := cli.Pay(
		context.Background(),
		*signer,
		[]types.ObjectId{*objId},
		[]types.Address{*recipient},
		[]decimal.Decimal{amount},
		nil,
		10000,
	)
	require.Nil(t, err)

	simulateCheck(t, cli, txnBytes, M1Account(t))
}

func TestClient_PaySui(t *testing.T) {
	cli := DevnetClient(t)
	signer := M1Address(t)
	recipient := Address
	objId, err := types.NewHexData(M1Coin1)
	require.NoError(t, err)

	amount := decimal.NewFromInt(10000)
	txnBytes, err := cli.PaySui(
		context.Background(),
		*signer,
		[]types.ObjectId{*objId},
		[]types.Address{*recipient},
		[]decimal.Decimal{amount},
		10000,
	)
	require.Nil(t, err)

	simulateCheck(t, cli, txnBytes, M1Account(t))
}

func TestClient_SplitCoin(t *testing.T) {
	cli := DevnetClient(t)
	signer := M1Address(t)
	objId, err := types.NewHexData(M1Coin1)
	require.NoError(t, err)
	splitCoins := []uint64{1e9} // 1SUI

	txnBytes, err := cli.SplitCoin(context.Background(), *signer, *objId, splitCoins, nil, 10000)
	require.Nil(t, err)

	simulateCheck(t, cli, txnBytes, M1Account(t))
}

func TestClient_SplitCoinEqual(t *testing.T) {
	cli := DevnetClient(t)
	signer := M1Address(t)
	objId, err := types.NewHexData(M1Coin1)
	require.NoError(t, err)

	txnBytes, err := cli.SplitCoinEqual(context.Background(), *signer, *objId, 2, nil, 10000)
	require.Nil(t, err)

	// simulateCheck(t, cli, txnBytes, M1Account(t))
	simulateAndSendTxn(t, cli, txnBytes, M1Account(t))
}

func TestClient_MergeCoins(t *testing.T) {
	cli := DevnetClient(t)
	signer := Address
	coin1, err := types.NewHexData(M1Coin1)
	require.NoError(t, err)
	coin2, err := types.NewHexData(M1Coin2)
	require.NoError(t, err)

	txnBytes, err := cli.MergeCoins(context.Background(), *signer, *coin1, *coin2, nil, 10000)
	require.Nil(t, err)

	simulateCheck(t, cli, txnBytes, M1Account(t))
}

func TestClient_Publish(t *testing.T) {
	t.Log("TestClient_Publish TODO")
	// cli := DevnetClient(t)

	// txnBytes, err := cli.Publish(context.Background(), *signer, *coin1, *coin2, nil, 10000)
	// require.Nil(t, err)
	// simulateCheck(t, cli, txnBytes, M1Account(t))
}

func TestClient_MoveCall(t *testing.T) {
	t.Log("TestClient_MoveCall TODO")
	// cli := DevnetClient(t)

	// txnBytes, err := cli.MoveCall(context.Background(), *signer, *coin1, *coin2, nil, 10000)
	// require.Nil(t, err)
	// simulateCheck(t, cli, txnBytes, M1Account(t))
}

func TestClient_BatchTransaction(t *testing.T) {
	t.Log("TestClient_BatchTransaction TODO")
	// cli := DevnetClient(t)

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
	acc *account.Account,
) *types.DryRunTransactionBlockResponse {
	if shouldSimulate() {
		simulate, err := cli.DryRunTransaction(context.Background(), txn)
		require.Nil(t, err)
		require.True(t, simulate.Effects.IsSuccess())
		return simulate
	}
	return nil
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
		require.True(t, simulate.Effects.IsSuccess())
	}
	if shouldExecute() {
		signature, err := acc.SignSecure(txn.TxBytes, sui_types.DefaultIntent())
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
