package client

import (
	"context"
	"encoding/json"
	"math"

	"github.com/coming-chat/go-sui/types"
	"github.com/shopspring/decimal"
)

const ROLLING_AVERAGE = 30

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

func (c *Client) RequestAddStake(ctx context.Context, signer types.Address, coins []types.ObjectId, amount uint64, validator types.Address, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	var resp types.TransactionBytes
	return &resp, c.CallContext(ctx, &resp, requestAddStake, signer, coins, amount, validator, gas, gasBudget)
}

func (c *Client) RequestWithdrawStake(ctx context.Context, signer types.Address, stakedSuiId types.ObjectId, gas *types.ObjectId, gasBudget uint64) (*types.TransactionBytes, error) {
	var resp types.TransactionBytes
	return &resp, c.CallContext(ctx, &resp, requestWithdrawStake, signer, stakedSuiId, gas, gasBudget)
}

// APY_e = (1 + epoch_rewards / stake)^365-1
// APY_e_30rollingaverage = average(APY_e,APY_e-1,…,APY_e-29);
// @return Average apy of all validator in the most recent epoch. key is the validator address, value is the average apy.
func (c *Client) GetAndCalculateRollingAverageApys(ctx context.Context, validatorCount int) (map[string]float64, error) {
	limit := uint(validatorCount) * ROLLING_AVERAGE
	if limit > QUERY_MAX_RESULT_LIMIT {
		limit = QUERY_MAX_RESULT_LIMIT
	}
	eventType := "0x3::validator_set::ValidatorEpochInfoEvent"
	events, err := c.QueryEvents(ctx, types.EventFilter{
		MoveEventType: &eventType,
	}, nil, &limit, true)
	if err != nil {
		return nil, err
	}

	apyMaps := make(map[string][]float64, 0)
	for _, event := range events.Data {
		address, apy := calculateAPY(event.ParsedJson)
		if address == "" {
			continue
		}
		if apyList, ok := apyMaps[address]; ok {
			if len(apyList) >= ROLLING_AVERAGE {
				continue // dont need more
			}
			apyList = append(apyList, apy)
			apyMaps[address] = apyList
		} else {
			apyMaps[address] = []float64{apy}
		}
	}

	apyResults := make(map[string]float64)
	for address, apys := range apyMaps {
		totalApy := float64(0)
		for _, apy := range apys {
			totalApy += apy
		}
		apyResults[address] = totalApy / float64(len(apys))
	}

	return apyResults, nil
}

func calculateAPY(parsedJson interface{}) (address string, apy float64) {
	data, err := json.Marshal(parsedJson)
	if err != nil {
		return
	}
	var model struct {
		ValidatorAddress  string          `json:"validator_address"`
		Stake             decimal.Decimal `json:"stake"`
		PoolStakingReward decimal.Decimal `json:"pool_staking_reward"`
	}
	err = json.Unmarshal(data, &model)
	if err != nil {
		return
	}
	if model.ValidatorAddress == "" {
		return
	}

	ratio, _ := model.PoolStakingReward.Div(model.Stake).Add(decimal.NewFromInt(1)).Float64()
	apy = math.Pow(ratio, 365.0) - 1
	if apy > 10000 {
		apy = 0
	}
	return model.ValidatorAddress, apy
}
