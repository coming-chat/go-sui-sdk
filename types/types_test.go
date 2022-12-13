package types

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAddressFromHex(t *testing.T) {
	addr, err := NewAddressFromHex("0x2")
	assert.Nil(t, err)

	t.Log(addr)
}

func TestObjectOwnerJsonENDE(t *testing.T) {
	dataStruct1 := struct {
		Owner *ObjectOwner `json:"owner"`
	}{}

	dataStruct2 := struct {
		Owner *ObjectOwner `json:"owner"`
	}{}
	jsonString1 := []byte(`{"owner":"Immutable"}`)

	jsonString2 := []byte(`{"owner":{"AddressOwner":"0x4a13f6340026019b39aaf0fb24b29e149c092a0e"}}`)

	err := json.Unmarshal(jsonString1, &dataStruct1)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(jsonString2, &dataStruct2)
	if err != nil {
		t.Fatal(err)
	}

	enData1, err := json.Marshal(dataStruct1)
	if err != nil {
		t.Fatal(err)
	}

	enData2, err := json.Marshal(dataStruct2)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(jsonString1, enData1) {
		t.Fatal("encode failed")
	}

	if !bytes.Equal(jsonString2, enData2) {
		t.Fatal("encode failed")
	}
}
