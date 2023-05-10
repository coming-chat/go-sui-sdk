package client

import (
	"context"

	"github.com/coming-chat/go-sui/types"
)

func (c *Client) GetLatestSuiSystemState(ctx context.Context) (*types.SuiSystemStateSummary, error) {
	var resp types.SuiSystemStateSummary
	return &resp, c.CallContext(ctx, &resp, getLatestSuiSystemState)
}

func (c *Client) GetValidatorsApy(ctx context.Context) (*types.ValidatorsApy, error) {
	var resp types.ValidatorsApy
	return &resp, c.CallContext(ctx, &resp, getValidatorsApy)
}

func (c *Client) GetStakes(ctx context.Context, owner suiAddress) ([]types.DelegatedStake, error) {
	var resp []types.DelegatedStake
	return resp, c.CallContext(ctx, &resp, getStakes, owner)
}

func (c *Client) GetStakesByIds(ctx context.Context, stakedSuiIds []suiObjectID) ([]types.DelegatedStake, error) {
	var resp []types.DelegatedStake
	return resp, c.CallContext(ctx, &resp, getStakesByIds, stakedSuiIds)
}

func (c *Client) RequestAddStake(ctx context.Context, signer suiAddress, coins []suiObjectID, amount types.SuiBigInt, validator suiAddress, gas *suiObjectID, gasBudget types.SuiBigInt) (*types.TransactionBytes, error) {
	var resp types.TransactionBytes
	return &resp, c.CallContext(ctx, &resp, requestAddStake, signer, coins, amount, validator, gas, gasBudget)
}

func (c *Client) RequestWithdrawStake(ctx context.Context, signer suiAddress, stakedSuiId suiObjectID, gas *suiObjectID, gasBudget types.SuiBigInt) (*types.TransactionBytes, error) {
	var resp types.TransactionBytes
	return &resp, c.CallContext(ctx, &resp, requestWithdrawStake, signer, stakedSuiId, gas, gasBudget)
}
