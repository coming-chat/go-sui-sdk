package types

import (
	"github.com/coming-chat/go-sui/v2/lib"
	"github.com/coming-chat/go-sui/v2/sui_types"
)

type ExecuteTransactionRequestType string

const (
	TxnRequestTypeWaitForEffectsCert    ExecuteTransactionRequestType = "WaitForEffectsCert"
	TxnRequestTypeWaitForLocalExecution ExecuteTransactionRequestType = "WaitForLocalExecution"
)

type EpochId = uint64

type GasCostSummary struct {
	ComputationCost         SafeSuiBigInt[uint64] `json:"computationCost"`
	StorageCost             SafeSuiBigInt[uint64] `json:"storageCost"`
	StorageRebate           SafeSuiBigInt[uint64] `json:"storageRebate"`
	NonRefundableStorageFee SafeSuiBigInt[uint64] `json:"nonRefundableStorageFee"`
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
	Owner     lib.TagJson[sui_types.Owner] `json:"owner"`
	Reference SuiObjectRef                 `json:"reference"`
}

type SuiTransactionBlockEffectsModifiedAtVersions struct {
	ObjectId       sui_types.ObjectID                      `json:"objectId"`
	SequenceNumber SafeSuiBigInt[sui_types.SequenceNumber] `json:"sequenceNumber"`
}

type SuiTransactionBlockEffectsV1 struct {
	/** The status of the execution */
	Status ExecutionStatus `json:"status"`
	/** The epoch when this transaction was executed */
	ExecutedEpoch SafeSuiBigInt[EpochId] `json:"executedEpoch"`
	/** The version that every modified (mutated or deleted) object had before it was modified by this transaction. **/
	ModifiedAtVersions []SuiTransactionBlockEffectsModifiedAtVersions `json:"modifiedAtVersions,omitempty"`
	GasUsed            GasCostSummary                                 `json:"gasUsed"`
	/** The object references of the shared objects used in this transaction. Empty if no shared objects were used. */
	SharedObjects []SuiObjectRef `json:"sharedObjects,omitempty"`
	/** The transaction digest */
	TransactionDigest sui_types.TransactionDigest `json:"transactionDigest"`
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
	EventsDigest *sui_types.TransactionEventsDigest `json:"eventsDigest,omitempty"`
	/** The set of transaction digests this transaction depends on */
	Dependencies []sui_types.TransactionDigest `json:"dependencies,omitempty"`
}

type SuiTransactionBlockEffects struct {
	V1 *SuiTransactionBlockEffectsV1 `json:"v1"`
}

func (t SuiTransactionBlockEffects) Tag() string {
	return "messageVersion"
}

func (t SuiTransactionBlockEffects) Content() string {
	return ""
}

func (t SuiTransactionBlockEffects) GasFee() int64 {
	if t.V1 == nil {
		return 0
	}
	fee := t.V1.GasUsed.StorageCost.Int64() -
		t.V1.GasUsed.StorageRebate.Int64() +
		t.V1.GasUsed.ComputationCost.Int64()
	return fee
}

func (t SuiTransactionBlockEffects) IsSuccess() bool {
	return t.V1.Status.Status == ExecutionStatusSuccess
}

const (
	SuiTransactionBlockKindSuiChangeEpoch             = "ChangeEpoch"
	SuiTransactionBlockKindSuiConsensusCommitPrologue = "ConsensusCommitPrologue"
	SuiTransactionBlockKindGenesis                    = "Genesis"
	SuiTransactionBlockKindProgrammableTransaction    = "ProgrammableTransaction"
)

type SuiTransactionBlockKind = lib.TagJson[TransactionBlockKind]

type TransactionBlockKind struct {
	/// A system transaction that will update epoch information on-chain.
	ChangeEpoch *SuiChangeEpoch `json:"ChangeEpoch,omitempty"`
	/// A system transaction used for initializing the initial state of the chain.
	Genesis *SuiGenesisTransaction `json:"Genesis,omitempty"`
	/// A system transaction marking the start of a series of transactions scheduled as part of a
	/// checkpoint
	ConsensusCommitPrologue *SuiConsensusCommitPrologue `json:"ConsensusCommitPrologue,omitempty"`
	/// A series of transactions where the results of one transaction can be used in future
	/// transactions
	ProgrammableTransaction *SuiProgrammableTransactionBlock `json:"ProgrammableTransaction,omitempty"`
	// .. more transaction types go here
}

func (t TransactionBlockKind) Tag() string {
	return "kind"
}

func (t TransactionBlockKind) Content() string {
	return ""
}

type SuiChangeEpoch struct {
	Epoch                 SafeSuiBigInt[EpochId] `json:"epoch"`
	StorageCharge         uint64                 `json:"storage_charge"`
	ComputationCharge     uint64                 `json:"computation_charge"`
	StorageRebate         uint64                 `json:"storage_rebate"`
	EpochStartTimestampMs uint64                 `json:"epoch_start_timestamp_ms"`
}

type SuiGenesisTransaction struct {
	Objects []sui_types.ObjectID `json:"objects"`
}

type SuiConsensusCommitPrologue struct {
	Epoch             uint64 `json:"epoch"`
	Round             uint64 `json:"round"`
	CommitTimestampMs uint64 `json:"commit_timestamp_ms"`
}

type SuiProgrammableTransactionBlock struct {
	Inputs []interface{} `json:"inputs"`
	/// The transactions to be executed sequentially. A failure in any transaction will
	/// result in the failure of the entire programmable transaction block.
	Commands []interface{} `json:"transactions"`
}

type SuiTransactionBlockDataV1 struct {
	Transaction SuiTransactionBlockKind `json:"transaction"`
	Sender      sui_types.SuiAddress    `json:"sender"`
	GasData     SuiGasData              `json:"gasData"`
}

type SuiTransactionBlockData struct {
	V1 *SuiTransactionBlockDataV1 `json:"v1,omitempty"`
}

func (t SuiTransactionBlockData) Tag() string {
	return "messageVersion"
}

func (t SuiTransactionBlockData) Content() string {
	return ""
}

type SuiTransactionBlock struct {
	Data         lib.TagJson[SuiTransactionBlockData] `json:"data"`
	TxSignatures []string                             `json:"txSignatures"`
}

type ObjectChange struct {
	Published *struct {
		PackageId sui_types.ObjectID                      `json:"packageId"`
		Version   SafeSuiBigInt[sui_types.SequenceNumber] `json:"version"`
		Digest    sui_types.ObjectDigest                  `json:"digest"`
		Nodules   []string                                `json:"nodules"`
	} `json:"published,omitempty"`
	/// Transfer objects to new address / wrap in another object
	Transferred *struct {
		Sender     sui_types.SuiAddress                    `json:"sender"`
		Recipient  ObjectOwner                             `json:"recipient"`
		ObjectType string                                  `json:"objectType"`
		ObjectId   sui_types.ObjectID                      `json:"objectId"`
		Version    SafeSuiBigInt[sui_types.SequenceNumber] `json:"version"`
		Digest     sui_types.ObjectDigest                  `json:"digest"`
	} `json:"transferred,omitempty"`
	/// Object mutated.
	Mutated *struct {
		Sender          sui_types.SuiAddress                    `json:"sender"`
		Owner           ObjectOwner                             `json:"owner"`
		ObjectType      string                                  `json:"objectType"`
		ObjectId        sui_types.ObjectID                      `json:"objectId"`
		Version         SafeSuiBigInt[sui_types.SequenceNumber] `json:"version"`
		PreviousVersion SafeSuiBigInt[sui_types.SequenceNumber] `json:"previousVersion"`
		Digest          sui_types.ObjectDigest                  `json:"digest"`
	} `json:"mutated,omitempty"`
	/// Delete object j
	Deleted *struct {
		Sender     sui_types.SuiAddress                    `json:"sender"`
		ObjectType string                                  `json:"objectType"`
		ObjectId   sui_types.ObjectID                      `json:"objectId"`
		Version    SafeSuiBigInt[sui_types.SequenceNumber] `json:"version"`
	} `json:"deleted,omitempty"`
	/// Wrapped object
	Wrapped *struct {
		Sender     sui_types.SuiAddress                    `json:"sender"`
		ObjectType string                                  `json:"objectType"`
		ObjectId   sui_types.ObjectID                      `json:"objectId"`
		Version    SafeSuiBigInt[sui_types.SequenceNumber] `json:"version"`
	} `json:"wrapped,omitempty"`
	/// New object creation
	Created *struct {
		Sender     sui_types.SuiAddress                    `json:"sender"`
		Owner      ObjectOwner                             `json:"owner"`
		ObjectType string                                  `json:"objectType"`
		ObjectId   sui_types.ObjectID                      `json:"objectId"`
		Version    SafeSuiBigInt[sui_types.SequenceNumber] `json:"version"`
		Digest     sui_types.ObjectDigest                  `json:"digest"`
	} `json:"created,omitempty"`
}

func (o ObjectChange) Tag() string {
	return "type"
}

func (o ObjectChange) Content() string {
	return ""
}

type BalanceChange struct {
	Owner    ObjectOwner `json:"owner"`
	CoinType string      `json:"coinType"`
	/* Coin balance change(positive means receive, negative means send) */
	Amount string `json:"amount"`
}

type SuiTransactionBlockResponse struct {
	Digest                  sui_types.TransactionDigest              `json:"digest"`
	Transaction             *SuiTransactionBlock                     `json:"transaction,omitempty"`
	RawTransaction          []byte                                   `json:"rawTransaction,omitempty"`
	Effects                 *lib.TagJson[SuiTransactionBlockEffects] `json:"effects,omitempty"`
	Events                  []SuiEvent                               `json:"events,omitempty"`
	TimestampMs             *SafeSuiBigInt[uint64]                   `json:"timestampMs,omitempty"`
	Checkpoint              *SafeSuiBigInt[CheckpointSequenceNumber] `json:"checkpoint,omitempty"`
	ConfirmedLocalExecution *bool                                    `json:"confirmedLocalExecution,omitempty"`
	ObjectChanges           []lib.TagJson[ObjectChange]              `json:"objectChanges,omitempty"`
	BalanceChanges          []BalanceChange                          `json:"balanceChanges,omitempty"`
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
	Effects lib.TagJson[SuiTransactionBlockEffects] `json:"effects"`
	Events  []SuiEvent                              `json:"events"`
	Results []ExecutionResultType                   `json:"results,omitempty"`
	Error   *string                                 `json:"error,omitempty"`
}

type TransactionFilter struct {
	Checkpoint   *sui_types.SequenceNumber `json:"Checkpoint,omitempty"`
	MoveFunction *struct {
		Package  sui_types.ObjectID `json:"package"`
		Module   string             `json:"module,omitempty"`
		Function string             `json:"function,omitempty"`
	} `json:"MoveFunction,omitempty"`
	InputObject      *sui_types.ObjectID   `json:"InputObject,omitempty"`
	ChangedObject    *sui_types.ObjectID   `json:"ChangedObject,omitempty"`
	FromAddress      *sui_types.SuiAddress `json:"FromAddress,omitempty"`
	ToAddress        *sui_types.SuiAddress `json:"ToAddress,omitempty"`
	FromAndToAddress *struct {
		From *sui_types.SuiAddress `json:"from"`
		To   *sui_types.SuiAddress `json:"to"`
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

type TransactionBlocksPage = Page[SuiTransactionBlockResponse, sui_types.TransactionDigest]

type DryRunTransactionBlockResponse struct {
	Effects        lib.TagJson[SuiTransactionBlockEffects] `json:"effects"`
	Events         []SuiEvent                              `json:"events"`
	ObjectChanges  []lib.TagJson[ObjectChange]             `json:"objectChanges"`
	BalanceChanges []BalanceChange                         `json:"balanceChanges"`
	Input          lib.TagJson[SuiTransactionBlockData]    `json:"input"`
}
