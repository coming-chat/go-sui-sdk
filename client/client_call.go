package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/coming-chat/go-sui/types"
)

func (c *Client) GetSuiCoinsOwnedByAddress(ctx context.Context, address types.Address) (types.Coins, error) {
	return c.GetCoinsOwnedByAddress(ctx, address, types.SuiCoinType)
}

func (c *Client) GetCoinsOwnedByAddress(ctx context.Context, address types.Address, coinType string) (coins types.Coins, err error) {
	defer func() {
		errP := recover()
		if errP != nil {
			err = fmt.Errorf("decode rpc data err: %v", errP)
			return
		}
	}()
	//TODO Get only 200 items for the time being, and then add parameters to get more.
	coinObjects, err := c.GetCoins(ctx, address, &coinType, nil, 200)
	if err != nil {
		return nil, err
	}
	for _, coin := range coinObjects.Data {
		coins = append(coins, types.Coin{
			Balance: uint64(coin.Balance),
			Type:    coin.CoinType,
			Owner:   &address,
			Reference: &types.ObjectRef{
				Digest:   coin.Digest,
				Version:  coin.Version,
				ObjectId: coin.CoinObjectId,
			},
		})
	}
	return coins, nil
}

// BatchGetObjectsOwnedByAddress @param filterType You can specify filtering out the specified resources, this will fetch all resources if it is not empty ""
func (c *Client) BatchGetObjectsOwnedByAddress(ctx context.Context, address types.Address, filterType string) ([]types.ObjectRead, error) {
	filterType = strings.TrimSpace(filterType)
	return c.BatchGetFilteredObjectsOwnedByAddress(ctx, address, func(oi types.ObjectInfo) bool {
		return filterType == "" || filterType == oi.Type
	})
}

func (c *Client) BatchGetFilteredObjectsOwnedByAddress(ctx context.Context, address types.Address, filter func(types.ObjectInfo) bool) ([]types.ObjectRead, error) {
	infos, err := c.GetObjectsOwnedByAddress(ctx, address)
	if err != nil {
		return nil, err
	}
	var elems []BatchElem
	for _, info := range infos {
		if filter != nil && filter(info) == false {
			// ignore objects if non-specified type
			continue
		}
		elems = append(elems, BatchElem{
			Method: "sui_getObject",
			Args:   []interface{}{info.ObjectId},
			Result: &types.ObjectRead{},
		})
	}
	if len(elems) == 0 {
		return []types.ObjectRead{}, nil
	}
	err = c.BatchCallContext(ctx, elems)
	if err != nil {
		return nil, err
	}
	objects := make([]types.ObjectRead, len(elems))
	for i, ele := range elems {
		if ele.Error != nil {
			return nil, ele.Error
		}
		objects[i] = *ele.Result.(*types.ObjectRead)
	}
	return objects, nil
}

func (c *Client) BatchTransaction(ctx context.Context, signer types.Address, txnParams []map[string]interface{}, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	err := c.CallContext(ctx, &resp, "sui_batchTransaction", signer, txnParams, gas, gasBudget)
	return &resp, err
}

func (c *Client) DryRunTransaction(ctx context.Context, tx *types.TransactionBytes) (*types.TransactionEffects, error) {
	resp := types.TransactionEffects{}
	err := c.CallContext(ctx, &resp, "sui_dryRunTransaction", tx.TxBytes)
	return &resp, err
}

func (c *Client) ExecuteTransaction(ctx context.Context, txn types.SignedTransaction, requestType types.ExecuteTransactionRequestType) (*types.ExecuteTransactionResponse, error) {
	resp := types.ExecuteTransactionResponse{}
	err := c.CallContext(ctx, &resp, "sui_executeTransaction", txn.TxBytes, txn.SigScheme, txn.Signature, txn.PublicKey, requestType)
	return &resp, err
}

func (c *Client) GetObject(ctx context.Context, objID types.ObjectId) (*types.ObjectRead, error) {
	resp := types.ObjectRead{}
	err := c.CallContext(ctx, &resp, "sui_getObject", objID)
	return &resp, err
}

func (c *Client) GetObjectsOwnedByAddress(ctx context.Context, address types.Address) ([]types.ObjectInfo, error) {
	var resp []types.ObjectInfo
	err := c.CallContext(ctx, &resp, "sui_getObjectsOwnedByAddress", address)
	return resp, err
}

func (c *Client) GetObjectsOwnedByObject(ctx context.Context, objID types.ObjectId) ([]types.ObjectInfo, error) {
	var resp []types.ObjectInfo
	err := c.CallContext(ctx, &resp, "sui_getObjectsOwnedByObject", objID)
	return resp, err
}

func (c *Client) GetRawObject(ctx context.Context, objID types.ObjectId) (*types.ObjectRead, error) {
	resp := types.ObjectRead{}
	err := c.CallContext(ctx, &resp, "sui_getRawObject", objID)
	return &resp, err
}

func (c *Client) GetTotalTransactionNumber(ctx context.Context) (uint64, error) {
	resp := uint64(0)
	err := c.CallContext(ctx, &resp, "sui_getTotalTransactionNumber")
	return resp, err
}

func (c *Client) GetTransactionsInRange(ctx context.Context, start, end uint64) ([]string, error) {
	var resp []string
	err := c.CallContext(ctx, &resp, "sui_getTransactionsInRange", start, end)
	return resp, err
}

func (c *Client) BatchGetTransaction(digests []string) (map[string]*types.TransactionResponse, error) {
	if len(digests) == 0 {
		return map[string]*types.TransactionResponse{}, nil
	}
	var elems []BatchElem
	results := make(map[string]*types.TransactionResponse)
	for _, v := range digests {
		results[v] = new(types.TransactionResponse)
		elems = append(elems, BatchElem{
			Method: "sui_getTransaction",
			Args:   []interface{}{v},
			Result: results[v],
		})
	}
	if err := c.BatchCall(elems); err != nil {
		return nil, err
	}
	return results, nil
}

func (c *Client) BatchGetObject(objects []types.ObjectId) (map[string]*types.ObjectRead, error) {
	if len(objects) == 0 {
		return map[string]*types.ObjectRead{}, nil
	}
	var elems []BatchElem
	results := make(map[string]*types.ObjectRead)
	for _, v := range objects {
		results[v.String()] = new(types.ObjectRead)
		elems = append(elems, BatchElem{
			Method: "sui_getObject",
			Args:   []interface{}{v},
			Result: results[v.String()],
		})
	}
	if err := c.BatchCall(elems); err != nil {
		return nil, err
	}
	return results, nil
}

func (c *Client) GetTransaction(ctx context.Context, digest string) (*types.TransactionResponse, error) {
	resp := types.TransactionResponse{}
	err := c.CallContext(ctx, &resp, "sui_getTransaction", digest)
	return &resp, err
}

// MergeCoins Create an unsigned transaction to merge multiple coins into one coin.
func (c *Client) MergeCoins(ctx context.Context, signer types.Address, primaryCoin, coinToMerge types.ObjectId, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	err := c.CallContext(ctx, &resp, "sui_mergeCoins", signer, primaryCoin, coinToMerge, gas, gasBudget)
	return &resp, err
}

// MoveCall Create an unsigned transaction to execute a Move call on the network, by calling the specified function in the module of a given package.
// TODO: not support param `typeArguments` yet.
// So now only methods with `typeArguments` are supported
func (c *Client) MoveCall(ctx context.Context, signer types.Address, packageId types.ObjectId, module, function string, typeArgs []string, arguments []any, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	err := c.CallContext(ctx, &resp, "sui_moveCall", signer, packageId, module, function, typeArgs, arguments, gas, gasBudget)
	return &resp, err
}

// SplitCoin Create an unsigned transaction to split a coin object into multiple coins.
func (c *Client) SplitCoin(ctx context.Context, signer types.Address, Coin types.ObjectId, splitAmounts []uint64, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	err := c.CallContext(ctx, &resp, "sui_splitCoin", signer, Coin, splitAmounts, gas, gasBudget)
	return &resp, err
}

// SplitCoinEqual Create an unsigned transaction to split a coin object into multiple equal-size coins.
func (c *Client) SplitCoinEqual(ctx context.Context, signer types.Address, Coin types.ObjectId, splitCount uint64, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	err := c.CallContext(ctx, &resp, "sui_splitCoinEqual", signer, Coin, splitCount, gas, gasBudget)
	return &resp, err
}

// TransferObject Create an unsigned transaction to transfer an object from one address to another. The object's type must allow public transfers
func (c *Client) TransferObject(ctx context.Context, signer, recipient types.Address, objID types.ObjectId, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	err := c.CallContext(ctx, &resp, "sui_transferObject", signer, objID, gas, gasBudget, recipient)
	return &resp, err
}

// TransferSui Create an unsigned transaction to send SUI coin object to a Sui address. The SUI object is also used as the gas object.
func (c *Client) TransferSui(ctx context.Context, signer, recipient types.Address, suiObjID types.ObjectId, amount, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	err := c.CallContext(ctx, &resp, "sui_transferSui", signer, suiObjID, gasBudget, recipient, amount)
	return &resp, err
}

// PayAllSui Create an unsigned transaction to send all SUI coins to one recipient.
func (c *Client) PayAllSui(ctx context.Context, signer, recipient types.Address, inputCoins []types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	err := c.CallContext(ctx, &resp, "sui_payAllSui", signer, inputCoins, recipient, gasBudget)
	return &resp, err
}

func (c *Client) GetCoinMetadata(ctx context.Context, coinType string) (*types.SuiCoinMetadata, error) {
	resp := types.SuiCoinMetadata{}
	err := c.CallContext(ctx, &resp, "sui_getCoinMetadata", coinType)
	return &resp, err
}

func (c *Client) Pay(ctx context.Context, signer types.Address, inputCoins []types.ObjectId, recipients []types.Address, amount []uint64, gas types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	err := c.CallContext(ctx, &resp, "sui_pay", signer, inputCoins, recipients, amount, gas, gasBudget)
	return &resp, err
}

func (c *Client) PaySui(ctx context.Context, signer types.Address, inputCoins []types.ObjectId, recipients []types.Address, amount []uint64, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	err := c.CallContext(ctx, &resp, "sui_paySui", signer, inputCoins, recipients, amount, gasBudget)
	return &resp, err
}

func (c *Client) GetAllBalances(ctx context.Context, address types.Address) ([]types.SuiCoinBalance, error) {
	var resp []types.SuiCoinBalance
	err := c.CallContext(ctx, &resp, "sui_getAllBalances", address)
	return resp, err
}

// GetBalance to use default sui coin(0x2::sui::SUI) when coinType is nil
func (c *Client) GetBalance(ctx context.Context, address types.Address, coinType *string) (*types.SuiCoinBalance, error) {
	resp := types.SuiCoinBalance{}
	return &resp, c.CallContext(ctx, &resp, "sui_getBalance", address, coinType)
}

func (c *Client) DevInspectMoveCall(ctx context.Context, signer types.Address, packageId types.ObjectId, module, function string, typeArgs []string, arguments []any) (*types.DevInspectResults, error) {
	var resp types.DevInspectResults
	err := c.CallContext(ctx, &resp, "sui_devInspectMoveCall", signer, packageId, module, function, typeArgs, arguments)
	return &resp, err
}

func (c *Client) DevInspectTransaction(ctx context.Context, txByte types.Base64Data) (*types.DevInspectResults, error) {
	var resp types.DevInspectResults
	err := c.CallContext(ctx, &resp, "sui_devInspectTransaction", txByte)
	return &resp, err
}

// GetCoins to use default sui coin(0x2::sui::SUI) when coinType is nil
// start with the first object when cursor is nil
func (c *Client) GetCoins(ctx context.Context, address types.Address, coinType *string, cursor *types.ObjectId, limit uint) (*types.CoinPage, error) {
	var resp types.CoinPage
	return &resp, c.CallContext(ctx, &resp, "sui_getCoins", address, coinType, cursor, limit)
}

// GetAllCoins
// start with the first object when cursor is nil
func (c *Client) GetAllCoins(ctx context.Context, address types.Address, cursor *types.ObjectId, limit uint) (*types.CoinPage, error) {
	var resp types.CoinPage
	return &resp, c.CallContext(ctx, &resp, "sui_getAllCoins", address, cursor, limit)
}

func (c *Client) GetTotalSupply(ctx context.Context, coinType string) (*types.Supply, error) {
	var resp types.Supply
	return &resp, c.CallContext(ctx, &resp, "sui_getTotalSupply", coinType)
}

func (c *Client) ExecuteTransactionSerializedSig(ctx context.Context, txn types.SignedTransactionSerializedSig, requestType types.ExecuteTransactionRequestType) (*types.ExecuteTransactionResponse, error) {
	resp := types.ExecuteTransactionResponse{}
	err := c.CallContext(ctx, &resp, "sui_executeTransactionSerializedSig", txn.TxBytes, txn.Signature, requestType)
	return &resp, err
}

/*
TODO
sui_getCommitteeInfo
sui_getDynamicFieldObject
sui_getDynamicFields
sui_getEvents
sui_getMoveFunctionArgTypes
sui_getNormalizedMoveFunction
sui_getNormalizedMoveModule
sui_getNormalizedMoveModulesByPackage
sui_getNormalizedMoveStruct
sui_getSuiSystemState
sui_getTransactionAuthSigners
sui_getTransactions
sui_subscribeEvent
sui_tryGetPastObject
sui_publish
*/
