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
	TransactionEffectsDigest string `json:"transactionEffectsDigest"`

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

type SuiCoinMetadata struct {
	Decimals    uint8    `json:"decimals"`
	Description string   `json:"description"`
	IconUrl     string   `json:"iconUrl,omitempty"`
	Id          ObjectId `json:"id"`
	Name        string   `json:"name"`
	Symbol      string   `json:"symbol"`
}

type SuiCoinBalance struct {
	CoinType        string `json:"coinType"`
	CoinObjectCount int64  `json:"coinObjectCount"`
	TotalBalance    int64  `json:"totalBalance"`
}

type DevInspectResults struct {
	Effects TransactionEffects `json:"effects"`
	Results DevInspectResult   `json:"results"`
}

type DevInspectResult struct {
	Err string `json:"Err,omitempty"`
	Ok  any    `json:"Ok,omitempty"` //Result_of_Array_of_Tuple_of_uint_and_SuiExecutionResult_or_String
}

type CoinPage struct {
	Data       []CoinObject `json:"data"`
	NextCursor string       `json:"nextCursor"`
}

type CoinObject struct {
	CoinType     string   `json:"coinType"`
	CoinObjectId ObjectId `json:"coinObjectId"`
	Version      int64    `json:"version"`
	Digest       string   `json:"digest"`
	Balance      int64    `json:"balance"`
}

type Supply struct {
	Value uint64 `json:"value"`
}
