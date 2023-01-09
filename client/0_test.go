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
	Address, _ = types.NewAddressFromHex("0xb08ae4d187ca0057baa1666fe43fb9d7f3693a9a")
	M1Mnemonic = os.Getenv("WalletSdkTestM1")
)

func TestnetClient(t *testing.T) *Client {
	c, err := Dial(TestnetRpcUrl)
	require.NoError(t, err)
	return c
}

func DevnetClient(t *testing.T) *Client {
	c, err := Dial(DevNetRpcUrl)

	coins, err := c.GetCoins(context.TODO(), *Address, nil, nil, 1)
	require.NoError(t, err)
	if len(coins.Data) == 0 {
		_, err = FaucetFundAccount(Address.String(), DevNetFaucetUrl)
		require.NoError(t, err)
	}
	return c
}

func M1Account(t *testing.T) *account.Account {
	a, err := account.NewAccountWithMnemonic(M1Mnemonic)
	require.NoError(t, err)
	return a
}
