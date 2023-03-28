package account

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var Mnemonic = os.Getenv("WalletSdkTestM1")

func TestMyAccouunt(t *testing.T) {
	account, err := NewAccountWithMnemonic(Mnemonic)
	require.Nil(t, err)

	t.Logf("addr = %v", account.Address)
}
