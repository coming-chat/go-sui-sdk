package client

import (
	"context"
	"testing"

	"github.com/coming-chat/go-sui/v2/account"
	"github.com/coming-chat/go-sui/v2/types"
	"github.com/stretchr/testify/require"
)

func TestAccountSignAndSend(t *testing.T) {
	// ManualTest_AccountSignAndSend(t)
}

func ManualTest_AccountSignAndSend(t *testing.T) {
	unsafeMnemonic := M1Mnemonic

	account, err := account.NewAccountWithMnemonic(unsafeMnemonic)
	require.Nil(t, err)
	t.Log(account.Address)

	cli := TestnetClient(t)
	signer := SuiAddressNoErr(account.Address)
	coins, err := cli.GetSuiCoinsOwnedByAddress(context.Background(), *signer)
	require.Nil(t, err)
	require.Greater(t, coins.TotalBalance().Int64(), SUI(0.01).Int64(), "insufficient balance")

	coinIds := make([]suiObjectID, len(coins))
	for i, c := range coins {
		coinIds[i] = c.CoinObjectId
	}
	gasBudget := types.NewSafeSuiBigInt(uint64(10000000))
	txn, err := cli.PayAllSui(context.Background(), *signer, *signer, coinIds, gasBudget)
	require.Nil(t, err)

	resp := executeTxn(t, cli, txn.TxBytes, account)
	t.Log("txn digest =", resp.Digest)
}
