package client

import (
	"context"
	"github.com/coming-chat/go-sui/types"
	"os"
	"testing"

	"github.com/coming-chat/go-sui/account"
	"github.com/stretchr/testify/require"
)

const (
	DevNetRpcUrl  = "https://fullnode.devnet.sui.io"
	TestnetRpcUrl = "https://fullnode.testnet.sui.io"
)

var (
	Address, _ = types.NewAddressFromHex("0x6fc6148816617c3c3eccb1d09e930f73f6712c9c")
	M1Mnemonic = os.Getenv("WalletSdkTestM1")
)

func TestnetClient(t *testing.T) *Client {
	c, err := Dial(TestnetRpcUrl)
	require.Nil(t, err)
	return c
}

func DevnetClient(t *testing.T) *Client {
	c, err := Dial(DevNetRpcUrl)

	coins, err := c.GetCoins(context.TODO(), *Address, nil, nil, 1)
	require.Nil(t, err)
	if len(coins.Data) == 0 {
		_, err = FaucetFundAccount(Address.String(), DevNetFaucetUrl)
		require.NoError(t, err)
	}
	return c
}

func M1Account(t *testing.T) *account.Account {
	a, err := account.NewAccountWithMnemonic(M1Mnemonic)
	require.Nil(t, err)
	return a
}
