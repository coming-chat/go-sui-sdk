package types

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	var all = ""
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
				All: &all,
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
			require.NoError(t1, err)
			t1.Logf("%#v", got)
		})
	}
}

func TestIsSameStringAddress(t *testing.T) {
	type args struct {
		addr1 string
		addr2 string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "same address",
			args: args{
				"0x00000123",
				"0x00000000000123",
			},
			want: true,
		},
		{
			name: "not same address",
			args: args{
				"0x123f",
				"0x00000000123",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSameStringAddress(tt.args.addr1, tt.args.addr2); got != tt.want {
				t.Errorf("IsSameStringAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
