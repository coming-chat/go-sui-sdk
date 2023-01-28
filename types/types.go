package types

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	SuiCoinType = "0x2::sui::SUI"
)

type Address = HexData

// NewAddressFromHex
/**
 * Creates Address from a hex string.
 * @param addr Hex string can be with a prefix or without a prefix,
 * e.g. '0x1aa' or '1aa'. Hex string will be left padded with 0s if too short.
 */
func NewAddressFromHex(addr string) (*Address, error) {
	if strings.HasPrefix(addr, "0x") || strings.HasPrefix(addr, "0X") {
		addr = addr[2:]
	}
	if len(addr)%2 != 0 {
		addr = "0" + addr
	}

	data, err := hex.DecodeString(addr)
	if err != nil {
		return nil, err
	}
	const addressLength = 20
	if len(data) > addressLength {
		return nil, fmt.Errorf("hex string is too long. Address's length is %v data", addressLength)
	}

	res := [addressLength]byte{}
	copy(res[addressLength-len(data):], data[:])
	return &Address{
		data: res[:],
	}, nil
}

// ShortString Returns the address with leading zeros trimmed, e.g. 0x2
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
	Digest   string   `json:"digest"`
	ObjectId ObjectId `json:"objectId"`
	Version  int64    `json:"version"`
}

type SignatureScheme string

const (
	SignatureSchemeEd25519   SignatureScheme = "ED25519"
	SignatureSchemeSecp256k1 SignatureScheme = "Secp256k1"
)

type SignatureSchemeSerialized byte

const (
	SignatureSchemeSerializedEd25519   SignatureSchemeSerialized = 0
	SignatureSchemeSerializedSecp256k1 SignatureSchemeSerialized = 1
)

type ExecuteTransactionRequestType string

const (
	TxnRequestTypeWaitForEffectsCert    ExecuteTransactionRequestType = "WaitForEffectsCert"
	TxnRequestTypeWaitForLocalExecution ExecuteTransactionRequestType = "WaitForLocalExecution"
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

type SignedTransactionSerializedSig struct {
	// transaction data bytes
	TxBytes *Base64Data `json:"tx_bytes"`

	// transaction signature
	Signature *Base64Data `json:"signature"`
}

type TransferObject struct {
	Recipient Address   `json:"recipient"`
	ObjectRef ObjectRef `json:"object_ref"`
}
type ModulePublish struct {
	Modules [][]byte `json:"modules"`
}
type MoveCall struct {
	Package  ObjectRef     `json:"package"`
	Module   string        `json:"module"`
	Function string        `json:"function"`
	TypeArgs []interface{} `json:"typeArguments"`
	Args     []interface{} `json:"arguments"`
}
type TransferSui struct {
	Recipient Address `json:"recipient"`
	Amount    uint64  `json:"amount"`
}
type Pay struct {
	Coins      []ObjectRef `json:"coins"`
	Recipients []Address   `json:"recipients"`
	Amounts    []uint64    `json:"amounts"`
}
type PaySui struct {
	Coins      []ObjectRef `json:"coins"`
	Recipients []Address   `json:"recipients"`
	Amounts    []uint64    `json:"amounts"`
}
type PayAllSui struct {
	Coins     []ObjectRef `json:"coins"`
	Recipient Address     `json:"recipient"`
}
type ChangeEpoch struct {
	Epoch             interface{} `json:"epoch"`
	StorageCharge     uint64      `json:"storage_charge"`
	ComputationCharge uint64      `json:"computation_charge"`
}

type SingleTransactionKind struct {
	TransferObject *TransferObject `json:"TransferObject,omitempty"`
	Publish        *ModulePublish  `json:"Publish,omitempty"`
	Call           *MoveCall       `json:"Call,omitempty"`
	TransferSui    *TransferSui    `json:"TransferSui,omitempty"`
	ChangeEpoch    *ChangeEpoch    `json:"ChangeEpoch,omitempty"`
	PaySui         *PaySui         `json:"PaySui,omitempty"`
	Pay            *Pay            `json:"Pay,omitempty"`
	PayAllSui      *PayAllSui      `json:"PayAllSui,omitempty"`
}

type SenderSignedData struct {
	Transactions []SingleTransactionKind `json:"transactions,omitempty"`

	Sender     *Address   `json:"sender"`
	GasPayment *ObjectRef `json:"gasPayment"`
	GasBudget  uint64     `json:"gasBudget"`
	// GasPrice     uint64      `json:"gasPrice"`
}

type OwnedObjectRef struct {
	Owner     *ObjectOwner `json:"owner"`
	Reference *ObjectRef   `json:"reference"`
}

type Event interface{}

type ObjectOwner struct {
	*ObjectOwnerInternal
	*string
}

type ObjectOwnerInternal struct {
	AddressOwner *Address `json:"AddressOwner,omitempty"`
	ObjectOwner  *Address `json:"ObjectOwner,omitempty"`
	SingleOwner  *Address `json:"SingleOwner,omitempty"`
	Shared       *struct {
		InitialSharedVersion int64 `json:"initial_shared_version"`
	} `json:"Shared,omitempty"`
}

func (o ObjectOwner) MarshalJSON() ([]byte, error) {
	if o.string != nil {
		data, err := json.Marshal(o.string)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if o.ObjectOwnerInternal != nil {
		data, err := json.Marshal(o.ObjectOwnerInternal)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, errors.New("nil value")
}

func (o *ObjectOwner) UnmarshalJSON(data []byte) error {
	if bytes.HasPrefix(data, []byte("\"")) {
		stringData := string(data[1 : len(data)-1])
		o.string = &stringData
		return nil
	}
	if bytes.HasPrefix(data, []byte("{")) {
		oOI := ObjectOwnerInternal{}
		err := json.Unmarshal(data, &oOI)
		if err != nil {
			return err
		}
		o.ObjectOwnerInternal = &oOI
		return nil
	}
	return errors.New("value not json")
}

type ObjectReadDetail struct {
	Data  map[string]interface{} `json:"data"`
	Owner *ObjectOwner           `json:"owner"`

	PreviousTransaction string     `json:"previousTransaction"`
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
	Digest   string       `json:"digest"`
	Type     string       `json:"type"`
	Owner    *ObjectOwner `json:"owner"`

	PreviousTransaction string `json:"previousTransaction"`
}

// IntentBytes See: sui/crates/sui-types/src/intent.rs
// This is currently hardcoded with [IntentScope::TransactionData = 0, Version::V0 = 0, AppId::Sui = 0]
var IntentBytes = []byte{0, 0, 0}

func (txn *TransactionBytes) SignWith(privateKey ed25519.PrivateKey) *SignedTransaction {
	signTx := bytes.NewBuffer(IntentBytes)
	signTx.Write(txn.TxBytes.Data())
	message := signTx.Bytes()
	signature := ed25519.Sign(privateKey, message)
	sign := Bytes(signature).GetBase64Data()
	publicKey := privateKey.Public().(ed25519.PublicKey)
	pub := Bytes(publicKey).GetBase64Data()

	return &SignedTransaction{
		TxBytes:   &txn.TxBytes,
		SigScheme: SignatureSchemeEd25519,
		Signature: &sign,
		PublicKey: &pub,
	}
}

func (txn *TransactionBytes) SignSerializedSigWith(privateKey ed25519.PrivateKey) *SignedTransactionSerializedSig {
	signTx := bytes.NewBuffer(IntentBytes)
	signTx.Write(txn.TxBytes.Data())
	message := signTx.Bytes()
	signatureData := bytes.NewBuffer([]byte{byte(SignatureSchemeSerializedEd25519)})
	signatureData.Write(ed25519.Sign(privateKey, message))
	signatureData.Write(privateKey.Public().(ed25519.PublicKey))
	signature := Base64Data{
		data: signatureData.Bytes(),
	}

	return &SignedTransactionSerializedSig{
		TxBytes:   &txn.TxBytes,
		Signature: &signature,
	}
}
