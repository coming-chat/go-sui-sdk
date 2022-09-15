package types

import (
	"errors"
	"fmt"
	"strings"
)

type ResourceType struct {
	Address    *Address
	ModuleName string
	FuncName   string

	SubType *ResourceType
}

func NewResourceType(str string) (*ResourceType, error) {
	if strings.Contains(str, "<") {
		return nil, errors.New("Not implemented")
	}
	parts := strings.Split(str, "::")
	if len(parts) != 3 {
		return nil, errors.New("Invalid type string literal.")
	}
	addr, err := NewAddressFromHex(parts[0])
	if err != nil {
		return nil, err
	}
	return &ResourceType{
		Address:    addr,
		ModuleName: parts[1],
		FuncName:   parts[2],
	}, nil
}

func (t *ResourceType) String() string {
	return fmt.Sprintf("%v::%v::%v", t.Address.String(), t.ModuleName, t.FuncName)
}

func (t *ResourceType) ShortString() string {
	return fmt.Sprintf("%v::%v::%v", t.Address.ShortString(), t.ModuleName, t.FuncName)
}
