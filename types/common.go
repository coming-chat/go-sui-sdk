package types

type TransactionDigest = string

type TransactionEffectsDigest = string

type TransactionEventDigest = string

type SequenceNumber = uint64

// export const ObjectId = string();
// export type ObjectId = Infer<typeof ObjectId>;

// export const SuiAddress = string();
// export type SuiAddress = Infer<typeof SuiAddress>;

type ObjectOwnerInternal struct {
	AddressOwner *Address `json:"AddressOwner,omitempty"`
	ObjectOwner  *Address `json:"ObjectOwner,omitempty"`
	SingleOwner  *Address `json:"SingleOwner,omitempty"`
	Shared       *struct {
		InitialSharedVersion uint64 `json:"initial_shared_version"`
	} `json:"Shared,omitempty"`
}

type ObjectOwner struct {
	*ObjectOwnerInternal
	*string
}
