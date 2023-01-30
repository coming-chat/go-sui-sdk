package types

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestTransactionQuery_MarshalJSON(t1 *testing.T) {
	type fields struct {
		All           *string
		MoveFunction  *MoveFunction
		InputObject   *ObjectId
		MutatedObject *ObjectId
		FromAddress   *Address
		ToAddress     *Address
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "test1",
			fields: fields{
				FromAddress: AddressFromHex(t1, "0x6fc6148816617c3c3eccb1d09e930f73f6712c9c"),
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := TransactionQuery{
				All:           tt.fields.All,
				MoveFunction:  tt.fields.MoveFunction,
				InputObject:   tt.fields.InputObject,
				MutatedObject: tt.fields.MutatedObject,
				FromAddress:   tt.fields.FromAddress,
				ToAddress:     tt.fields.ToAddress,
			}
			got, err := json.Marshal(t)
			if !tt.wantErr(t1, err, fmt.Sprintf("MarshalJSON()")) {
				return
			}
			assert.Equalf(t1, tt.want, got, "MarshalJSON()")
		})
	}
}
