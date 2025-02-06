package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFaucetFundAccount_Devnet(t *testing.T) {
	// addr := M1Account(t).Address
	addr := "0xd77955e670f42c1bc5e94b9e68e5fe9bdbed9134d784f2a14dfe5fc1b24b5d9f"

	res, err := FaucetFundAccount(addr, DevNetFaucetUrl)
	require.Nil(t, err)
	t.Log("hash = ", res)
}

// func TestFaucetFundAccount_Testnet(t *testing.T) {
// 	addr := "0xd77955e670f42c1bc5e94b9e68e5fe9bdbed9134d784f2a14dfe5fc1b24b5d9f"
// 	res, err := FaucetFundAccount(addr, TestNetFaucetUrl)
// 	require.Nil(t, err)
// }
