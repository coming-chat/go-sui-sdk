package types

import "errors"

var (
	ErrNoCoinsFound        = errors.New("no coins found")
	ErrInsufficientBalance = errors.New("insufficient account balance")

	ErrNeedMergeCoin    = errors.New("no coins of such a large amount were found to execute this transaction")
	ErrNeedSplitGasCoin = errors.New("missing an extra coin to use as the transaction fee")

	ErrCoinsNotMatchRequest = errors.New("coins not match request")
	ErrCoinsNeedMoreObject  = errors.New("you should get more SUI coins and try again")
)
