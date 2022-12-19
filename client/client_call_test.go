package client

import (
	"context"
	"github.com/coming-chat/go-sui/account"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/coming-chat/go-sui/types"
	"github.com/stretchr/testify/require"
)

var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:    3,
		IdleConnTimeout: 30 * time.Second,
	},
	Timeout: 30 * time.Second,
}

func TestClient_GetTransactionsInRange(t *testing.T) {
	type fields struct {
		idCounter uint32
		rpcUrl    string
		client    *http.Client
	}
	type args struct {
		ctx   context.Context
		start uint64
		end   uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "test for devnet",
			fields: fields{
				idCounter: 0,
				rpcUrl:    DevnetRpcUrl,
				client:    client,
			},
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
			c := &Client{
				idCounter: tt.fields.idCounter,
				rpcUrl:    tt.fields.rpcUrl,
				client:    tt.fields.client,
			}
			got, err := c.GetTransactionsInRange(tt.args.ctx, tt.args.start, tt.args.end)
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
	type fields struct {
		idCounter uint32
		rpcUrl    string
		client    *http.Client
	}
	type args struct {
		digests []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "test for devnet transaction",
			fields: fields{
				idCounter: 1,
				rpcUrl:    DevnetRpcUrl,
				client:    client,
			},
			args: args{
				digests: []string{"TkLw7eH9NtKh6pSb7evL8EcCf7RDMEsJ3VU7FqJRpf8"},
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				idCounter: tt.fields.idCounter,
				rpcUrl:    tt.fields.rpcUrl,
				client:    tt.fields.client,
			}
			got, err := c.BatchGetTransaction(tt.args.digests)
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
	type fields struct {
		idCounter uint32
		rpcUrl    string
		client    *http.Client
	}
	type args struct {
		objects []types.ObjectId
	}
	var (
		o1, _ = types.NewHexData("0x523203b287a2c1df0a707a6b563aa7d29bd216d6")
		o2, _ = types.NewHexData("0xb1e550000000000000000000000000000000008c")
		o3, _ = types.NewHexData("0xb1e550000000000000000000000000000000005a")
	)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "test for devnet",
			fields: fields{
				idCounter: 1,
				rpcUrl:    DevnetRpcUrl,
				client:    client,
			},
			args: args{
				objects: []types.ObjectId{*o1, *o2, *o3},
			},
			want:    3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				idCounter: tt.fields.idCounter,
				rpcUrl:    tt.fields.rpcUrl,
				client:    tt.fields.client,
			}
			got, err := c.BatchGetObject(tt.args.objects)
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
	type fields struct {
		idCounter uint32
		rpcUrl    string
		client    *http.Client
	}
	type args struct {
		ctx   context.Context
		objID types.ObjectId
	}
	var (
		o, _ = types.NewHexData("0xb1e55000000000000000000000000000000000ca")
	)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "test for devnet",
			fields: fields{
				idCounter: 1,
				rpcUrl:    DevnetRpcUrl,
				client:    client,
			},
			args: args{
				ctx:   context.TODO(),
				objID: *o,
			},
			want:    3,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				idCounter: tt.fields.idCounter,
				rpcUrl:    tt.fields.rpcUrl,
				client:    tt.fields.client,
			}
			got, err := c.GetObject(tt.args.ctx, tt.args.objID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%+v", got)
		})
	}
}

func TestClient_DryRunTransaction(t *testing.T) {
	c, err := Dial(DevnetRpcUrl)
	if err != nil {
		t.Logf("%e", err)
	}
	type args struct {
		ctx context.Context
		tx  *types.TransactionBytes
	}
	signer, err := types.NewAddressFromHex("0x3503a04d1e0de4f58d10484122d6dc42abfbe291")
	if err != nil {
		t.Logf("%e", err)
	}
	coins, err := c.GetSuiCoinsOwnedByAddress(context.Background(), *signer)
	if err != nil {
		t.Logf("%e", err)
	}
	coin, err := coins.PickCoinNoLess(1000)
	if err != nil {
		t.Fatal(err)
	}
	typeArgs := []string{}
	objectId, _ := types.NewHexData("0x00e2cd853b00a1531b5a5579156a174891543e50")
	arguments := []any{
		objectId,
		[]byte("13e8531463853d9a3ff017d140be14a9357f6b1d::coins::BTC"),
	}
	packageId, err := types.NewHexData("0xe558bd8e7a6a88a405ffd93cc71ecf1ade45686c")
	if err != nil {
		t.Logf("%e", err)
	}
	tx, err := c.MoveCall(context.Background(), *signer, *packageId, "interfaces", "get_dola_token_liquidity", typeArgs, arguments, &coin.Reference.ObjectId, 1000)
	if err != nil {
		t.Logf("%e", err)
	}
	tests := []struct {
		name string
		args args
		// want    *types.TransactionEffects
		wantErr bool
	}{
		{
			name: "dry run",
			args: args{
				ctx: context.Background(),
				tx:  tx,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := c.DryRunTransaction(tt.args.ctx, tt.args.tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.DryRunTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_GetObjectsOwnedByAddress(t *testing.T) {
	address := "0x6c5d2cd6e62734f61b4e318e58cbfd1c4b99dfaf"
	cli := DevnetClient(t)

	addr, err := types.NewAddressFromHex(address)
	require.Nil(t, err)
	objects, err := cli.GetObjectsOwnedByAddress(context.Background(), *addr)
	require.Nil(t, err)
	t.Log(objects)

	objectsss, err := cli.BatchGetObjectsOwnedByAddress(context.Background(), *addr, "0xe6e1a26c0be45fc0ec73521c2d3ca268a843a89b::capy::Capy")
	require.Nil(t, err)
	t.Log(objectsss)
}

func TestClient_GetSuiCoinsOwnedByAddress(t *testing.T) {
	addr, err := types.NewAddressFromHex("0x6fc6148816617c3c3eccb1d09e930f73f6712c9c")
	if err != nil {
		t.Error(err)
	}
	type fields struct {
		idCounter uint32
		rpcUrl    string
		client    *http.Client
	}
	type args struct {
		ctx     context.Context
		address types.Address
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "case 1",
			fields: fields{
				idCounter: 1,
				rpcUrl:    DevnetRpcUrl,
				client:    client,
			},
			args: args{
				ctx:     context.TODO(),
				address: *addr,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				idCounter: tt.fields.idCounter,
				rpcUrl:    tt.fields.rpcUrl,
				client:    tt.fields.client,
			}
			got, err := c.GetSuiCoinsOwnedByAddress(tt.args.ctx, tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSuiCoinsOwnedByAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("coin data: %v", got)
		})
	}
}

func TestClient_CallDmens(t *testing.T) {
	chatGPTSuiAccount, err := account.NewAccountWithMnemonic(os.Getenv("mnemonic"))
	if err != nil {
		t.Fatal(err)
	}
	signer, err := types.NewAddressFromHex(chatGPTSuiAccount.Address)
	if err != nil {
		t.Fatal(err)
	}
	rpc, err := Dial(DevnetRpcUrl)
	if err != nil {
		t.Fatal(err)
	}
	coins, err := rpc.GetSuiCoinsOwnedByAddress(context.Background(), *signer)
	if err != nil {
		t.Fatal(err)
	}
	coin, err := coins.PickCoinNoLess(1000)
	if err != nil {
		t.Fatal(err)
	}
	typeArgs := []string{}
	arguments := []any{
		os.Getenv("dmensId"),
		1,
		3,
		"test",
		os.Getenv("relyId"),
	}
	packageId, err := types.NewHexData("0x8e27cccd2250bedd6a437efb3d5abcc3f1d60f04")
	if err != nil {
		t.Fatal(err)
	}
	tx, err := rpc.MoveCall(
		context.Background(),
		*signer,
		*packageId,
		"dmens",
		"post_with_ref",
		typeArgs,
		arguments,
		&coin.Reference.ObjectId,
		1000,
	)
	if err != nil {
		t.Fatal(err)
	}

	signedTxn := tx.SignWith(chatGPTSuiAccount.PrivateKey)
	result, err := rpc.ExecuteTransaction(context.Background(), *signedTxn, types.TxnRequestTypeWaitForLocalExecution)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}
