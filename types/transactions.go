package types

import (
	"github.com/shopspring/decimal"
)

type ExecuteTransactionRequestType string

const (
	TxnRequestTypeWaitForEffectsCert    ExecuteTransactionRequestType = "WaitForEffectsCert"
	TxnRequestTypeWaitForLocalExecution ExecuteTransactionRequestType = "WaitForLocalExecution"
)

type EpochId = string

type GasCostSummary struct {
	ComputationCost         decimal.Decimal `json:"computationCost"`
	StorageCost             decimal.Decimal `json:"storageCost"`
	StorageRebate           decimal.Decimal `json:"storageRebate"`
	NonRefundableStorageFee decimal.Decimal `json:"nonRefundableStorageFee"`
}

const (
	ExecutionStatusSuccess = "success"
	ExecutionStatusFailure = "failure"
)

type ExecutionStatus struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type OwnedObjectRef struct {
	Owner     ObjectOwner  `json:"owner"`
	Reference SuiObjectRef `json:"reference"`
}

type TransactionEffectsModifiedAtVersions struct {
	ObjectId       ObjectId       `json:"objectId"`
	SequenceNumber SequenceNumber `json:"sequenceNumber"`
}

type TransactionEffects struct {
	MessageVersion string `json:"messageVersion"`

	/** The status of the execution */
	Status ExecutionStatus `json:"status"`
	/** The epoch when this transaction was executed */
	ExecutedEpoch EpochId `json:"executedEpoch"`
	/** The version that every modified (mutated or deleted) object had before it was modified by this transaction. **/
	ModifiedAtVersions []TransactionEffectsModifiedAtVersions `json:"modifiedAtVersions,omitempty"`
	GasUsed            GasCostSummary                         `json:"gasUsed"`
	/** The object references of the shared objects used in this transaction. Empty if no shared objects were used. */
	SharedObjects []SuiObjectRef `json:"sharedObjects,omitempty"`
	/** The transaction digest */
	TransactionDigest TransactionDigest `json:"transactionDigest"`
	/** ObjectRef and owner of new objects created */
	Created []OwnedObjectRef `json:"created,omitempty"`
	/** ObjectRef and owner of mutated objects, including gas object */
	Mutated []OwnedObjectRef `json:"mutated,omitempty"`
	/**
	 * ObjectRef and owner of objects that are unwrapped in this transaction.
	 * Unwrapped objects are objects that were wrapped into other objects in the past,
	 * and just got extracted out.
	 */
	Unwrapped []OwnedObjectRef `json:"unwrapped,omitempty"`
	/** Object Refs of objects now deleted (the old refs) */
	Deleted []SuiObjectRef `json:"deleted,omitempty"`
	/** Object Refs of objects now deleted (the old refs) */
	UnwrappedThenDeleted []SuiObjectRef `json:"unwrapped_then_deleted,omitempty"`
	/** Object refs of objects now wrapped in other objects */
	Wrapped []SuiObjectRef `json:"wrapped,omitempty"`
	/**
	 * The updated gas object reference. Have a dedicated field for convenient access.
	 * It's also included in mutated.
	 */
	GasObject OwnedObjectRef `json:"gasObject"`
	/** The events emitted during execution. Note that only successful transactions emit events */
	EventsDigest *TransactionEventDigest `json:"eventsDigest,omitempty"`
	/** The set of transaction digests this transaction depends on */
	Dependencies []TransactionDigest `json:"dependencies,omitempty"`
}

func (te *TransactionEffects) GasFee() uint64 {
	feeInt := te.GasUsed.StorageCost.Sub(te.GasUsed.StorageRebate).Add(te.GasUsed.ComputationCost)
	return feeInt.BigInt().Uint64()
}

func (te *TransactionEffects) IsSuccess() bool {
	return te.Status.Status == ExecutionStatusSuccess
}

type TransactionEvents = []SuiEvent

const (
	SuiTransactionBlockKindSuiChangeEpoch             = "ChangeEpoch"
	SuiTransactionBlockKindSuiConsensusCommitPrologue = "ConsensusCommitPrologue"
	SuiTransactionBlockKindGenesis                    = "Genesis"
	SuiTransactionBlockKindProgrammableTransaction    = "ProgrammableTransaction"
)

type SuiTransactionBlockKind interface{}

type SuiTransactionBlockData struct {
	MessageVersion string                  `json:"messageVersion"`
	Transaction    SuiTransactionBlockKind `json:"transaction"`
	Sender         Address                 `json:"sender"`
	GasData        SuiGasData              `json:"gasData"`
}

type SuiTransactionBlock struct {
	Data         SuiTransactionBlockData `json:"data"`
	TxSignatures []string                `json:"txSignatures"`
}

type SuiObjectChange interface{}

type BalanceChange struct {
	Owner    ObjectOwner `json:"owner"`
	CoinType string      `json:"coinType"`
	/* Coin balance change(positive means receive, negative means send) */
	Amount string `json:"amount"`
}

type SuiTransactionBlockResponse struct {
	Digest                  TransactionDigest    `json:"digest,omitempty"`
	Transaction             *SuiTransactionBlock `json:"transaction,omitempty"`
	Effects                 *TransactionEffects  `json:"effects,omitempty"`
	Events                  TransactionEvents    `json:"events,omitempty"`
	TimestampMs             *int64               `json:"timestampMs,omitempty"`
	Checkpoint              *int64               `json:"checkpoint,omitempty"`
	ConfirmedLocalExecution *bool                `json:"confirmedLocalExecution,omitempty"`
	ObjectChanges           []SuiObjectChange    `json:"objectChanges,omitempty"`
	BalanceChanges          []BalanceChange      `json:"balanceChanges,omitempty"`
	/* Errors that occurred in fetching/serializing the transaction. */
	Errors []string `json:"errors,omitempty"`
}

type ReturnValueType interface{}
type MutableReferenceOutputType interface{}
type ExecutionResultType struct {
	MutableReferenceOutputs []MutableReferenceOutputType `json:"mutableReferenceOutputs,omitempty"`
	ReturnValues            []ReturnValueType            `json:"returnValues,omitempty"`
}

type DevInspectResults struct {
	Effects TransactionEffects    `json:"effects"`
	Events  TransactionEvents     `json:"events"`
	Results []ExecutionResultType `json:"results,omitempty"`
	Error   string                `json:"error,omitempty"`
}

type TransactionFilter struct {
	Checkpoint   *SequenceNumber `json:"Checkpoint,omitempty"`
	MoveFunction *struct {
		Package  ObjectId `json:"package"`
		Module   string   `json:"module,omitempty"`
		Function string   `json:"function,omitempty"`
	} `json:"MoveFunction,omitempty"`
	InputObject      *ObjectId `json:"InputObject,omitempty"`
	ChangedObject    *ObjectId `json:"ChangedObject,omitempty"`
	FromAddress      *Address  `json:"FromAddress,omitempty"`
	ToAddress        *Address  `json:"ToAddress,omitempty"`
	FromAndToAddress *struct {
		From *Address `json:"from"`
		To   *Address `json:"to"`
	} `json:"FromAndToAddress,omitempty"`
	TransactionKind *string `json:"TransactionKind,omitempty"`
}

type SuiTransactionBlockResponseOptions struct {
	/* Whether to show transaction input data. Default to be false. */
	ShowInput bool `json:"showInput,omitempty"`
	/* Whether to show transaction effects. Default to be false. */
	ShowEffects bool `json:"showEffects,omitempty"`
	/* Whether to show transaction events. Default to be false. */
	ShowEvents bool `json:"showEvents,omitempty"`
	/* Whether to show object changes. Default to be false. */
	ShowObjectChanges bool `json:"showObjectChanges,omitempty"`
	/* Whether to show coin balance changes. Default to be false. */
	ShowBalanceChanges bool `json:"showBalanceChanges,omitempty"`
}

type SuiTransactionBlockResponseQuery struct {
	Filter  *TransactionFilter                  `json:"filter,omitempty"`
	Options *SuiTransactionBlockResponseOptions `json:"options,omitempty"`
}

type PaginatedTransactionResponse struct {
	Data        []SuiTransactionBlockResponse `json:"data,omitempty"`
	NextCursor  *TransactionDigest            `json:"nextCursor,omitempty"`
	HasNextPage bool                          `json:"hasNextPage"`
}

type DryRunTransactionBlockResponse struct {
	Effects        TransactionEffects `json:"effects"`
	Events         TransactionEvents  `json:"events"`
	ObjectChanges  []SuiObjectChange  `json:"objectChanges,omitempty"`
	BalanceChanges []BalanceChange    `json:"balanceChanges,omitempty"`
}
