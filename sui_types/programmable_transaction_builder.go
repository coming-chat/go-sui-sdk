package sui_types

import (
	"errors"
	"fmt"
	"github.com/coming-chat/go-sui/lib"
	"github.com/coming-chat/go-sui/move_types"
	"github.com/fardream/go-bcs/bcs"
	"github.com/mitchellh/hashstructure/v2"
	"strconv"
)

type BuilderArg struct {
	Object              *ObjectID
	Pure                *[]uint8
	ForcedNonUniquePure *uint
}

func (b BuilderArg) String() string {
	hash, err := hashstructure.Hash(b, hashstructure.FormatV2, nil)
	if err != nil {
		panic(err)
	}
	return strconv.FormatUint(hash, 10)
}

type ProgrammableTransactionBuilder struct {
	Inputs         map[string]CallArg //maybe has hash clash
	InputsKeyOrder []BuilderArg
	Commands       []Command
}

func NewProgrammableTransactionBuilder() *ProgrammableTransactionBuilder {
	return &ProgrammableTransactionBuilder{
		Inputs: make(map[string]CallArg),
	}
}

func (p *ProgrammableTransactionBuilder) Finish() ProgrammableTransaction {
	var inputs []CallArg
	for _, v := range p.InputsKeyOrder {
		inputs = append(inputs, p.Inputs[v.String()])
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
	_, ok := p.Inputs[key.String()]
	p.Inputs[key.String()] = value
	if !ok {
		p.InputsKeyOrder = append(p.InputsKeyOrder, key)
		return uint16(len(p.InputsKeyOrder) - 1)
	}
	for i, v := range p.InputsKeyOrder {
		if v.String() == key.String() {
			return uint16(i)
		}
	}
	return 0
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
	}.String()]; ok {
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

func (p *ProgrammableTransactionBuilder) PaySui(
	recipients []SuiAddress,
	amounts []uint64,
) error {
	return p.PayMulInternal(
		recipients, amounts, Argument{
			GasCoin: &lib.EmptyEnum{},
		},
	)
}

func (p *ProgrammableTransactionBuilder) PayAllSui(recipient SuiAddress) error {
	recArg, err := p.Pure(recipient)
	if err != nil {
		return err
	}
	p.Command(
		Command{
			TransferObjects: &struct {
				Arguments []Argument
				Argument  Argument
			}{Arguments: []Argument{{GasCoin: &lib.EmptyEnum{}}}, Argument: recArg},
		},
	)
	return nil
}

func (p *ProgrammableTransactionBuilder) Pay(
	coins []*ObjectRef,
	recipients []SuiAddress,
	amounts []uint64,
) error {
	if len(coins) == 0 {
		return errors.New("coins is empty")
	}
	coinArg, err := p.Obj(
		ObjectArg{
			ImmOrOwnedObject: coins[0],
		},
	)
	coins = coins[1:]
	if err != nil {
		return err
	}
	var mergeArgs []Argument
	for _, v := range coins {
		mergeCoin, err := p.Obj(
			ObjectArg{
				ImmOrOwnedObject: v,
			},
		)
		if err != nil {
			return err
		}
		mergeArgs = append(mergeArgs, mergeCoin)
	}
	if len(mergeArgs) != 0 {
		p.Command(
			Command{
				MergeCoins: &struct {
					Argument  Argument
					Arguments []Argument
				}{Argument: coinArg, Arguments: mergeArgs},
			},
		)
	}
	return p.PayMulInternal(recipients, amounts, coinArg)
}

func (p *ProgrammableTransactionBuilder) PayMulInternal(
	recipients []SuiAddress,
	amounts []uint64, coin Argument,
) error {
	if len(recipients) != len(amounts) {
		return fmt.Errorf(
			"recipients and amounts mismatch. Got %d recipients but %d amounts",
			len(recipients),
			len(amounts),
		)
	}
	if len(amounts) == 0 {
		return nil
	}
	var (
		amtArgs              []Argument
		recipientMap         = make(map[SuiAddress][]int)
		recipientMapKeyIndex []SuiAddress
	)
	for i := 0; i < len(amounts); i++ {
		amt, err := p.Pure(amounts[i])
		if err != nil {
			return err
		}
		recipientMap[recipients[i]] = append(recipientMap[recipients[i]], i)
		if len(recipientMap[recipients[i]]) == 1 {
			recipientMapKeyIndex = append(recipientMapKeyIndex, recipients[i])
		}
		amtArgs = append(amtArgs, amt)
	}
	splitCoinResult := p.Command(
		Command{
			SplitCoins: &struct {
				Argument  Argument
				Arguments []Argument
			}{Argument: coin, Arguments: amtArgs},
		},
	)
	if splitCoinResult.Result == nil {
		return errors.New("self.command should always give a Argument::Result")
	}
	for _, v := range recipientMapKeyIndex {
		recArg, err := p.Pure(v)
		if err != nil {
			return err
		}
		var coins []Argument
		for _, j := range recipientMap[v] {
			sdCoin := Argument{
				NestedResult: &struct {
					Result1 uint16
					Result2 uint16
				}{Result1: *splitCoinResult.Result, Result2: uint16(j)},
			}
			coins = append(coins, sdCoin)
		}
		p.Command(
			Command{
				TransferObjects: &struct {
					Arguments []Argument
					Argument  Argument
				}{Arguments: coins, Argument: recArg},
			},
		)
	}
	return nil
}
