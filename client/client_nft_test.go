package client

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMintNFT(t *testing.T) {
	cli := DevnetClient(t)

	var (
		timeNow = time.Now().Format("06-01-02 15:04")
		nftName = "ComingChat NFT at " + timeNow
		nftDesc = "This is a NFT created by ComingChat"
		nftUrl  = "https://coming.chat/favicon.ico"
	)
	coins, err := cli.GetSuiCoinsOwnedByAddress(context.TODO(), *Address)
	require.NoError(t, err)

	firstCoin, err := coins.PickCoinNoLess(12000)
	require.NoError(t, err)

	txnBytes, err := cli.MintNFT(context.TODO(), *Address, nftName, nftDesc, nftUrl, &firstCoin.Reference.ObjectId, 12000)
	require.NoError(t, err)
	t.Log(txnBytes.TxBytes)

	response, err := cli.DevInspectTransaction(context.TODO(), txnBytes.TxBytes)
	require.NoError(t, err)
	t.Logf("%#v", response)
}

func TestGetDevNFTs(t *testing.T) {
	cli := DevnetClient(t)

	nfts, err := cli.GetNFTsOwnedByAddress(context.TODO(), *Address)
	require.NoError(t, err)
	for _, nft := range nfts {
		t.Log(nft.Details)
	}
}
