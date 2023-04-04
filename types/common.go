package types

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type TransactionDigest = string

type TransactionEffectsDigest = string

type TransactionEventDigest = string

type SequenceNumber = uint64

// export const ObjectId = string();
// export type ObjectId = Infer<typeof ObjectId>;

// export const SuiAddress = string();
// export type SuiAddress = Infer<typeof SuiAddress>;

type ObjectOwnerInternal struct {
	AddressOwner *Address `json:"AddressOwner,omitempty"`
	ObjectOwner  *Address `json:"ObjectOwner,omitempty"`
	SingleOwner  *Address `json:"SingleOwner,omitempty"`
	Shared       *struct {
		InitialSharedVersion uint64 `json:"initial_shared_version"`
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
	tmp := make(map[string]interface{})
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	rv := reflect.ValueOf(t).Elem().Field(0)
	v, ok := tmp[t.Data.Tag()]
	if !ok {
		return fmt.Errorf("no such tag: %s in json data %v", t.Data.Tag(), tmp)
	}
	for i := 0; i < rv.Type().NumField(); i++ {
		if !strings.Contains(rv.Type().Field(i).Tag.Get("json"), v.(string)) {
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
			value, has := tmp[t.Data.Content()]
			if !has {
				return fmt.Errorf("json data [%v] get content key [%s] failed", tmp, t.Data.Content())
			}
			jsonData, err = json.Marshal(value)
			if err != nil {
				return err
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
