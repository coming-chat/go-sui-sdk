package sui_types

//
//import (
//	"errors"
//	"github.com/coming-chat/go-sui/types"
//	"github.com/fardream/go-bcs/bcs"
//	"reflect"
//)
//
//type MoveCallArg []any
//
//func (m MoveCallArg) GetMoveCallArgs() ([]*CallArg, error) {
//	var callArgs []*CallArg
//	for _, v := range m {
//		arg, err := ArgToCallArg(reflect.ValueOf(v))
//		if err != nil {
//			return nil, err
//		}
//		callArgs = append(callArgs, arg)
//	}
//	return callArgs, nil
//}
//
//func ArgToCallArg(value reflect.Value) (*CallArg, error) {
//	if !value.CanInterface() {
//		return nil, errors.New("this field is not exported")
//	}
//	switch value.Type().Kind() {
//	case reflect.Pointer:
//		return ArgToCallArg(reflect.Indirect(value))
//	case reflect.Slice, reflect.Array:
//		var objectArgs []*ObjectArg
//		for i := 0; i < value.Len(); i++ {
//			arg, err := ArgToCallArg(value.Index(i))
//			if err != nil {
//				return nil, err
//			}
//			if arg.Object != nil {
//				objectArgs = append(objectArgs, arg.Object)
//			} else {
//				break
//			}
//		}
//		if len(objectArgs) != 0 {
//			return &CallArg{
//				ObjVec: objectArgs,
//			}, nil
//		}
//	case reflect.Struct:
//		if objectRef, ok := (value.Interface()).(types.ObjectRef); ok {
//			return &CallArg{
//				Object: &ObjectArg{
//					ImmOrOwnedObject: &objectRef,
//				},
//			}, nil
//		}
//		if sharedObject, ok := (value.Interface()).(SharedObject); ok {
//			return &CallArg{
//				Object: &ObjectArg{
//					SharedObject: &sharedObject,
//				},
//			}, nil
//		}
//	default:
//	}
//
//	data, err := bcs.Marshal(value.Interface())
//	if err != nil {
//		return nil, err
//	}
//	return &CallArg{
//		Pure: &data,
//	}, nil
//}
