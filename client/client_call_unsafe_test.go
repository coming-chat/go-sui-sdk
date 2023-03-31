package client

import (
	"context"
	"testing"

	"github.com/coming-chat/go-sui/account"
	"github.com/coming-chat/go-sui/types"
	"github.com/stretchr/testify/require"
)

func shouldSimulate() bool {
	if shouldExecute() {
		return true
	}
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

func TestClient_SplitCoin(t *testing.T) {
	cli := DevnetClient(t)
	signer := Address
	objId, err := types.NewHexData("0x11462c88e74bb00079e3c043efb664482ee4551744ee691c7623b98503cb3f4d") // 0.2 SUI
	require.NoError(t, err)
	splitCoins := []uint64{150000000} // 0.15 SUI

	txnBytes, err := cli.SplitCoin(context.Background(), *signer, *objId, splitCoins, nil, 10000)
	require.Nil(t, err)
	t.Log(txnBytes)

	simulateAndSendTxn(t, cli, txnBytes, M1Account(t))
}

func simulateAndSendTxn(t *testing.T, cli *Client, txn *types.TransactionBytes, acc *account.Account) *types.SuiTransactionBlockResponse {
	if shouldSimulate() {
		simulate, err := cli.DryRunTransaction(context.Background(), txn)
		require.Nil(t, err)
		t.Log(simulate)
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
