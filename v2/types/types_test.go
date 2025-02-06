package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAddressFromHex(t *testing.T) {
	addr := AddressFromHex(t, "0x2")
	require.Equal(t, addr.ShortString(), "0x2")
	t.Log(addr)
}

func TestObjectOwnerJsonENDE(t *testing.T) {
	{
		var dataStruct struct {
			Owner *ObjectOwner `json:"owner"`
		}
		jsonString := []byte(`{"owner":"Immutable"}`)

		err := json.Unmarshal(jsonString, &dataStruct)
		require.Nil(t, err)
		enData, err := json.Marshal(dataStruct)
		require.Nil(t, err)
		require.Equal(t, jsonString, enData)
	}
	{
		var dataStruct struct {
			Owner *ObjectOwner `json:"owner"`
		}
		jsonString := []byte(`{"owner":{"AddressOwner":"0xfb1f678fcfe31c7c1924319e49614ffbe3a984842ceed559aa2d772e60a2ef8f"}}`)

		err := json.Unmarshal(jsonString, &dataStruct)
		require.Nil(t, err)
		enData, err := json.Marshal(dataStruct)
		require.Nil(t, err)
		require.Equal(t, jsonString, enData)
	}
}

func TestTransactionQuery_MarshalJSON(t1 *testing.T) {
	// var all = ""
	// type fields struct {
	// 	All           *string
	// 	MoveFunction  *MoveFunction
	// 	InputObject   *ObjectId
	// 	MutatedObject *ObjectId
	// 	FromAddress   *Address
	// 	ToAddress     *Address
	// }
	// tests := []struct {
	// 	name    string
	// 	fields  fields
	// 	want    []byte
	// 	wantErr assert.ErrorAssertionFunc
	// }{
	// 	{
	// 		name: "test1",
	// 		fields: fields{
	// 			All: &all,
	// 		},
	// 	},
	// }
	// for _, tt := range tests {
	// 	t1.Run(tt.name, func(t1 *testing.T) {
	// 		t := TransactionQuery{
	// 			All:           tt.fields.All,
	// 			MoveFunction:  tt.fields.MoveFunction,
	// 			InputObject:   tt.fields.InputObject,
	// 			MutatedObject: tt.fields.MutatedObject,
	// 			FromAddress:   tt.fields.FromAddress,
	// 			ToAddress:     tt.fields.ToAddress,
	// 		}
	// 		got, err := json.Marshal(t)
	// 		require.NoError(t1, err)
	// 		t1.Logf("%#v", got)
	// 	})
	// }
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
				"0x000000123",
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
