package client

import (
	"context"

	"github.com/coming-chat/go-sui/types"
)

const (
	packageIdStake          = "0x03"
	systemState             = "0x05"
	moduleNameStake         = "sui_system"
	funcNameAddStakeMulCoin = "request_add_stake_mul_coin"
	funcNameWithdrawStake   = "request_withdraw_stake"
)

func (c *Client) GetLatestSuiSystemState(ctx context.Context) (*types.SuiSystemStateSummary, error) {
	var resp types.SuiSystemStateSummary
	return &resp, c.CallContext(ctx, &resp, getLatestSuiSystemState)
}

func (c *Client) GetStakes(ctx context.Context, owner types.Address) ([]types.DelegatedStake, error) {
	var resp []types.DelegatedStake
	return resp, c.CallContext(ctx, &resp, getStakes, owner)
}

func (c *Client) GetStakesByIds(ctx context.Context, stakedSuiIds []types.ObjectId) ([]types.DelegatedStake, error) {
	var resp []types.DelegatedStake
	return resp, c.CallContext(ctx, &resp, getStakesByIds, stakedSuiIds)
}

// TODO: fix params
func (c *Client) RequestAddDelegation(ctx context.Context, signer types.Address, coins []types.ObjectId, amount uint64, validator types.Address, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	packageId, _ := types.NewHexData(packageIdStake)
	systemId, _ := types.NewHexData(systemState)
	return c.MoveCall(ctx,
		signer,
		*packageId,
		moduleNameStake,
		funcNameAddStakeMulCoin,
		[]string{},
		[]any{
			systemId,
			coins,
			amount,
			validator,
		},
		gas, gasBudget)
}

func (c *Client) RequestWithdrawDelegation(ctx context.Context, signer types.Address, stakedSui types.ObjectId, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	packageId, _ := types.NewHexData(packageIdStake)
	systemId, _ := types.NewHexData(systemState)
	return c.MoveCall(ctx,
		signer,
		*packageId,
		moduleNameStake,
		funcNameWithdrawStake,
		[]string{},
		[]any{
			systemId,
			stakedSui,
		},
		gas, gasBudget)
}
