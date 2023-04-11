package types

import (
	"errors"
	"math/big"
	"sort"

	"github.com/shopspring/decimal"
)

// type LockedBalance struct {
// 	EpochId int64 `json:"epochId"`
// 	Number  int64 `json:"number"`
// }

type Coin struct {
	CoinType     string                `json:"coinType"`
	CoinObjectId ObjectId              `json:"coinObjectId"`
	Version      SuiBigInt             `json:"version"`
	Digest       TransactionDigest     `json:"digest"`
	Balance      SafeSuiBigInt[uint64] `json:"balance"`

	LockedUntilEpoch    *SafeSuiBigInt[uint64] `json:"lockedUntilEpoch,omitempty"`
	PreviousTransaction TransactionDigest      `json:"previousTransaction"`
}

type CoinPage = Page[Coin, ObjectId]

type Balance struct {
	CoinType        string                              `json:"coinType"`
	CoinObjectCount uint64                              `json:"coinObjectCount"`
	TotalBalance    SuiBigInt                           `json:"totalBalance"`
	LockedBalance   map[SafeSuiBigInt[uint64]]SuiBigInt `json:"lockedBalance"`
}

type Supply struct {
	Value decimal.Decimal `json:"value"`
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
		total = total.Add(total, big.NewInt(coin.Balance.Int64()))
	}
	return total
}

func (cs Coins) PickCoinNoLess(amount uint64) (*Coin, error) {
	for i, coin := range cs {
		if coin.Balance.Uint64() >= amount {
			cs = append(cs[:i], cs[i+1:]...)
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
		if cs[i].Balance.Uint64() < gasAmount {
			continue
		}

		if nil == gasCoin || gasCoin.Balance.Uint64() > cs[i].Balance.Uint64() {
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
					return sortedCoins[i].Balance.Uint64() < sortedCoins[j].Balance.Uint64()
				} else {
					return sortedCoins[i].Balance.Uint64() >= sortedCoins[j].Balance.Uint64()
				}
			},
		)
	}

	result := make(Coins, 0)
	total := big.NewInt(0)
	for _, coin := range sortedCoins {
		result = append(result, coin)
		total = new(big.Int).Add(total, big.NewInt(coin.Balance.Int64()))
		if total.Cmp(amount) >= 0 {
			return result, nil
		}
	}

	return nil, ErrCoinsNeedMoreObject
}
