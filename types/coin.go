package types

import (
	"errors"
	"github.com/coming-chat/go-sui/sui_types"
	"math/big"
	"sort"
)

const SUI_COIN_TYPE = "0x2::sui::SUI"

const MAX_INPUT_COUNT_MERGE = 256 - 1
const MAX_INPUT_COUNT_STAKE = 512 - 1

// type LockedBalance struct {
// 	EpochId int64 `json:"epochId"`
// 	Number  int64 `json:"number"`
// }

type Coin struct {
	CoinType     string                   `json:"coinType"`
	CoinObjectId sui_types.ObjectID       `json:"coinObjectId"`
	Version      sui_types.SequenceNumber `json:"version"`
	Digest       sui_types.ObjectDigest   `json:"digest"`
	Balance      SafeSuiBigInt[uint64]    `json:"balance"`

	LockedUntilEpoch    *SafeSuiBigInt[uint64]      `json:"lockedUntilEpoch,omitempty"`
	PreviousTransaction sui_types.TransactionDigest `json:"previousTransaction"`
}

func (c *Coin) Reference() *sui_types.ObjectRef {
	return &sui_types.ObjectRef{
		Digest:   c.Digest,
		Version:  c.Version,
		ObjectId: c.CoinObjectId,
	}
}

type CoinPage = Page[Coin, sui_types.ObjectID]

type Balance struct {
	CoinType        string                              `json:"coinType"`
	CoinObjectCount uint64                              `json:"coinObjectCount"`
	TotalBalance    SuiBigInt                           `json:"totalBalance"`
	LockedBalance   map[SafeSuiBigInt[uint64]]SuiBigInt `json:"lockedBalance"`
}

type Supply struct {
	Value SafeSuiBigInt[uint64] `json:"value"`
}

type PickedCoins struct {
	Coins        []Coin
	TotalAmount  big.Int
	TargetAmount big.Int

	// max coin value except the picked coins, It may be used to help set the gas budget.
	RemainingMaxCoinValue uint64
}

func (cs *PickedCoins) Count() int {
	return len(cs.Coins)
}

// Only have one coin, and the coin's amount is equal to the target amount
func (cs *PickedCoins) OnlyOneAndAmountMatch() bool {
	return len(cs.Coins) == 1 && cs.TotalAmount.Cmp(&cs.TargetAmount) == 0
}

func (cs *PickedCoins) CoinIds() []sui_types.ObjectID {
	coinIds := make([]sui_types.ObjectID, len(cs.Coins))
	for idx, coin := range cs.Coins {
		coinIds[idx] = coin.CoinObjectId
	}
	return coinIds
}

// PickupCoins
// Select coins that match the target amount.
// @param inputCoins queried page coin datas
// @param targetAmount total amount of coins to be selected from inputCoins
// @param limit the max number of coins selected, default is `MAX_INPUT_COUNT_MERGE`
// @param reserveGasCoin Only valid when coin is SUI, Do we need to keep at least one coin unselected?
// @throw ErrNoCoinsFound If the count of input coins is 0.
// @throw ErrInsufficientBalance If the input coins are all that is left and the total amount is less than the target amount.
// @throw ErrNeedMergeCoin If there are many coins, but the total amount of coins limited is less than the target amount.
// @throw ErrNeedSplitGasCoin If the coin to be selected is SUI, the total amount of left all coins is greater than the target amount, but cannot reserved another gas coin.
func PickupCoins(inputCoins *CoinPage, targetAmount big.Int, limit int, reserveGasCoin bool) (*PickedCoins, error) {
	inputCount := len(inputCoins.Data)
	if inputCount <= 0 {
		return nil, ErrNoCoinsFound
	}
	if limit == 0 {
		limit = MAX_INPUT_COUNT_MERGE
	}
	coins := inputCoins.Data
	// sort by balance descend
	sort.Slice(
		coins, func(i, j int) bool {
			return coins[i].Balance.Uint64() > coins[j].Balance.Uint64()
		},
	)

	// First find a coin with a value that is exactly equal to the target amount.
	for idx, coin := range coins {
		if coin.Balance.Uint64() == targetAmount.Uint64() {
			maxGas := uint64(0)
			if idx == 0 && len(coins) > 1 {
				maxGas = coins[1].Balance.Uint64()
			} else {
				maxGas = coins[0].Balance.Uint64()
			}
			return &PickedCoins{
				Coins:        Coins{coin},
				TotalAmount:  targetAmount,
				TargetAmount: targetAmount,

				RemainingMaxCoinValue: maxGas,
			}, nil
		}
		if coin.Balance.Uint64() < targetAmount.Uint64() {
			break
		}
	}

	total := big.NewInt(0)
	pickedCoins := []Coin{}
	maxGas := uint64(0)
	for idx, coin := range coins {
		total = total.Add(total, big.NewInt(0).SetUint64(coin.Balance.Uint64()))
		pickedCoins = append(pickedCoins, coin)
		if total.Cmp(&targetAmount) >= 0 {
			if idx+1 < len(coins) {
				maxGas = coins[idx+1].Balance.Uint64()
			}
			break
		}
	}

	if total.Cmp(&targetAmount) < 0 {
		if inputCoins.HasNextPage {
			return nil, ErrNeedMergeCoin
		} else {
			return nil, ErrInsufficientBalance
		}
	}
	if limit < len(pickedCoins) {
		return nil, ErrNeedMergeCoin
	}
	isLeftCoin := inputCoins.HasNextPage || inputCount > len(pickedCoins)
	isSUI := coins[0].CoinType == SUI_COIN_TYPE
	if isSUI && reserveGasCoin && !isLeftCoin {
		return nil, ErrNeedSplitGasCoin
	}

	return &PickedCoins{
		Coins:        pickedCoins,
		TotalAmount:  *total,
		TargetAmount: targetAmount,

		RemainingMaxCoinValue: maxGas,
	}, nil
}

type Coins []Coin

func (cs Coins) TotalBalance() *big.Int {
	total := big.NewInt(0)
	for _, coin := range cs {
		total = total.Add(total, big.NewInt(0).SetUint64(coin.Balance.Uint64()))
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

const (
	PickSmaller = iota // pick smaller coins to match amount
	PickBigger         // pick bigger coins to match amount
	PickByOrder        // pick coins by coins order to match amount
)

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
		total = new(big.Int).Add(total, big.NewInt(0).SetUint64(coin.Balance.Uint64()))
		if total.Cmp(amount) >= 0 {
			return result, nil
		}
	}

	return nil, ErrCoinsNeedMoreObject
}
