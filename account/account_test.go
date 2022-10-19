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
	t.Logf("addr = %v", account.Address) // 0x6c5d2cd6e62734f61b4e318e58cbfd1c4b99dfaf
}
