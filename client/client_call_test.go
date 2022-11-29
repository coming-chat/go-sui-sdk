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
		digests []types.Digest
	}
	var (
		d1, _ = types.NewBase64Data("2Tpxt7+bP0XP1h0s4j0vWyYYFJQGusuctv5LEuUiSPY=")
		d2, _ = types.NewBase64Data("JV657cXeBAmezgkjIMWASF5ugJR6OhYVfSJZumhUrlE=")
		d3, _ = types.NewBase64Data("kkgD8Vu6g2YLEa0jG6Teh9xo3A78OocWkUb27InA968=")
	)
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
				digests: []types.Digest{*d1, *d2, *d3},
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
		o1, _ = types.NewHexData("0x582d1e989991cd4255ac3d2ba5ac7db15d3077ba")
		o2, _ = types.NewHexData("0x5abee7585fcb043a30e827a5bad42132c7a243ca")
		o3, _ = types.NewHexData("0x8b5bfe60fe69e40f1565802d41c3950955e8fead")
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
