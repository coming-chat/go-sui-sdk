package client

import (
	"context"

	"github.com/coming-chat/go-sui/types"
)

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
