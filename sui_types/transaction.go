package sui_types

import "github.com/coming-chat/go-sui/lib"

func NewProgrammableAllowSponsor(
	sender SuiAddress,
	gasPayment []*ObjectRef,
	pt ProgrammableTransaction,
	gasBudge,
	gasPrice uint64,
	sponsor SuiAddress,
) TransactionData {
	kind := TransactionKind{
		ProgrammableTransaction: &pt,
	}
	return newWithGasCoinsAllowSponsor(kind, sender, gasPayment, gasBudge, gasPrice, sponsor)
}

func NewProgrammable(
	sender SuiAddress,
	gasPayment []*ObjectRef,
	pt ProgrammableTransaction,
	gasBudget uint64,
	gasPrice uint64,
) TransactionData {
	return NewProgrammableAllowSponsor(sender, gasPayment, pt, gasBudget, gasPrice, sender)
}

func newWithGasCoinsAllowSponsor(
	kind TransactionKind,
	sender SuiAddress,
	gasPayment []*ObjectRef,
	gasBudget uint64,
	gasPrice uint64,
	gasSponsor SuiAddress,
) TransactionData {
	return TransactionData{
		V1: &TransactionDataV1{
			Kind:   kind,
			Sender: sender,
			GasData: GasData{
				Price:   gasPrice,
				Owner:   gasSponsor,
				Payment: gasPayment,
				Budget:  gasBudget,
			},
			Expiration: TransactionExpiration{
				None: &lib.EmptyEnum{},
			},
		},
	}
}
