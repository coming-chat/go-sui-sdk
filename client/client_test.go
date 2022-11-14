package client

import (
	"context"
	"testing"

	"github.com/coming-chat/go-sui/types"
	"github.com/stretchr/testify/require"
)

func TestClient_SplitCoin(t *testing.T) {
	account := M1Account(t)
	client := DevnetClient(t)

	signer, err := types.NewAddressFromHex(account.Address)
	require.Nil(t, err)
	coins, err := client.GetSuiCoinsOwnedByAddress(context.Background(), *signer)
	require.Nil(t, err)

	firstCoin, err := coins.PickCoinNoLess(100000)
	require.Nil(t, err)
	everyAmount := firstCoin.Balance / 2
	amounts := []uint64{everyAmount, everyAmount}

	txn, err := client.SplitCoin(context.Background(), *signer, firstCoin.Reference.ObjectId, amounts, nil, 100000)
	require.Nil(t, err)

	t.Log(txn.TxBytes.String())

	signedTxn := txn.SignWith(account.PrivateKey)
	resp, err := client.ExecuteTransaction(context.Background(), *signedTxn, types.TxnRequestTypeWaitForLocalExecution)
	require.Nil(t, err)
	t.Log("hash =", resp.TransactionDigest())
}

func TestClient_SplitCoinEqual(t *testing.T) {
	account := M1Account(t)
	client := DevnetClient(t)

	signer, err := types.NewAddressFromHex(account.Address)
	require.Nil(t, err)
	coins, err := client.GetSuiCoinsOwnedByAddress(context.Background(), *signer)
	require.Nil(t, err)

	firstCoin, err := coins.PickCoinNoLess(100000)
	require.Nil(t, err)

	txn, err := client.SplitCoinEqual(context.Background(), *signer, firstCoin.Reference.ObjectId, 2, nil, 100000)
	require.Nil(t, err)

	t.Log(txn.TxBytes.String())

	signedTxn := txn.SignWith(account.PrivateKey)
	resp, err := client.ExecuteTransaction(context.Background(), *signedTxn, types.TxnRequestTypeWaitForLocalExecution)
	require.Nil(t, err)
	t.Log("hash =", resp.TransactionDigest())
}

func TestTransaction(t *testing.T) {
	digest := "gTYc+0O3nl1m2uCf36HEOWMFgEAe/eyyKBKDpJT6wV0="

	dig, err := types.NewBase64Data(digest)
	require.NoError(t, err)

	cli := DevnetClient(t)
	resp, err := cli.GetTransaction(context.Background(), *dig)
	require.NoError(t, err)

	t.Log(resp)
}

func TestBatchCall_GetObject(t *testing.T) {
	objKeys := []string{
		"0x055a53f0f5e8f711c04d31e8d7ae7d021f2e0171",
		"0x1911b14cf5b5a356258a5379a871aaecb8b88a50",
		"0x0000000000000000000000000000000000000005",
		"0xb540b59c0f6fe006aad6df91776aba721a18ecc4",
	}

	elems := make([]BatchElem, len(objKeys))
	for i := 0; i < len(objKeys); i++ {
		ele := BatchElem{
			Method: "sui_getObject",
			Args:   []interface{}{objKeys[i]},
			Result: &types.ObjectRead{},
		}
		elems[i] = ele
	}

	cli := DevnetClient(t)
	err := cli.BatchCall(elems)
	require.NoError(t, err)

	for _, ele := range elems {
		t.Log("res = ", ele.Result)
		if ele.Error != nil {
			t.Log("❗️❗️❗️", ele.Error)
		}
	}
	t.Log("")
}

func TestBatchGetObjectsOwnedByAddress(t *testing.T) {
	addr, err := types.NewAddressFromHex("0x6c5d2cd6e62734f61b4e318e58cbfd1c4b99dfaf")
	require.NoError(t, err)

	cli := TestnetClient(t)
	coins, err := cli.GetSuiCoinsOwnedByAddress(context.Background(), *addr)
	require.NoError(t, err)

	t.Log(coins.TotalBalance().String())
}
