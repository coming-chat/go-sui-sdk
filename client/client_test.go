package client

import (
	"testing"

	"github.com/coming-chat/go-sui/types"
	"github.com/stretchr/testify/assert"
)

const DevnetRpcUrl = "https://gateway.devnet.sui.io:443"

func TestClient_Call(t *testing.T) {
	client, err := Dial(DevnetRpcUrl)
	assert.Nil(t, err)

	txn := types.TransactionBytes{}
	params := []interface{}{
		"0xbb8f7e72ae99d371020a1ccfe703bfb64a8a430f",
		"0x36d3176a796e167ffcbd823c94718e7db56b955f",
		[]int{40000, 5000, 5000},
		"0x9f662fec10f77b5cfd1bed5ffa53232b8a62a982",
		2000,
	}
	err = client.Call(&txn, "sui_splitCoin", params...)
	assert.Nil(t, err)

	t.Log(txn)
	t.Log(txn.TxBytes.String())
}
