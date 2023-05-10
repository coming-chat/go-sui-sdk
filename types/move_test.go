package types

import (
	"reflect"
	"testing"

	"github.com/coming-chat/go-sui/sui_types"
	"github.com/stretchr/testify/require"
)

func AddressFromHex(t *testing.T, hex string) *suiAddress {
	addr, err := sui_types.NewAddressFromHex(hex)
	require.NoError(t, err)
	return addr
}

func TestNewResourceType(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    *ResourceType
		wantErr bool
	}{
		{
			name: "sample",
			str:  "0x23::coin::Xxxx",
			want: &ResourceType{AddressFromHex(t, "0x23"), "coin", "Xxxx", nil},
		},
		{
			name: "three level",
			str:  "0xabc::Coin::Xxxx<0x789::AAA::ppp<0x111::mod3::func3>>",
			want: &ResourceType{AddressFromHex(t, "0xabc"), "Coin", "Xxxx",
				&ResourceType{AddressFromHex(t, "0x789"), "AAA", "ppp",
					&ResourceType{AddressFromHex(t, "0x111"), "mod3", "func3", nil},
				},
			},
		},
		{
			name:    "error addrress",
			str:     "0x123abcg::coin::Xxxx",
			wantErr: true,
		},
		{
			name:    "error format",
			str:     "0x1::m1::f1<0x2::m2::f2>x",
			wantErr: true,
		},
		{
			name:    "error format2",
			str:     "0x1::m1::f1<<0x3::m3::f3>0x2::m2::f2>",
			wantErr: true,
		},
		{
			name:    "error format3",
			str:     "<0x3::m3::f3>0x1::m1::f1<0x2::m2::f2>",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewResourceType(tt.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewResourceType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResourceType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceType_String(t *testing.T) {
	typeString := "0x1::mmm1::fff1<0x123abcdef::mm2::ff3>"

	resourceType, err := NewResourceType(typeString)
	require.NoError(t, err)
	res := "0x0000000000000000000000000000000000000000000000000000000000000001::mmm1::fff1<0x0000000000000000000000000000000000000000000000000000000123abcdef::mm2::ff3>"
	require.Equal(t, resourceType.String(), res)
}

func TestResourceType_ShortString(t *testing.T) {
	tests := []struct {
		name string
		arg  *ResourceType
		want string
	}{
		{
			arg:  &ResourceType{AddressFromHex(t, "0x1"), "m1", "f1", nil},
			want: "0x1::m1::f1",
		},
		{
			arg: &ResourceType{AddressFromHex(t, "0x1"), "m1", "f1",
				&ResourceType{AddressFromHex(t, "2"), "m2", "f2",
					&ResourceType{AddressFromHex(t, "0x123abcdef"), "m3", "f3", nil}}},
			want: "0x1::m1::f1<0x2::m2::f2<0x123abcdef::m3::f3>>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arg.ShortString(); got != tt.want {
				t.Errorf("ResourceType.ShortString() = %v, want %v", got, tt.want)
			}
		})
	}
}
