package sui_types

import (
	"github.com/coming-chat/go-sui/v2/lib"
)

type IntentScope struct {
	TransactionData         *lib.EmptyEnum // Used for a user signature on a transaction data.
	TransactionEffects      *lib.EmptyEnum // Used for an authority signature on transaction effects.
	CheckpointSummary       *lib.EmptyEnum // Used for an authority signature on a checkpoint summary.
	PersonalMessage         *lib.EmptyEnum // Used for a user signature on a personal message.
	SenderSignedTransaction *lib.EmptyEnum // Used for an authority signature on a user signed transaction.
	ProofOfPossession       *lib.EmptyEnum // Used as a signature representing an authority's proof of possesion of its authority protocol key.
	HeaderDigest            *lib.EmptyEnum // Used for narwhal authority signature on header digest.
}

func (i IntentScope) IsBcsEnum() {
}

type IntentVersion struct {
	V0 *lib.EmptyEnum
}

func (i IntentVersion) IsBcsEnum() {
}

type AppId struct {
	Sui     *lib.EmptyEnum
	Narwhal *lib.EmptyEnum
}

func (a AppId) IsBcsEnum() {
}

type Intent struct {
	Scope   IntentScope
	Version IntentVersion
	AppId   AppId
}

func DefaultIntent() Intent {
	return Intent{
		Scope: IntentScope{
			TransactionData: &lib.EmptyEnum{},
		},
		Version: IntentVersion{
			V0: &lib.EmptyEnum{},
		},
		AppId: AppId{
			Sui: &lib.EmptyEnum{},
		},
	}
}

type IntentValue interface {
	TransactionData | ~[]byte
}

type IntentMessage[T IntentValue] struct {
	Intent Intent
	Value  T
}

func NewIntentMessage[T IntentValue](intent Intent, value T) IntentMessage[T] {
	return IntentMessage[T]{
		Intent: intent,
		Value:  value,
	}
}
