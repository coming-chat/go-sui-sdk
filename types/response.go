package types

type AuthSignInfo interface{}

type CertifiedTransaction struct {
	TransactionDigest string        `json:"transactionDigest"`
	TxSignature       string        `json:"txSignature"`
	AuthSignInfo      *AuthSignInfo `json:"authSignInfo"`

	Data *SenderSignedData `json:"data"`
}

type GasCostSummary struct {
	ComputationCost uint64 `json:"computationCost"`
	StorageCost     uint64 `json:"storageCost"`
	StorageRebate   uint64 `json:"storageRebate"`
}

const (
	TransactionStatusSuccess = "success"
	TransactionStatusFailure = "failure"
)

type TransactionStatus struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type TransactionEffects struct {
	Status TransactionStatus `json:"status"`

	TransactionDigest string          `json:"transactionDigest"`
	GasUsed           *GasCostSummary `json:"gasUsed"`
	GasObject         *OwnedObjectRef `json:"gasObject"`
	Events            []Event         `json:"events,omitempty"`
	Dependencies      []string        `json:"dependencies,omitempty"`

	// SharedObjects []ObjectRef      `json:"sharedObjects"`
	Created   []OwnedObjectRef `json:"created,omitempty"`
	Mutated   []OwnedObjectRef `json:"mutated,omitempty"`
	Unwrapped []OwnedObjectRef `json:"unwrapped,omitempty"`
	Deleted   []ObjectRef      `json:"deleted,omitempty"`
	Wrapped   []ObjectRef      `json:"wrapped,omitempty"`
}

func (te *TransactionEffects) GasFee() uint64 {
	return te.GasUsed.StorageCost - te.GasUsed.StorageRebate + te.GasUsed.ComputationCost
}

type ParsedTransactionResponse interface{}

// TxnRequestTypeImmediateReturn       ExecuteTransactionRequestType = "ImmediateReturn"
// TxnRequestTypeWaitForTxCert         ExecuteTransactionRequestType = "WaitForTxCert"
// TxnRequestTypeWaitForEffectsCert    ExecuteTransactionRequestType = "WaitForEffectsCert"

type TransactionResponse struct {
	Certificate *CertifiedTransaction     `json:"certificate"`
	Effects     *TransactionEffects       `json:"effects"`
	ParsedData  ParsedTransactionResponse `json:"parsed_data,omitempty"`
	TimestampMs uint64                    `json:"timestamp_ms,omitempty"`
}

type ExecuteTransactionEffects struct {
	TransactionEffectsDigest Digest `json:"transactionEffectsDigest"`

	Effects      TransactionEffects `json:"effects"`
	AuthSignInfo *AuthSignInfo      `json:"authSignInfo"`
}

type ExecuteTransactionResponse struct {
	ImmediateReturn *struct {
		TransactionDigest string `json:"tx_digest"`
	} `json:"ImmediateReturn,omitempty"`

	TxCert *struct {
		Certificate CertifiedTransaction `json:"certificate"`
	} `json:"TxCert,omitempty"`

	EffectsCert *struct {
		Certificate CertifiedTransaction      `json:"certificate"`
		Effects     ExecuteTransactionEffects `json:"effects"`

		ConfirmedLocalExecution bool `json:"confirmed_local_execution"`
	} `json:"EffectsCert,omitempty"`
}

func (r *ExecuteTransactionResponse) TransactionDigest() string {
	switch {
	case r.ImmediateReturn != nil:
		return r.ImmediateReturn.TransactionDigest
	case r.TxCert != nil:
		return r.TxCert.Certificate.TransactionDigest
	case r.EffectsCert != nil:
		return r.EffectsCert.Certificate.TransactionDigest
	}
	return ""
}
