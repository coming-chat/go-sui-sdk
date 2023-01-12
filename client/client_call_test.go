package client

import (
	"context"
	"github.com/coming-chat/go-sui/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClient_GetTransactionsInRange(t *testing.T) {
	chain := DevnetClient(t)
	type args struct {
		ctx   context.Context
		start uint64
		end   uint64
	}
	tests := []struct {
		name    string
		chain   *Client
		args    args
		want    int
		wantErr bool
	}{
		{
			name:  "test for devnet",
			chain: chain,
			args: args{
				ctx:   context.TODO(),
				start: 0,
				end:   10,
			},
			want:    10,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.chain.GetTransactionsInRange(tt.args.ctx, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTransactionsInRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("GetTransactionsInRange() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_BatchGetTransaction(t *testing.T) {
	chain := DevnetClient(t)
	coins, err := chain.GetCoins(context.TODO(), *Address, nil, nil, 1)
	require.NoError(t, err)
	object, err := chain.GetObject(context.TODO(), coins.Data[0].CoinObjectId)
	require.NoError(t, err)
	type args struct {
		digests []string
	}
	tests := []struct {
		name    string
		chain   *Client
		args    args
		want    int
		wantErr bool
	}{
		{
			name:  "test for devnet transaction",
			chain: chain,
			args: args{
				digests: []string{object.Details.PreviousTransaction},
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.chain.BatchGetTransaction(tt.args.digests)
			if (err != nil) != tt.wantErr {
				t.Errorf("BatchGetTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("BatchGetTransaction() got = %v, want %v", got, tt.want)
			}
			t.Logf("%+v", got)
		})
	}
}

func TestClient_BatchGetObject(t *testing.T) {
	type args struct {
		objects []types.ObjectId
	}
	chain := DevnetClient(t)
	coins, err := chain.GetCoins(context.TODO(), *Address, nil, nil, 1)
	require.NoError(t, err)

	tests := []struct {
		name    string
		chain   *Client
		args    args
		want    int
		wantErr bool
	}{
		{
			name:  "test for devnet",
			chain: chain,
			args: args{
				objects: []types.ObjectId{coins.Data[0].CoinObjectId},
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.chain.BatchGetObject(tt.args.objects)
			if (err != nil) != tt.wantErr {
				t.Errorf("BatchGetObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("BatchGetObject() got = %v, want %v", got, tt.want)
			}
			t.Logf("%+v", got)
		})
	}
}

func TestClient_GetObject(t *testing.T) {
	type args struct {
		ctx   context.Context
		objID types.ObjectId
	}
	chain := DevnetClient(t)
	coins, err := chain.GetCoins(context.TODO(), *Address, nil, nil, 1)
	require.NoError(t, err)

	tests := []struct {
		name    string
		chain   *Client
		args    args
		want    int
		wantErr bool
	}{
		{
			name:  "test for devnet",
			chain: chain,
			args: args{
				ctx:   context.TODO(),
				objID: coins.Data[0].CoinObjectId,
			},
			want:    3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.chain.GetObject(tt.args.ctx, tt.args.objID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%+v", got)
		})
	}
}

func TestClient_DryRunTransaction(t *testing.T) {
	chain := DevnetClient(t)
	coins, err := chain.GetCoins(context.TODO(), *Address, nil, nil, 1)
	require.NoError(t, err)
	tx, err := chain.TransferSui(context.TODO(), *Address, *Address, coins.Data[0].CoinObjectId, 1000, 1000)
	require.NoError(t, err)
	type args struct {
		ctx context.Context
		tx  *types.TransactionBytes
	}
	tests := []struct {
		name  string
		args  args
		chain *Client
		// want    *types.TransactionEffects
		wantErr bool
	}{
		{
			name:  "dry run",
			chain: chain,
			args: args{
				ctx: context.TODO(),
				tx:  tx,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.chain.DryRunTransaction(tt.args.ctx, tt.args.tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.DryRunTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%#v", result)
		})
	}
}

func TestClient_GetObjectsOwnedByAddress(t *testing.T) {
	cli := DevnetClient(t)

	objects, err := cli.GetObjectsOwnedByAddress(context.TODO(), *Address)
	require.NoError(t, err)
	t.Log(objects)

	filterObject, err := cli.BatchGetObjectsOwnedByAddress(context.TODO(), *Address, types.SuiCoinType)
	require.NoError(t, err)
	t.Log(filterObject)
}

func TestClient_GetSuiCoinsOwnedByAddress(t *testing.T) {
	chain := DevnetClient(t)
	type args struct {
		ctx     context.Context
		address types.Address
	}
	tests := []struct {
		name    string
		chain   *Client
		args    args
		wantErr bool
	}{
		{
			name:  "case 1",
			chain: chain,
			args: args{
				ctx:     context.TODO(),
				address: *Address,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.chain.GetSuiCoinsOwnedByAddress(tt.args.ctx, tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSuiCoinsOwnedByAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("coin data: %v", got)
		})
	}
}

func TestClient_GetCoinMetadata(t *testing.T) {
	chain := DevnetClient(t)
	metadata, err := chain.GetCoinMetadata(context.TODO(), types.SuiCoinType)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", metadata)
}

// TestClient_Pay need another coin type(not default sui coin)
//func TestClient_Pay(t *testing.T) {
//	chain := DevnetClient(t)
//	coins, err := chain.GetCoins(context.TODO(), *Address, nil, nil, 1)
//	require.NoError(t, err)
//	inputCoins := []types.ObjectId{coins.Data[0].CoinObjectId}
//
//	tx, err := chain.Pay(context.TODO(), *Address, inputCoins, []types.Address{*Address}, []uint64{1000}, coins.Data[len(coins.Data)-1].CoinObjectId, 2000)
//	require.NoError(t, err)
//	t.Logf("%#v", tx)
//	inspectResult, err := chain.DevInspectTransaction(context.TODO(), tx.TxBytes)
//	require.NoError(t, err)
//	t.Logf("%#v", inspectResult)
//}

func TestClient_PaySui(t *testing.T) {
	chain := DevnetClient(t)

	recipients := []types.Address{*Address}

	coins, err := chain.GetCoins(context.TODO(), *Address, nil, nil, 1)
	require.NoError(t, err)
	inputCoins := []types.ObjectId{coins.Data[0].CoinObjectId}

	tx, err := chain.PaySui(context.TODO(), *Address, inputCoins, recipients, []uint64{1000}, 1000)
	require.NoError(t, err)
	t.Logf("%#v", tx)

	inspectResult, err := chain.DevInspectTransaction(context.TODO(), tx.TxBytes)
	require.NoError(t, err)
	if inspectResult.Effects.Status.Error != "" {
		t.Fatalf("%#v", inspectResult)
	}
	t.Logf("%#v", inspectResult)
}

func TestClient_GetAllBalances(t *testing.T) {
	chain := DevnetClient(t)
	balances, err := chain.GetAllBalances(context.TODO(), *Address)
	require.NoError(t, err)
	t.Logf("%#v", balances)
}

func TestClient_GetBalance(t *testing.T) {
	chain := DevnetClient(t)
	balance, err := chain.GetBalance(context.TODO(), *Address, nil)
	require.NoError(t, err)
	t.Logf("%#v", balance)
}

func TestClient_DevInspectMoveCall(t *testing.T) {
	chain := DevnetClient(t)

	packageId, err := types.NewHexData("0x2ac78d1e11aabe14a9b22a9d574ec10079bc10b0")
	require.NoError(t, err)

	devInspectResults, err := chain.DevInspectMoveCall(
		context.TODO(),
		*Address,
		*packageId,
		"dmens",
		"post",
		[]string{},
		[]any{
			"0x1a933e6c29a0113905449a5de595c214579adff7",
			0,
			0,
			"hello",
		},
	)
	require.NoError(t, err)
	if devInspectResults.Effects.Status.Error != "" {
		t.Fatalf("%#v", devInspectResults)
	}
	t.Logf("%T", devInspectResults)
}

func TestClient_DevInspectTransaction(t *testing.T) {
	chain := DevnetClient(t)

	signer, err := types.NewAddressFromHex("0x6fc6148816617c3c3eccb1d09e930f73f6712c9c")
	require.NoError(t, err)
	coin, err := types.NewHexData("0x451d7adec18632a43c1d7948e504be0ee07ef1e9")
	require.NoError(t, err)
	txnBytes, err := chain.MintNFT(
		context.TODO(),
		*signer,
		"ComingChat NFT",
		"This is a NFT created by ComingChat",
		"https://coming.chat/favicon.ico",
		coin,
		12000,
	)
	require.NoError(t, err)
	devInspectResults, err := chain.DevInspectTransaction(context.TODO(), txnBytes.TxBytes)
	require.NoError(t, err)
	if devInspectResults.Effects.Status.Error != "" {
		t.Fatalf("%#v", devInspectResults)
	}
	t.Logf("%#v", devInspectResults)
}

func TestClient_GetCoins(t *testing.T) {
	chain := DevnetClient(t)
	defaultCoinType := types.SuiCoinType
	coins, err := chain.GetCoins(context.TODO(), *Address, &defaultCoinType, nil, 1)
	require.NoError(t, err)
	t.Logf("%#v", coins)
}

func TestClient_SplitCoin(t *testing.T) {
	cli := DevnetClient(t)

	coins, err := cli.GetSuiCoinsOwnedByAddress(context.TODO(), *Address)
	require.NoError(t, err)

	firstCoin, err := coins.PickCoinNoLess(100)
	require.NoError(t, err)
	everyAmount := firstCoin.Balance / 2
	amounts := []uint64{everyAmount, everyAmount}

	txn, err := cli.SplitCoin(context.TODO(), *Address, firstCoin.Reference.ObjectId, amounts, nil, 1000)
	require.NoError(t, err)

	t.Log(txn.TxBytes.String())

	inspectTxResult, err := cli.DevInspectTransaction(context.TODO(), txn.TxBytes)
	require.NoError(t, err)
	if inspectTxResult.Effects.Status.Error != "" {
		t.Fatalf("%#v", inspectTxResult)
	}
	t.Logf("%#v", inspectTxResult)
}

func TestClient_SplitCoinEqual(t *testing.T) {
	cli := DevnetClient(t)

	coins, err := cli.GetSuiCoinsOwnedByAddress(context.TODO(), *Address)
	require.NoError(t, err)

	firstCoin, err := coins.PickCoinNoLess(1000)
	require.NoError(t, err)

	getCoins, err := cli.GetCoins(context.TODO(), *Address, nil, nil, 0)
	require.NoError(t, err)

	txn, err := cli.SplitCoinEqual(context.TODO(), *Address, firstCoin.Reference.ObjectId, 2, &getCoins.Data[len(getCoins.Data)-1].CoinObjectId, 10000)
	require.NoError(t, err)

	t.Log(txn.TxBytes.String())

	inspectRes, err := cli.DevInspectTransaction(context.TODO(), txn.TxBytes)
	require.NoError(t, err)

	if inspectRes.Effects.Status.Error != "" {
		t.Fatalf("%#v", inspectRes)
	}
	t.Logf("%#v", inspectRes)
}

func TestGetTransaction(t *testing.T) {
	cli := DevnetClient(t)
	transactions, err := cli.GetTransactionsInRange(context.TODO(), 0, 1)
	require.NoError(t, err)
	require.NotEmpty(t, transactions)
	resp, err := cli.GetTransaction(context.TODO(), transactions[0])
	require.NoError(t, err)

	t.Logf("%#v", resp)
}

func TestBatchCall_GetObject(t *testing.T) {
	cli := DevnetClient(t)
	coins, err := cli.GetCoins(context.TODO(), *Address, nil, nil, 3)
	require.NoError(t, err)
	var objKeys []string
	for _, v := range coins.Data {
		objKeys = append(objKeys, v.CoinObjectId.String())
	}

	elems := make([]BatchElem, len(objKeys))
	for i := 0; i < len(objKeys); i++ {
		ele := BatchElem{
			Method: "sui_getObject",
			Args:   []interface{}{objKeys[i]},
			Result: &types.ObjectRead{},
		}
		elems[i] = ele
	}

	err = cli.BatchCall(elems)
	require.NoError(t, err)

	t.Logf("%#v", elems)
}

func TestBatchGetObjectsOwnedByAddress(t *testing.T) {
	cli := DevnetClient(t)
	coins, err := cli.GetSuiCoinsOwnedByAddress(context.TODO(), *Address)
	require.NoError(t, err)

	t.Logf("%#v", coins)
}

func TestClient_GetAllCoins(t *testing.T) {
	chain := DevnetClient(t)
	type args struct {
		ctx     context.Context
		address types.Address
		cursor  *types.ObjectId
		limit   uint
	}
	tests := []struct {
		name    string
		chain   *Client
		args    args
		want    *types.CoinPage
		wantErr bool
	}{
		{
			name:  "test case 1",
			chain: chain,
			args: args{
				ctx:     context.TODO(),
				address: *Address,
				cursor:  nil,
				limit:   3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := chain.GetAllCoins(tt.args.ctx, tt.args.address, tt.args.cursor, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllCoins() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%#v", got)
		})
	}
}

func TestClient_GetTotalSupply(t *testing.T) {
	chain := DevnetClient(t)
	type args struct {
		ctx      context.Context
		coinType string
	}
	tests := []struct {
		name    string
		chain   *Client
		args    args
		want    uint64
		wantErr bool
	}{
		{
			name:  "test 1",
			chain: chain,
			args: args{
				context.TODO(),
				types.SuiCoinType,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.chain.GetTotalSupply(tt.args.ctx, tt.args.coinType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTotalSupply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%d", got)
		})
	}
}
