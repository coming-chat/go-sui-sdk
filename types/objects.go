package types

import (
	"encoding/json"
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

type SuiParsedData interface{}
type SuiRawData interface{}

type SuiObjectData struct {
	ObjectId ObjectId       `json:"objectId"`
	Version  SequenceNumber `json:"version"`
	Digest   ObjectDigest   `json:"digest"`
	/**
	 * Type of the object, default to be undefined unless SuiObjectDataOptions.showType is set to true
	 */
	Type *string `json:"type,omitempty"`
	/**
	 * Move object content or package content, default to be undefined unless SuiObjectDataOptions.showContent is set to true
	 */
	Content SuiParsedData `json:"content,omitempty"`
	/**
	 * Move object content or package content in BCS bytes, default to be undefined unless SuiObjectDataOptions.showBcs is set to true
	 */
	Bcs SuiRawData `json:"bcs,omitempty"`
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
	StorageRebate *int64 `json:"storageRebate,omitempty"`
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
	UnKnown *struct{} `json:"unKnown"`
}

func (e *SuiObjectResponseError) UnmarshalJSON(data []byte) error {
	type orError struct {
		Code     string          `json:"code"`
		ObjectId *ObjectId       `json:"object_id,omitempty"`
		Version  *SequenceNumber `json:"version,omitempty"`
		Digest   *ObjectDigest   `json:"digest,omitempty"`
	}
	var tmp = orError{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	switch tmp.Code {
	case "notExists":
		e.NotExists = &struct {
			ObjectId ObjectId `json:"object_id"`
		}{
			ObjectId: *tmp.ObjectId,
		}
	case "deleted":
		e.Deleted = &struct {
			ObjectId ObjectId       `json:"object_id"`
			Version  SequenceNumber `json:"version"`
			Digest   ObjectDigest   `json:"digest"`
		}{
			ObjectId: *tmp.ObjectId,
			Version:  *tmp.Version,
			Digest:   *tmp.Digest,
		}
	case "unknown":
		e.UnKnown = &struct{}{}
	default:
	}
	return nil
}

type SuiObjectResponse struct {
	Data  *SuiObjectData          `json:"data,omitempty"`
	Error *SuiObjectResponseError `json:"error,omitempty"`
}

type CheckpointedObjectId struct {
	ObjectId     ObjectId `json:"objectId"`
	AtCheckpoint *int     `json:"atCheckpoint"`
}

type PaginatedObjectsResponse struct {
	Data        []SuiObjectResponse   `json:"data,omitempty"`
	NextCursor  *CheckpointedObjectId `json:"nextCursor,omitempty"`
	HasNextPage bool                  `json:"hasNextPage"`
}

type SuiObjectDataFilter struct {
	Package    *ObjectId   `json:"Package,omitempty"`
	MoveModule *MoveModule `json:"MoveModule,omitempty"`
	StructType string      `json:"StructType,omitempty"`
}

type SuiObjectResponseQuery struct {
	Filter  *SuiObjectDataFilter  `json:"filter,omitempty"`
	Options *SuiObjectDataOptions `json:"options,omitempty"`
}
