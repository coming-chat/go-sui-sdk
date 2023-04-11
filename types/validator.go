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
	Principal         SuiBigInt   `json:"principal"`
	Status            StakeStatus `json:"status"`
	EstimatedReward   *SuiBigInt  `json:"estimatedReward,omitempty"`
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

	VotingPower             SuiBigInt `json:"votingPower"`
	GasPrice                SuiBigInt `json:"gasPrice"`
	CommissionRate          SuiBigInt `json:"commissionRate"`
	NextEpochStake          SuiBigInt `json:"nextEpochStake"`
	NextEpochGasPrice       SuiBigInt `json:"nextEpochGasPrice"`
	NextEpochCommissionRate SuiBigInt `json:"nextEpochCommissionRate"`
	StakingPoolId           string    `json:"stakingPoolId"`

	StakingPoolActivationEpoch   *SuiBigInt `json:"stakingPoolActivationEpoch,omitempty"`
	StakingPoolDeactivationEpoch *SuiBigInt `json:"stakingPoolDeactivationEpoch,omitempty"`

	StakingPoolSuiBalance    SuiBigInt `json:"stakingPoolSuiBalance"`
	RewardsPool              SuiBigInt `json:"rewardsPool"`
	PoolTokenBalance         SuiBigInt `json:"poolTokenBalance"`
	PendingStake             SuiBigInt `json:"pendingStake"`
	PendingPoolTokenWithdraw SuiBigInt `json:"pendingPoolTokenWithdraw"`
	PendingTotalSuiWithdraw  SuiBigInt `json:"pendingTotalSuiWithdraw"`
	ExchangeRatesId          string    `json:"exchangeRatesId"`
	ExchangeRatesSize        SuiBigInt `json:"exchangeRatesSize"`
}

func (v *SuiValidatorSummary) CalculateAPY(epoch uint64) float64 {
	var (
		stakingPoolSuiBalance      = v.StakingPoolSuiBalance
		stakingPoolActivationEpoch = v.StakingPoolActivationEpoch
		poolTokenBalance           = v.PoolTokenBalance
	)

	// If the staking pool is active then we calculate its APY. Or if staking started in epoch 0
	if stakingPoolActivationEpoch.IsZero() {
		numEpochsParticipated := epoch - stakingPoolActivationEpoch.BigInt().Uint64()
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
	Epoch                                 SuiBigInt             `json:"epoch"`
	ProtocolVersion                       SuiBigInt             `json:"protocolVersion"`
	SystemStateVersion                    SuiBigInt             `json:"systemStateVersion"`
	StorageFundTotalObjectStorageRebates  SuiBigInt             `json:"storageFundTotalObjectStorageRebates"`
	StorageFundNonRefundableBalance       SuiBigInt             `json:"storageFundNonRefundableBalance"`
	ReferenceGasPrice                     SuiBigInt             `json:"referenceGasPrice"`
	SafeMode                              bool                  `json:"safeMode"`
	SafeModeStorageRewards                SuiBigInt             `json:"safeModeStorageRewards"`
	SafeModeComputationRewards            SuiBigInt             `json:"safeModeComputationRewards"`
	SafeModeStorageRebates                SuiBigInt             `json:"safeModeStorageRebates"`
	SafeModeNonRefundableStorageFee       SuiBigInt             `json:"safeModeNonRefundableStorageFee"`
	EpochStartTimestampMs                 SuiBigInt             `json:"epochStartTimestampMs"`
	EpochDurationMs                       SuiBigInt             `json:"epochDurationMs"`
	StakeSubsidyStartEpoch                SuiBigInt             `json:"stakeSubsidyStartEpoch"`
	MaxValidatorCount                     SuiBigInt             `json:"maxValidatorCount"`
	MinValidatorJoiningStake              SuiBigInt             `json:"minValidatorJoiningStake"`
	ValidatorLowStakeThreshold            SuiBigInt             `json:"validatorLowStakeThreshold"`
	ValidatorVeryLowStakeThreshold        SuiBigInt             `json:"validatorVeryLowStakeThreshold"`
	ValidatorLowStakeGracePeriod          SuiBigInt             `json:"validatorLowStakeGracePeriod"`
	StakeSubsidyBalance                   SuiBigInt             `json:"stakeSubsidyBalance"`
	StakeSubsidyDistributionCounter       SuiBigInt             `json:"stakeSubsidyDistributionCounter"`
	StakeSubsidyCurrentDistributionAmount SuiBigInt             `json:"stakeSubsidyCurrentDistributionAmount"`
	StakeSubsidyPeriodLength              SuiBigInt             `json:"stakeSubsidyPeriodLength"`
	StakeSubsidyDecreaseRate              uint16                `json:"stakeSubsidyDecreaseRate"`
	TotalStake                            SuiBigInt             `json:"totalStake"`
	ActiveValidators                      []SuiValidatorSummary `json:"activeValidators"`
	PendingActiveValidatorsId             ObjectId              `json:"pendingActiveValidatorsId"`
	PendingActiveValidatorsSize           SuiBigInt             `json:"pendingActiveValidatorsSize"`
	PendingRemovals                       []uint64              `json:"pendingRemovals"`
	StakingPoolMappingsId                 ObjectId              `json:"stakingPoolMappingsId"`
	StakingPoolMappingsSize               SuiBigInt             `json:"stakingPoolMappingsSize"`
	InactivePoolsId                       ObjectId              `json:"inactivePoolsId"`
	InactivePoolsSize                     SuiBigInt             `json:"inactivePoolsSize"`
	ValidatorCandidatesId                 ObjectId              `json:"validatorCandidatesId"`
	ValidatorCandidatesSize               SuiBigInt             `json:"validatorCandidatesSize"`
	AtRiskValidators                      []interface{}         `json:"atRiskValidators"`
	ValidatorReportRecords                []interface{}         `json:"validatorReportRecords"`
}
