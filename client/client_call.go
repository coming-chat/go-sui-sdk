package client

import (
	"context"
	"strings"

	"github.com/coming-chat/go-sui/v2/lib"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
)

// NOTE: This copys the query limit from our Rust JSON RPC backend, this needs to be kept in sync!
const QUERY_MAX_RESULT_LIMIT = 1000

type suiAddress = sui_types.SuiAddress
type suiObjectID = sui_types.ObjectID
type suiDigest = sui_types.TransactionDigest
type suiBase64Data = lib.Base64Data

// MARK - Getter Function

// GetBalance to use default sui coin(0x2::sui::SUI) when coinType is empty
func (c *Client) GetBalance(ctx context.Context, owner suiAddress, coinType string) (*types.Balance, error) {
	resp := types.Balance{}
	if coinType == "" {
		return &resp, c.CallContext(ctx, &resp, getBalance, owner)
	} else {
		return &resp, c.CallContext(ctx, &resp, getBalance, owner, coinType)
	}
}

func (c *Client) GetAllBalances(ctx context.Context, owner suiAddress) ([]types.Balance, error) {
	var resp []types.Balance
	return resp, c.CallContext(ctx, &resp, getAllBalances, owner)
}

// GetSuiCoinsOwnedByAddress This function will retrieve a maximum of 200 coins.
func (c *Client) GetSuiCoinsOwnedByAddress(ctx context.Context, address suiAddress) (types.Coins, error) {
	coinType := types.SuiCoinType
	page, err := c.GetCoins(ctx, address, &coinType, nil, 200)
	if err != nil {
		return nil, err
	}
	return page.Data, nil
}

// GetCoins to use default sui coin(0x2::sui::SUI) when coinType is nil
// start with the first object when cursor is nil
func (c *Client) GetCoins(
	ctx context.Context,
	owner suiAddress,
	coinType *string,
	cursor *suiObjectID,
	limit uint,
) (*types.CoinPage, error) {
	var resp types.CoinPage
	return &resp, c.CallContext(ctx, &resp, getCoins, owner, coinType, cursor, limit)
}

// GetAllCoins
// start with the first object when cursor is nil
func (c *Client) GetAllCoins(
	ctx context.Context,
	owner suiAddress,
	cursor *suiObjectID,
	limit uint,
) (*types.CoinPage, error) {
	var resp types.CoinPage
	return &resp, c.CallContext(ctx, &resp, getAllCoins, owner, cursor, limit)
}

func (c *Client) GetCoinMetadata(ctx context.Context, coinType string) (*types.SuiCoinMetadata, error) {
	var resp types.SuiCoinMetadata
	return &resp, c.CallContext(ctx, &resp, getCoinMetadata, coinType)
}

func (c *Client) GetObject(
	ctx context.Context,
	objID suiObjectID,
	options *types.SuiObjectDataOptions,
) (*types.SuiObjectResponse, error) {
	var resp types.SuiObjectResponse
	return &resp, c.CallContext(ctx, &resp, getObject, objID, options)
}

func (c *Client) MultiGetObjects(
	ctx context.Context,
	objIDs []suiObjectID,
	options *types.SuiObjectDataOptions,
) ([]types.SuiObjectResponse, error) {
	var resp []types.SuiObjectResponse
	return resp, c.CallContext(ctx, &resp, multiGetObjects, objIDs, options)
}

// address : <SuiAddress> - the owner's Sui address
// query : <ObjectResponseQuery> - the objects query criteria.
// cursor : <CheckpointedObjectID> - An optional paging cursor. If provided, the query will start from the next item after the specified cursor. Default to start from the first item if not specified.
// limit : <uint> - Max number of items returned per page, default to [QUERY_MAX_RESULT_LIMIT_OBJECTS] if is 0
func (c *Client) GetOwnedObjects(
	ctx context.Context,
	address suiAddress,
	query *types.SuiObjectResponseQuery,
	cursor *types.CheckpointedObjectId,
	limit *uint,
) (*types.ObjectsPage, error) {
	var resp types.ObjectsPage
	return &resp, c.CallContext(ctx, &resp, getOwnedObjects, address, query, cursor, limit)
}

func (c *Client) GetTotalSupply(ctx context.Context, coinType string) (*types.Supply, error) {
	var resp types.Supply
	return &resp, c.CallContext(ctx, &resp, getTotalSupply, coinType)
}

func (c *Client) GetTotalTransactionBlocks(ctx context.Context) (string, error) {
	var resp string
	return resp, c.CallContext(ctx, &resp, getTotalTransactionBlocks)
}

func (c *Client) GetLatestCheckpointSequenceNumber(ctx context.Context) (string, error) {
	var resp string
	return resp, c.CallContext(ctx, &resp, getLatestCheckpointSequenceNumber)
}

// BatchGetObjectsOwnedByAddress @param filterType You can specify filtering out the specified resources, this will fetch all resources if it is not empty ""
func (c *Client) BatchGetObjectsOwnedByAddress(
	ctx context.Context,
	address suiAddress,
	options types.SuiObjectDataOptions,
	filterType string,
) ([]types.SuiObjectResponse, error) {
	filterType = strings.TrimSpace(filterType)
	return c.BatchGetFilteredObjectsOwnedByAddress(
		ctx, address, options, func(sod *types.SuiObjectData) bool {
			return filterType == "" || filterType == *sod.Type
		},
	)
}

func (c *Client) BatchGetFilteredObjectsOwnedByAddress(
	ctx context.Context,
	address suiAddress,
	options types.SuiObjectDataOptions,
	filter func(*types.SuiObjectData) bool,
) ([]types.SuiObjectResponse, error) {
	query := types.SuiObjectResponseQuery{
		Options: &types.SuiObjectDataOptions{
			ShowType: true,
		},
	}
	filteringObjs, err := c.GetOwnedObjects(ctx, address, &query, nil, nil)
	if err != nil {
		return nil, err
	}
	objIds := make([]suiObjectID, 0)
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

func (c *Client) GetTransactionBlock(
	ctx context.Context,
	digest suiDigest,
	options types.SuiTransactionBlockResponseOptions,
) (*types.SuiTransactionBlockResponse, error) {
	resp := types.SuiTransactionBlockResponse{}
	return &resp, c.CallContext(ctx, &resp, getTransactionBlock, digest, options)
}

func (c *Client) GetReferenceGasPrice(ctx context.Context) (*types.SafeSuiBigInt[uint64], error) {
	var resp types.SafeSuiBigInt[uint64]
	return &resp, c.CallContext(ctx, &resp, getReferenceGasPrice)
}

func (c *Client) GetEvents(ctx context.Context, digest suiDigest) ([]types.SuiEvent, error) {
	var resp []types.SuiEvent
	return resp, c.CallContext(ctx, &resp, getEvents, digest)
}

func (c *Client) TryGetPastObject(
	ctx context.Context,
	objectId suiObjectID,
	version uint64,
	options *types.SuiObjectDataOptions,
) (*types.SuiPastObjectResponse, error) {
	var resp types.SuiPastObjectResponse
	return &resp, c.CallContext(ctx, &resp, tryGetPastObject, objectId, version, options)
}

func (c *Client) DevInspectTransactionBlock(
	ctx context.Context,
	senderAddress suiAddress,
	txByte suiBase64Data,
	gasPrice *types.SafeSuiBigInt[uint64],
	epoch *uint64,
) (*types.DevInspectResults, error) {
	var resp types.DevInspectResults
	return &resp, c.CallContext(ctx, &resp, devInspectTransactionBlock, senderAddress, txByte, gasPrice, epoch)
}

func (c *Client) DryRunTransaction(
	ctx context.Context,
	txBytes suiBase64Data,
) (*types.DryRunTransactionBlockResponse, error) {
	var resp types.DryRunTransactionBlockResponse
	return &resp, c.CallContext(ctx, &resp, dryRunTransactionBlock, txBytes)
}

func (c *Client) ExecuteTransactionBlock(
	ctx context.Context, txBytes suiBase64Data, signatures []any,
	options *types.SuiTransactionBlockResponseOptions, requestType types.ExecuteTransactionRequestType,
) (*types.SuiTransactionBlockResponse, error) {
	resp := types.SuiTransactionBlockResponse{}
	return &resp, c.CallContext(ctx, &resp, executeTransactionBlock, txBytes, signatures, options, requestType)
}

// TransferObject Create an unsigned transaction to transfer an object from one address to another. The object's type must allow public transfers
func (c *Client) TransferObject(
	ctx context.Context,
	signer, recipient suiAddress,
	objID suiObjectID,
	gas *suiObjectID,
	gasBudget types.SafeSuiBigInt[uint64],
) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, transferObject, signer, objID, gas, gasBudget, recipient)
}

// TransferSui Create an unsigned transaction to send SUI coin object to a Sui address. The SUI object is also used as the gas object.
func (c *Client) TransferSui(
	ctx context.Context, signer, recipient suiAddress, suiObjID suiObjectID, amount,
	gasBudget types.SafeSuiBigInt[uint64],
) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, transferSui, signer, suiObjID, gasBudget, recipient, amount)
}

// PayAllSui Create an unsigned transaction to send all SUI coins to one recipient.
func (c *Client) PayAllSui(
	ctx context.Context,
	signer, recipient suiAddress,
	inputCoins []suiObjectID,
	gasBudget types.SafeSuiBigInt[uint64],
) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, payAllSui, signer, inputCoins, recipient, gasBudget)
}

func (c *Client) Pay(
	ctx context.Context,
	signer suiAddress,
	inputCoins []suiObjectID,
	recipients []suiAddress,
	amount []types.SafeSuiBigInt[uint64],
	gas *suiObjectID,
	gasBudget types.SafeSuiBigInt[uint64],
) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, pay, signer, inputCoins, recipients, amount, gas, gasBudget)
}

func (c *Client) PaySui(
	ctx context.Context,
	signer suiAddress,
	inputCoins []suiObjectID,
	recipients []suiAddress,
	amount []types.SafeSuiBigInt[uint64],
	gasBudget types.SafeSuiBigInt[uint64],
) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, paySui, signer, inputCoins, recipients, amount, gasBudget)
}

// SplitCoin Create an unsigned transaction to split a coin object into multiple coins.
func (c *Client) SplitCoin(
	ctx context.Context,
	signer suiAddress,
	Coin suiObjectID,
	splitAmounts []types.SafeSuiBigInt[uint64],
	gas *suiObjectID,
	gasBudget types.SafeSuiBigInt[uint64],
) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, splitCoin, signer, Coin, splitAmounts, gas, gasBudget)
}

// SplitCoinEqual Create an unsigned transaction to split a coin object into multiple equal-size coins.
func (c *Client) SplitCoinEqual(
	ctx context.Context,
	signer suiAddress,
	Coin suiObjectID,
	splitCount types.SafeSuiBigInt[uint64],
	gas *suiObjectID,
	gasBudget types.SafeSuiBigInt[uint64],
) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, splitCoinEqual, signer, Coin, splitCount, gas, gasBudget)
}

// MergeCoins Create an unsigned transaction to merge multiple coins into one coin.
func (c *Client) MergeCoins(
	ctx context.Context,
	signer suiAddress,
	primaryCoin, coinToMerge suiObjectID,
	gas *suiObjectID,
	gasBudget types.SafeSuiBigInt[uint64],
) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, mergeCoins, signer, primaryCoin, coinToMerge, gas, gasBudget)
}

func (c *Client) Publish(
	ctx context.Context,
	sender suiAddress,
	compiledModules []*suiBase64Data,
	dependencies []suiObjectID,
	gas suiObjectID,
	gasBudget uint,
) (*types.TransactionBytes, error) {
	var resp types.TransactionBytes
	return &resp, c.CallContext(ctx, &resp, publish, sender, compiledModules, dependencies, gas, gasBudget)
}

// MoveCall Create an unsigned transaction to execute a Move call on the network, by calling the specified function in the module of a given package.
// TODO: not support param `typeArguments` yet.
// So now only methods with `typeArguments` are supported
// TODO: execution_mode : <SuiTransactionBlockBuilderMode>
func (c *Client) MoveCall(
	ctx context.Context,
	signer suiAddress,
	packageId suiObjectID,
	module, function string,
	typeArgs []string,
	arguments []any,
	gas *suiObjectID,
	gasBudget types.SafeSuiBigInt[uint64],
) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(
		ctx,
		&resp,
		moveCall,
		signer,
		packageId,
		module,
		function,
		typeArgs,
		arguments,
		gas,
		gasBudget,
	)
}

// TODO: execution_mode : <SuiTransactionBlockBuilderMode>
func (c *Client) BatchTransaction(
	ctx context.Context,
	signer suiAddress,
	txnParams []map[string]interface{},
	gas *suiObjectID,
	gasBudget uint64,
) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	return &resp, c.CallContext(ctx, &resp, batchTransaction, signer, txnParams, gas, gasBudget)
}

func (c *Client) QueryTransactionBlocks(
	ctx context.Context, query types.SuiTransactionBlockResponseQuery,
	cursor *suiDigest, limit *uint, descendingOrder bool,
) (*types.TransactionBlocksPage, error) {
	resp := types.TransactionBlocksPage{}
	return &resp, c.CallContext(ctx, &resp, queryTransactionBlocks, query, cursor, limit, descendingOrder)
}

func (c *Client) QueryEvents(
	ctx context.Context, query types.EventFilter, cursor *types.EventId, limit *uint,
	descendingOrder bool,
) (*types.EventPage, error) {
	var resp types.EventPage
	return &resp, c.CallContext(ctx, &resp, queryEvents, query, cursor, limit, descendingOrder)
}

func (c *Client) GetDynamicFields(
	ctx context.Context, parentObjectId suiObjectID, cursor *suiObjectID,
	limit *uint,
) (*types.DynamicFieldPage, error) {
	var resp types.DynamicFieldPage
	return &resp, c.CallContext(ctx, &resp, getDynamicFields, parentObjectId, cursor, limit)
}

func (c *Client) GetDynamicFieldObject(
	ctx context.Context, parentObjectId suiObjectID,
	name sui_types.DynamicFieldName,
) (*types.SuiObjectResponse, error) {
	var resp types.SuiObjectResponse
	return &resp, c.CallContext(ctx, &resp, getDynamicFieldObject, parentObjectId, name)
}
