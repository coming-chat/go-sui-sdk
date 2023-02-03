package types

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/fardream/go-bcs/bcs"
	"strings"
	// "github.com/btcsuite/btcutil/base58"
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

func (a HexData) Data() []byte {
	return a.data
}
func (a HexData) Length() int {
	return len(a.data)
}
func (a HexData) String() string {
	return "0x" + hex.EncodeToString(a.data)
}

func (a HexData) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

func (a *HexData) UnmarshalJSON(data []byte) error {
	str := ""
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}
	tmp, err := NewHexData(str)
	if err == nil {
		a.data = tmp.data
	}
	return err
}

func (a HexData) MarshalBCS() ([]byte, error) {
	return a.data, nil
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

func (h Base64Data) MarshalBCS() ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	en := bcs.NewEncoder(buffer)
	err := en.Encode(h.data)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// type Base58Data struct {
// 	data []byte
// }

// func NewBase58Data(str string) (*Base58Data, error) {
// 	data := base58.Decode(str)
// 	return &Base58Data{data}, nil
// }

// func (h Base58Data) Data() []byte {
// 	return h.data
// }
// func (h Base58Data) Length() int {
// 	return len(h.data)
// }
// func (h Base58Data) String() string {
// 	return base58.Encode(h.data)
// }

// func (h Base58Data) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(h.String())
// }

// func (h *Base58Data) UnmarshalJSON(data []byte) error {
// 	str := ""
// 	err := json.Unmarshal(data, &str)
// 	if err != nil {
// 		return err
// 	}
// 	tmp, err := NewBase58Data(str)
// 	if err == nil {
// 		h.data = tmp.data
// 	}
// 	return err
// }
