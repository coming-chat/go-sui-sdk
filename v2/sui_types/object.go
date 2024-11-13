package sui_types

import "github.com/coming-chat/go-sui/v2/lib"

const (
	ObjectStartVersion = SequenceNumber(1)
)

type Data struct {
	Move    *MoveObject
	Package *MovePackage
}

func (d Data) IsBcsEnum() {

}

type MoveObjectType = MoveObjectType_

type MoveObject struct {
	Type              MoveObjectType
	HasPublicTransfer bool
	Version           SequenceNumber
	Contents          []uint8
}

type Owner struct {
	AddressOwner *SuiAddress `json:"AddressOwner"`
	ObjectOwner  *SuiAddress `json:"ObjectOwner"`
	Shared       *struct {
		InitialSharedVersion SequenceNumber `json:"initial_shared_version"`
	} `json:"Shared,omitempty"`
	Immutable *lib.EmptyEnum `json:"Immutable,omitempty"`
}

func (o Owner) IsBcsEnum() {
}

func (o Owner) Tag() string {
	return ""
}

func (o Owner) Content() string {
	return ""
}

type Object struct {
	Data                Data
	Owner               Owner
	PreviousTransaction TransactionDigest
	StorageRebate       uint64
}
