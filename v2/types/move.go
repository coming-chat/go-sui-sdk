package types

import (
	"errors"
	"fmt"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"strings"
)

type ResourceType struct {
	Address    *sui_types.SuiAddress
	ModuleName string
	FuncName   string

	SubType *ResourceType
}

func NewResourceType(str string) (*ResourceType, error) {
	ltIdx := strings.Index(str, "<")
	var subType *ResourceType
	var err error
	if ltIdx != -1 {
		gtIdx := strings.LastIndex(str, ">")
		if gtIdx != len(str)-1 {
			return nil, errors.New("invalid type string literal")
		}
		subType, err = NewResourceType(str[ltIdx+1 : gtIdx])
		if err != nil {
			return nil, err
		}
		str = str[:ltIdx]
	}

	parts := strings.Split(str, "::")
	if len(parts) != 3 {
		return nil, errors.New("invalid type string literal")
	}
	addr, err := sui_types.NewAddressFromHex(parts[0])
	if err != nil {
		return nil, err
	}
	return &ResourceType{
		Address:    addr,
		ModuleName: parts[1],
		FuncName:   parts[2],
		SubType:    subType,
	}, nil
}

func (t *ResourceType) String() string {
	if t.SubType != nil {
		return fmt.Sprintf("%v::%v::%v<%v>", t.Address.String(), t.ModuleName, t.FuncName, t.SubType.String())
	} else {
		return fmt.Sprintf("%v::%v::%v", t.Address.String(), t.ModuleName, t.FuncName)
	}
}

func (t *ResourceType) ShortString() string {
	if t.SubType != nil {
		return fmt.Sprintf("%v::%v::%v<%v>", t.Address.ShortString(), t.ModuleName, t.FuncName, t.SubType.ShortString())
	} else {
		return fmt.Sprintf("%v::%v::%v", t.Address.ShortString(), t.ModuleName, t.FuncName)
	}
}
