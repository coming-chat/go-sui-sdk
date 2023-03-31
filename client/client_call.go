package client

import (
	"context"
	"github.com/shopspring/decimal"
	"strings"

	"github.com/coming-chat/go-sui/types"
)

// GetBalance to use default sui coin(0x2::sui::SUI) when coinType is empty
func (c *Client) GetBalance(ctx context.Context, owner types.Address, coinType string) (*types.CoinBalance, error) {
	resp := types.CoinBalance{}
	if coinType == "" {
		return &resp, c.CallContext(ctx, &resp, getBalance, owner)
	} else {
		return &resp, c.CallContext(ctx, &resp, getBalance, owner, coinType)
	}
}

func (c *Client) GetAllBalances(ctx context.Context, owner types.Address) ([]types.CoinBalance, error) {
	var resp []types.CoinBalance
	return resp, c.CallContext(ctx, &resp, getAllBalances, owner)
}

// GetSuiCoinsOwnedByAddress This function will retrieve a maximum of 200 coins.
func (c *Client) GetSuiCoinsOwnedByAddress(ctx context.Context, address types.Address) (types.Coins, error) {
	coinType := types.SuiCoinType
	page, err := c.GetCoins(ctx, address, &coinType, nil, 200)
	if err != nil {
		return nil, err
	}
	return page.Data, nil
}

// GetCoins to use default sui coin(0x2::sui::SUI) when coinType is nil
// start with the first object when cursor is nil
func (c *Client) GetCoins(ctx context.Context, owner types.Address, coinType *string, cursor *types.ObjectId, limit uint) (*types.PaginatedCoins, error) {
	var resp types.PaginatedCoins
	return &resp, c.CallContext(ctx, &resp, getCoins, owner, coinType, cursor, limit)
}

// GetAllCoins
// start with the first object when cursor is nil
func (c *Client) GetAllCoins(ctx context.Context, owner types.Address, cursor *types.ObjectId, limit uint) (*types.PaginatedCoins, error) {
	var resp types.PaginatedCoins
	return &resp, c.CallContext(ctx, &resp, getAllCoins, owner, cursor, limit)
}

func (c *Client) GetCoinMetadata(ctx context.Context, coinType string) (*types.SuiCoinMetadata, error) {
	var resp types.SuiCoinMetadata
	return &resp, c.CallContext(ctx, &resp, getCoinMetadata, coinType)
}

func (c *Client) GetObject(ctx context.Context, objID types.ObjectId, options *types.SuiObjectDataOptions) (*types.SuiObjectResponse, error) {
	var resp types.SuiObjectResponse
	return &resp, c.CallContext(ctx, &resp, getObject, objID, options)
}

func (c *Client) MultiGetObjects(ctx context.Context, objIDs []types.ObjectId, options *types.SuiObjectDataOptions) ([]types.SuiObjectResponse, error) {
	var resp []types.SuiObjectResponse
	return resp, c.CallContext(ctx, &resp, multiGetObjects, objIDs, options)
}

// address : <SuiAddress> - the owner's Sui address
// query : <ObjectResponseQuery> - the objects query criteria.
// cursor : <CheckpointedObjectID> - An optional paging cursor. If provided, the query will start from the next item after the specified cursor. Default to start from the first item if not specified.
// limit : <uint> - Max number of items returned per page, default to [QUERY_MAX_RESULT_LIMIT_OBJECTS] if is 0
func (c *Client) GetOwnedObjects(ctx context.Context, address types.Address, query *types.SuiObjectResponseQuery, cursor *types.CheckpointedObjectId, limit uint) (*types.PaginatedObjectsResponse, error) {
	var resp types.PaginatedObjectsResponse
	if limit > 0 {
		return &resp, c.CallContext(ctx, &resp, getOwnedObjects, address, query, cursor, limit)
	} else {
		return &resp, c.CallContext(ctx, &resp, getOwnedObjects, address, query, cursor)
	}
}

func (c *Client) GetTotalSupply(ctx context.Context, coinType string) (*types.CoinSupply, error) {
	var resp types.CoinSupply
	return &resp, c.CallContext(ctx, &resp, getTotalSupply, coinType)
}

func (c *Client) GetTotalTransactionBlocks(ctx context.Context) (string, error) {
	var resp string
	return resp, c.CallContext(ctx, &resp, getTotalTransactionBlocks)
}

// BatchGetObjectsOwnedByAddress @param filterType You can specify filtering out the specified resources, this will fetch all resources if it is not empty ""
func (c *Client) BatchGetObjectsOwnedByAddress(ctx context.Context, address types.Address, options types.SuiObjectDataOptions, filterType string) ([]types.SuiObjectResponse, error) {
	filterType = strings.TrimSpace(filterType)
	return c.BatchGetFilteredObjectsOwnedByAddress(ctx, address, options, func(sod *types.SuiObjectData) bool {
		return filterType == "" || filterType == *sod.Type
	})
}

func (c *Client) BatchGetFilteredObjectsOwnedByAddress(ctx context.Context, address types.Address, options types.SuiObjectDataOptions, filter func(*types.SuiObjectData) bool) ([]types.SuiObjectResponse, error) {
	query := types.SuiObjectResponseQuery{
		Options: &types.SuiObjectDataOptions{
			ShowType: true,
		},
	}
	filteringObjs, err := c.GetOwnedObjects(ctx, address, &query, nil, 0)
	if err != nil {
		return nil, err
	}
	objIds := make([]types.ObjectId, 0)
	for _, obj := range filteringObjs.Data {
		if obj.Data == nil {
			continue // error obj
		}
		if filter != nil && filter(obj.Data) == false {
			continue // ignore objects if non-specified type
		}
		objIds = append(objIds, obj.Data.ObjectId)
	}

	return c.MultiGetObjects(ctx, objIds, &options)
}

func (c *Client) GetTransactionBlock(ctx context.Context, digest types.TransactionDigest, options types.SuiTransactionBlockResponseOptions) (*types.SuiTransactionBlockResponse, error) {
	resp := types.SuiTransactionBlockResponse{}
	return &resp, c.CallContext(ctx, &resp, getTransactionBlock, digest, options)
}

func (c *Client) GetReferenceGasPrice(ctx context.Context) (*decimal.Decimal, error) {
	var resp decimal.Decimal
	return &resp, c.CallContext(ctx, &resp, getReferenceGasPrice)
}

func (c *Client) DevInspectTransactionBlock(ctx context.Context, senderAddress types.Address, txByte types.Base64Data, gasPrice *decimal.Decimal, epoch *uint64) (*types.DevInspectResults, error) {
	var resp types.DevInspectResults
	return &resp, c.CallContext(ctx, &resp, devInspectTransactionBlock, senderAddress, txByte, gasPrice, epoch)
}

// MARK - Unmigrated

func (c *Client) BatchTransaction(ctx context.Context, signer types.Address, txnParams []map[string]interface{}, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, batchTransaction, signer, txnParams, gas, gasBudget)
}

//func (c *Client) DryRunTransaction(ctx context.Context, tx *types.TransactionBytes) (*types.TransactionEffects, error) {
//	resp := types.TransactionEffects{}
//	return &resp, c.CallContext(ctx, &resp, dryRunTransaction, tx.TxBytes)
//}

//func (c *Client) ExecuteTransaction(ctx context.Context, txn types.SignedTransactionSerializedSig, requestType types.ExecuteTransactionRequestType) (*types.ExecuteTransactionResponse, error) {
//	resp := types.ExecuteTransactionResponse{}
//	return &resp, c.CallContext(ctx, &resp, executeTransaction, txn.TxBytes, txn.Signature, requestType)
//}

//func (c *Client) BatchGetTransaction(digests []string) (map[string]*types.TransactionResponse, error) {
//	if len(digests) == 0 {
//		return map[string]*types.TransactionResponse{}, nil
//	}
//	var elems []BatchElem
//	results := make(map[string]*types.TransactionResponse)
//	for _, v := range digests {
//		results[v] = new(types.TransactionResponse)
//		elems = append(elems, BatchElem{
//			Method: getTransaction,
//			Args:   []interface{}{v},
//			Result: results[v],
//		})
//	}
//	return results, c.BatchCall(elems)
//}

// MergeCoins Create an unsigned transaction to merge multiple coins into one coin.
func (c *Client) MergeCoins(ctx context.Context, signer types.Address, primaryCoin, coinToMerge types.ObjectId, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, mergeCoins, signer, primaryCoin, coinToMerge, gas, gasBudget)
}

// MoveCall Create an unsigned transaction to execute a Move call on the network, by calling the specified function in the module of a given package.
// TODO: not support param `typeArguments` yet.
// So now only methods with `typeArguments` are supported
func (c *Client) MoveCall(ctx context.Context, signer types.Address, packageId types.ObjectId, module, function string, typeArgs []string, arguments []any, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, moveCall, signer, packageId, module, function, typeArgs, arguments, gas, gasBudget)
}

// SplitCoin Create an unsigned transaction to split a coin object into multiple coins.
func (c *Client) SplitCoin(ctx context.Context, signer types.Address, Coin types.ObjectId, splitAmounts []uint64, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, splitCoin, signer, Coin, splitAmounts, gas, gasBudget)
}

// SplitCoinEqual Create an unsigned transaction to split a coin object into multiple equal-size coins.
func (c *Client) SplitCoinEqual(ctx context.Context, signer types.Address, Coin types.ObjectId, splitCount uint64, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, splitCoinEqual, signer, Coin, splitCount, gas, gasBudget)
}

// TransferObject Create an unsigned transaction to transfer an object from one address to another. The object's type must allow public transfers
func (c *Client) TransferObject(ctx context.Context, signer, recipient types.Address, objID types.ObjectId, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, transferObject, signer, objID, gas, gasBudget, recipient)
}

// TransferSui Create an unsigned transaction to send SUI coin object to a Sui address. The SUI object is also used as the gas object.
func (c *Client) TransferSui(ctx context.Context, signer, recipient types.Address, suiObjID types.ObjectId, amount,
	gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, transferSui, signer, suiObjID, gasBudget, recipient, amount)
}

// PayAllSui Create an unsigned transaction to send all SUI coins to one recipient.
func (c *Client) PayAllSui(ctx context.Context, signer, recipient types.Address, inputCoins []types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, payAllSui, signer, inputCoins, recipient, gasBudget)

}

func (c *Client) Pay(ctx context.Context, signer types.Address, inputCoins []types.ObjectId, recipients []types.Address, amount []uint64, gas types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, pay, signer, inputCoins, recipients, amount, gas, gasBudget)
}

func (c *Client) PaySui(ctx context.Context, signer types.Address, inputCoins []types.ObjectId, recipients []types.Address, amount []uint64, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, paySui, signer, inputCoins, recipients, amount, gasBudget)
}

//func (c *Client) ExecuteTransactionSerializedSig(ctx context.Context, txn types.SignedTransactionSerializedSig, requestType types.ExecuteTransactionRequestType) (*types.ExecuteTransactionResponse, error) {
//	resp := types.ExecuteTransactionResponse{}
//	return &resp, c.CallContext(ctx, &resp, executeTransactionSerializedSig, txn.TxBytes, txn.Signature, requestType)
//}

func (c *Client) Publish(ctx context.Context, address types.Address, compiledModules []*types.Base64Data, gas types.ObjectId, gasBudget uint) (*types.TransactionBytes, error) {
	var resp types.TransactionBytes
	return &resp, c.CallContext(ctx, &resp, publish, address, compiledModules, gas, gasBudget)
}

func (c *Client) TryGetPastObject(ctx context.Context, objectId types.ObjectId, version uint64) (*types.ObjectRead, error) {
	resp := types.ObjectRead{}
	return &resp, c.CallContext(ctx, &resp, tryGetPastObject, objectId, version)
}

func (c *Client) GetEvents(ctx context.Context, eventQuery types.EventQuery, cursor *types.EventID, limit uint, descendingOrder bool) (*types.EventPage, error) {
	resp := types.EventPage{}
	return &resp, c.CallContext(ctx, &resp, getEvents, eventQuery, cursor, limit, descendingOrder)
}
