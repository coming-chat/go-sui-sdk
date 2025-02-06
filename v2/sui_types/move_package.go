package sui_types

type MovePackage struct {
	Id              ObjectID
	Version         SequenceNumber
	ModuleMap       map[string][]uint8
	TypeOriginTable []TypeOrigin
	LinkageTable    map[ObjectID]UpgradeInfo
}

type TypeOrigin struct {
	ModuleName string
	StructName string
	Package    ObjectID
}

type UpgradeInfo struct {
	UpgradedId      ObjectID
	UpgradedVersion SequenceNumber
}
