package types

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strings"
)

type BytesData interface {
	Data() []byte
	Length() int
	String() string
}

type Bytes []byte

func (b Bytes) GetHexData() HexData {
	return HexData{b}
}
func (b Bytes) GetBase64Data() Base64Data {
	return Base64Data{b}
}

type HexData struct {
	data []byte
}

func NewHexData(str string) (*HexData, error) {
	if strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X") {
		str = str[2:]
	}
	data, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}
	return &HexData{data}, nil
}

func (h HexData) Data() []byte {
	return h.data
}
func (h HexData) Length() int {
	return len(h.data)
}
func (h HexData) String() string {
	return "0x" + hex.EncodeToString(h.data)
}

func (h HexData) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

func (h *HexData) UnmarshalJSON(data []byte) error {
	str := ""
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	tmp, err := NewHexData(str)
	if err == nil {
		h.data = tmp.data
	}
	return err
}

type Base64Data struct {
	data []byte
}

func NewBase64Data(str string) (*Base64Data, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}
	return &Base64Data{data}, nil
}

func (h Base64Data) Data() []byte {
	return h.data
}
func (h Base64Data) Length() int {
	return len(h.data)
}
func (h Base64Data) String() string {
	return base64.StdEncoding.EncodeToString(h.data)
}

func (h Base64Data) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

func (h *Base64Data) UnmarshalJSON(data []byte) error {
	str := ""
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	tmp, err := NewBase64Data(str)
	if err == nil {
		h.data = tmp.data
	}
	return err
}
