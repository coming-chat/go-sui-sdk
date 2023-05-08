package lib

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type TagJsonType interface {
	Tag() string
	Content() string
}

type TagJson[T TagJsonType] struct {
	Data T
}

func (t *TagJson[T]) UnmarshalJSON(data []byte) error {
	rv := reflect.ValueOf(t).Elem().Field(0)
	if t.Data.Tag() == "" {
		var tmp string
		err := json.Unmarshal(data, &tmp)
		if err != nil {
			return err
		}
		for i := 0; i < rv.Type().NumField(); i++ {
			if rv.Field(i).Type().Name() == tmp && rv.Field(i).IsNil() {
				rv.Field(i).Set(reflect.New(rv.Field(i).Type().Elem()))
				break
			}
		}
		return nil
	}
	tmp := make(map[string]json.RawMessage)
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
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
