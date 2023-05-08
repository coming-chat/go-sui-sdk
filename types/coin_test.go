package types

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func balanceObject(val uint64) SafeSuiBigInt[uint64] {
	return NewSafeSuiBigInt(val)
}

func TestCoins_PickSUICoinsWithGas(t *testing.T) {
	// coins 1,2,3,4,5
	testCoins := Coins{
		{Balance: balanceObject(3)},
		{Balance: balanceObject(5)},
		{Balance: balanceObject(1)},
		{Balance: balanceObject(4)},
		{Balance: balanceObject(2)},
	}
	type args struct {
		amount     *big.Int
		gasAmount  uint64
		pickMethod int
	}
	tests := []struct {
		name    string
		cs      Coins
		args    args
		want    Coins
		want1   *Coin
		wantErr bool
	}{
		{
			name: "case success 1",
			cs:   testCoins,
			args: args{
				amount:     new(big.Int),
				gasAmount:  0,
				pickMethod: PickSmaller,
			},
			want:    nil,
			want1:   nil,
			wantErr: false,
		},
		{
			name: "case success 2",
			cs:   testCoins,
			args: args{
				amount:     big.NewInt(1),
				gasAmount:  2,
				pickMethod: PickSmaller,
			},
			want:    Coins{{Balance: balanceObject(1)}},
			want1:   &Coin{Balance: balanceObject(2)},
			wantErr: false,
		},
		{
			name: "case success 3",
			cs:   testCoins,
			args: args{
				amount:     big.NewInt(4),
				gasAmount:  2,
				pickMethod: PickSmaller,
			},
			want:    Coins{{Balance: balanceObject(1)}, {Balance: balanceObject(3)}},
			want1:   &Coin{Balance: balanceObject(2)},
			wantErr: false,
		},
		{
			name: "case success 4",
			cs:   testCoins,
			args: args{
				amount:     big.NewInt(6),
				gasAmount:  2,
				pickMethod: PickSmaller,
			},
			want:    Coins{{Balance: balanceObject(1)}, {Balance: balanceObject(3)}, {Balance: balanceObject(4)}},
			want1:   &Coin{Balance: balanceObject(2)},
			wantErr: false,
		},
		{
			name: "case error 1",
			cs:   testCoins,
			args: args{
				amount:     big.NewInt(6),
				gasAmount:  6,
				pickMethod: PickSmaller,
			},
			want:    Coins{},
			want1:   nil,
			wantErr: true,
		},
		{
			name: "case error 1",
			cs:   testCoins,
			args: args{
				amount:     big.NewInt(100),
				gasAmount:  3,
				pickMethod: PickSmaller,
			},
			want:    Coins{},
			want1:   &Coin{Balance: balanceObject(3)},
			wantErr: true,
		},
		{
			name: "case bigger 1",
			cs:   testCoins,
			args: args{
				amount:     big.NewInt(3),
				gasAmount:  3,
				pickMethod: PickBigger,
			},
			want:    Coins{{Balance: balanceObject(5)}},
			want1:   &Coin{Balance: balanceObject(3)},
			wantErr: false,
		},
		{
			name: "case order 1",
			cs:   testCoins,
			args: args{
				amount:     big.NewInt(3),
				gasAmount:  3,
				pickMethod: PickByOrder,
			},
			want:    Coins{{Balance: balanceObject(5)}},
			want1:   &Coin{Balance: balanceObject(3)},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.cs.PickSUICoinsWithGas(tt.args.amount, tt.args.gasAmount, tt.args.pickMethod)
			if (err != nil) != tt.wantErr {
				t.Errorf("Coins.PickSUICoinsWithGas() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != 0 && len(tt.want) != 0 {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Coins.PickSUICoinsWithGas() got = %v, want %v", got, tt.want)
				}
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Coins.PickSUICoinsWithGas() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCoins_PickCoins(t *testing.T) {
	// coins 1,2,3,4,5
	testCoins := Coins{
		{Balance: balanceObject(3)},
		{Balance: balanceObject(5)},
		{Balance: balanceObject(1)},
		{Balance: balanceObject(4)},
		{Balance: balanceObject(2)},
	}
	type args struct {
		amount     *big.Int
		pickMethod int
	}
	tests := []struct {
		name    string
		cs      Coins
		args    args
		want    Coins
		wantErr bool
	}{
		{
			name:    "smaller 1",
			cs:      testCoins,
			args:    args{amount: big.NewInt(2), pickMethod: PickSmaller},
			want:    Coins{{Balance: balanceObject(1)}, {Balance: balanceObject(2)}},
			wantErr: false,
		},
		{
			name:    "smaller 2",
			cs:      testCoins,
			args:    args{amount: big.NewInt(4), pickMethod: PickSmaller},
			want:    Coins{{Balance: balanceObject(1)}, {Balance: balanceObject(2)}, {Balance: balanceObject(3)}},
			wantErr: false,
		},
		{
			name:    "bigger 1",
			cs:      testCoins,
			args:    args{amount: big.NewInt(2), pickMethod: PickBigger},
			want:    Coins{{Balance: balanceObject(5)}},
			wantErr: false,
		},
		{
			name:    "bigger 2",
			cs:      testCoins,
			args:    args{amount: big.NewInt(6), pickMethod: PickBigger},
			want:    Coins{{Balance: balanceObject(5)}, {Balance: balanceObject(4)}},
			wantErr: false,
		},
		{
			name:    "pick by order 1",
			cs:      testCoins,
			args:    args{amount: big.NewInt(6), pickMethod: PickByOrder},
			want:    Coins{{Balance: balanceObject(3)}, {Balance: balanceObject(5)}},
			wantErr: false,
		},
		{
			name:    "pick by order 2",
			cs:      testCoins,
			args:    args{amount: big.NewInt(15), pickMethod: PickByOrder},
			want:    testCoins,
			wantErr: false,
		},
		{
			name:    "pick error",
			cs:      testCoins,
			args:    args{amount: big.NewInt(16), pickMethod: PickByOrder},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cs.PickCoins(tt.args.amount, tt.args.pickMethod)
			if (err != nil) != tt.wantErr {
				t.Errorf("Coins.PickCoins() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Coins.PickCoins() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPickupCoins(t *testing.T) {
	coin := func(n uint64) Coin {
		return Coin{Balance: balanceObject(n), CoinType: SUI_COIN_TYPE}
	}

	type args struct {
		inputCoins   *CoinPage
		targetAmount big.Int
		limit        int
		gasAmount    uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *PickedCoins
		wantErr error
	}{
		{
			name: "need sort",
			args: args{
				inputCoins: &Page[Coin, HexData]{
					Data: []Coin{
						coin(1e3), coin(1e5), coin(1e2), coin(1e4),
					},
				},
				targetAmount: *big.NewInt(1e5 + 1e3),
				limit:        100,
			},
			want: &PickedCoins{
				Coins: []Coin{
					coin(1e5), coin(1e4),
				},
				TotalAmount:  *big.NewInt(1e5 + 1e4),
				TargetAmount: *big.NewInt(1e5 + 1e3),
				GasCoins: []Coin{
					coin(1e3),
				},
				GasTotalAmount: 1e3,
			},
		},
		{
			name: "ErrNoCoinsFound",
			args: args{
				inputCoins: &Page[Coin, HexData]{
					Data: []Coin{},
				},
				targetAmount: *big.NewInt(101000),
			},
			wantErr: ErrNoCoinsFound,
		},
		{
			name: "ErrInsufficientBalance",
			args: args{
				inputCoins: &Page[Coin, HexData]{
					Data: []Coin{
						coin(1e5), coin(1e6), coin(1e4),
					},
				},
				targetAmount: *big.NewInt(1e9),
			},
			wantErr: ErrInsufficientBalance,
		},
		{
			name: "ErrNeedMergeCoin 1",
			args: args{
				inputCoins: &Page[Coin, HexData]{
					Data: []Coin{
						coin(1e5), coin(1e6), coin(1e4),
					},
					HasNextPage: true,
				},
				targetAmount: *big.NewInt(1e9),
			},
			wantErr: ErrNeedMergeCoin,
		},
		{
			name: "ErrNeedMergeCoin 2",
			args: args{
				inputCoins: &Page[Coin, HexData]{
					Data: []Coin{
						coin(1e5), coin(1e6), coin(1e4), coin(1e5),
					},
					HasNextPage: false,
				},
				targetAmount: *big.NewInt(1e6 + 2e5 + 1e3),
				limit:        3,
			},
			wantErr: ErrNeedMergeCoin,
		},
		{
			name: "ErrNeedSplitGasCoin",
			args: args{
				inputCoins: &Page[Coin, HexData]{
					Data: []Coin{
						{Balance: balanceObject(1e5), CoinType: SUI_COIN_TYPE},
						{Balance: balanceObject(1e5), CoinType: SUI_COIN_TYPE},
					},
					HasNextPage: false,
				},
				targetAmount: *big.NewInt(1e5 + 1),
				gasAmount:    1e4,
			},
			wantErr: ErrNeedSplitGasCoin,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PickupCoins(tt.args.inputCoins, tt.args.targetAmount, tt.args.limit, tt.args.gasAmount)
			if tt.wantErr != nil {
				require.Equal(t, err, tt.wantErr)
			} else {
				require.Equal(t, got, tt.want)
			}
		})
	}
}
