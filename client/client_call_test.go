package client

import (
	"context"
	"github.com/coming-chat/go-sui/types"
	"net/http"
	"testing"
	"time"
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
