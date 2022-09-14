package types

type Address = HexData
type ObjectId = HexData
type Digest = Base64Data

type InputObjectKind map[string]interface{}

type TransactionBytes struct {
	// the gas object to be used
	Gas ObjectRef `json:"gas"`

	// objects to be used in this transaction
	InputObjects []InputObjectKind `json:"inputObjects"`

	// transaction data bytes
	TxBytes Base64Data `json:"txBytes"`
}

type ObjectRef struct {
	Digest   Digest   `json:"digest"`
	ObjectId ObjectId `json:"objectId"`
	Version  int      `json:"version"`
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

type CertifiedTransaction map[string]interface{}

type TransactionEffects map[string]interface{}

type ParsedTransactionResponse interface{}

type TransactionResponse struct {
	Certificate []CertifiedTransaction    `json:"certificate"`
	Effects     []TransactionEffects      `json:"effects"`
	ParsedData  ParsedTransactionResponse `json:"parsed_data,omitempty"`
	TimestampMs uint64                    `json:"timestamp_ms,omitempty"`
}

type ObjectOwner struct {
	AddressOwner *Address `json:"AddressOwner"`
}

type ObjectReadDetail struct {
	Data  map[string]interface{} `json:"data"`
	Owner *ObjectOwner           `json:"owner"`

	PreviousTransaction *Digest    `json:"previousTransaction"`
	StorageRebate       int        `json:"storageRebate"`
	Reference           *ObjectRef `json:"reference"`
}

type ObjectRead struct {
	Details *ObjectReadDetail `json:"details"`
	Status  string            `json:"status"`
}

type ObjectInfo struct {
	ObjectId *ObjectId    `json:"objectId"`
	Version  int          `json:"version"`
	Digest   *Digest      `json:"digest"`
	Type     string       `json:"type"`
	Owner    *ObjectOwner `json:"owner"`

	PreviousTransaction *Digest `json:"previousTransaction"`
}
