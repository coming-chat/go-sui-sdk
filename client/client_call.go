package client

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/coming-chat/go-sui/types"
)

func (c *Client) GetSuiCoinsOwnedByAddress(ctx context.Context, address types.Address) (types.Coins, error) {
	coinType := "0x2::coin::Coin<0x2::sui::SUI>"
	coinObjects, err := c.BatchGetObjectsOwnedByAddress(ctx, address, coinType)
	if err != nil {
		return nil, err
	}

	type coinData struct {
		Fields struct {
			Balance uint64 `json:"balance"`
		} `json:"fields"`
	}
	coins := types.Coins{}
	for _, coin := range coinObjects {
		if coin.Status != types.ObjectStatusExists {
			continue
		}
		bytes, err := json.Marshal(coin.Details.Data)
		if err != nil {
			return nil, err
		}
		coindata := coinData{}
		err = json.Unmarshal(bytes, &coindata)
		if err != nil {
			return nil, err
		}

		coins = append(coins, types.Coin{
			Balance:             coindata.Fields.Balance,
			Type:                coinType,
			Owner:               coin.Details.Owner,
			PreviousTransaction: coin.Details.PreviousTransaction,
			Reference:           coin.Details.Reference,
		})
	}

	return coins, nil
}

// @param filterType You can specify filtering out the specified resources, this will fetch all resources if it is not empty ""
func (c *Client) BatchGetObjectsOwnedByAddress(ctx context.Context, address types.Address, filterType string) ([]types.ObjectRead, error) {
	infos, err := c.GetObjectsOwnedByAddress(ctx, address)
	if err != nil {
		return nil, err
	}
	filterType = strings.TrimSpace(filterType)
	elems := []BatchElem{}
	for _, info := range infos {
		if filterType != "" && filterType != info.Type {
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

func (c *Client) ExecuteTransaction(ctx context.Context, txn types.SignedTransaction) (*types.TransactionResponse, error) {
	resp := types.TransactionResponse{}
	err := c.CallContext(ctx, &resp, "sui_executeTransaction", txn.TxBytes, txn.SigScheme, txn.Signature, txn.PublicKey)
	return &resp, err
}

func (c *Client) GetObject(ctx context.Context, objID types.ObjectId) (*types.ObjectRead, error) {
	resp := types.ObjectRead{}
	err := c.CallContext(ctx, &resp, "sui_getObject", objID)
	return &resp, err
}

func (c *Client) GetObjectsOwnedByAddress(ctx context.Context, address types.Address) ([]types.ObjectInfo, error) {
	resp := []types.ObjectInfo{}
	err := c.CallContext(ctx, &resp, "sui_getObjectsOwnedByAddress", address)
	return resp, err
}

func (c *Client) GetObjectsOwnedByObject(ctx context.Context, objID types.ObjectId) ([]types.ObjectInfo, error) {
	resp := []types.ObjectInfo{}
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

func (c *Client) GetTransaction(ctx context.Context, digest types.Base64Data) (*types.TransactionResponse, error) {
	resp := types.TransactionResponse{}
	err := c.CallContext(ctx, &resp, "sui_getTransaction", digest)
	return &resp, err
}

// Create an unsigned transaction to merge multiple coins into one coin.
func (c *Client) MergeCoins(ctx context.Context, signer types.Address, primaryCoin, coinToMerge types.ObjectId, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	err := c.CallContext(ctx, &resp, "sui_mergeCoins", signer, primaryCoin, coinToMerge, gas, gasBudget)
	return &resp, err
}

// Create an unsigned transaction to split a coin object into multiple coins.
func (c *Client) SplitCoin(ctx context.Context, signer types.Address, Coin types.ObjectId, splitAmounts []uint64, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	err := c.CallContext(ctx, &resp, "sui_splitCoin", signer, Coin, splitAmounts, gas, gasBudget)
	return &resp, err
}

// Create an unsigned transaction to split a coin object into multiple equal-size coins.
func (c *Client) SplitCoinEqual(ctx context.Context, signer types.Address, Coin types.ObjectId, splitCount uint64, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	err := c.CallContext(ctx, &resp, "sui_splitCoinEqual", signer, Coin, splitCount, gas, gasBudget)
	return &resp, err
}

// Create an unsigned transaction to transfer an object from one address to another. The object's type must allow public transfers
func (c *Client) TransferObject(ctx context.Context, signer, recipient types.Address, objID, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	err := c.CallContext(ctx, &resp, "sui_transferObject", signer, objID, gas, gasBudget, recipient)
	return &resp, err
}

// Create an unsigned transaction to send SUI coin object to a Sui address. The SUI object is also used as the gas object.
func (c *Client) TransferSui(ctx context.Context, signer, recipient types.Address, suiObjID types.ObjectId, amount, gasBudget uint64) (*types.TransactionBytes, error) {
	resp := types.TransactionBytes{}
	err := c.CallContext(ctx, &resp, "sui_transferSui", signer, suiObjID, gasBudget, recipient, amount)
	return &resp, err
}
