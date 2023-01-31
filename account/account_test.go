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

	t.Logf("addr = %v", account.Address)
}
