package account

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var Mnemonic = os.Getenv("WalletSdkTestM1")

func TestMyAccouunt(t *testing.T) {
	account, err := NewAccountWithMnemonic(Mnemonic)
	assert.Nil(t, err)

	t.Logf("pri = %x", account.PrivateKey[:32])
	t.Logf("pub = %x", account.PublicKey)
	t.Logf("addr = %v", account.Address) // 0x0bd43fc3aa4f62e8943d16f66beb7546fafb2bac
}
