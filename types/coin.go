package types

import (
	"errors"
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
	CoinType     string                `json:"coinType"`
	CoinObjectId ObjectId              `json:"coinObjectId"`
	Version      SuiBigInt             `json:"version"`
	Digest       TransactionDigest     `json:"digest"`
	Balance      SafeSuiBigInt[uint64] `json:"balance"`

	LockedUntilEpoch    *SafeSuiBigInt[uint64] `json:"lockedUntilEpoch,omitempty"`
	PreviousTransaction TransactionDigest      `json:"previousTransaction"`
}

func (c *Coin) Reference() *ObjectRef {
	return &ObjectRef{
		Digest:   c.Digest,
		Version:  c.Version,
		ObjectId: c.CoinObjectId,
	}
}

func (c *Coin) IsSUI() bool {
	return c.CoinType == SUI_COIN_TYPE
}

type CoinPage = Page[Coin, ObjectId]

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

	// There may be at least one coin even if gasAmount is 0
	GasCoins       []Coin
	GasTotalAmount uint64
	GasAmount      uint64
}

func (cs *PickedCoins) Count() int {
	return len(cs.Coins)
}

// Only have one coin, and the coin's amount is equal to the target amount
func (cs *PickedCoins) OnlyOneAndAmountMatch() bool {
	return len(cs.Coins) == 1 && cs.TotalAmount.Cmp(&cs.TargetAmount) == 0
}

func (cs *PickedCoins) CoinIds() []ObjectId {
	coinIds := make([]ObjectId, len(cs.Coins))
	for idx, coin := range cs.Coins {
		coinIds[idx] = coin.CoinObjectId
	}
	return coinIds
}

func (cs *PickedCoins) MaxGasCoin() *Coin {
	if cs.GasCoins == nil || len(cs.GasCoins) <= 0 {
		return nil
	}
	return &cs.GasCoins[0]
}

func (cs *PickedCoins) GasEnough() bool {
	return cs.GasTotalAmount >= cs.GasAmount
}

// return ((totalAmount + gasTotalAmount) - (targetAmount + gasAmount)) >= 0
func (cs *PickedCoins) AllAmountEnough() bool {
	return big.NewInt(0).Sub(&cs.TotalAmount, &cs.TargetAmount).Int64()-int64(cs.GasAmount)+int64(cs.GasTotalAmount) >= 0
}

// PickupCoins
// Select coins that match the target amount.
// @param inputCoins queried page coin datas
// @param targetAmount total amount of coins to be selected from inputCoins
// @param limit the max number of coins selected, default is `MAX_INPUT_COUNT_MERGE`
// @param gasAmount Only valid when coin is SUI, Amount of gas need to be selected, set 0 if no need separate gas coin
//
// @throw ErrNoCoinsFound If the count of input coins is 0.
// @throw ErrInsufficientBalance If the input coins are all that is left and the total amount is less than the target amount.
// @throw ErrNeedMergeCoin If there are many coins, but the total amount of coins limited is less than the target amount.
// @throw ErrNeedSplitGasCoin If the coin to be selected is SUI, the total amount of left all coins is greater than the target amount, but cannot reserved another gas coin.
func PickupCoins(inputCoins *CoinPage, targetAmount big.Int, limit int, gasAmount uint64) (*PickedCoins, error) {
	if inputCoins == nil || len(inputCoins.Data) <= 0 {
		return nil, ErrNoCoinsFound
	}
	if limit <= 0 {
		limit = MAX_INPUT_COUNT_MERGE
	}
	coins := inputCoins.Data
	// sort by balance descend
	sort.Slice(coins, func(i, j int) bool {
		return coins[i].Balance.Uint64() > coins[j].Balance.Uint64()
	})

	var res *PickedCoins

out:
	for {
		// First find a coin with a value that is exactly equal to the target amount.
		for idx, coin := range coins {
			if coin.Balance.Uint64() == targetAmount.Uint64() {
				res = &PickedCoins{
					Coins:        Coins{coin},
					TotalAmount:  targetAmount,
					TargetAmount: targetAmount,
				}
				coins = append(coins[:idx], coins[idx+1:]...)
				break out
			}
			if coin.Balance.Uint64() < targetAmount.Uint64() {
				break
			}
		}

		total := big.NewInt(0)
		pickedCoins := []Coin{}
		for idx, coin := range coins {
			total = total.Add(total, big.NewInt(0).SetUint64(coin.Balance.data))
			pickedCoins = append(pickedCoins, coin)
			if total.Cmp(&targetAmount) >= 0 {
				res = &PickedCoins{
					Coins:        pickedCoins,
					TotalAmount:  *total,
					TargetAmount: targetAmount,
				}
				coins = coins[idx+1:]
				break
			}
		}
		if res == nil { // This means that the totalAmount < targetAmount
			if inputCoins.HasNextPage {
				return nil, ErrNeedMergeCoin
			} else {
				return nil, ErrInsufficientBalance
			}
		}
		if limit < len(pickedCoins) {
			return nil, ErrNeedMergeCoin
		}
		break out
	}
	isSUI := res.Coins[0].IsSUI()
	if !isSUI {
		return res, nil
	}
	if gasAmount != 0 && len(coins) == 0 && !inputCoins.HasNextPage { // There is no gas coin.
		return nil, ErrNeedSplitGasCoin
	}

	totalGas := uint64(0)
	pickedGas := []Coin{}
	for _, coin := range coins {
		totalGas = totalGas + coin.Balance.data
		pickedGas = append(pickedGas, coin)
		if totalGas >= gasAmount {
			break
		}
	}
	// The gas fee is estimated and cannot be accurately equal, so no comparison is made
	res.GasCoins = pickedGas
	res.GasTotalAmount = totalGas
	res.GasAmount = gasAmount
	return res, nil
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
