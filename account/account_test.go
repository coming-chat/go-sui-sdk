package account

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const Mnemonic = "crack coil okay hotel glue embark all employ east impact stomach cigar"

func TestAccount(t *testing.T) {
	account, err := NewAccountWithMnemonic(Mnemonic)
	assert.Nil(t, err)
	assert.Equal(t, account.Address, "0xbb8f7e72ae99d371020a1ccfe703bfb64a8a430f")

	t.Logf("pri = %x", account.PrivateKey[:32])
	t.Logf("pub = %x", account.PublicKey)
	t.Logf("addr = %v", account.Address)
}

func TestMyAccouunt(t *testing.T) {
	mnemonic := os.Getenv("WalletSdkTestM1")
	account, err := NewAccountWithMnemonic(mnemonic)
	assert.Nil(t, err)

	t.Logf("pri = %x", account.PrivateKey[:32])
	t.Logf("pub = %x", account.PublicKey)
	t.Logf("addr = %v", account.Address) // 0x0bd43fc3aa4f62e8943d16f66beb7546fafb2bac
}
