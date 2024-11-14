package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/coming-chat/go-sui/v2/lib"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"strings"
)

const (
	SuiCoinType   = "0x2::sui::SUI"
	DevNetRpcUrl  = "https://fullnode.devnet.sui.io"
	TestnetRpcUrl = "https://fullnode.testnet.sui.io"
	MainnetRpcUrl = "https://fullnode.mainnet.sui.io"
)

// ShortString Returns the address with leading zeros trimmed, e.g. 0x2

type InputObjectKind map[string]interface{}

type TransactionBytes struct {
	// the gas object to be used
	Gas []sui_types.ObjectRef `json:"gas"`

	// objects to be used in this transaction
	InputObjects []InputObjectKind `json:"inputObjects"`

	// transaction data bytes
	TxBytes lib.Base64Data `json:"txBytes"`
}

type TransferObject struct {
	Recipient sui_types.SuiAddress `json:"recipient"`
	ObjectRef sui_types.ObjectRef  `json:"object_ref"`
}
type ModulePublish struct {
	Modules [][]byte `json:"modules"`
}
type MoveCall struct {
	Package  sui_types.ObjectID `json:"package"`
	Module   string             `json:"module"`
	Function string             `json:"function"`
	TypeArgs []interface{}      `json:"typeArguments"`
	Args     []interface{}      `json:"arguments"`
}
type TransferSui struct {
	Recipient sui_types.SuiAddress `json:"recipient"`
	Amount    uint64               `json:"amount"`
}
type Pay struct {
	Coins      []sui_types.ObjectRef  `json:"coins"`
	Recipients []sui_types.SuiAddress `json:"recipients"`
	Amounts    []uint64               `json:"amounts"`
}
type PaySui struct {
	Coins      []sui_types.ObjectRef  `json:"coins"`
	Recipients []sui_types.SuiAddress `json:"recipients"`
	Amounts    []uint64               `json:"amounts"`
}
type PayAllSui struct {
	Coins     []sui_types.ObjectRef `json:"coins"`
	Recipient sui_types.SuiAddress  `json:"recipient"`
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

	Sender     *sui_types.SuiAddress `json:"sender"`
	GasPayment *sui_types.ObjectRef  `json:"gasPayment"`
	GasBudget  uint64                `json:"gasBudget"`
	// GasPrice     uint64      `json:"gasPrice"`
}

type TimeRange struct {
	StartTime uint64 `json:"startTime"` // left endpoint of time interval, milliseconds since epoch, inclusive
	EndTime   uint64 `json:"endTime"`   // right endpoint of time interval, milliseconds since epoch, exclusive
}

type MoveModule struct {
	Package sui_types.ObjectID `json:"package"`
	Module  string             `json:"module"`
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

func IsSameStringAddress(addr1, addr2 string) bool {
	if strings.HasPrefix(addr1, "0x") {
		addr1 = addr1[2:]
	}
	if strings.HasPrefix(addr2, "0x") {
		addr2 = addr2[2:]
	}
	addr1 = strings.TrimLeft(addr1, "0")
	return strings.TrimLeft(addr1, "0") == strings.TrimLeft(addr2, "0")
}
