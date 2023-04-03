package types

import (
	"github.com/shopspring/decimal"
)

type Balance struct {
	Value decimal.Decimal `json:"value"`
}

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

	StakingPoolActivationEpoch   int64 `json:"stakingPoolActivationEpoch"`
	StakingPoolDeactivationEpoch int64 `json:"stakingPoolDeactivationEpoch"`

	StakingPoolSuiBalance    decimal.Decimal `json:"stakingPoolSuiBalance"`
	RewardsPool              decimal.Decimal `json:"rewardsPool"`
	PoolTokenBalance         decimal.Decimal `json:"poolTokenBalance"`
	PendingStake             decimal.Decimal `json:"pendingStake"`
	PendingPoolTokenWithdraw decimal.Decimal `json:"pendingPoolTokenWithdraw"`
	PendingTotalSuiWithdraw  decimal.Decimal `json:"pendingTotalSuiWithdraw"`
	ExchangeRatesId          string          `json:"exchangeRatesId"`
	ExchangeRatesSize        int64           `json:"exchangeRatesSize"`
}

func (v *SuiValidatorSummary) CalculateAPY(epoch int64) float64 {
	var (
		stakingPoolSuiBalance      = v.StakingPoolSuiBalance
		stakingPoolActivationEpoch = v.StakingPoolActivationEpoch
		poolTokenBalance           = v.PoolTokenBalance
	)

	// If the staking pool is active then we calculate its APY. Or if staking started in epoch 0
	if stakingPoolActivationEpoch == 0 {
		numEpochsParticipated := epoch - stakingPoolActivationEpoch
		pow1 := stakingPoolSuiBalance.Sub(poolTokenBalance).Div(poolTokenBalance).Add(decimal.NewFromInt(1))
		pow2 := decimal.NewFromInt(365).Div(decimal.NewFromInt(numEpochsParticipated))
		apy := pow1.Pow(pow2).Sub(decimal.NewFromInt(1)).Mul(decimal.NewFromInt(100))
		apyValue, _ := apy.Round(4).Float64()
		if apyValue > 100000 {
			return 0
		} else {
			return apyValue
		}
	} else {
		return 0
	}
}

type SuiSystemStateSummary struct {
	Epoch                                 int64                 `json:"epoch"`
	ProtocolVersion                       int64                 `json:"protocolVersion"`
	SystemStateVersion                    int64                 `json:"systemStateVersion"`
	StorageFundTotalObjectStorageRebates  int64                 `json:"storageFundTotalObjectStorageRebates"`
	StorageFundNonRefundableBalance       int64                 `json:"storageFundNonRefundableBalance"`
	ReferenceGasPrice                     int64                 `json:"referenceGasPrice"`
	SafeMode                              bool                  `json:"safeMode"`
	SafeModeStorageRewards                int64                 `json:"safeModeStorageRewards"`
	SafeModeComputationRewards            int64                 `json:"safeModeComputationRewards"`
	SafeModeStorageRebates                int64                 `json:"safeModeStorageRebates"`
	SafeModeNonRefundableStorageFee       int64                 `json:"safeModeNonRefundableStorageFee"`
	EpochStartTimestampMs                 int64                 `json:"epochStartTimestampMs"`
	EpochDurationMs                       int64                 `json:"epochDurationMs"`
	StakeSubsidyStartEpoch                int64                 `json:"stakeSubsidyStartEpoch"`
	MaxValidatorCount                     int64                 `json:"maxValidatorCount"`
	MinValidatorJoiningStake              decimal.Decimal       `json:"minValidatorJoiningStake"`
	ValidatorLowStakeThreshold            decimal.Decimal       `json:"validatorLowStakeThreshold"`
	ValidatorVeryLowStakeThreshold        decimal.Decimal       `json:"validatorVeryLowStakeThreshold"`
	ValidatorLowStakeGracePeriod          int64                 `json:"validatorLowStakeGracePeriod"`
	StakeSubsidyBalance                   decimal.Decimal       `json:"stakeSubsidyBalance"`
	StakeSubsidyDistributionCounter       int64                 `json:"stakeSubsidyDistributionCounter"`
	StakeSubsidyCurrentDistributionAmount decimal.Decimal       `json:"stakeSubsidyCurrentDistributionAmount"`
	StakeSubsidyPeriodLength              int64                 `json:"stakeSubsidyPeriodLength"`
	StakeSubsidyDecreaseRate              int64                 `json:"stakeSubsidyDecreaseRate"`
	TotalStake                            decimal.Decimal       `json:"totalStake"`
	ActiveValidators                      []SuiValidatorSummary `json:"activeValidators"`
	PendingActiveValidatorsId             string                `json:"pendingActiveValidatorsId"`
	PendingActiveValidatorsSize           int64                 `json:"pendingActiveValidatorsSize"`
	PendingRemovals                       []int64               `json:"pendingRemovals"`
	StakingPoolMappingsId                 string                `json:"stakingPoolMappingsId"`
	StakingPoolMappingsSize               int64                 `json:"stakingPoolMappingsSize"`
	InactivePoolsId                       string                `json:"inactivePoolsId"`
	InactivePoolsSize                     int64                 `json:"inactivePoolsSize"`
	ValidatorCandidatesId                 string                `json:"validatorCandidatesId"`
	ValidatorCandidatesSize               int64                 `json:"validatorCandidatesSize"`
	AtRiskValidators                      []interface{}         `json:"atRiskValidators"`
	ValidatorReportRecords                []interface{}         `json:"validatorReportRecords"`
}
