package types

type EventId struct {
	TxDigest TransactionDigest `json:"txDigest"`
	EventSeq SequenceNumber    `json:"eventSeq"`
}

type SuiEvent struct {
	Id EventId `json:"id"`
	// Move package where this event was emitted.
	PackageId ObjectId `json:"packageId"`
	// Move module where this event was emitted.
	TransactionModule string `json:"transactionModule"`
	// Sender's Sui address.
	Sender Address `json:"sender"`
	// Move event type.
	Type string `json:"type"`
	// Parsed json value of the event
	ParsedJson interface{} `json:"parsedJson,omitempty"`
	// Base 58 encoded bcs bytes of the move event
	Bcs         string `json:"bcs"`
	TimestampMs *int64 `json:"timestampMs,omitempty"`
}
