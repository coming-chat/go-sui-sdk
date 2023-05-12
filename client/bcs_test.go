package client

import (
	"context"
	"testing"

	"github.com/coming-chat/go-sui/sui_types"
	"github.com/coming-chat/go-sui/types"
	"github.com/fardream/go-bcs/bcs"
	"github.com/stretchr/testify/require"
)

func TestBCS_TransferObject(t *testing.T) {
	sender, err := sui_types.NewAddressFromHex("0xd77955e670f42c1bc5e94b9e68e5fe9bdbed9134d784f2a14dfe5fc1b24b5d9f")
	require.NoError(t, err)
	recipient := sender
	gasBudget := SUI(0.01).Uint64()

	cli := TestnetClient(t)
	coins := GetCoins(t, cli, *sender, 2)
	coin, gas := coins[0], coins[1]

	gasPrice := uint64(1000)
	// gasPrice, err := cli.GetReferenceGasPrice(context.Background())

	// build with BCS
	ptb := sui_types.NewProgrammableTransactionBuilder()
	err = ptb.TransferObject(*recipient, []*sui_types.ObjectRef{coin.Reference()})
	require.NoError(t, err)
	pt := ptb.Finish()
	tx := sui_types.NewProgrammable(
		*sender, []*sui_types.ObjectRef{
			gas.Reference(),
		},
		pt, gasBudget, gasPrice,
	)
	txBytesBCS, err := bcs.Marshal(tx)
	require.NoError(t, err)

	// build with remote rpc
	txn, err := cli.TransferObject(context.Background(), *sender, *recipient,
		coin.CoinObjectId,
		&gas.CoinObjectId,
		types.NewSafeSuiBigInt(gasBudget))
	require.NoError(t, err)
	txBytesRemote := txn.TxBytes.Data()

	require.Equal(t, txBytesBCS, txBytesRemote)
}

func TestBCS_TransferSui(t *testing.T) {
	sender, err := sui_types.NewAddressFromHex("0xd77955e670f42c1bc5e94b9e68e5fe9bdbed9134d784f2a14dfe5fc1b24b5d9f")
	require.NoError(t, err)
	recipient := sender
	amount := SUI(0.001).Uint64()
	gasBudget := SUI(0.01).Uint64()

	cli := TestnetClient(t)
	coin := GetCoins(t, cli, *sender, 1)[0]

	gasPrice := uint64(1000)
	// gasPrice, err := cli.GetReferenceGasPrice(context.Background())

	// build with BCS
	ptb := sui_types.NewProgrammableTransactionBuilder()
	err = ptb.TransferSui(*recipient, &amount)
	require.NoError(t, err)
	pt := ptb.Finish()
	tx := sui_types.NewProgrammable(
		*sender, []*sui_types.ObjectRef{
			coin.Reference(),
		},
		pt, gasBudget, gasPrice,
	)
	txBytesBCS, err := bcs.Marshal(tx)
	require.NoError(t, err)

	// build with remote rpc
	txn, err := cli.TransferSui(context.Background(), *sender, *recipient, coin.CoinObjectId,
		types.NewSafeSuiBigInt(amount),
		types.NewSafeSuiBigInt(gasBudget))
	require.NoError(t, err)
	txBytesRemote := txn.TxBytes.Data()

	require.Equal(t, txBytesBCS, txBytesRemote)
}

func TestBCS_PaySui(t *testing.T) {
	sender, err := sui_types.NewAddressFromHex("0xd77955e670f42c1bc5e94b9e68e5fe9bdbed9134d784f2a14dfe5fc1b24b5d9f")
	require.NoError(t, err)
	recipient := sender
	recipient2, _ := sui_types.NewAddressFromHex("0x123456")
	amount := SUI(0.001).Uint64()
	gasBudget := SUI(0.01).Uint64()

	cli := TestnetClient(t)
	coin := GetCoins(t, cli, *sender, 1)[0]

	gasPrice := uint64(1000)
	// gasPrice, err := cli.GetReferenceGasPrice(context.Background())

	// build with BCS
	ptb := sui_types.NewProgrammableTransactionBuilder()
	err = ptb.PaySui([]suiAddress{*recipient, *recipient2}, []uint64{amount, amount})
	require.NoError(t, err)
	pt := ptb.Finish()
	tx := sui_types.NewProgrammable(
		*sender, []*sui_types.ObjectRef{
			coin.Reference(),
		},
		pt, gasBudget, gasPrice,
	)
	txBytesBCS, err := bcs.Marshal(tx)
	require.NoError(t, err)

	// FIXME: Fail when there are two equal recipients
	resp := simulateCheck(t, cli, &types.TransactionBytes{
		TxBytes: txBytesBCS,
	}, true)
	gasFee := resp.Effects.Data.GasFee()
	t.Log(gasFee)

	// build with remote rpc
	// txn, err := cli.PaySui(context.Background(), *sender, []suiObjectID{coin.CoinObjectId},
	// 	[]suiAddress{*recipient, *recipient2},
	// 	[]types.SafeSuiBigInt[uint64]{types.NewSafeSuiBigInt(amount), types.NewSafeSuiBigInt(amount)},
	// 	types.NewSafeSuiBigInt(gasBudget))
	// require.NoError(t, err)
	// txBytesRemote := txn.TxBytes.Data()

	// XXX: Fail when there are multiple recipients
	// require.Equal(t, txBytesBCS, txBytesRemote)
}

func TestBCS_PayAllSui(t *testing.T) {
	sender, err := sui_types.NewAddressFromHex("0xd77955e670f42c1bc5e94b9e68e5fe9bdbed9134d784f2a14dfe5fc1b24b5d9f")
	require.NoError(t, err)
	recipient := sender
	gasBudget := SUI(0.01).Uint64()

	cli := TestnetClient(t)
	coins := GetCoins(t, cli, *sender, 2)
	coin, coin2 := coins[0], coins[1]

	gasPrice := uint64(1000)
	// gasPrice, err := cli.GetReferenceGasPrice(context.Background())

	// build with BCS
	ptb := sui_types.NewProgrammableTransactionBuilder()
	err = ptb.PayAllSui(*recipient)
	require.NoError(t, err)
	pt := ptb.Finish()
	tx := sui_types.NewProgrammable(
		*sender, []*sui_types.ObjectRef{
			coin.Reference(),
			coin2.Reference(),
		},
		pt, gasBudget, gasPrice,
	)
	txBytesBCS, err := bcs.Marshal(tx)
	require.NoError(t, err)

	// build with remote rpc
	txn, err := cli.PayAllSui(context.Background(), *sender, *recipient,
		[]suiObjectID{
			coin.CoinObjectId, coin2.CoinObjectId,
		},
		types.NewSafeSuiBigInt(gasBudget))
	require.NoError(t, err)
	txBytesRemote := txn.TxBytes.Data()

	require.Equal(t, txBytesBCS, txBytesRemote)
}

func TestBCS_Pay(t *testing.T) {
	sender, err := sui_types.NewAddressFromHex("0xd77955e670f42c1bc5e94b9e68e5fe9bdbed9134d784f2a14dfe5fc1b24b5d9f")
	require.NoError(t, err)
	recipient := sender
	recipient2, _ := sui_types.NewAddressFromHex("0x123456")
	amount := SUI(0.001).Uint64()
	gasBudget := SUI(0.01).Uint64()

	cli := TestnetClient(t)
	coins := GetCoins(t, cli, *sender, 2)
	coin, gas := coins[0], coins[1]

	gasPrice := uint64(1000)
	// gasPrice, err := cli.GetReferenceGasPrice(context.Background())

	// build with BCS
	ptb := sui_types.NewProgrammableTransactionBuilder()
	err = ptb.Pay(
		[]*sui_types.ObjectRef{coin.Reference()},
		[]suiAddress{*recipient, *recipient2},
		[]uint64{amount, amount})
	require.NoError(t, err)
	pt := ptb.Finish()
	tx := sui_types.NewProgrammable(
		*sender, []*sui_types.ObjectRef{
			gas.Reference(),
		},
		pt, gasBudget, gasPrice,
	)
	txBytesBCS, err := bcs.Marshal(tx)
	require.NoError(t, err)

	// FIXME: Fail when there are two equal recipients
	resp := simulateCheck(t, cli, &types.TransactionBytes{
		TxBytes: txBytesBCS,
	}, true)
	gasfee := resp.Effects.Data.GasFee()
	t.Log(gasfee)

	// build with remote rpc
	// txn, err := cli.Pay(context.Background(), *sender,
	// 	[]suiObjectID{coin.CoinObjectId},
	// 	[]suiAddress{*recipient, *recipient2},
	// 	[]types.SafeSuiBigInt[uint64]{types.NewSafeSuiBigInt(amount), types.NewSafeSuiBigInt(amount)},
	// 	&gas.CoinObjectId,
	// 	types.NewSafeSuiBigInt(gasBudget))
	// require.NoError(t, err)
	// txBytesRemote := txn.TxBytes.Data()

	// XXX: Fail when there are multiple recipients
	// require.Equal(t, txBytesBCS, txBytesRemote)
}

func GetCoins(t *testing.T, cli *Client, sender suiAddress, needCount int) []types.Coin {
	coins, err := cli.GetCoins(context.Background(), sender, nil, nil, uint(needCount))
	require.NoError(t, err)
	require.True(t, len(coins.Data) >= needCount)
	return coins.Data
}
