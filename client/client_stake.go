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
