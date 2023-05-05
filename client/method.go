package client

const (
	SuiXPrefix   = "suix_"
	SuiPrefix    = "sui_"
	UnsafePrefix = "unsafe_"
)

type Method interface {
	String() string
}

type SuiMethod string

func (s SuiMethod) String() string {
	return SuiPrefix + string(s)
}

type SuiXMethod string

func (s SuiXMethod) String() string {
	return SuiXPrefix + string(s)
}

type UnsafeMethod string

func (u UnsafeMethod) String() string {
	return UnsafePrefix + string(u)
}

const (
	devInspectTransactionBlock        SuiMethod    = "devInspectTransactionBlock"
	dryRunTransactionBlock            SuiMethod    = "dryRunTransactionBlock"
	executeTransactionBlock           SuiMethod    = "executeTransactionBlock"
	getCheckpoint                     SuiMethod    = "getCheckpoint"
	getCheckpoints                    SuiMethod    = "getCheckpoints"
	getEvents                         SuiMethod    = "getEvents"
	getLatestCheckpointSequenceNumber SuiMethod    = "getLatestCheckpointSequenceNumber"
	getMoveFunctionArgTypes           SuiMethod    = "getMoveFunctionArgTypes"
	getNormalizedMoveFunction         SuiMethod    = "getNormalizedMoveFunction"
	getNormalizedMoveModule           SuiMethod    = "getNormalizedMoveModule"
	getNormalizedMoveModulesByPackage SuiMethod    = "getNormalizedMoveModulesByPackage"
	getNormalizedMoveStruct           SuiMethod    = "getNormalizedMoveStruct"
	getObject                         SuiMethod    = "getObject"
	getTotalTransactionBlocks         SuiMethod    = "getTotalTransactionBlocks"
	getTransactionBlock               SuiMethod    = "getTransactionBlock"
	multiGetObjects                   SuiMethod    = "multiGetObjects"
	multiGetTransactionBlocks         SuiMethod    = "multiGetTransactionBlocks"
	tryGetPastObject                  SuiMethod    = "tryGetPastObject"
	tryMultiGetPastObjects            SuiMethod    = "tryMultiGetPastObjects"
	getAllBalances                    SuiXMethod   = "getAllBalances"
	getAllCoins                       SuiXMethod   = "getAllCoins"
	getBalance                        SuiXMethod   = "getBalance"
	getCoinMetadata                   SuiXMethod   = "getCoinMetadata"
	getCoins                          SuiXMethod   = "getCoins"
	getCommitteeInfo                  SuiXMethod   = "getCommitteeInfo"
	getCurrentEpoch                   SuiXMethod   = "getCurrentEpoch"
	getDynamicFieldObject             SuiXMethod   = "getDynamicFieldObject"
	getDynamicFields                  SuiXMethod   = "getDynamicFields"
	getEpochs                         SuiXMethod   = "getEpochs"
	getLatestSuiSystemState           SuiXMethod   = "getLatestSuiSystemState"
	getMoveCallMetrics                SuiXMethod   = "getMoveCallMetrics"
	getNetworkMetrics                 SuiXMethod   = "getNetworkMetrics"
	getOwnedObjects                   SuiXMethod   = "getOwnedObjects"
	getReferenceGasPrice              SuiXMethod   = "getReferenceGasPrice"
	getStakes                         SuiXMethod   = "getStakes"
	getStakesByIds                    SuiXMethod   = "getStakesByIds"
	getTotalSupply                    SuiXMethod   = "getTotalSupply"
	getValidatorsApy                  SuiXMethod   = "getValidatorsApy"
	queryEvents                       SuiXMethod   = "queryEvents"
	queryObjects                      SuiXMethod   = "queryObjects"
	queryTransactionBlocks            SuiXMethod   = "queryTransactionBlocks"
	subscribeEvent                    SuiXMethod   = "subscribeEvent"
	batchTransaction                  UnsafeMethod = "batchTransaction"
	mergeCoins                        UnsafeMethod = "mergeCoins"
	moveCall                          UnsafeMethod = "moveCall"
	pay                               UnsafeMethod = "pay"
	payAllSui                         UnsafeMethod = "payAllSui"
	paySui                            UnsafeMethod = "paySui"
	publish                           UnsafeMethod = "publish"
	requestAddStake                   UnsafeMethod = "requestAddStake"
	requestWithdrawStake              UnsafeMethod = "requestWithdrawStake"
	splitCoin                         UnsafeMethod = "splitCoin"
	splitCoinEqual                    UnsafeMethod = "splitCoinEqual"
	transferObject                    UnsafeMethod = "transferObject"
	transferSui                       UnsafeMethod = "transferSui"
)
