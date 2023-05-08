package sui_types

import (
	"encoding/hex"
	"fmt"
	"github.com/coming-chat/go-sui/lib"
	"github.com/coming-chat/go-sui/move_types"
	"strings"
)

type SuiAddress = lib.HexData

type SequenceNumber = uint64

// NewAddressFromHex
/**
 * Creates SuiAddress from a hex string.
 * @param addr Hex string can be with a prefix or without a prefix,
 * e.g. '0x1aa' or '1aa'. Hex string will be left padded with 0s if too short.
 */
func NewAddressFromHex(addr string) (*SuiAddress, error) {
	if strings.HasPrefix(addr, "0x") || strings.HasPrefix(addr, "0X") {
		addr = addr[2:]
	}
	if len(addr)%2 != 0 {
		addr = "0" + addr
	}

	data, err := hex.DecodeString(addr)
	if err != nil {
		return nil, err
	}
	const addressLength = 32
	if len(data) > addressLength {
		return nil, fmt.Errorf("hex string is too long. SuiAddress's length is %v data", addressLength)
	}

	res := [addressLength]byte{}
	copy(res[addressLength-len(data):], data[:])
	address := SuiAddress(res[:])
	return &address, nil
}

type ObjectID = lib.HexData

// ObjectRef for BCS, need to keep this order
type ObjectRef struct {
	ObjectId ObjectID       `json:"objectId"`
	Version  SequenceNumber `json:"version"`
	Digest   ObjectDigest   `json:"digest"`
}

type MoveObjectType_ struct {
	Other     *move_types.StructTag
	GasCoin   *lib.EmptyEnum
	StakedSui *lib.EmptyEnum
	Coin      *move_types.TypeTag
}
