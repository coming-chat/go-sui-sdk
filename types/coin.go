package types

import (
	"errors"
	"github.com/shopspring/decimal"
	"math/big"
	"sort"
)

// type LockedBalance struct {
// 	EpochId int64 `json:"epochId"`
// 	Number  int64 `json:"number"`
// }

type CoinStruct struct {
	CoinType     string            `json:"coinType"`
	CoinObjectId ObjectId          `json:"coinObjectId"`
	Version      uint64            `json:"version"`
	Digest       TransactionDigest `json:"digest"`
	Balance      int64             `json:"balance"`

	LockedUntilEpoch    *int              `json:"lockedUntilEpoch"`
	PreviousTransaction TransactionDigest `json:"previousTransaction"`
}

type PaginatedCoins struct {
	Data        []CoinStruct `json:"data"`
	NextCursor  *ObjectId    `json:"nextCursor,omitempty"`
	HasNextPage bool         `json:"hasNextPage"`
}

type Balance struct {
	CoinType        string                              `json:"coinType"`
	CoinObjectCount uint64                              `json:"coinObjectCount"`
	TotalBalance    decimal.Decimal                     `json:"totalBalance"`
	LockedBalance   map[decimal.Decimal]decimal.Decimal `json:"lockedBalance"`
}

type CoinSupply struct {
	Value uint64 `json:"value"`
}

var ErrCoinsNotMatchRequest error
var ErrCoinsNeedMoreObject error

const (
	PickSmaller = iota // pick smaller coins to match amount
	PickBigger         // pick bigger coins to match amount
	PickByOrder        // pick coins by coins order to match amount
)

// type Coin struct {
// 	Balance   uint64     `json:"balance"`
// 	Type      string     `json:"type"`
// 	Owner     *Address   `json:"owner"`
// 	Reference *ObjectRef `json:"reference"`
// }

type Coin = CoinStruct

func (c *Coin) Reference() *ObjectRef {
	return &ObjectRef{
		Digest:   c.Digest,
		Version:  c.Version,
		ObjectId: c.CoinObjectId,
	}
}

type Coins []Coin

func init() {
	ErrCoinsNotMatchRequest = errors.New("coins not match request")
	ErrCoinsNeedMoreObject = errors.New("you should get more SUI coins and try again")
}

func (cs Coins) TotalBalance() *big.Int {
	total := big.NewInt(0)
	for _, coin := range cs {
		total = total.Add(total, big.NewInt(coin.Balance))
	}
	return total
}

func (cs Coins) PickCoinNoLess(amount uint64) (*Coin, error) {
	for _, coin := range cs {
		if uint64(coin.Balance) >= amount {
			return &coin, nil
		}
	}
	if len(cs) <= 3 {
		return nil, errors.New("insufficient balance")
	}
	return nil, errors.New("no coin is enough to cover the gas")
}

// PickSUICoinsWithGas pick coins, which sum >= amount, and pick a gas coin >= gasAmount which not in coins
// if not satisfated amount/gasAmount, an ErrCoinsNotMatchRequest/ErrCoinsNeedMoreObject error will return
// if gasAmount == 0, a nil gasCoin will return
// pickMethod, see PickSmaller|PickBigger|PickByOrder
func (cs Coins) PickSUICoinsWithGas(amount *big.Int, gasAmount uint64, pickMethod int) (Coins, *Coin, error) {
	if gasAmount == 0 {
		res, err := cs.PickCoins(amount, pickMethod)
		return res, nil, err
	}

	if amount.Cmp(big.NewInt(0)) == 0 && gasAmount == 0 {
		return make(Coins, 0), nil, nil
	} else if len(cs) == 0 {
		return cs, nil, ErrCoinsNeedMoreObject
	}

	// find smallest to match gasAmount
	var gasCoin *Coin
	var selectIndex int
	for i := range cs {
		if uint64(cs[i].Balance) < gasAmount {
			continue
		}

		if nil == gasCoin || gasCoin.Balance > cs[i].Balance {
			gasCoin = &cs[i]
			selectIndex = i
		}
	}
	if nil == gasCoin {
		return nil, nil, ErrCoinsNotMatchRequest
	}

	lastCoins := make(Coins, 0, len(cs)-1)
	lastCoins = append(lastCoins, cs[0:selectIndex]...)
	lastCoins = append(lastCoins, cs[selectIndex+1:]...)
	pickCoins, err := lastCoins.PickCoins(amount, pickMethod)
	return pickCoins, gasCoin, err
}

// PickCoins pick coins, which sum >= amount,
// pickMethod, see PickSmaller|PickBigger|PickByOrder
// if not satisfated amount, an ErrCoinsNeedMoreObject error will return
func (cs Coins) PickCoins(amount *big.Int, pickMethod int) (Coins, error) {
	var sortedCoins Coins
	if pickMethod == PickByOrder {
		sortedCoins = cs
	} else {
		sortedCoins = make(Coins, len(cs))
		copy(sortedCoins, cs)
		sort.Slice(
			sortedCoins, func(i, j int) bool {
				if pickMethod == PickSmaller {
					return sortedCoins[i].Balance < sortedCoins[j].Balance
				} else {
					return sortedCoins[i].Balance >= sortedCoins[j].Balance
				}
			},
		)
	}

	result := make(Coins, 0)
	total := big.NewInt(0)
	for _, coin := range sortedCoins {
		result = append(result, coin)
		total = new(big.Int).Add(total, big.NewInt(coin.Balance))
		if total.Cmp(amount) >= 0 {
			return result, nil
		}
	}

	return nil, ErrCoinsNeedMoreObject
}
