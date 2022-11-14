package client

import (
	"context"
	"testing"
	"time"

	"github.com/coming-chat/go-sui/types"
	"github.com/stretchr/testify/require"
)

func TestMintNFT(t *testing.T) {
	account := M1Account(t)
	client := TestnetClient(t)

	signer, err := types.NewAddressFromHex(account.Address)
	require.Nil(t, err)
	gasBudget := uint64(12000)

	var (
		timeNow = time.Now().Format("06-01-02 15:04")
		nftName = "ComingChat NFT at " + timeNow
		nftDesc = "This is a NFT created by ComingChat"
		nftUrl  = "https://coming.chat/favicon.ico"
	)
	txnBytes, err := client.MintNFT(context.Background(), *signer, nftName, nftDesc, nftUrl, nil, gasBudget)
	require.Nil(t, err)
	t.Log(txnBytes.TxBytes)

	// sign & send
	signedTxn := txnBytes.SignWith(account.PrivateKey)
	response, err := client.ExecuteTransaction(context.Background(), *signedTxn, types.TxnRequestTypeWaitForLocalExecution)
	require.Nil(t, err)
	t.Log("hash: ", response.TransactionDigest())
}

func TestGetDevNFTs(t *testing.T) {
	account := M1Account(t)
	address, err := types.NewAddressFromHex(account.Address)
	require.Nil(t, err)

	client := DevnetClient(t)

	nfts, err := client.GetNFTsOwnedByAddress(context.Background(), *address)
	require.Nil(t, err)
	for _, nft := range nfts {
		t.Log(nft.Details)
	}
}
