package types

import (
	"encoding/hex"
	"fmt"
	"strings"
)

type Address = HexData

/**
 * Creates Address from a hex string.
 * @param addr Hex string can be with a prefix or without a prefix,
 *   e.g. '0x1aa' or '1aa'. Hex string will be left padded with 0s if too short.
 */
func NewAddressFromHex(addr string) (*Address, error) {
	if strings.HasPrefix(addr, "0x") || strings.HasPrefix(addr, "0X") {
		addr = addr[2:]
	}
	if len(addr)%2 != 0 {
		addr = "0" + addr
	}

	bytes, err := hex.DecodeString(addr)
	if err != nil {
		return nil, err
	}
	const addressLength = 20
	if len(bytes) > addressLength {
		return nil, fmt.Errorf("Hex string is too long. Address's length is %v bytes.", addressLength)
	}

	res := Address{}
	copy(res.data[addressLength-len(bytes):], bytes[:])
	return &res, nil
}

// Returns the address with leading zeros trimmed, e.g. 0x2
func (a Address) ShortString() string {
	return "0x" + strings.TrimLeft(hex.EncodeToString(a.data), "0")
}

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
	AddressOwner *Address `json:"AddressOwner,omitempty"`
	ObjectOwner  *Address `json:"ObjectOwner,omitempty"`
	SingleOwner  *Address `json:"SingleOwner,omitempty"`
}

type ObjectReadDetail struct {
	Data  map[string]interface{} `json:"data"`
	Owner *ObjectOwner           `json:"owner"`

	PreviousTransaction *Digest    `json:"previousTransaction"`
	StorageRebate       int        `json:"storageRebate"`
	Reference           *ObjectRef `json:"reference"`
}

type ObjectStatus string

const (
	ObjectStatusExists    ObjectStatus = "Exists"
	ObjectStatusNotExists ObjectStatus = "NotExists"
	ObjectStatusDeleted   ObjectStatus = "Deleted"
)

type ObjectRead struct {
	Details *ObjectReadDetail `json:"details"`
	Status  ObjectStatus      `json:"status"`
}

type ObjectInfo struct {
	ObjectId *ObjectId    `json:"objectId"`
	Version  int          `json:"version"`
	Digest   *Digest      `json:"digest"`
	Type     string       `json:"type"`
	Owner    *ObjectOwner `json:"owner"`

	PreviousTransaction *Digest `json:"previousTransaction"`
}
