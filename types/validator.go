package types

import (
	"math"

	"github.com/shopspring/decimal"
)

type StakeStatus = string

const (
	StakeStatusActive   = "Active"
	StakeStatusPending  = "Pending"
	StakeStatusUnstaked = "Unstaked"
)

type StakeObject struct {
	StakedSuiId       ObjectId    `json:"stakedSuiId"`
	StakeRequestEpoch EpochId     `json:"stakeRequestEpoch"`
	StakeActiveEpoch  EpochId     `json:"stakeActiveEpoch"`
	Principal         uint64      `json:"principal"`
	Status            StakeStatus `json:"status"`
	EstimatedReward   *uint64     `json:"estimatedReward,omitempty"`
}

type DelegatedStake struct {
	ValidatorAddress Address       `json:"validatorAddress"`
	StakingPool      ObjectId      `json:"stakingPool"`
	Stakes           []StakeObject `json:"stakes"`
}

type SuiValidatorSummary struct {
	SuiAddress             Address `json:"suiAddress"`
	ProtocolPubkeyBytes    string  `json:"protocolPubkeyBytes"`
	NetworkPubkeyBytes     string  `json:"networkPubkeyBytes"`
	WorkerPubkeyBytes      string  `json:"workerPubkeyBytes"`
	ProofOfPossessionBytes string  `json:"proofOfPossessionBytes"`
	OperationCapId         string  `json:"operationCapId"`
	Name                   string  `json:"name"`
	Description            string  `json:"description"`
	ImageUrl               string  `json:"imageUrl"`
	ProjectUrl             string  `json:"projectUrl"`
	P2pAddress             string  `json:"p2pAddress"`
	NetAddress             string  `json:"netAddress"`
	PrimaryAddress         string  `json:"primaryAddress"`
	WorkerAddress          string  `json:"workerAddress"`

	NextEpochProtocolPubkeyBytes string `json:"nextEpochProtocolPubkeyBytes"`
	NextEpochProofOfPossession   string `json:"nextEpochProofOfPossession"`
	NextEpochNetworkPubkeyBytes  string `json:"nextEpochNetworkPubkeyBytes"`
	NextEpochWorkerPubkeyBytes   string `json:"nextEpochWorkerPubkeyBytes"`
	NextEpochNetAddress          string `json:"nextEpochNetAddress"`
	NextEpochP2pAddress          string `json:"nextEpochP2pAddress"`
	NextEpochPrimaryAddress      string `json:"nextEpochPrimaryAddress"`
	NextEpochWorkerAddress       string `json:"nextEpochWorkerAddress"`

	VotingPower             int64           `json:"votingPower"`
	GasPrice                int64           `json:"gasPrice"`
	CommissionRate          int64           `json:"commissionRate"`
	NextEpochStake          decimal.Decimal `json:"nextEpochStake"`
	NextEpochGasPrice       int64           `json:"nextEpochGasPrice"`
	NextEpochCommissionRate int64           `json:"nextEpochCommissionRate"`
	StakingPoolId           string          `json:"stakingPoolId"`

	StakingPoolActivationEpoch   uint64 `json:"stakingPoolActivationEpoch"`
	StakingPoolDeactivationEpoch int64  `json:"stakingPoolDeactivationEpoch"`

	StakingPoolSuiBalance    decimal.Decimal `json:"stakingPoolSuiBalance"`
	RewardsPool              decimal.Decimal `json:"rewardsPool"`
	PoolTokenBalance         decimal.Decimal `json:"poolTokenBalance"`
	PendingStake             decimal.Decimal `json:"pendingStake"`
	PendingPoolTokenWithdraw decimal.Decimal `json:"pendingPoolTokenWithdraw"`
	PendingTotalSuiWithdraw  decimal.Decimal `json:"pendingTotalSuiWithdraw"`
	ExchangeRatesId          string          `json:"exchangeRatesId"`
	ExchangeRatesSize        int64           `json:"exchangeRatesSize"`
}

func (v *SuiValidatorSummary) CalculateAPY(epoch uint64) float64 {
	var (
		stakingPoolSuiBalance      = v.StakingPoolSuiBalance
		stakingPoolActivationEpoch = v.StakingPoolActivationEpoch
		poolTokenBalance           = v.PoolTokenBalance
	)

	// If the staking pool is active then we calculate its APY. Or if staking started in epoch 0
	if stakingPoolActivationEpoch == 0 {
		numEpochsParticipated := epoch - stakingPoolActivationEpoch
		pow1, _ := stakingPoolSuiBalance.Sub(poolTokenBalance).Div(poolTokenBalance).Add(decimal.NewFromInt(1)).Float64()
		pow2, _ := decimal.NewFromInt(365).Div(decimal.NewFromInt(int64(numEpochsParticipated))).Float64()
		apy := (math.Pow(pow1, pow2) - 1) * 100
		if apy > 100000 {
			return 0
		} else {
			return apy
		}
	} else {
		return 0
	}
}

type SuiSystemStateSummary struct {
	Epoch                                 uint64                `json:"epoch"`
	ProtocolVersion                       uint64                `json:"protocolVersion"`
	SystemStateVersion                    uint64                `json:"systemStateVersion"`
	StorageFundTotalObjectStorageRebates  uint64                `json:"storageFundTotalObjectStorageRebates"`
	StorageFundNonRefundableBalance       uint64                `json:"storageFundNonRefundableBalance"`
	ReferenceGasPrice                     uint64                `json:"referenceGasPrice"`
	SafeMode                              bool                  `json:"safeMode"`
	SafeModeStorageRewards                uint64                `json:"safeModeStorageRewards"`
	SafeModeComputationRewards            uint64                `json:"safeModeComputationRewards"`
	SafeModeStorageRebates                uint64                `json:"safeModeStorageRebates"`
	SafeModeNonRefundableStorageFee       uint64                `json:"safeModeNonRefundableStorageFee"`
	EpochStartTimestampMs                 uint64                `json:"epochStartTimestampMs"`
	EpochDurationMs                       uint64                `json:"epochDurationMs"`
	StakeSubsidyStartEpoch                uint64                `json:"stakeSubsidyStartEpoch"`
	MaxValidatorCount                     uint64                `json:"maxValidatorCount"`
	MinValidatorJoiningStake              uint64                `json:"minValidatorJoiningStake"`
	ValidatorLowStakeThreshold            uint64                `json:"validatorLowStakeThreshold"`
	ValidatorVeryLowStakeThreshold        uint64                `json:"validatorVeryLowStakeThreshold"`
	ValidatorLowStakeGracePeriod          uint64                `json:"validatorLowStakeGracePeriod"`
	StakeSubsidyBalance                   uint64                `json:"stakeSubsidyBalance"`
	StakeSubsidyDistributionCounter       uint64                `json:"stakeSubsidyDistributionCounter"`
	StakeSubsidyCurrentDistributionAmount uint64                `json:"stakeSubsidyCurrentDistributionAmount"`
	StakeSubsidyPeriodLength              uint64                `json:"stakeSubsidyPeriodLength"`
	StakeSubsidyDecreaseRate              uint16                `json:"stakeSubsidyDecreaseRate"`
	TotalStake                            uint64                `json:"totalStake"`
	ActiveValidators                      []SuiValidatorSummary `json:"activeValidators"`
	PendingActiveValidatorsId             string                `json:"pendingActiveValidatorsId"`
	PendingActiveValidatorsSize           uint64                `json:"pendingActiveValidatorsSize"`
	PendingRemovals                       []uint64              `json:"pendingRemovals"`
	StakingPoolMappingsId                 string                `json:"stakingPoolMappingsId"`
	StakingPoolMappingsSize               uint64                `json:"stakingPoolMappingsSize"`
	InactivePoolsId                       string                `json:"inactivePoolsId"`
	InactivePoolsSize                     uint64                `json:"inactivePoolsSize"`
	ValidatorCandidatesId                 string                `json:"validatorCandidatesId"`
	ValidatorCandidatesSize               uint64                `json:"validatorCandidatesSize"`
	AtRiskValidators                      []interface{}         `json:"atRiskValidators"`
	ValidatorReportRecords                []interface{}         `json:"validatorReportRecords"`
}
