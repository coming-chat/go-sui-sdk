package move_types

import "github.com/coming-chat/go-sui/lib"

type StructTag struct {
	Address    AccountAddress
	Module     Identifier
	Name       Identifier
	TypeParams []TypeTag
}

type TypeTag struct {
	Bool    *lib.EmptyEnum
	U8      *lib.EmptyEnum
	U64     *lib.EmptyEnum
	U128    *lib.EmptyEnum
	Address *lib.EmptyEnum
	Signer  *lib.EmptyEnum
	Vector  *TypeTag
	Struct  *StructTag
	U16     *lib.EmptyEnum
	U32     *lib.EmptyEnum
	U256    *lib.EmptyEnum
}

func (t TypeTag) IsBcsEnum() {
}
