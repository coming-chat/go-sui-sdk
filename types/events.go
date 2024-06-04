package types

import "github.com/W3Tools/go-sui-sdk/v2/sui_types"

type EventId struct {
	TxDigest sui_types.TransactionDigest `json:"txDigest"`
	EventSeq SafeSuiBigInt[uint64]       `json:"eventSeq"`
}

type SuiEvent struct {
	Id EventId `json:"id"`
	// Move package where this event was emitted.
	PackageId sui_types.ObjectID `json:"packageId"`
	// Move module where this event was emitted.
	TransactionModule string `json:"transactionModule"`
	// Sender's Sui sui_types.address.
	Sender sui_types.SuiAddress `json:"sender"`
	// Move event type.
	Type string `json:"type"`
	// Parsed json value of the event
	ParsedJson interface{} `json:"parsedJson,omitempty"`
	// Base 58 encoded bcs bytes of the move event
	Bcs         string                 `json:"bcs"`
	TimestampMs *SafeSuiBigInt[uint64] `json:"timestampMs,omitempty"`
}

type EventFilter struct {
	All *[]EventFilter `json:"All,omitempty"`

	// Events emitted from the specified transaction.
	// digest of the transaction, as base-64 encoded string
	Transaction *sui_types.TransactionDigest `json:"Transaction,omitempty"`

	// Events emitted from the specified Move module.
	MoveModule *MoveModule `json:"MoveModule,omitempty"`

	// Events emitted, defined on the specified Move module.
	MoveEventModule *MoveModule `json:"MoveEventModule,omitempty"`

	// Move struct name of the event
	MoveEvent *string `json:"MoveEvent,omitempty"`

	// Type of event described in Events section
	EventType *string `json:"EventType,omitempty"`

	// Query by sender address
	Sender *sui_types.SuiAddress `json:"Sender,omitempty"`

	// Query by recipient
	Recipient *Recipient `json:"Recipient,omitempty"`

	// Return events associated with the given object
	Object *sui_types.ObjectID `json:"Object,omitempty"`

	// Return events emitted in [start_time, end_time] interval
	TimeRange *TimeRange `json:"TimeRange,omitempty"`
}

type EventPage = Page[SuiEvent, EventId]
