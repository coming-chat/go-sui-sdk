package client

import (
	"context"

	"github.com/coming-chat/go-sui/types"
)

func (c *Client) GetDelegatedStakes(ctx context.Context, owner types.Address) ([]types.DelegatedStake, error) {
	var resp []types.DelegatedStake
	return resp, c.CallContext(ctx, &resp, "sui_getDelegatedStakes", owner)
}

func (c *Client) GetValidators(ctx context.Context) ([]types.ValidatorMetadata, error) {
	var resp []types.ValidatorMetadata
	return resp, c.CallContext(ctx, &resp, "sui_getValidators")
}

func (c *Client) GetSuiSystemState(ctx context.Context) (*types.SuiSystemState, error) {
	var resp types.SuiSystemState
	return &resp, c.CallContext(ctx, &resp, "sui_getSuiSystemState")
}

func (c *Client) RequestAddDelegation(ctx context.Context, signer types.Address, coins []types.ObjectId, amount uint64, validator types.Address, gas types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	var resp types.TransactionBytes
	return &resp, c.CallContext(ctx, &resp, "sui_requestAddDelegation", signer, coins, amount, validator, gas, gasBudget)
}

func (c *Client) RequestSwitchDelegation(ctx context.Context, signer types.Address, delegation, stakedSui types.ObjectId, newValidator types.Address, gas types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	var resp types.TransactionBytes
	return &resp, c.CallContext(ctx, &resp, "sui_requestSwitchDelegation", signer, delegation, stakedSui, newValidator, gas, gasBudget)
}

func (c *Client) RequestWithdrawDelegation(ctx context.Context, signer types.Address, delegation, stakedSui types.ObjectId, gas types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	var resp types.TransactionBytes
	return &resp, c.CallContext(ctx, &resp, "sui_requestWithdrawDelegation", signer, delegation, stakedSui, gas, gasBudget)
}
