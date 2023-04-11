package types

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/shopspring/decimal"
)

type TransactionDigest = string

type TransactionEffectsDigest = string

type TransactionEventDigest = string

type SequenceNumber = uint64

type SuiBigInt = decimal.Decimal

type SafeSuiBigInt[T ~int64 | ~uint64] struct {
	data T
}

func (s *SafeSuiBigInt[T]) UnmarshalJSON(data []byte) error {
	num := decimal.NewFromInt(0)
	err := num.UnmarshalJSON(data)
	if err != nil {
		return err
	}

	if num.BigInt().IsInt64() {
		s.data = T(num.BigInt().Int64())
		return nil
	}

	if num.BigInt().IsUint64() {
		s.data = T(num.BigInt().Uint64())
		return nil
	}
	return fmt.Errorf("json data [%s] is not T", string(data))
}

func (s SafeSuiBigInt[T]) MarshalJSON() ([]byte, error) {
	return decimal.NewFromInt(int64(s.data)).MarshalJSON()
}

func (s SafeSuiBigInt[T]) Int64() int64 {
	return int64(s.data)
}

func (s SafeSuiBigInt[T]) Uint64() uint64 {
	return uint64(s.data)
}

// export const ObjectId = string();
// export type ObjectId = Infer<typeof ObjectId>;

// export const SuiAddress = string();
// export type SuiAddress = Infer<typeof SuiAddress>;

type ObjectOwnerInternal struct {
	AddressOwner *Address `json:"AddressOwner,omitempty"`
	ObjectOwner  *Address `json:"ObjectOwner,omitempty"`
	SingleOwner  *Address `json:"SingleOwner,omitempty"`
	Shared       *struct {
		InitialSharedVersion SequenceNumber `json:"initial_shared_version"`
	} `json:"Shared,omitempty"`
}

type ObjectOwner struct {
	*ObjectOwnerInternal
	*string
}

type TagJsonType interface {
	Tag() string
	Content() string
}

type TagJson[T TagJsonType] struct {
	Data T
}

func (t *TagJson[T]) UnmarshalJSON(data []byte) error {
	tmp := make(map[string]json.RawMessage)
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	rv := reflect.ValueOf(t).Elem().Field(0)
	v, ok := tmp[t.Data.Tag()]
	if !ok {
		return fmt.Errorf("no such tag: %s in json data %v", t.Data.Tag(), tmp)
	}
	var subType string
	err = json.Unmarshal(v, &subType)
	if err != nil {
		return fmt.Errorf("the tag [%s] value is not string", t.Data.Tag())
	}
	for i := 0; i < rv.Type().NumField(); i++ {
		if !strings.Contains(rv.Type().Field(i).Tag.Get("json"), subType) {
			continue
		}
		if rv.Field(i).Kind() != reflect.Pointer {
			return fmt.Errorf("field %s not pointer", rv.Field(i).Type().Name())
		}
		if rv.Field(i).IsNil() {
			rv.Field(i).Set(reflect.New(rv.Field(i).Type().Elem()))
		}
		jsonData := data
		if t.Data.Content() != "" {
			jsonData, ok = tmp[t.Data.Content()]
			if !ok {
				return fmt.Errorf("json data [%v] get content key [%s] failed", tmp, t.Data.Content())
			}
		}
		err = json.Unmarshal(jsonData, rv.Field(i).Interface())
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("no tag[%s] value <%s> in struct fields", t.Data.Tag(), v)
}

type Page[T SuiTransactionBlockResponse | SuiEvent | Coin | SuiObjectResponse, C TransactionDigest | EventId | ObjectId] struct {
	Data        []T  `json:"data"`
	NextCursor  *C   `json:"nextCursor,omitempty"`
	HasNextPage bool `json:"hasNextPage"`
}
