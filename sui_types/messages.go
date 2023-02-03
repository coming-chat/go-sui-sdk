package sui_types

import "github.com/coming-chat/go-sui/types"

type TransactionData struct {
	Kind       TransactionKind
	Sender     types.Address
	GasPayment types.ObjectRef
	GasPrice   uint64
	GasBudget  uint64
}

type TransactionKind struct {
	Single *SingleTransactionKind
	Batch  []*SingleTransactionKind
}

func (t TransactionKind) IsBcsEnum() {
}

type SingleTransactionKind struct {
	TransferObject *TransferObject
	Publish        *MoveModulePublish
	Call           *MoveCall
	TransferSui    *TransferSui
	Pay            *Pay
	PaySui         *PaySui
	PayAllSui      *PayAllSui
	ChangeEpoch    *ChangeEpoch
	Genesis        *GenesisTransaction
}

func (s SingleTransactionKind) IsBcsEnum() {
}

type TransferObject struct {
	Recipient types.Address
	ObjectRef types.ObjectRef
}

type MoveModulePublish struct {
	Modules [][]byte
}

type MoveCall struct {
	Package       types.ObjectId
	Module        string
	Function      string
	TypeArguments []*TypeTag
	Arguments     []*CallArg
}

type TransferSui struct {
	Recipient types.Address
	Amount    *uint64 `bcs:"optional"`
}

type Pay struct {
	Coins      []*types.ObjectRef
	Recipients []*types.Address
	Amounts    []*uint64
}

type PaySui = Pay

type PayAllSui struct {
	Coins     []*types.ObjectRef
	Recipient types.Address
}

type ChangeEpoch struct {
}

type GenesisTransaction struct {
}

type TypeTag struct {
	//#[serde(rename = "bool", alias = "Bool")]
	//Bool,
	//#[serde(rename = "u8", alias = "U8")]
	//U8,
	//#[serde(rename = "u64", alias = "U64")]
	//U64,
	//#[serde(rename = "u128", alias = "U128")]
	//U128,
	//#[serde(rename = "address", alias = "Address")]
	//Address,
	//#[serde(rename = "signer", alias = "Signer")]
	//Signer,
	//#[serde(rename = "vector", alias = "Vector")]
	//Vector(Box<TypeTag>),
	//#[serde(rename = "struct", alias = "Struct")]
	//Struct(Box<StructTag>),
	//
	//// NOTE: Added in bytecode version v6, do not reorder!
	//#[serde(rename = "u16", alias = "U16")]
	//U16,
	//#[serde(rename = "u32", alias = "U32")]
	//U32,
	//#[serde(rename = "u256", alias = "U256")]
	//U256,
}

func (t TypeTag) IsBcsEnum() {
}

type CallArg struct {
	Pure   *[]byte
	Object *ObjectArg
	ObjVec []*ObjectArg
}

func (c CallArg) IsBcsEnum() {
}

type ObjectArg struct {
	ImmOrOwnedObject *types.ObjectRef
	SharedObject     *SharedObject
}

func (o ObjectArg) IsBcsEnum() {
}

type SharedObject struct {
	Id                   types.ObjectId
	InitialSharedVersion uint64
	Mutable              bool `bcs:"-"`
}
