package types

type AuthSignInfo interface{}

type CertifiedTransaction struct {
	TransactionDigest string        `json:"transactionDigest"`
	TxSignature       string        `json:"txSignature"`
	AuthSignInfo      *AuthSignInfo `json:"authSignInfo"`

	Data *SenderSignedData `json:"data"`
}

type ParsedTransactionResponse interface{}

type ExecuteTransactionEffects struct {
	TransactionEffectsDigest string `json:"transactionEffectsDigest"`

	Effects      TagJson[SuiTransactionBlockEffects] `json:"effects"`
	AuthSignInfo *AuthSignInfo                       `json:"authSignInfo"`
}

type ExecuteTransactionResponse struct {
	Certificate CertifiedTransaction      `json:"certificate"`
	Effects     ExecuteTransactionEffects `json:"effects"`

	ConfirmedLocalExecution bool `json:"confirmed_local_execution"`
}

func (r *ExecuteTransactionResponse) TransactionDigest() string {
	return r.Certificate.TransactionDigest
}

type SuiCoinMetadata struct {
	Decimals    uint8    `json:"decimals"`
	Description string   `json:"description"`
	IconUrl     string   `json:"iconUrl,omitempty"`
	Id          ObjectId `json:"id"`
	Name        string   `json:"name"`
	Symbol      string   `json:"symbol"`
}

type DevInspectResult struct {
	Err string `json:"Err,omitempty"`
	Ok  any    `json:"Ok,omitempty"` //Result_of_Array_of_Tuple_of_uint_and_SuiExecutionResult_or_String
}
