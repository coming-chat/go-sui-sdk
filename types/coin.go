package types

import "math/big"

type Coin struct {
	Balance             uint64       `json:"balance"`
	Type                string       `json:"type"`
	Owner               *ObjectOwner `json:"owner"`
	PreviousTransaction *Digest      `json:"previousTransaction"`
	Reference           *ObjectRef   `json:"reference"`
}

type Coins []Coin

func (cs Coins) TotalBalance() *big.Int {
	total := big.NewInt(0)
	for _, coin := range cs {
		total = total.Add(total, big.NewInt(0).SetUint64(coin.Balance))
	}
	return total
}
