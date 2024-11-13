package lib

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcutil/base58"
	"strings"
)

type BytesData interface {
	~[]byte
	Data() []byte
	Length() int
	String() string
}

type Bytes []byte

func (b Bytes) GetHexData() HexData {
	return HexData(b)
}
func (b Bytes) GetBase64Data() Base64Data {
	return Base64Data(b)
}

type HexData []byte

func NewHexData(str string) (*HexData, error) {
	if strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X") {
		str = str[2:]
	}
	data, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}
	hexData := HexData(data)
	return &hexData, nil
}

func (h HexData) Data() []byte {
	return h
}
func (h HexData) Length() int {
	return len(h)
}
func (h HexData) String() string {
	return "0x" + hex.EncodeToString(h)
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
		*h = *tmp
	}
	return err
}

func (h HexData) ShortString() string {
	return "0x" + strings.TrimLeft(hex.EncodeToString(h), "0")
}

type Base64Data []byte

func NewBase64Data(str string) (*Base64Data, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}
	b64 := Base64Data(data)
	return &b64, nil
}

func (h Base64Data) Data() []byte {
	return h
}
func (h Base64Data) Length() int {
	return len(h)
}
func (h Base64Data) String() string {
	return base64.StdEncoding.EncodeToString(h)
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
		*h = *tmp
	}
	return err
}

type Base58 []byte

func NewBase58(str string) (*Base58, error) {
	data := Base58(base58.Decode(str))
	return &data, nil
}

func (b Base58) Data() []byte {
	return b
}
func (b Base58) Length() int {
	return len(b)
}
func (b Base58) String() string {
	return base58.Encode(b)
}

func (b Base58) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

func (b *Base58) UnmarshalJSON(data []byte) error {
	str := ""
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	tmp, err := NewBase58(str)
	if err == nil {
		*b = *tmp
	}
	return err
}
