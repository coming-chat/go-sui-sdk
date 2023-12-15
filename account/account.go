package account

import (
	"encoding/base64"
	"encoding/hex"

	"github.com/coming-chat/go-aptos/crypto/derivation"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/blake2b"
)

const (
	ADDRESS_LENGTH = 64
)

type Account struct {
	KeyPair sui_types.SuiKeyPair
	Address string
}

func NewAccount(scheme sui_types.SignatureScheme, seed []byte) *Account {
	suiKeyPair := sui_types.NewSuiKeyPair(scheme, seed)
	tmp := []byte{scheme.Flag()}
	tmp = append(tmp, suiKeyPair.PublicKey()...)
	addrBytes := blake2b.Sum256(tmp)
	address := "0x" + hex.EncodeToString(addrBytes[:])[:ADDRESS_LENGTH]

	return &Account{
		KeyPair: suiKeyPair,
		Address: address,
	}
}

func NewAccountWithKeystore(keystore string) (*Account, error) {
	ksByte, err := base64.StdEncoding.DecodeString(keystore)
	if err != nil {
		return nil, err
	}
	scheme, err := sui_types.NewSignatureScheme(ksByte[0])
	if err != nil {
		return nil, err
	}
	return NewAccount(scheme, ksByte[1:]), nil
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
	scheme, err := sui_types.NewSignatureScheme(0)
	if err != nil {
		return nil, err
	}
	return NewAccount(scheme, key.Key), nil
}

func (a *Account) Sign(data []byte) []byte {
	switch a.KeyPair.Flag() {
	case 0:
		return a.KeyPair.Ed25519.Sign(data)
	default:
		return []byte{}
	}
}

func (a *Account) SignSecureWithoutEncode(txnBytes []byte, intent sui_types.Intent) (sui_types.Signature, error) {
	message := sui_types.NewIntentMessage(intent, bcsBytes(txnBytes))
	signature, err := sui_types.NewSignatureSecure(message, &a.KeyPair)
	if err != nil {
		return sui_types.Signature{}, err
	}
	return signature, nil
}

type bcsBytes []byte

func (b bcsBytes) MarshalBCS() ([]byte, error) {
	return b, nil
}
