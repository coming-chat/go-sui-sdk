package move_types

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
)

const SuiAddressLen = 32

type AccountAddress [SuiAddressLen]uint8

func NewAccountAddressHex(str string) (*AccountAddress, error) {
	if strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X") {
		str = str[2:]
	}
	if len(str)%2 != 0 {
		str = "0" + str
	}
	data, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}
	if len(data) > SuiAddressLen {
		return nil, errors.New("the len is invalid")
	}
	var accountAddress AccountAddress
	copy(accountAddress[SuiAddressLen-len(data):], data[:])
	return &accountAddress, nil
}

func (a AccountAddress) Data() []byte {
	return a[:]
}
func (a AccountAddress) Length() int {
	return len(a)
}
func (a AccountAddress) String() string {
	return "0x" + hex.EncodeToString(a[:])
}

func (a AccountAddress) ShortString() string {
	return "0x" + strings.TrimLeft(hex.EncodeToString(a[:]), "0")
}

func (a AccountAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

func (a *AccountAddress) UnmarshalJSON(data []byte) error {
	str := ""
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	tmp, err := NewAccountAddressHex(str)
	if err == nil {
		*a = *tmp
	}
	return err
}

func (a AccountAddress) MarshalBCS() ([]byte, error) {
	return a[:], nil
}
