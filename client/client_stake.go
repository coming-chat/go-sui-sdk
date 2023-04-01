package client

import (
	"context"

	"github.com/coming-chat/go-sui/types"
)

const (
	packageIdStake          = "0x03"
	systemState             = "0x05"
	moduleNameStake         = "sui_system"
	funcNameAddStake        = "request_add_stake"
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

func (c *Client) RequestAddDelegation(ctx context.Context, signer types.Address, coin types.ObjectId, validator types.Address, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	packageId, _ := types.NewHexData(packageIdStake)
	systemId, _ := types.NewHexData(systemState)
	return c.MoveCall(ctx,
		signer,
		*packageId,
		moduleNameStake,
		funcNameAddStake,
		[]string{},
		[]any{
			systemId,
			coin,
			validator,
		},
		gas, gasBudget)
}

// TODO: fix params
func (c *Client) RequestAddDelegationMulCoin(ctx context.Context, signer types.Address, coins []types.ObjectId, amount uint64, validator types.Address, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	packageId, _ := types.NewHexData(packageIdStake)
	systemId, _ := types.NewHexData(systemState)
	return c.MoveCall(ctx,
		signer,                  // 0xd77955e670f42c1bc5e94b9e68e5fe9bdbed9134d784f2a14dfe5fc1b24b5d9f
		*packageId,              // 0x3
		moduleNameStake,         // sui_system
		funcNameAddStakeMulCoin, // request_add_stake_mul_coin
		[]string{},
		[]any{
			systemId,  // 0x5
			coins,     // [0x0501ebf5518912e380e8b3b68f93548418fb1bce59ed025f68ad6d236f012f92]
			amount,    // uint64(10000) ???
			validator, // 0x8ce890590fed55c37d44a043e781ad94254b413ee079a53fb5c037f7a6311304
		},
		gas, gasBudget)
}

func (c *Client) RequestWithdrawDelegation(ctx context.Context, signer types.Address, stakedSuiId types.ObjectId, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
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
			stakedSuiId,
		},
		gas, gasBudget)
}
