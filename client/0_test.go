package client

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
	"github.com/shopspring/decimal"

	"github.com/coming-chat/go-sui/v2/account"
	"github.com/stretchr/testify/require"
)

var (
	M1Mnemonic = os.Getenv("WalletSdkTestM1")
	Address, _ = sui_types.NewAddressFromHex("0x7e875ea78ee09f08d72e2676cf84e0f1c8ac61d94fa339cc8e37cace85bebc6e")
)

func MainnetClient(t *testing.T) *Client {
	c, err := Dial(types.MainnetRpcUrl)
	require.NoError(t, err)
	return c
}

func TestnetClient(t *testing.T) *Client {
	c, err := Dial(types.TestnetRpcUrl)
	require.NoError(t, err)
	return c
}

func DevnetClient(t *testing.T) *Client {
	c, err := Dial(types.DevNetRpcUrl)
	require.NoError(t, err)

	balance, err := c.GetBalance(context.Background(), *Address, types.SUI_COIN_TYPE)
	require.NoError(t, err)
	if balance.TotalBalance.BigInt().Uint64() < SUI(0.3).Uint64() {
		_, err = FaucetFundAccount(Address.String(), DevNetFaucetUrl)
		require.NoError(t, err)
	}
	return c
}

func ChainClient(t *testing.T) *Client {
	suiEnv := os.Getenv("SUI_NETWORK")
	switch suiEnv {
	case "testnet":
		return TestnetClient(t)
	case "devnet":
		return DevnetClient(t)
	case "":
		fallthrough
	default:
		return TestnetClient(t)
	}
}

func M1Account(t *testing.T) *account.Account {
	a, err := account.NewAccountWithMnemonic(M1Mnemonic)
	require.NoError(t, err)
	return a
}

func M1Address(t *testing.T) *suiAddress {
	return Address
}

func Signer(t *testing.T) *account.Account {
	return M1Account(t)
}

type SUI float64

func (s SUI) Int64() int64 {
	return int64(s * 1e9)
}
func (s SUI) Uint64() uint64 {
	return uint64(s * 1e9)
}
func (s SUI) Decimal() decimal.Decimal {
	return decimal.NewFromInt(s.Int64())
}
func (s SUI) String() string {
	return strconv.FormatInt(s.Int64(), 10)
}

func SuiAddressNoErr(str string) *suiAddress {
	s, _ := sui_types.NewAddressFromHex(str)
	return s
}
