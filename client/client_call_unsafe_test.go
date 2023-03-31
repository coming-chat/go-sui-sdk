package client

import (
	"context"
	"testing"

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

func TestClient_TransferObject(t *testing.T) {
	cli := DevnetClient(t)
	signer := Address
	recipient := Address
	objId, err := types.NewHexData("0x11462c88e74bb00079e3c043efb664482ee4551744ee691c7623b98503cb3f4d") // 0.2 SUI
	require.NoError(t, err)

	txnBytes, err := cli.TransferObject(context.Background(), *signer, *recipient, *objId, nil, 10000)
	require.Nil(t, err)

	simulateCheck(t, cli, txnBytes, M1Account(t))
}

func TestClient_TransferSui(t *testing.T) {
	cli := DevnetClient(t)
	signer := Address
	recipient := Address
	objId, err := types.NewHexData("0x11462c88e74bb00079e3c043efb664482ee4551744ee691c7623b98503cb3f4d") // 0.2 SUI
	require.NoError(t, err)

	txnBytes, err := cli.TransferSui(context.Background(), *signer, *recipient, *objId, 100000, 10000)
	require.Nil(t, err)

	simulateCheck(t, cli, txnBytes, M1Account(t))
}

func TestClient_PayAllSui(t *testing.T) {
	cli := DevnetClient(t)
	signer := Address
	recipient := Address
	objId, err := types.NewHexData("0x11462c88e74bb00079e3c043efb664482ee4551744ee691c7623b98503cb3f4d") // 0.2 SUI
	require.NoError(t, err)
	coin2, err := types.NewHexData("0x0fe1d3981da6954ed97e98e715ba41c54e67e4b461b2420a1e082b93d5700871") // 0.2 SUI
	require.NoError(t, err)

	txnBytes, err := cli.PayAllSui(context.Background(), *signer, *recipient, []types.ObjectId{*objId, *coin2}, 10000)
	require.Nil(t, err)

	simulateCheck(t, cli, txnBytes, M1Account(t))
}

func TestClient_Pay(t *testing.T) {
	cli := DevnetClient(t)
	signer := Address
	recipient := Address
	objId, err := types.NewHexData("0x11462c88e74bb00079e3c043efb664482ee4551744ee691c7623b98503cb3f4d") // 0.2 SUI
	require.NoError(t, err)

	amount := decimal.NewFromInt(10000)
	txnBytes, err := cli.Pay(context.Background(), *signer, []types.ObjectId{*objId}, []types.Address{*recipient}, []decimal.Decimal{amount}, nil, 10000)
	require.Nil(t, err)

	simulateCheck(t, cli, txnBytes, M1Account(t))
}

func TestClient_PaySui(t *testing.T) {
	cli := DevnetClient(t)
	signer := Address
	recipient := Address
	objId, err := types.NewHexData("0x11462c88e74bb00079e3c043efb664482ee4551744ee691c7623b98503cb3f4d") // 0.2 SUI
	require.NoError(t, err)

	amount := decimal.NewFromInt(10000)
	txnBytes, err := cli.PaySui(context.Background(), *signer, []types.ObjectId{*objId}, []types.Address{*recipient}, []decimal.Decimal{amount}, 10000)
	require.Nil(t, err)

	simulateCheck(t, cli, txnBytes, M1Account(t))
}

func TestClient_SplitCoin(t *testing.T) {
	cli := DevnetClient(t)
	signer := Address
	objId, err := types.NewHexData("0x11462c88e74bb00079e3c043efb664482ee4551744ee691c7623b98503cb3f4d") // 0.2 SUI
	require.NoError(t, err)
	splitCoins := []uint64{150000000} // 0.15 SUI

	txnBytes, err := cli.SplitCoin(context.Background(), *signer, *objId, splitCoins, nil, 10000)
	require.Nil(t, err)

	simulateCheck(t, cli, txnBytes, M1Account(t))
}

func TestClient_SplitCoinEqual(t *testing.T) {
	cli := DevnetClient(t)
	signer := Address
	objId, err := types.NewHexData("0x11462c88e74bb00079e3c043efb664482ee4551744ee691c7623b98503cb3f4d") // 0.2 SUI
	require.NoError(t, err)

	txnBytes, err := cli.SplitCoinEqual(context.Background(), *signer, *objId, 2, nil, 10000)
	require.Nil(t, err)

	simulateCheck(t, cli, txnBytes, M1Account(t))
}

func TestClient_MergeCoins(t *testing.T) {
	cli := DevnetClient(t)
	signer := Address
	coin1, err := types.NewHexData("0x0b7b4b8474f22334aaf5688ed1ca05ef2877d8e695f8dc541db3935a52eace79") // 0.2 SUI
	require.NoError(t, err)
	coin2, err := types.NewHexData("0x0fe1d3981da6954ed97e98e715ba41c54e67e4b461b2420a1e082b93d5700871") // 0.2 SUI
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
// @return Always nil
func simulateCheck(t *testing.T, cli *Client, txn *types.TransactionBytes, acc *account.Account) *types.SuiTransactionBlockResponse {
	if shouldSimulate() {
		simulate, err := cli.DryRunTransaction(context.Background(), txn)
		require.Nil(t, err)
		require.True(t, simulate.Effects.IsSuccess())
	}
	return nil
}

func simulateAndSendTxn(t *testing.T, cli *Client, txn *types.TransactionBytes, acc *account.Account) *types.SuiTransactionBlockResponse {
	if shouldExecute() || shouldSimulate() {
		simulate, err := cli.DryRunTransaction(context.Background(), txn)
		require.Nil(t, err)
		require.True(t, simulate.Effects.IsSuccess())
	}
	if shouldExecute() {
		signedTxn := txn.SignSerializedSigWith(acc.PrivateKey)
		options := types.SuiTransactionBlockResponseOptions{
			ShowEffects: true,
		}
		resp, err := cli.ExecuteSignedTransaction(context.Background(), *signedTxn, &options, types.TxnRequestTypeWaitForLocalExecution)
		require.NoError(t, err)
		t.Log(resp)
		return resp
	}
	return nil
}
