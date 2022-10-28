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
	client := DevnetClient(t)

	signer, err := types.NewAddressFromHex(account.Address)
	require.Nil(t, err)
	gasBudget := uint64(100000)

	var (
		timeNow = time.Now().Format("06-01-02 15:04")
		nftName = "ComingChat NFT at " + timeNow
		nftDesc = "This is a NFT created by ComingChat"
		nftUrl  = "https://coming.chat/favicon.ico"
	)
	txnBytes, err := client.MintDevnetNFT(context.Background(), *signer, nftName, nftDesc, nftUrl, nil, gasBudget)
	require.Nil(t, err)
	t.Log(txnBytes.TxBytes)

	// sign & send
	signedTxn := txnBytes.SignWith(account.PrivateKey)
	response, err := client.ExecuteTransaction(context.Background(), *signedTxn, types.TxnRequestTypeWaitForLocalExecution)
	require.Nil(t, err)
	t.Log("hash: ", response.TransactionDigest())
}

func TestGetDevNFTs(t *testing.T) {
	address, err := types.NewAddressFromHex("0xbb8f7e72ae99d371020a1ccfe703bfb64a8a430f")
	require.Nil(t, err)

	client := DevnetClient(t)

	nfts, err := client.GetDevnetNFTOwnedByAddress(context.Background(), *address)
	require.Nil(t, err)
	for _, nft := range nfts {
		t.Log(nft.Details)
	}
}
