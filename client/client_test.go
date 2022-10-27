package client

import (
	"context"
	"testing"

	"github.com/coming-chat/go-sui/types"
	"github.com/stretchr/testify/require"
)

func TestClient_Call(t *testing.T) {
	account := M1Account(t)
	client := DevnetClient(t)

	signer, err := types.NewAddressFromHex(account.Address)
	require.Nil(t, err)
	coins, err := client.GetSuiCoinsOwnedByAddress(context.Background(), *signer)
	require.Nil(t, err)

	firstCoin := coins[0]
	everyAmount := firstCoin.Balance / 3
	lastAmount := everyAmount
	if lastAmount > 30000 {
		lastAmount = everyAmount - 30000
	}
	amounts := []uint64{everyAmount, everyAmount, lastAmount}
	gasCoin := coins[1]

	txn, err := client.SplitCoin(context.Background(), *signer, firstCoin.Reference.ObjectId, amounts, gasCoin.Reference.ObjectId, 100000)
	require.Nil(t, err)

	t.Log(txn)
	t.Log(txn.TxBytes.String())

	signedTxn := txn.SignWith(account.PrivateKey)
	resp, err := client.ExecuteTransaction(context.Background(), *signedTxn, types.TxnRequestTypeWaitForLocalExecution)
	require.Nil(t, err)
	t.Log(resp.TransactionDigest())
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
	addr, err := types.NewAddressFromHex("0xbb8f7e72ae99d371020a1ccfe703bfb64a8a430f")
	require.NoError(t, err)

	cli := DevnetClient(t)
	coins, err := cli.GetSuiCoinsOwnedByAddress(context.Background(), *addr)
	require.NoError(t, err)

	t.Log(coins.TotalBalance().String())
}
