package client

import (
	"context"
	"testing"

	"github.com/coming-chat/go-sui/types"
	"github.com/stretchr/testify/require"
)

const DevnetRpcUrl = "https://gateway.devnet.sui.io:443"

func TestClient_Call(t *testing.T) {
	client := DevnetClient(t)

	txn := types.TransactionBytes{}
	params := []interface{}{
		"0xbb8f7e72ae99d371020a1ccfe703bfb64a8a430f",
		"0x36d3176a796e167ffcbd823c94718e7db56b955f",
		[]int{40000, 5000, 5000},
		"0x9f662fec10f77b5cfd1bed5ffa53232b8a62a982",
		2000,
	}
	err := client.Call(&txn, "sui_splitCoin", params...)
	require.NoError(t, err)

	t.Log(txn)
	t.Log(txn.TxBytes.String())
}

func DevnetClient(t *testing.T) *Client {
	c, err := Dial(DevnetRpcUrl)
	require.NoError(t, err)
	return c
}

func TestTransaction(t *testing.T) {
	// digest := "2yhXOzBqTsOpcWNZKCSsKySaUTJUVgpGyrzhQVu7PcM="
	digest := "4nMHqXi60PLxj/DxLCWwkiO3L41kIz89qMDEpStRdP8="

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
