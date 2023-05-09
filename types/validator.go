package types

import (
	"encoding/json"
	"fmt"
	"github.com/coming-chat/go-sui/lib"
	"github.com/coming-chat/go-sui/sui_types"
	"math"
	"reflect"
	"strings"

	"github.com/shopspring/decimal"
)

type StakeStatus = lib.TagJson[Status]

type Status struct {
	Pending *struct{} `json:"Pending,omitempty"`
	Active  *struct {
		EstimatedReward SafeSuiBigInt[uint64] `json:"estimatedReward"`
	} `json:"Active,omitempty"`
	Unstaked *struct{} `json:"Unstaked,omitempty"`
}

func (s Status) Tag() string {
	return "status"
}

func (s Status) Content() string {
	return ""
}

const (
	StakeStatusActive   = "Active"
	StakeStatusPending  = "Pending"
	StakeStatusUnstaked = "Unstaked"
)

type Stake struct {
	StakedSuiId       sui_types.ObjectID     `json:"stakedSuiId"`
	StakeRequestEpoch SafeSuiBigInt[EpochId] `json:"stakeRequestEpoch"`
	StakeActiveEpoch  SafeSuiBigInt[EpochId] `json:"stakeActiveEpoch"`
	Principal         SafeSuiBigInt[uint64]  `json:"principal"`
	StakeStatus       *StakeStatus           `json:"-,flatten"`
}

type JsonFlatten[T Stake] struct {
	Data T
}

func (s *JsonFlatten[T]) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &s.Data)
	if err != nil {
		return err
	}
	rv := reflect.ValueOf(s).Elem().Field(0)
	for i := 0; i < rv.Type().NumField(); i++ {
		tag := rv.Type().Field(i).Tag.Get("json")
		if strings.Contains(tag, "flatten") {
			if rv.Field(i).Kind() != reflect.Pointer {
				return fmt.Errorf("field %s not pointer", rv.Field(i).Type().Name())
			}
			if rv.Field(i).IsNil() {
				rv.Field(i).Set(reflect.New(rv.Field(i).Type().Elem()))
			}
			err = json.Unmarshal(data, rv.Field(i).Interface())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type DelegatedStake struct {
	ValidatorAddress sui_types.SuiAddress `json:"validatorAddress"`
	StakingPool      sui_types.ObjectID   `json:"stakingPool"`
	Stakes           []JsonFlatten[Stake] `json:"stakes"`
}

type SuiValidatorSummary struct {
	SuiAddress             sui_types.SuiAddress `json:"suiAddress"`
	ProtocolPubkeyBytes    lib.Base64Data       `json:"protocolPubkeyBytes"`
	NetworkPubkeyBytes     lib.Base64Data       `json:"networkPubkeyBytes"`
	WorkerPubkeyBytes      lib.Base64Data       `json:"workerPubkeyBytes"`
	ProofOfPossessionBytes lib.Base64Data       `json:"proofOfPossessionBytes"`
	OperationCapId         sui_types.ObjectID   `json:"operationCapId"`
	Name                   string               `json:"name"`
	Description            string               `json:"description"`
	ImageUrl               string               `json:"imageUrl"`
	ProjectUrl             string               `json:"projectUrl"`
	P2pAddress             string               `json:"p2pAddress"`
	NetAddress             string               `json:"netAddress"`
	PrimaryAddress         string               `json:"primaryAddress"`
	WorkerAddress          string               `json:"workerAddress"`

	NextEpochProtocolPubkeyBytes lib.Base64Data `json:"nextEpochProtocolPubkeyBytes"`
	NextEpochProofOfPossession   lib.Base64Data `json:"nextEpochProofOfPossession"`
	NextEpochNetworkPubkeyBytes  lib.Base64Data `json:"nextEpochNetworkPubkeyBytes"`
	NextEpochWorkerPubkeyBytes   lib.Base64Data `json:"nextEpochWorkerPubkeyBytes"`
	NextEpochNetAddress          string         `json:"nextEpochNetAddress"`
	NextEpochP2pAddress          string         `json:"nextEpochP2pAddress"`
	NextEpochPrimaryAddress      string         `json:"nextEpochPrimaryAddress"`
	NextEpochWorkerAddress       string         `json:"nextEpochWorkerAddress"`

	VotingPower             SafeSuiBigInt[uint64] `json:"votingPower"`
	GasPrice                SafeSuiBigInt[uint64] `json:"gasPrice"`
	CommissionRate          SafeSuiBigInt[uint64] `json:"commissionRate"`
	NextEpochStake          SafeSuiBigInt[uint64] `json:"nextEpochStake"`
	NextEpochGasPrice       SafeSuiBigInt[uint64] `json:"nextEpochGasPrice"`
	NextEpochCommissionRate SafeSuiBigInt[uint64] `json:"nextEpochCommissionRate"`
	StakingPoolId           sui_types.ObjectID    `json:"stakingPoolId"`

	StakingPoolActivationEpoch   SafeSuiBigInt[uint64] `json:"stakingPoolActivationEpoch"`
	StakingPoolDeactivationEpoch SafeSuiBigInt[uint64] `json:"stakingPoolDeactivationEpoch"`

	StakingPoolSuiBalance    SafeSuiBigInt[uint64] `json:"stakingPoolSuiBalance"`
	RewardsPool              SafeSuiBigInt[uint64] `json:"rewardsPool"`
	PoolTokenBalance         SafeSuiBigInt[uint64] `json:"poolTokenBalance"`
	PendingStake             SafeSuiBigInt[uint64] `json:"pendingStake"`
	PendingPoolTokenWithdraw SafeSuiBigInt[uint64] `json:"pendingPoolTokenWithdraw"`
	PendingTotalSuiWithdraw  SafeSuiBigInt[uint64] `json:"pendingTotalSuiWithdraw"`
	ExchangeRatesId          sui_types.ObjectID    `json:"exchangeRatesId"`
	ExchangeRatesSize        SafeSuiBigInt[uint64] `json:"exchangeRatesSize"`
}

func (v *SuiValidatorSummary) CalculateAPY(epoch uint64) float64 {
	var (
		stakingPoolSuiBalance      = v.StakingPoolSuiBalance
		stakingPoolActivationEpoch = v.StakingPoolActivationEpoch
		poolTokenBalance           = v.PoolTokenBalance
	)

	// If the staking pool is active then we calculate its APY. Or if staking started in epoch 0
	if stakingPoolActivationEpoch.Uint64() == 0 {
		numEpochsParticipated := epoch - stakingPoolActivationEpoch.Uint64()
		pow1, _ := decimal.NewFromInt(stakingPoolSuiBalance.Int64()).Sub(decimal.NewFromInt(stakingPoolSuiBalance.Int64())).
			Div(decimal.NewFromInt(poolTokenBalance.Int64())).Add(
			decimal.
				NewFromInt(1),
		).Float64()
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

type TypeName []sui_types.SuiAddress
type SuiSystemStateSummary struct {
	Epoch                                 SafeSuiBigInt[uint64]   `json:"epoch"`
	ProtocolVersion                       SafeSuiBigInt[uint64]   `json:"protocolVersion"`
	SystemStateVersion                    SafeSuiBigInt[uint64]   `json:"systemStateVersion"`
	StorageFundTotalObjectStorageRebates  SafeSuiBigInt[uint64]   `json:"storageFundTotalObjectStorageRebates"`
	StorageFundNonRefundableBalance       SafeSuiBigInt[uint64]   `json:"storageFundNonRefundableBalance"`
	ReferenceGasPrice                     SafeSuiBigInt[uint64]   `json:"referenceGasPrice"`
	SafeMode                              bool                    `json:"safeMode"`
	SafeModeStorageRewards                SafeSuiBigInt[uint64]   `json:"safeModeStorageRewards"`
	SafeModeComputationRewards            SafeSuiBigInt[uint64]   `json:"safeModeComputationRewards"`
	SafeModeStorageRebates                SafeSuiBigInt[uint64]   `json:"safeModeStorageRebates"`
	SafeModeNonRefundableStorageFee       SafeSuiBigInt[uint64]   `json:"safeModeNonRefundableStorageFee"`
	EpochStartTimestampMs                 SafeSuiBigInt[uint64]   `json:"epochStartTimestampMs"`
	EpochDurationMs                       SafeSuiBigInt[uint64]   `json:"epochDurationMs"`
	StakeSubsidyStartEpoch                SafeSuiBigInt[uint64]   `json:"stakeSubsidyStartEpoch"`
	MaxValidatorCount                     SafeSuiBigInt[uint64]   `json:"maxValidatorCount"`
	MinValidatorJoiningStake              SafeSuiBigInt[uint64]   `json:"minValidatorJoiningStake"`
	ValidatorLowStakeThreshold            SafeSuiBigInt[uint64]   `json:"validatorLowStakeThreshold"`
	ValidatorVeryLowStakeThreshold        SafeSuiBigInt[uint64]   `json:"validatorVeryLowStakeThreshold"`
	ValidatorLowStakeGracePeriod          SafeSuiBigInt[uint64]   `json:"validatorLowStakeGracePeriod"`
	StakeSubsidyBalance                   SafeSuiBigInt[uint64]   `json:"stakeSubsidyBalance"`
	StakeSubsidyDistributionCounter       SafeSuiBigInt[uint64]   `json:"stakeSubsidyDistributionCounter"`
	StakeSubsidyCurrentDistributionAmount SafeSuiBigInt[uint64]   `json:"stakeSubsidyCurrentDistributionAmount"`
	StakeSubsidyPeriodLength              SafeSuiBigInt[uint64]   `json:"stakeSubsidyPeriodLength"`
	StakeSubsidyDecreaseRate              uint16                  `json:"stakeSubsidyDecreaseRate"`
	TotalStake                            SafeSuiBigInt[uint64]   `json:"totalStake"`
	ActiveValidators                      []SuiValidatorSummary   `json:"activeValidators"`
	PendingActiveValidatorsId             sui_types.ObjectID      `json:"pendingActiveValidatorsId"`
	PendingActiveValidatorsSize           SafeSuiBigInt[uint64]   `json:"pendingActiveValidatorsSize"`
	PendingRemovals                       []SafeSuiBigInt[uint64] `json:"pendingRemovals"`
	StakingPoolMappingsId                 sui_types.ObjectID      `json:"stakingPoolMappingsId"`
	StakingPoolMappingsSize               SafeSuiBigInt[uint64]   `json:"stakingPoolMappingsSize"`
	InactivePoolsId                       sui_types.ObjectID      `json:"inactivePoolsId"`
	InactivePoolsSize                     SafeSuiBigInt[uint64]   `json:"inactivePoolsSize"`
	ValidatorCandidatesId                 sui_types.ObjectID      `json:"validatorCandidatesId"`
	ValidatorCandidatesSize               SafeSuiBigInt[uint64]   `json:"validatorCandidatesSize"`
	AtRiskValidators                      interface{}             `json:"atRiskValidators"`
	ValidatorReportRecords                interface{}             `json:"validatorReportRecords"`
}

type ValidatorsApy struct {
	Epoch SafeSuiBigInt[EpochId] `json:"epoch"`
	Apys  []struct {
		Address string  `json:"address"`
		Apy     float64 `json:"apy"`
	} `json:"apys"`
}

func (apys *ValidatorsApy) ApyMap() map[string]float64 {
	res := make(map[string]float64)
	for _, apy := range apys.Apys {
		res[apy.Address] = apy.Apy
	}
	return res
}
