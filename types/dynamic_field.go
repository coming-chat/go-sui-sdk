package types

import "github.com/coming-chat/go-sui/sui_types"

type DynamicFieldInfo struct {
	Name sui_types.DynamicFieldName `json:"name"`
	//Base58
	BcsName    string                     `json:"bcsName"`
	Type       sui_types.DynamicFieldType `json:"type"`
	ObjectType string                     `json:"objectType"`
	ObjectId   ObjectId                   `json:"objectId"`
	Version    SequenceNumber             `json:"version"`
	Digest     ObjectDigest               `json:"digest"`
}

type DynamicFieldPage Page[DynamicFieldInfo, ObjectId]
