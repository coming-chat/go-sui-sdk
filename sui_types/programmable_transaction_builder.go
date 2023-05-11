package sui_types

import (
	"errors"
	"fmt"
	"github.com/coming-chat/go-sui/lib"
	"github.com/coming-chat/go-sui/move_types"
	"github.com/fardream/go-bcs/bcs"
)

type BuilderArg struct {
	Object              *ObjectID
	Pure                *[]uint8
	ForcedNonUniquePure *uint
}

type ProgrammableTransactionBuilder struct {
	Inputs         map[BuilderArg]CallArg
	InputsKeyOrder []BuilderArg
	Commands       []Command
}

func NewProgrammableTransactionBuilder() *ProgrammableTransactionBuilder {
	return &ProgrammableTransactionBuilder{
		Inputs: make(map[BuilderArg]CallArg),
	}
}

func (p *ProgrammableTransactionBuilder) Finish() ProgrammableTransaction {
	var inputs []CallArg
	for _, v := range p.InputsKeyOrder {
		inputs = append(inputs, p.Inputs[v])
	}
	return ProgrammableTransaction{
		Inputs:   inputs,
		Commands: p.Commands,
	}
}

func (p *ProgrammableTransactionBuilder) ForceSeparatePure(value any) (Argument, error) {
	pureData, err := bcs.Marshal(value)
	if err != nil {
		return Argument{}, err
	}
	return p.pureBytes(pureData, true), nil
}

func (p *ProgrammableTransactionBuilder) pureBytes(bytes []byte, forceSeparate bool) Argument {
	var arg BuilderArg
	if forceSeparate {
		length := uint(len(p.Inputs))
		arg = BuilderArg{
			ForcedNonUniquePure: &length,
		}
	} else {
		arg = BuilderArg{
			Pure: &bytes,
		}
	}
	i := p.insertFull(
		arg, CallArg{
			Pure: &bytes,
		},
	)
	return Argument{
		Input: &i,
	}

}

func (p *ProgrammableTransactionBuilder) insertFull(key BuilderArg, value CallArg) uint16 {
	p.Inputs[key] = value
	p.InputsKeyOrder = append(p.InputsKeyOrder, key)
	return uint16(len(p.InputsKeyOrder) - 1)
}
func (p *ProgrammableTransactionBuilder) Pure(value any) (Argument, error) {
	pureData, err := bcs.Marshal(value)
	if err != nil {
		return Argument{}, err
	}
	return p.pureBytes(pureData, false), nil
}

func (p *ProgrammableTransactionBuilder) Obj(objArg ObjectArg) (Argument, error) {
	id := objArg.id()
	var oj ObjectArg
	if oldValue, ok := p.Inputs[BuilderArg{
		Object: &id,
	}]; ok {
		var oldObjArg ObjectArg
		switch {
		case oldValue.Pure != nil:
			return Argument{}, errors.New("invariant violation! object has Pure argument")
		case oldValue.Object != nil:
			oldObjArg = *oldValue.Object
		}

		switch {
		case oldObjArg.SharedObject.InitialSharedVersion == objArg.SharedObject.InitialSharedVersion:
			if oldObjArg.id() != objArg.id() {
				return Argument{}, errors.New("invariant violation! object has id does not match call arg")
			}
			oj = ObjectArg{
				SharedObject: &struct {
					Id                   ObjectID
					InitialSharedVersion SequenceNumber
					Mutable              bool
				}{
					Id:                   id,
					InitialSharedVersion: objArg.SharedObject.InitialSharedVersion,
					Mutable:              oldObjArg.SharedObject.Mutable || objArg.SharedObject.Mutable,
				},
			}
		default:
			if oldObjArg != objArg {
				return Argument{}, fmt.Errorf(
					"mismatched Object argument kind for object %s. "+
						"%v is not compatible with %v", id.String(), oldValue, objArg,
				)
			}
			oj = objArg
		}
	} else {
		oj = objArg
	}
	i := p.insertFull(
		BuilderArg{
			Object: &id,
		}, CallArg{
			Object: &oj,
		},
	)
	return Argument{
		Input: &i,
	}, nil
}

func (p *ProgrammableTransactionBuilder) Input(callArg CallArg) (Argument, error) {
	switch {
	case callArg.Pure != nil:
		return p.pureBytes(*callArg.Pure, false), nil
	case callArg.Object != nil:
		return p.Obj(*callArg.Object)
	default:
		return Argument{}, errors.New("this callArg is nil")
	}
}

func (p *ProgrammableTransactionBuilder) MakeObjList(objs []ObjectArg) error {
	var objArgs []Argument
	for _, v := range objs {
		objArg, err := p.Obj(v)
		if err != nil {
			return err
		}
		objArgs = append(objArgs, objArg)
	}
	p.Command(
		Command{
			MakeMoveVec: &struct {
				TypeTag   *move_types.TypeTag `bcs:"optional"`
				Arguments []Argument
			}{TypeTag: nil, Arguments: objArgs},
		},
	)
	return nil
}

func (p *ProgrammableTransactionBuilder) Command(command Command) Argument {
	p.Commands = append(p.Commands, command)
	i := uint16(len(p.Commands)) - 1
	return Argument{
		Result: &i,
	}
}

func (p *ProgrammableTransactionBuilder) TransferObject(
	recipient SuiAddress,
	objectRefs []*ObjectRef,
) error {
	recArg, err := p.Pure(recipient)
	if err != nil {
		return err
	}
	var objArgs []Argument
	for _, v := range objectRefs {
		objArg, err := p.Obj(
			ObjectArg{
				ImmOrOwnedObject: v,
			},
		)
		if err != nil {
			return err
		}
		objArgs = append(objArgs, objArg)
	}
	p.Command(
		Command{
			TransferObjects: &struct {
				Arguments []Argument
				Argument  Argument
			}{Arguments: objArgs, Argument: recArg},
		},
	)
	return nil
}

func (p *ProgrammableTransactionBuilder) TransferSui(recipient SuiAddress, amount *uint64) error {
	recArg, err := p.Pure(recipient)
	if err != nil {
		return err
	}
	var coinArg Argument
	if amount == nil {
		coinArg = Argument{
			GasCoin: &lib.EmptyEnum{},
		}
	} else {
		amtArg, err := p.Pure(*amount)
		if err != nil {
			return err
		}
		coinArg = p.Command(
			Command{
				SplitCoins: &struct {
					Argument  Argument
					Arguments []Argument
				}{
					Argument: Argument{
						GasCoin: &lib.EmptyEnum{},
					}, Arguments: []Argument{
						amtArg,
					},
				},
			},
		)
	}
	p.Command(
		Command{
			TransferObjects: &struct {
				Arguments []Argument
				Argument  Argument
			}{
				Arguments: []Argument{
					coinArg,
				}, Argument: recArg,
			},
		},
	)
	return nil
}

func (p *ProgrammableTransactionBuilder) MoveCall(
	packageID ObjectID,
	module move_types.Identifier,
	function move_types.Identifier,
	typeArguments []move_types.TypeTag,
	callArgs []CallArg,
) error {
	var arguments []Argument
	for _, v := range callArgs {
		argument, err := p.Input(v)
		if err != nil {
			return err
		}
		arguments = append(arguments, argument)
	}
	p.Command(
		Command{
			MoveCall: &ProgrammableMoveCall{
				Package:       packageID,
				Module:        module,
				Function:      function,
				TypeArguments: typeArguments,
				Arguments:     arguments,
			},
		},
	)
	return nil
}
