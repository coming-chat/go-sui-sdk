package types

import (
	"math/big"
	"reflect"
	"testing"
)

func TestCoins_PickSUICoinsWithGas(t *testing.T) {
	// coins 1,2,3,4,5
	testCoins := Coins{
		{
			Balance: 3,
		},
		{
			Balance: 5,
		},
		{
			Balance: 1,
		},
		{
			Balance: 4,
		},
		{
			Balance: 2,
		},
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
			want:    Coins{{Balance: 1}},
			want1:   &Coin{Balance: 2},
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
			want:    Coins{{Balance: 1}, {Balance: 3}},
			want1:   &Coin{Balance: 2},
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
			want:    Coins{{Balance: 1}, {Balance: 3}, {Balance: 4}},
			want1:   &Coin{Balance: 2},
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
			want1:   &Coin{Balance: 3},
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
			want:    Coins{{Balance: 5}},
			want1:   &Coin{Balance: 3},
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
			want:    Coins{{Balance: 5}},
			want1:   &Coin{Balance: 3},
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
		{
			Balance: 3,
		},
		{
			Balance: 5,
		},
		{
			Balance: 1,
		},
		{
			Balance: 4,
		},
		{
			Balance: 2,
		},
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
			want:    Coins{{Balance: 1}, {Balance: 2}},
			wantErr: false,
		},
		{
			name:    "smaller 2",
			cs:      testCoins,
			args:    args{amount: big.NewInt(4), pickMethod: PickSmaller},
			want:    Coins{{Balance: 1}, {Balance: 2}, {Balance: 3}},
			wantErr: false,
		},
		{
			name:    "bigger 1",
			cs:      testCoins,
			args:    args{amount: big.NewInt(2), pickMethod: PickBigger},
			want:    Coins{{Balance: 5}},
			wantErr: false,
		},
		{
			name:    "bigger 2",
			cs:      testCoins,
			args:    args{amount: big.NewInt(6), pickMethod: PickBigger},
			want:    Coins{{Balance: 5}, {Balance: 4}},
			wantErr: false,
		},
		{
			name:    "pick by order 1",
			cs:      testCoins,
			args:    args{amount: big.NewInt(6), pickMethod: PickByOrder},
			want:    Coins{{Balance: 3}, {Balance: 5}},
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
