package types

import (
	"bytes"
	"encoding/json"
	"math"
	"math/big"
)

type Balance struct {
	Value uint64 `json:"Value"`
}

type StakedSui struct {
	Id struct {
		Id ObjectId `json:"id"`
	} `json:"id"`
	ValidatorAddress       Address `json:"validator_address"`
	PoolStartingEpoch      uint64  `json:"pool_starting_epoch"`
	DelegationRequestEpoch uint64  `json:"delegation_request_epoch"`
	Principal              Balance `json:"principal"`
	SuiTokenLock           int     `json:"sui_token_lock,omitempty"`
}

type ActiveFields struct {
	Id struct {
		Id ObjectId `json:"id"`
	} `json:"id"`
	StakedSuiId        ObjectId `json:"staked_sui_id"`
	PrincipalSuiAmount uint64   `json:"principal_sui_amount"`
	PoolTokens         Balance  `json:"pool_tokens"`
}

type ActiveDelegationStatus struct {
	Active ActiveFields `json:"Active"`
}

type DelegatedStake struct {
	StakedSui        StakedSui   `json:"staked_sui"`
	DelegationStatus interface{} `json:"delegation_status"` // "Pending" or ActiveDelegationStatus
}

type ValidatorMetadata struct {
	SuiAddress              string `json:"sui_address"`
	PubkeyBytes             []byte `json:"pubkey_bytes"`
	NetworkPubkeyBytes      []byte `json:"network_pubkey_bytes"`
	WorkerPubkeyBytes       []byte `json:"worker_pubkey_bytes"`
	ProofOfPossessionBytes  []byte `json:"proof_of_possession_bytes"`
	Name                    string `json:"name"`
	Description             string `json:"description"`
	ImageUrl                string `json:"image_url"`
	ProjectUrl              string `json:"project_url"`
	NetAddress              []byte `json:"net_address"`
	ConsensusAddress        []byte `json:"consensus_address"`
	WorkerAddress           []byte `json:"worker_address"`
	NextEpochStake          uint64 `json:"next_epoch_stake"`
	NextEpochDelegation     uint64 `json:"next_epoch_delegation"`
	NextEpochGasPrice       uint64 `json:"next_epoch_gas_price"`
	NextEpochCommissionRate uint64 `json:"next_epoch_commission_rate"`
}

type StakingPool struct {
	ValidatorAddress      Address     `json:"validator_address"`
	StartingEpoch         uint64      `json:"starting_epoch"`
	SuiBalance            uint64      `json:"sui_balance"`
	RewardsPool           Balance     `json:"rewards_pool"`
	DelegationTokenSupply Supply      `json:"delegation_token_supply"`
	PendingDelegations    interface{} `json:"pending_delegations"` //: LinkedTable<ObjectID>,
	PendingWithdraws      interface{} `json:"pending_withdraws"`   //: TableVec,
}

type Validator struct {
	Metadata              ValidatorMetadata `json:"metadata"`
	VotingPower           uint64            `json:"voting_power"`
	StakeAmount           uint64            `json:"stake_amount"`
	PendingStake          uint64            `json:"pending_stake"`
	PendingWithdraw       uint64            `json:"pending_withdraw"`
	GasPrice              uint64            `json:"gas_price"`
	DelegationStakingPool StakingPool       `json:"delegation_staking_pool"`
	CommissionRate        uint64            `json:"commission_rate"`
}

type ValidatorSet struct {
	ValidatorStake            uint64              `json:"validator_stake"`
	DelegationStake           uint64              `json:"delegation_stake"`
	ActiveValidators          []Validator         `json:"active_validators"`
	PendingValidators         []Validator         `json:"pending_validators"`
	PendingRemovals           []uint64            `json:"pending_removals"`
	NextEpochValidators       []ValidatorMetadata `json:"next_epoch_validators"`
	PendingDelegationSwitches interface{}         `json:"pending_delegation_switches"` //: VecMap<ValidatorPair, TableVec>,
}

type SuiSystemState struct {
	Info                   interface{}  `json:"info"` //: UID
	Epoch                  uint64       `json:"epoch"`
	Validators             ValidatorSet `json:"validators"`
	Treasury_cap           Supply       `json:"treasury_cap"`
	StorageFund            Balance      `json:"storage_fund"`
	Parameters             interface{}  `json:"parameters"` //: SystemParameters,
	ReferenceGasPrice      uint64       `json:"reference_gas_price"`
	ValidatorReportRecords interface{}  `json:"validator_report_records"` //: VecMap<SuiAddress, VecSet<SuiAddress>>,
	StakeSubsidy           interface{}  `json:"stake_subsidy"`            //: StakeSubsidy,
	SafeMode               bool         `json:"safe_mode"`
	EpochStartTimestampMs  uint64       `json:"epoch_start_timestamp_ms"`
}

func (v *Validator) CalculateAPY(epoch uint64) float64 {
	p := v.DelegationStakingPool
	if epoch < p.StartingEpoch {
		return 0
	}

	numEpochsParticipated := epoch - p.StartingEpoch
	if p.DelegationTokenSupply.Value == 0 || numEpochsParticipated == 0 {
		return 0
	}
	pow1, _ := big.NewFloat(0).Quo(
		big.NewFloat(float64(p.SuiBalance)),
		big.NewFloat(float64(p.DelegationTokenSupply.Value)),
	).Float64()
	pow2, _ := big.NewFloat(0).Quo(
		big.NewFloat(365.0),
		big.NewFloat(float64(numEpochsParticipated)),
	).Float64()

	apy := math.Pow(pow1, pow2) - 1
	if apy > 100000 {
		return 0
	} else {
		return apy
	}
}

func (d *DelegatedStake) CalculateEarnAmount(activeValidators []Validator) (earn uint64, validator *Validator) {
	for i, v := range activeValidators {
		if bytes.Compare(v.DelegationStakingPool.ValidatorAddress.data, d.StakedSui.ValidatorAddress.data) == 0 {
			validator = &activeValidators[i]
			break
		}
	}
	if validator == nil {
		return
	}

	statusBytes, err := json.Marshal(d.DelegationStatus)
	if err != nil {
		return
	}
	var status ActiveDelegationStatus
	err = json.Unmarshal(statusBytes, &status)
	if err != nil {
		return
	}

	poolTokens := status.Active.PoolTokens.Value
	principal := status.Active.PrincipalSuiAmount
	delegationTokenSupply := validator.DelegationStakingPool.DelegationTokenSupply.Value
	suiBalance := validator.DelegationStakingPool.SuiBalance

	// currentSuiWorth = ┗ poolTokens * suiBalance / delegationTokenSupply ┛  (round down)
	currentSuiWorth := big.NewInt(0).Mul(big.NewInt(int64(poolTokens)), big.NewInt(int64(suiBalance)))
	currentSuiWorth = currentSuiWorth.Quo(currentSuiWorth, big.NewInt(int64(delegationTokenSupply))) // Quo implements truncated division.

	earnInt := big.NewInt(0).Sub(currentSuiWorth, big.NewInt(int64(principal)))
	return earnInt.Uint64(), validator
}
