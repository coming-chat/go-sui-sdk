package client

import (
	"context"

	"testing"

	"github.com/coming-chat/go-sui/types"
	"github.com/stretchr/testify/require"
)

func TestGetDelegatedStakes(t *testing.T) {
	cli := TestnetClient(t)

	acc := M1Account(t)
	addr, _ := types.NewAddressFromHex(acc.Address)

	res, err := cli.GetDelegatedStakes(context.Background(), *addr)
	require.Nil(t, err)
	t.Log(res)
}

func TestGetValidators(t *testing.T) {
	cli := TestnetClient(t)

	res, err := cli.GetValidators(context.Background())
	require.Nil(t, err)
	// t.Log(res)
	for _, validator := range res {
		if string(validator.Name) == "Chainode Tech" {
			t.Log(validator)
		}
	}
}

func TestGetSuiSystemState(t *testing.T) {
	cli := TestnetClient(t)

	res, err := cli.GetSuiSystemState(context.Background())
	require.Nil(t, err)

	for _, v := range res.Validators.ActiveValidators {
		t.Logf("%v, %v\n", string(v.Metadata.Name), v.CalculateAPY(res.Epoch))
	}
}
