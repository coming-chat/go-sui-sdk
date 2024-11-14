package crypto

import "crypto/ed25519"

type KeyPair interface {
	Sign(msg []byte) []byte
	PublicKey() []byte
	PrivateKey() []byte
}

type Ed25519KeyPair struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

func NewEd25519KeyPair(privateKey ed25519.PrivateKey) *Ed25519KeyPair {
	return &Ed25519KeyPair{
		privateKey: privateKey,
		publicKey:  privateKey.Public().(ed25519.PublicKey),
	}
}

func (e *Ed25519KeyPair) Sign(msg []byte) []byte {
	return ed25519.Sign(e.privateKey, msg)
}

func (e *Ed25519KeyPair) PublicKey() []byte {
	return e.publicKey
}

func (e *Ed25519KeyPair) PrivateKey() []byte {
	return e.privateKey
}
