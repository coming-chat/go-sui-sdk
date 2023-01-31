package account

import (
	"crypto/ed25519"
	"encoding/hex"

	"github.com/coming-chat/go-aptos/crypto/derivation"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/sha3"
)

const (
	SIGNATURE_SCHEME_FLAG_ED25519 = 0x0

	ADDRESS_LENGTH = 40
)

type Account struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
	Address    string
}

func NewAccount(seed []byte) *Account {
	privateKey := ed25519.NewKeyFromSeed(seed[:])
	publicKey := privateKey.Public().(ed25519.PublicKey)

	tmp := []byte{SIGNATURE_SCHEME_FLAG_ED25519}
	tmp = append(tmp, publicKey...)
	addrBytes := sha3.Sum256(tmp)
	address := "0x" + hex.EncodeToString(addrBytes[:])[:ADDRESS_LENGTH]

	return &Account{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    address,
	}
}

func NewAccountWithMnemonic(mnemonic string) (*Account, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}
	key, err := derivation.DeriveForPath("m/44'/784'/0'/0'/0'", seed)
	if err != nil {
		return nil, err
	}
	return NewAccount(key.Key), nil
}

// GetOldVersionPrivateKeyWithMnemonic
// Deprecated: sui change the plan, replace with NewAccountWithMnemonic
func GetOldVersionPrivateKeyWithMnemonic(mnemonic string) ([]byte, error) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, err
	}
	return seed[:ed25519.SeedSize], nil
}

func (a *Account) Sign(data []byte) []byte {
	return ed25519.Sign(a.PrivateKey, data)
}
