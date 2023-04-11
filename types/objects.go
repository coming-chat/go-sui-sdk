package types

import (
	"github.com/shopspring/decimal"
)

type ObjectDigest = string

type SuiObjectRef struct {
	/** Base64 string representing the object digest */
	Digest TransactionDigest `json:"digest"`
	/** Hex code as string representing the object id */
	ObjectId string `json:"objectId"`
	/** Object version */
	Version int64 `json:"version"`
}

type SuiGasData struct {
	Payment []SuiObjectRef `json:"payment"`
	/** Gas Object's owner */
	Owner  string `json:"owner"`
	Price  int64  `json:"price"`
	Budget int64  `json:"budget"`
}

type SuiParsedData struct {
	MoveObject *SuiParsedMoveObject `json:"moveObject,omitempty"`
	Package    *SuiMovePackage      `json:"package,omitempty"`
}

func (p SuiParsedData) Tag() string {
	return "dataType"
}

func (p SuiParsedData) Content() string {
	return ""
}

type SuiMovePackage struct {
	Disassembled map[string]interface{} `json:"disassembled"`
}

type SuiParsedMoveObject struct {
	Type              string `json:"type"`
	HasPublicTransfer bool   `json:"hasPublicTransfer"`
	Fields            any    `json:"fields"`
}

type SuiRawData struct {
	MoveObject *SuiRawMoveObject  `json:"moveObject,omitempty"`
	Package    *SuiRawMovePackage `json:"package,omitempty"`
}

func (r SuiRawData) Tag() string {
	return "dataType"
}

func (r SuiRawData) Content() string {
	return ""
}

type SuiRawMoveObject struct {
	Type              string         `json:"type"`
	HasPublicTransfer bool           `json:"hasPublicTransfer"`
	Version           SequenceNumber `json:"version"`
	BcsBytes          Base64Data     `json:"bcsBytes"`
}

type SuiRawMovePackage struct {
	Id              ObjectId              `json:"id"`
	Version         SequenceNumber        `json:"version"`
	ModuleMap       map[string]Base64Data `json:"moduleMap"`
	TypeOriginTable []TypeOrigin          `json:"typeOriginTable"`
	LinkageTable    map[string]UpgradeInfo
}

type UpgradeInfo struct {
	UpgradedId      ObjectId
	UpgradedVersion SequenceNumber
}

type TypeOrigin struct {
	ModuleName string   `json:"moduleName"`
	StructName string   `json:"structName"`
	Package    ObjectId `json:"package"`
}

type SuiObjectData struct {
	ObjectId ObjectId        `json:"objectId"`
	Version  decimal.Decimal `json:"version"`
	Digest   ObjectDigest    `json:"digest"`
	/**
	 * Type of the object, default to be undefined unless SuiObjectDataOptions.showType is set to true
	 */
	Type *string `json:"type,omitempty"`
	/**
	 * Move object content or package content, default to be undefined unless SuiObjectDataOptions.showContent is set to true
	 */
	Content *TagJson[SuiParsedData] `json:"content,omitempty"`
	/**
	 * Move object content or package content in BCS bytes, default to be undefined unless SuiObjectDataOptions.showBcs is set to true
	 */
	Bcs *TagJson[SuiRawData] `json:"bcs,omitempty"`
	/**
	 * The owner of this object. Default to be undefined unless SuiObjectDataOptions.showOwner is set to true
	 */
	Owner *ObjectOwner `json:"owner,omitempty"`
	/**
	 * The digest of the transaction that created or last mutated this object.
	 * Default to be undefined unless SuiObjectDataOptions.showPreviousTransaction is set to true
	 */
	PreviousTransaction *TransactionDigest `json:"previousTransaction,omitempty"`
	/**
	 * The amount of SUI we would rebate if this object gets deleted.
	 * This number is re-calculated each time the object is mutated based on
	 * the present storage gas price.
	 * Default to be undefined unless SuiObjectDataOptions.showStorageRebate is set to true
	 */
	StorageRebate *decimal.Decimal `json:"storageRebate,omitempty"`
	/**
	 * Display metadata for this object, default to be undefined unless SuiObjectDataOptions.showDisplay is set to true
	 * This can also be None if the struct type does not have Display defined
	 * See more details in https://forums.sui.io/t/nft-object-display-proposal/4872
	 */
	Display interface{} `json:"display,omitempty"`
}

type SuiObjectDataOptions struct {
	/* Whether to fetch the object type, default to be false */
	ShowType bool `json:"showType,omitempty"`
	/* Whether to fetch the object content, default to be false */
	ShowContent bool `json:"showContent,omitempty"`
	/* Whether to fetch the object content in BCS bytes, default to be false */
	ShowBcs bool `json:"showBcs,omitempty"`
	/* Whether to fetch the object owner, default to be false */
	ShowOwner bool `json:"showOwner,omitempty"`
	/* Whether to fetch the previous transaction digest, default to be false */
	ShowPreviousTransaction bool `json:"showPreviousTransaction,omitempty"`
	/* Whether to fetch the storage rebate, default to be false */
	ShowStorageRebate bool `json:"showStorageRebate,omitempty"`
	/* Whether to fetch the display metadata, default to be false */
	ShowDisplay bool `json:"showDisplay,omitempty"`
}

type SuiObjectResponseError struct {
	NotExists *struct {
		ObjectId ObjectId `json:"object_id"`
	} `json:"notExists,omitempty"`
	Deleted *struct {
		ObjectId ObjectId       `json:"object_id"`
		Version  SequenceNumber `json:"version"`
		Digest   ObjectDigest   `json:"digest"`
	} `json:"deleted,omitempty"`
	UnKnown      *struct{} `json:"unKnown"`
	DisplayError *struct {
		Error string `json:"error"`
	} `json:"displayError"`
}

func (e SuiObjectResponseError) Tag() string {
	return "code"
}

func (e SuiObjectResponseError) Content() string {
	return ""
}

type SuiObjectResponse struct {
	Data  *SuiObjectData                   `json:"data,omitempty"`
	Error *TagJson[SuiObjectResponseError] `json:"error,omitempty"`
}

type CheckpointedObjectId struct {
	ObjectId     ObjectId `json:"objectId"`
	AtCheckpoint *int     `json:"atCheckpoint"`
}

type PaginatedObjectsResponse struct {
	Data []SuiObjectResponse `json:"data,omitempty"`
	// NextCursor  *CheckpointedObjectId `json:"nextCursor,omitempty"`
	NextCursor  interface{} `json:"nextCursor,omitempty"` // will update when testnet after 0.30.0
	HasNextPage bool        `json:"hasNextPage"`
}

type ObjectsPage = Page[SuiObjectResponse, ObjectId]

type SuiObjectDataFilter struct {
	Package    *ObjectId   `json:"Package,omitempty"`
	MoveModule *MoveModule `json:"MoveModule,omitempty"`
	StructType string      `json:"StructType,omitempty"`
}

type SuiObjectResponseQuery struct {
	Filter  *SuiObjectDataFilter  `json:"filter,omitempty"`
	Options *SuiObjectDataOptions `json:"options,omitempty"`
}

type SuiPastObjectResponse = TagJson[SuiPastObject]

type SuiPastObject struct {
	/// The object exists and is found with this version
	VersionFound *SuiObjectData `json:"VersionFound,omitempty"`
	/// The object does not exist
	ObjectNotExists *ObjectId `json:"ObjectNotExists,omitempty"`
	/// The object is found to be deleted with this version
	ObjectDeleted *SuiObjectRef `json:"ObjectDeleted,omitempty"`
	/// The object exists but not found with this version
	VersionNotFound *struct{ ObjectId SequenceNumber } `json:"VersionNotFound,omitempty"`
	/// The asked object version is higher than the latest
	VersionTooHigh *struct {
		ObjectId      ObjectId       `json:"object_id"`
		AskedVersion  SequenceNumber `json:"asked_version"`
		LatestVersion SequenceNumber `json:"latest_version"`
	} `json:"VersionTooHigh,omitempty"`
}

func (s SuiPastObject) Tag() string {
	return "status"
}

func (s SuiPastObject) Content() string {
	return "details"
}
