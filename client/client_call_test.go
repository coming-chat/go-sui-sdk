package client

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/coming-chat/go-sui/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

//func TestClient_BatchGetTransaction(t *testing.T) {
//	chain := DevnetClient(t)
//	coins, err := chain.GetCoins(context.TODO(), *Address, nil, nil, 1)
//	require.NoError(t, err)
//	object, err := chain.GetObject(context.TODO(), coins.Data[0].CoinObjectId, nil)
//	require.NoError(t, err)
//	type args struct {
//		digests []string
//	}
//	tests := []struct {
//		name    string
//		chain   *Client
//		args    args
//		want    int
//		wantErr bool
//	}{
//		{
//			name:  "test for devnet transaction",
//			chain: chain,
//			args: args{
//				digests: []string{*object.Data.PreviousTransaction},
//			},
//			want:    1,
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := tt.chain.BatchGetTransaction(tt.args.digests)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("BatchGetTransaction() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if len(got) != tt.want {
//				t.Errorf("BatchGetTransaction() got = %v, want %v", got, tt.want)
//			}
//			t.Logf("%+v", got)
//		})
//	}
//}

func TestClient_DryRunTransaction(t *testing.T) {
	chain := DevnetClient(t)
	coins, err := chain.GetSuiCoinsOwnedByAddress(context.TODO(), *Address)
	require.NoError(t, err)
	coin, err := coins.PickCoinNoLess(2000)
	require.NoError(t, err)
	tx, err := chain.TransferSui(context.TODO(), *Address, *Address, coin.CoinObjectId, 1000, 1000)
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
		t.Run(
			tt.name, func(t *testing.T) {
				result, err := tt.chain.DryRunTransaction(tt.args.ctx, tt.args.tx)
				if (err != nil) != tt.wantErr {
					t.Errorf("Client.DryRunTransaction() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				t.Logf("%#v", result)
			},
		)
	}
}

// TestClient_ExecuteTransactionSerializedSig
// This test case will affect the real coin in the test case of account
// temporary disabled
//func TestClient_ExecuteTransactionSerializedSig(t *testing.T) {
//	chain := DevnetClient(t)
//	coins, err := chain.GetSuiCoinsOwnedByAddress(context.TODO(), *Address)
//	require.NoError(t, err)
//	coin, err := coins.PickCoinNoLess(2000)
//	require.NoError(t, err)
//	tx, err := chain.TransferSui(context.TODO(), *Address, *Address, coin.Reference.ObjectId, 1000, 1000)
//	require.NoError(t, err)
//	account := M1Account(t)
//	signedTx := tx.SignSerializedSigWith(account.PrivateKey)
//	txResult, err := chain.ExecuteTransactionSerializedSig(context.TODO(), *signedTx, types.TxnRequestTypeWaitForEffectsCert)
//	require.NoError(t, err)
//	t.Logf("%#v", txResult)
//}

//func TestClient_ExecuteTransaction(t *testing.T) {
//	chain := DevnetClient(t)
//	coins, err := chain.GetSuiCoinsOwnedByAddress(context.TODO(), *Address)
//	require.NoError(t, err)
//	coin, err := coins.PickCoinNoLess(2000)
//	require.NoError(t, err)
//	tx, err := chain.TransferSui(context.TODO(), *Address, *Address, coin.Reference.ObjectId, 1000, 1000)
//	require.NoError(t, err)
//	account := M1Account(t)
//	signedTx := tx.SignSerializedSigWith(account.PrivateKey)
//	txResult, err := chain.ExecuteTransaction(context.TODO(), *signedTx, types.TxnRequestTypeWaitForEffectsCert)
//	require.NoError(t, err)
//	t.Logf("%#v", txResult)
//}

func TestClient_BatchGetObjectsOwnedByAddress(t *testing.T) {
	cli := DevnetClient(t)

	options := types.SuiObjectDataOptions{
		ShowType:    true,
		ShowContent: true,
	}
	coinType := fmt.Sprintf("0x2::coin::Coin<%v>", types.SuiCoinType)
	filterObject, err := cli.BatchGetObjectsOwnedByAddress(context.TODO(), *Address, options, coinType)
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
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := tt.chain.GetSuiCoinsOwnedByAddress(tt.args.ctx, tt.args.address)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetSuiCoinsOwnedByAddress() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				t.Logf("coin data: %v", got)
			},
		)
	}
}

func TestClient_GetCoinMetadata(t *testing.T) {
	chain := DevnetClient(t)
	metadata, err := chain.GetCoinMetadata(context.TODO(), types.SuiCoinType)
	require.Nil(t, err)
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

//func TestClient_PaySui(t *testing.T) {
//	chain := DevnetClient(t)
//
//	recipients := []types.Address{*Address}
//
//	coins, err := chain.GetSuiCoinsOwnedByAddress(context.TODO(), *Address)
//	require.NoError(t, err)
//	coin, err := coins.PickCoinNoLess(2000)
//	require.NoError(t, err)
//	inputCoins := []types.ObjectId{coin.CoinObjectId}
//
//	tx, err := chain.PaySui(context.TODO(), *Address, inputCoins, recipients, []uint64{1000}, 1000)
//	require.NoError(t, err)
//	t.Logf("%#v", tx)
//
//	inspectResult, err := chain.DryRunTransaction(context.TODO(), tx)
//	require.NoError(t, err)
//	if inspectResult.Status.Error != "" {
//		t.Fatalf("%#v", inspectResult)
//	}
//	t.Logf("%#v", inspectResult)
//}

func TestClient_GetAllBalances(t *testing.T) {
	chain := DevnetClient(t)
	balances, err := chain.GetAllBalances(context.TODO(), *Address)
	require.NoError(t, err)
	t.Logf("%#v", balances)
}

func TestClient_GetBalance(t *testing.T) {
	chain := DevnetClient(t)
	balance, err := chain.GetBalance(context.TODO(), *Address, "")
	require.NoError(t, err)
	t.Logf("%#v", balance)
}

//func TestClient_DevInspectMoveCall(t *testing.T) {
//	chain := DevnetClient(t)
//
//	packageId, err := types.NewHexData("0xb08873e9b44960657723604e4f6bc70c2d1c2b50")
//	require.NoError(t, err)
//
//	devInspectResults, err := chain.DevInspectMoveCall(
//		context.TODO(),
//		*Address,
//		*packageId,
//		"profile",
//		"register",
//		[]string{},
//		[]any{
//			"0xae71509d1be0c751bbced577bd1598e617161c29",
//			"",
//			"",
//		},
//	)
//	require.NoError(t, err)
//	if devInspectResults.Effects.Status.Error != "" {
//		t.Fatalf("%#v", devInspectResults)
//	}
//	t.Logf("%T", devInspectResults)
//}

//func TestClient_DevInspectTransactionBlock(t *testing.T) {
//	chain := DevnetClient(t)
//	packageId, err := types.NewAddressFromHex("0x2")
//	require.NoError(t, err)
//	require.NoError(t, err)
//	arg := sui_types.MoveCallArg{
//		"ComingChat NFT",
//		"This is a NFT created by ComingChat",
//		"https://coming.chat/favicon.ico",
//	}
//	args, err := arg.GetMoveCallArgs()
//	require.NoError(t, err)
//	tKind := sui_types.TransactionKind{
//		Single: &sui_types.SingleTransactionKind{
//			Call: &sui_types.MoveCall{
//				Package:       *packageId,
//				Module:        "devnet_nft",
//				Function:      "mint",
//				TypeArguments: []*sui_types.TypeTag{},
//				Arguments:     args,
//			},
//		},
//	}
//	txBytes, err := bcs.Marshal(tKind)
//	require.NoError(t, err)
//
//	devInspectResults, err := chain.DevInspectTransactionBlock(context.TODO(), *Address, types.Bytes(txBytes).GetBase64Data(), nil, nil)
//	require.NoError(t, err)
//	if devInspectResults.Effects.Status.Error != "" {
//		t.Fatalf("%#v", devInspectResults)
//	}
//	t.Logf("%#v", devInspectResults)
//}

func TestClient_GetCoins(t *testing.T) {
	chain := DevnetClient(t)
	defaultCoinType := types.SuiCoinType
	coins, err := chain.GetCoins(context.TODO(), *Address, &defaultCoinType, nil, 1)
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
		want    *types.PaginatedCoins
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
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := chain.GetAllCoins(tt.args.ctx, tt.args.address, tt.args.cursor, tt.args.limit)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetAllCoins() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				t.Logf("%#v", got)
			},
		)
	}
}

//func TestClient_SplitCoin(t *testing.T) {
//	cli := DevnetClient(t)
//
//	coins, err := cli.GetSuiCoinsOwnedByAddress(context.TODO(), *Address)
//	require.NoError(t, err)
//
//	firstCoin, err := coins.PickCoinNoLess(100)
//	require.NoError(t, err)
//	everyAmount := uint64(firstCoin.Balance) / 2
//	amounts := []uint64{everyAmount, everyAmount}
//
//	txn, err := cli.SplitCoin(context.TODO(), *Address, firstCoin.CoinObjectId, amounts, nil, 1000)
//	require.NoError(t, err)
//
//	t.Log(txn.TxBytes.String())
//
//	inspectTxResult, err := cli.DryRunTransaction(context.TODO(), txn)
//	require.NoError(t, err)
//	if inspectTxResult.Status.Error != "" {
//		t.Fatalf("%#v", inspectTxResult)
//	}
//	t.Logf("%#v", inspectTxResult)
//}

//func TestClient_SplitCoinEqual(t *testing.T) {
//	cli := DevnetClient(t)
//
//	coins, err := cli.GetSuiCoinsOwnedByAddress(context.TODO(), *Address)
//	require.NoError(t, err)
//
//	firstCoin, err := coins.PickCoinNoLess(1000)
//	require.NoError(t, err)
//
//	getCoins, err := cli.GetCoins(context.TODO(), *Address, nil, nil, 0)
//	require.NoError(t, err)
//
//	txn, err := cli.SplitCoinEqual(context.TODO(), *Address, firstCoin.CoinObjectId, 2, &getCoins.Data[len(getCoins.Data)-1].CoinObjectId, 1000)
//	require.NoError(t, err)
//
//	t.Log(txn.TxBytes.String())
//
//	inspectRes, err := cli.DryRunTransaction(context.TODO(), txn)
//	require.NoError(t, err)
//
//	if inspectRes.Status.Error != "" {
//		t.Fatalf("%#v", inspectRes)
//	}
//	t.Logf("%#v", inspectRes)
//}

func TestClient_GetTransaction(t *testing.T) {
	cli := DevnetClient(t)
	digest := "5rMRjX2HWFcWeeNUvMBmpBEa44zsVV7JSNayrGwhVRPy"
	resp, err := cli.GetTransactionBlock(
		context.Background(), digest, types.SuiTransactionBlockResponseOptions{
			ShowInput: false,
		},
	)
	require.NoError(t, err)
	t.Logf("%#v", resp)
}

func TestBatchCall_GetObject(t *testing.T) {
	cli := DevnetClient(t)

	if false {
		// get sepcified object
		idstr := "0x4ad2f0a918a241d6a19573212aeb56947bb9255a14e921a7ec78b262536826f0"
		objId, err := types.NewHexData(idstr)
		require.Nil(t, err)
		obj, err := cli.GetObject(
			context.Background(), *objId, &types.SuiObjectDataOptions{
				ShowType:    true,
				ShowContent: true,
			},
		)
		require.Nil(t, err)
		t.Log(obj.Data)
	}

	coins, err := cli.GetCoins(context.TODO(), *Address, nil, nil, 3)
	require.NoError(t, err)
	if len(coins.Data) == 0 {
		return
	}
	objId := coins.Data[0].CoinObjectId
	obj, err := cli.GetObject(context.Background(), objId, nil)
	require.Nil(t, err)
	t.Log(obj.Data)
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
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := tt.chain.GetObject(tt.args.ctx, tt.args.objID, nil)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetObject() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				t.Logf("%+v", got)
			},
		)
	}
}

func TestClient_MultiGetObjects(t *testing.T) {
	chain := DevnetClient(t)
	coins, err := chain.GetCoins(context.TODO(), *Address, nil, nil, 1)
	require.NoError(t, err)
	if len(coins.Data) != 0 {
		t.Log("Warning: No Object Id for test.")
		return
	}

	obj := coins.Data[0].CoinObjectId
	objs := []types.ObjectId{obj, obj}
	resp, err := chain.MultiGetObjects(
		context.Background(), objs, &types.SuiObjectDataOptions{
			ShowType:  true,
			ShowOwner: true,
		},
	)
	require.Nil(t, err)
	require.Equal(t, len(objs), len(resp))
	require.Equal(t, resp[0], resp[1])
}

func TestClient_GetOwnedObjects(t *testing.T) {
	cli := DevnetClient(t)

	obj, err := types.NewHexData("0x02")
	require.Nil(t, err)
	query := types.SuiObjectResponseQuery{
		Filter: &types.SuiObjectDataFilter{
			Package: obj,
			// StructType: "0x2::coin::Coin<0x2::sui::SUI>",
		},
		Options: &types.SuiObjectDataOptions{
			ShowType: true,
		},
	}

	objs, err := cli.GetOwnedObjects(context.Background(), *Address, &query, nil, 0)
	require.Nil(t, err)
	t.Log(objs.Data)
}

func TestBatchGetObjectsOwnedByAddress(t *testing.T) {
	cli := DevnetClient(t)
	coins, err := cli.GetSuiCoinsOwnedByAddress(context.TODO(), *Address)
	require.NoError(t, err)

	t.Logf("%#v", coins)
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
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := tt.chain.GetTotalSupply(tt.args.ctx, tt.args.coinType)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetTotalSupply() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				t.Logf("%d", got)
			},
		)
	}
}
func TestClient_GetTotalTransactionBlocks(t *testing.T) {
	cli := DevnetClient(t)
	res, err := cli.GetTotalTransactionBlocks(context.Background())
	require.Nil(t, err)
	t.Log(res)
}

//func TestClient_Publish(t *testing.T) {
//	chain := DevnetClient(t)
//	dmens, err := types.NewBase64Data(DmensDmensB64)
//	require.NoError(t, err)
//	profile, err := types.NewBase64Data(DmensProfileB64)
//	require.NoError(t, err)
//	coins, err := chain.GetSuiCoinsOwnedByAddress(context.TODO(), *Address)
//	require.NoError(t, err)
//	coin, err := coins.PickCoinNoLess(30000)
//	require.NoError(t, err)
//	type args struct {
//		ctx             context.Context
//		address         types.Address
//		compiledModules []*types.Base64Data
//		gas             types.ObjectId
//		gasBudget       uint
//	}
//	tests := []struct {
//		name    string
//		client  *Client
//		args    args
//		want    *types.TransactionBytes
//		wantErr bool
//	}{
//		{
//			name:   "test for dmens publish",
//			client: chain,
//			args: args{
//				ctx:             context.TODO(),
//				address:         *Address,
//				compiledModules: []*types.Base64Data{dmens, profile},
//				gas:             coin.CoinObjectId,
//				gasBudget:       30000,
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := tt.client.Publish(tt.args.ctx, tt.args.address, tt.args.compiledModules, tt.args.gas, tt.args.gasBudget)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			t.Logf("%#v", got)
//
//			txResult, err := tt.client.DryRunTransaction(context.TODO(), got)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//
//			t.Logf("%#v", txResult)
//		})
//	}
//}

func TestClient_TryGetPastObject(t *testing.T) {
	cli := DevnetClient(t)
	objId, err := types.NewHexData("0x11462c88e74bb00079e3c043efb664482ee4551744ee691c7623b98503cb3f4d")
	require.Nil(t, err)
	data, err := cli.TryGetPastObject(context.Background(), *objId, 903, nil)
	require.Nil(t, err)
	t.Log(data)
}

func TestClient_GetEvents(t *testing.T) {
	cli := DevnetClient(t)
	digest := "bWEVPGbA81GDJ4655fFuiabV11Z2gSgJyqfURXyNL6G"
	res, err := cli.GetEvents(context.Background(), digest)
	require.NoError(t, err)
	t.Log(res)
}

func TestClient_GetReferenceGasPrice(t *testing.T) {
	cli := DevnetClient(t)
	gasPrice, err := cli.GetReferenceGasPrice(context.Background())
	require.Nil(t, err)
	t.Logf("current gas price = %v", gasPrice)
}

func TestClient_DevInspectTransactionBlock(t *testing.T) {
	chain := DevnetClient(t)
	price, err := chain.GetReferenceGasPrice(context.TODO())
	require.NoError(t, err)
	coins, err := chain.GetSuiCoinsOwnedByAddress(context.TODO(), *Address)
	require.NoError(t, err)
	coin, err := coins.PickCoinNoLess(1000)
	require.NoError(t, err)
	tx, err := chain.TransferSui(context.TODO(), *Address, *Address, coin.CoinObjectId, 1000, 1000)
	require.NoError(t, err)
	type args struct {
		ctx           context.Context
		senderAddress types.Address
		txByte        types.Base64Data
		gasPrice      *decimal.Decimal
		epoch         *uint64
		chain         *Client
	}
	tests := []struct {
		name    string
		args    args
		want    *types.DevInspectResults
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				ctx:           context.TODO(),
				senderAddress: *Address,
				txByte:        tx.TxBytes,
				epoch:         nil,
				gasPrice:      price,
				chain:         chain,
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := tt.args.chain.DevInspectTransactionBlock(
					tt.args.ctx,
					tt.args.senderAddress,
					tt.args.txByte,
					tt.args.gasPrice,
					tt.args.epoch,
				)
				if (err != nil) != tt.wantErr {
					t.Errorf("DevInspectTransactionBlock() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("DevInspectTransactionBlock() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestClient_QueryTransactionBlocks(t *testing.T) {
	cli := DevnetClient(t)
	limit := uint(10)
	type args struct {
		ctx             context.Context
		query           types.SuiTransactionBlockResponseQuery
		cursor          *types.TransactionDigest
		limit           *uint
		descendingOrder bool
	}
	tests := []struct {
		name    string
		args    args
		want    *types.TransactionBlocksPage
		wantErr bool
	}{
		{
			name: "test for queryTransactionBlocks",
			args: args{
				ctx: context.TODO(),
				query: types.SuiTransactionBlockResponseQuery{
					Filter: &types.TransactionFilter{
						FromAddress: Address,
					},
					Options: &types.SuiTransactionBlockResponseOptions{
						ShowInput:   true,
						ShowEffects: true,
					},
				},
				cursor:          nil,
				limit:           &limit,
				descendingOrder: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := cli.QueryTransactionBlocks(
					tt.args.ctx,
					tt.args.query,
					tt.args.cursor,
					tt.args.limit,
					tt.args.descendingOrder,
				)
				if (err != nil) != tt.wantErr {
					t.Errorf("QueryTransactionBlocks() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				t.Logf("%#v", got)
			},
		)
	}
}

func TestClient_QueryEvents(t *testing.T) {
	cli := DevnetClient(t)
	limit := uint(10)
	type args struct {
		ctx             context.Context
		query           types.EventFilter
		cursor          *types.EventId
		limit           *uint
		descendingOrder bool
	}
	tests := []struct {
		name    string
		args    args
		want    *types.EventPage
		wantErr bool
	}{
		{
			name: "test for query events",
			args: args{
				ctx: context.TODO(),
				query: types.EventFilter{
					Sender: Address,
				},
				cursor:          nil,
				limit:           &limit,
				descendingOrder: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := cli.QueryEvents(
					tt.args.ctx,
					tt.args.query,
					tt.args.cursor,
					tt.args.limit,
					tt.args.descendingOrder,
				)
				if (err != nil) != tt.wantErr {
					t.Errorf("QueryEvents() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("QueryEvents() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
