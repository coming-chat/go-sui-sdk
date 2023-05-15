package client

import (
	"context"

	"github.com/coming-chat/go-sui/lib"
	"github.com/coming-chat/go-sui/sui_types"
	"github.com/coming-chat/go-sui/sui_types/sui_system_state"
	"github.com/coming-chat/go-sui/types"
	"github.com/fardream/go-bcs/bcs"
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

func BCS_RequestAddStake(signer suiAddress, coins []*sui_types.ObjectRef, amount types.SafeSuiBigInt[uint64], validator suiAddress, gasPrice, gasBudget uint64) ([]byte, error) {
	// build with BCS
	ptb := sui_types.NewProgrammableTransactionBuilder()
	amtArg, err := ptb.Pure(amount.Uint64())
	if err != nil {
		return nil, err
	}
	arg0, err := ptb.Obj(sui_types.SuiSystemMutObj)
	if err != nil {
		return nil, err
	}
	arg1 := ptb.Command(sui_types.Command{
		SplitCoins: &struct {
			Argument  sui_types.Argument
			Arguments []sui_types.Argument
		}{
			Argument:  sui_types.Argument{GasCoin: &lib.EmptyEnum{}},
			Arguments: []sui_types.Argument{amtArg},
		},
	}) // the coin is split result argument
	arg2, err := ptb.Pure(validator)
	if err != nil {
		return nil, err
	}

	ptb.Command(sui_types.Command{
		MoveCall: &sui_types.ProgrammableMoveCall{
			Package:  *sui_types.SuiSystemAddress,
			Module:   sui_system_state.SuiSystemModuleName,
			Function: sui_types.AddStakeFunName,
			Arguments: []sui_types.Argument{
				arg0, arg1, arg2,
			},
		},
	})
	pt := ptb.Finish()
	tx := sui_types.NewProgrammable(
		signer, coins, pt, gasBudget, gasPrice,
	)
	return bcs.Marshal(tx)
}
