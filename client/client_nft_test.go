package client

import (
	"context"
	"os"
	"testing"

	"github.com/coming-chat/go-sui/account"
	"github.com/coming-chat/go-sui/types"
	"github.com/stretchr/testify/require"
)

var Mnemonic = os.Getenv("WalletSdkTestM1")

func TestMintNFT(t *testing.T) {
	account, err := account.NewAccountWithMnemonic(Mnemonic)
	require.Nil(t, err)

	client := DevnetClient(t)

	signer, err := types.NewAddressFromHex(account.Address)
	require.Nil(t, err)
	gasObj, err := types.NewHexData("0x1dde86ffbc05ab0964b70f07029e65b8d74b4f66")
	require.Nil(t, err)

	const (
		nftName = "ComingChat NFT"
		nftDesc = "This is a NFT created by ComingChat"
		nftUrl  = "https://coming.chat/favicon.ico"
	)
	txnBytes, err := client.MintDevnetNFT(context.Background(), *signer, nftName, nftDesc, nftUrl, gasObj, 100000)
	require.Nil(t, err)
	t.Log(txnBytes.TxBytes)

	// sign & send
	signedTxn := txnBytes.SignWith(account.PrivateKey)
	response, err := client.ExecuteTransaction(context.Background(), *signedTxn)
	require.Nil(t, err)
	t.Log("hash: ", response.Certificate.TransactionDigest)
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
