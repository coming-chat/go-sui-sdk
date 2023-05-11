package sui_types

import (
	"github.com/coming-chat/go-sui/lib"
	"github.com/coming-chat/go-sui/move_types"
	"github.com/fardream/go-bcs/bcs"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferSui(t *testing.T) {
	ptb := NewProgrammableTransactionBuilder()
	recipient, err := NewAddressFromHex("0x7e875ea78ee09f08d72e2676cf84e0f1c8ac61d94fa339cc8e37cace85bebc6e")
	require.NoError(t, err)
	amount := uint64(100000)
	err = ptb.TransferSui(*recipient, &amount)
	require.NoError(t, err)
	pt := ptb.Finish()
	kind := TransactionKind{
		ProgrammableTransaction: &pt,
	}
	digest, err := NewDigest("HvbE2UZny6cP4KukaXetmj4jjpKTDTjVo23XEcu7VgSn")
	require.NoError(t, err)
	objectId, err := move_types.NewAccountAddressHex("0x13c1c3d0e15b4039cec4291c75b77c972c10c8e8e70ab4ca174cf336917cb4db")
	require.NoError(t, err)
	tx := TransactionData{
		V1: &TransactionDataV1{
			Kind:   kind,
			Sender: *recipient,
			GasData: GasData{
				Price: 1000,
				Owner: *recipient,
				Payment: []*ObjectRef{
					{
						ObjectId: *objectId,
						Version:  14924029,
						Digest:   *digest,
					},
				},
				Budget: 10000000,
			},
			Expiration: TransactionExpiration{
				None: &lib.EmptyEnum{},
			},
		},
	}
	txByte, err := bcs.Marshal(tx)
	require.NoError(t, err)
	t.Logf("%x", txByte)
}
