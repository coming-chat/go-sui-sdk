package sui_types

import "github.com/coming-chat/go-sui/v2/move_types"

const (
	StakingPoolModuleName = move_types.Identifier("staking_pool")
	StakedSuiStructName   = move_types.Identifier("StakedSui")

	AddStakeMulCoinFunName = move_types.Identifier("request_add_stake_mul_coin")
	AddStakeFunName        = move_types.Identifier("request_add_stake")
	WithdrawStakeFunName   = move_types.Identifier("request_withdraw_stake")
)
