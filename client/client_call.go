package client

import (
	"context"
	"strings"

	"github.com/coming-chat/go-sui/types"
)

// GetBalance to use default sui coin(0x2::sui::SUI) when coinType is empty
func (c *Client) GetBalance(ctx context.Context, owner types.Address, coinType string) (*types.CoinBalance, error) {
	resp := types.CoinBalance{}
	if coinType == "" {
		return &resp, c.CallContext(ctx, &resp, "suix_getBalance", owner)
	} else {
		return &resp, c.CallContext(ctx, &resp, "suix_getBalance", owner, coinType)
	}
}

func (c *Client) GetAllBalances(ctx context.Context, owner types.Address) ([]types.CoinBalance, error) {
	var resp []types.CoinBalance
	return resp, c.CallContext(ctx, &resp, "suix_getAllBalances", owner)
}

func (c *Client) GetSuiCoinsOwnedByAddress(ctx context.Context, address types.Address) (types.Coins, error) {
	cointype := types.SuiCoinType
	page, err := c.GetCoins(ctx, address, &cointype, nil, 200)
	if err != nil {
		return nil, err
	}
	return page.Data, nil
}

// GetCoins to use default sui coin(0x2::sui::SUI) when coinType is nil
// start with the first object when cursor is nil
func (c *Client) GetCoins(ctx context.Context, owner types.Address, coinType *string, cursor *types.ObjectId, limit uint) (*types.PaginatedCoins, error) {
	var resp types.PaginatedCoins
	return &resp, c.CallContext(ctx, &resp, "suix_getCoins", owner, coinType, cursor, limit)
}

// GetAllCoins
// start with the first object when cursor is nil
func (c *Client) GetAllCoins(ctx context.Context, owner types.Address, cursor *types.ObjectId, limit uint) (*types.PaginatedCoins, error) {
	var resp types.PaginatedCoins
	return &resp, c.CallContext(ctx, &resp, "suix_getAllCoins", owner, cursor, limit)
}

// GetAllCoins
// start with the first object when cursor is nil
func (c *Client) GetCoinMetadata(ctx context.Context, coinType string) (*types.SuiCoinMetadata, error) {
	var resp types.SuiCoinMetadata
	return &resp, c.CallContext(ctx, &resp, "suix_getCoinMetadata", coinType)
}

func (c *Client) GetObject(ctx context.Context, objID types.ObjectId, options *types.SuiObjectDataOptions) (*types.SuiObjectResponse, error) {
	var resp types.SuiObjectResponse
	return &resp, c.CallContext(ctx, &resp, "sui_getObject", objID, options)
}

func (c *Client) MultiGetObjects(ctx context.Context, objIDs []types.ObjectId, options *types.SuiObjectDataOptions) ([]types.SuiObjectResponse, error) {
	var resp []types.SuiObjectResponse
	return resp, c.CallContext(ctx, &resp, "sui_multiGetObjects", objIDs, options)
}

// address : <SuiAddress> - the owner's Sui address
// query : <ObjectResponseQuery> - the objects query criteria.
// cursor : <CheckpointedObjectID> - An optional paging cursor. If provided, the query will start from the next item after the specified cursor. Default to start from the first item if not specified.
// limit : <uint> - Max number of items returned per page, default to [QUERY_MAX_RESULT_LIMIT_OBJECTS] if is 0
func (c *Client) GetOwnedObjects(ctx context.Context, address types.Address, query *types.SuiObjectResponseQuery, cursor *types.CheckpointedObjectId, limit uint) (*types.PaginatedObjectsResponse, error) {
	var resp types.PaginatedObjectsResponse
	if limit > 0 {
		return &resp, c.CallContext(ctx, &resp, "suix_getOwnedObjects", address, query, cursor, limit)
	} else {
		return &resp, c.CallContext(ctx, &resp, "suix_getOwnedObjects", address, query, cursor)
	}
}

func (c *Client) GetTotalSupply(ctx context.Context, coinType string) (*types.CoinSupply, error) {
	var resp types.CoinSupply
	return &resp, c.CallContext(ctx, &resp, "suix_getTotalSupply", coinType)
}

func (c *Client) GetTotalTransactionBlocks(ctx context.Context) (string, error) {
	var resp string
	return resp, c.CallContext(ctx, &resp, "sui_getTotalTransactionBlocks")
}

// MARK - Unmigrated

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
	return &resp, c.CallContext(ctx, &resp, "sui_batchTransaction", signer, txnParams, gas, gasBudget)
}

func (c *Client) DryRunTransaction(ctx context.Context, tx *types.TransactionBytes) (*types.TransactionEffects, error) {
	resp := types.TransactionEffects{}
	return &resp, c.CallContext(ctx, &resp, "sui_dryRunTransaction", tx.TxBytes)
}

func (c *Client) ExecuteTransaction(ctx context.Context, txn types.SignedTransactionSerializedSig, requestType types.ExecuteTransactionRequestType) (*types.ExecuteTransactionResponse, error) {
	resp := types.ExecuteTransactionResponse{}
	return &resp, c.CallContext(ctx, &resp, "sui_executeTransaction", txn.TxBytes, txn.Signature, requestType)
}

func (c *Client) GetObjectsOwnedByAddress(ctx context.Context, address types.Address) ([]types.ObjectInfo, error) {
	var resp []types.ObjectInfo
	return resp, c.CallContext(ctx, &resp, "suix_getOwnedObjects", address)
}

func (c *Client) GetObjectsOwnedByObject(ctx context.Context, objID types.ObjectId) ([]types.ObjectInfo, error) {
	var resp []types.ObjectInfo
	return resp, c.CallContext(ctx, &resp, "sui_getObjectsOwnedByObject", objID)
}

func (c *Client) GetRawObject(ctx context.Context, objID types.ObjectId) (*types.ObjectRead, error) {
	resp := types.ObjectRead{}
	return &resp, c.CallContext(ctx, &resp, "sui_getRawObject", objID)
}

func (c *Client) GetTotalTransactionNumber(ctx context.Context) (uint64, error) {
	resp := uint64(0)
	return resp, c.CallContext(ctx, &resp, "sui_getTotalTransactionNumber")
}

func (c *Client) GetTransactionsInRange(ctx context.Context, start, end uint64) ([]string, error) {
	var resp []string
	return resp, c.CallContext(ctx, &resp, "sui_getTransactionsInRange", start, end)
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
	return results, c.BatchCall(elems)
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
	return results, c.BatchCall(elems)
}

func (c *Client) GetTransaction(ctx context.Context, digest string) (*types.TransactionResponse, error) {
	resp := types.TransactionResponse{}
	return &resp, c.CallContext(ctx, &resp, "sui_getTransaction", digest)
}

// MergeCoins Create an unsigned transaction to merge multiple coins into one coin.
func (c *Client) MergeCoins(ctx context.Context, signer types.Address, primaryCoin, coinToMerge types.ObjectId, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, "sui_mergeCoins", signer, primaryCoin, coinToMerge, gas, gasBudget)
}

// MoveCall Create an unsigned transaction to execute a Move call on the network, by calling the specified function in the module of a given package.
// TODO: not support param `typeArguments` yet.
// So now only methods with `typeArguments` are supported
func (c *Client) MoveCall(ctx context.Context, signer types.Address, packageId types.ObjectId, module, function string, typeArgs []string, arguments []any, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, "sui_moveCall", signer, packageId, module, function, typeArgs, arguments, gas, gasBudget)
}

// SplitCoin Create an unsigned transaction to split a coin object into multiple coins.
func (c *Client) SplitCoin(ctx context.Context, signer types.Address, Coin types.ObjectId, splitAmounts []uint64, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, "sui_splitCoin", signer, Coin, splitAmounts, gas, gasBudget)
}

// SplitCoinEqual Create an unsigned transaction to split a coin object into multiple equal-size coins.
func (c *Client) SplitCoinEqual(ctx context.Context, signer types.Address, Coin types.ObjectId, splitCount uint64, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, "sui_splitCoinEqual", signer, Coin, splitCount, gas, gasBudget)
}

// TransferObject Create an unsigned transaction to transfer an object from one address to another. The object's type must allow public transfers
func (c *Client) TransferObject(ctx context.Context, signer, recipient types.Address, objID types.ObjectId, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, "sui_transferObject", signer, objID, gas, gasBudget, recipient)
}

// TransferSui Create an unsigned transaction to send SUI coin object to a Sui address. The SUI object is also used as the gas object.
func (c *Client) TransferSui(ctx context.Context, signer, recipient types.Address, suiObjID types.ObjectId, amount, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, "sui_transferSui", signer, suiObjID, gasBudget, recipient, amount)
}

// PayAllSui Create an unsigned transaction to send all SUI coins to one recipient.
func (c *Client) PayAllSui(ctx context.Context, signer, recipient types.Address, inputCoins []types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, "sui_payAllSui", signer, inputCoins, recipient, gasBudget)

}

func (c *Client) Pay(ctx context.Context, signer types.Address, inputCoins []types.ObjectId, recipients []types.Address, amount []uint64, gas types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, "sui_pay", signer, inputCoins, recipients, amount, gas, gasBudget)
}

func (c *Client) PaySui(ctx context.Context, signer types.Address, inputCoins []types.ObjectId, recipients []types.Address, amount []uint64, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, "sui_paySui", signer, inputCoins, recipients, amount, gasBudget)
}

func (c *Client) DevInspectTransaction(ctx context.Context, senderAddress types.Address, txByte types.Base64Data, gasPrice *uint64, epoch *uint64) (*types.DevInspectResults, error) {
	var resp types.DevInspectResults
	return &resp, c.CallContext(ctx, &resp, "sui_devInspectTransaction", senderAddress, txByte, gasPrice, epoch)
}

func (c *Client) ExecuteTransactionSerializedSig(ctx context.Context, txn types.SignedTransactionSerializedSig, requestType types.ExecuteTransactionRequestType) (*types.ExecuteTransactionResponse, error) {
	resp := types.ExecuteTransactionResponse{}
	return &resp, c.CallContext(ctx, &resp, "sui_executeTransactionSerializedSig", txn.TxBytes, txn.Signature, requestType)
}

func (c *Client) Publish(ctx context.Context, address types.Address, compiledModules []*types.Base64Data, gas types.ObjectId, gasBudget uint) (*types.TransactionBytes, error) {
	var resp types.TransactionBytes
	return &resp, c.CallContext(ctx, &resp, "sui_publish", address, compiledModules, gas, gasBudget)
}

func (c *Client) GetTransactions(ctx context.Context, transactionQuery types.TransactionQuery, cursor *string, limit uint, descendingOrder bool) (*types.TransactionsPage, error) {
	var resp types.TransactionsPage
	return &resp, c.CallContext(ctx, &resp, "sui_getTransactions", transactionQuery, cursor, limit, descendingOrder)
}

func (c *Client) TryGetPastObject(ctx context.Context, objectId types.ObjectId, version uint64) (*types.ObjectRead, error) {
	resp := types.ObjectRead{}
	return &resp, c.CallContext(ctx, &resp, "sui_tryGetPastObject", objectId, version)
}

func (c *Client) GetEvents(ctx context.Context, eventQuery types.EventQuery, cursor *types.EventID, limit uint, descendingOrder bool) (*types.EventPage, error) {
	resp := types.EventPage{}
	return &resp, c.CallContext(ctx, &resp, "sui_getEvents", eventQuery, cursor, limit, descendingOrder)
}

func (c *Client) GetReferenceGasPrice(ctx context.Context) (uint64, error) {
	resp := uint64(0)
	return resp, c.CallContext(ctx, &resp, "sui_getReferenceGasPrice")
}

/*
TODO
sui_getCommitteeInfo
sui_getDynamicFieldObject
sui_getDynamicFields
sui_getMoveFunctionArgTypes
sui_getNormalizedMoveFunction
sui_getNormalizedMoveModule
sui_getNormalizedMoveModulesByPackage
sui_getNormalizedMoveStruct
sui_getTransactionAuthSigners
sui_subscribeEvent
*/
