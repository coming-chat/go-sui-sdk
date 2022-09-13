package types

type TransactionBytes struct {
	// the gas object to be used
	Gas ObjectRef `json:"gas"`

	// objects to be used in this transaction
	InputObjects []string `json:"inputObjects"` // TODO: type

	// transaction data bytes
	TxBytes Base64Data `json:"txBytes"`
}

type ObjectRef struct {
	Digest   Base64Data `json:"digest"`
	ObjectId HexData    `json:"objectId"`
	Version  int        `json:"version"`
}

type SignatureScheme string

const (
	SignatureSchemeEd25519   SignatureScheme = "ED25519"
	SignatureSchemeSecp256k1 SignatureScheme = "Secp256k1"
)

type SignedTransaction struct {
	// transaction data bytes
	TxBytes *Base64Data `json:"tx_bytes"`

	// Flag of the signature scheme that is used.
	SigScheme SignatureScheme `json:"sig_scheme"`

	// transaction signature
	Signature *Base64Data `json:"signature"`

	// signer's public key
	PublicKey *Base64Data `json:"pub_key"`
}
