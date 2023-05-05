package sui_types

type EmptyEnum struct {
}

func (e EmptyEnum) MarshalBCS() ([]byte, error) {
	return []byte{}, nil
}

type IntentScope struct {
	TransactionData         *EmptyEnum // Used for a user signature on a transaction data.
	TransactionEffects      *EmptyEnum // Used for an authority signature on transaction effects.
	CheckpointSummary       *EmptyEnum // Used for an authority signature on a checkpoint summary.
	PersonalMessage         *EmptyEnum // Used for a user signature on a personal message.
	SenderSignedTransaction *EmptyEnum // Used for an authority signature on a user signed transaction.
	ProofOfPossession       *EmptyEnum // Used as a signature representing an authority's proof of possesion of its authority protocol key.
	HeaderDigest            *EmptyEnum // Used for narwhal authority signature on header digest.
}

func (i IntentScope) IsBcsEnum() {
}

type IntentVersion struct {
	V0 *EmptyEnum
}

func (i IntentVersion) IsBcsEnum() {
}

type AppId struct {
	Sui     *EmptyEnum
	Narwhal *EmptyEnum
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
			TransactionData: &EmptyEnum{},
		},
		Version: IntentVersion{
			V0: &EmptyEnum{},
		},
		AppId: AppId{
			Sui: &EmptyEnum{},
		},
	}
}

type IntentValue interface {
	// | TransactionData
	~[]byte
	MarshalBCS() ([]byte, error)
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
