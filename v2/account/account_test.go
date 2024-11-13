package account

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/stretchr/testify/require"
)

var Mnemonic = os.Getenv("WalletSdkTestM1")

func TestMyAccouunt(t *testing.T) {
	account, err := NewAccountWithMnemonic(Mnemonic)
	require.Nil(t, err)

	t.Logf("addr = %v", account.Address)
}

func Test_Signature_Marshal_Unmarshal(t *testing.T) {
	account, err := NewAccountWithMnemonic(Mnemonic)
	require.Nil(t, err)

	msg := "Coming chat is very good jopfpzf"
	msgBytes := []byte(msg)

	signature1, err := account.SignSecureWithoutEncode(msgBytes, sui_types.DefaultIntent())
	require.Nil(t, err)

	marshaedData, err := json.Marshal(signature1)
	require.Nil(t, err)

	var signature2 sui_types.Signature
	err = json.Unmarshal(marshaedData, &signature2)
	require.Nil(t, err)

	require.Equal(t, signature1, signature2)
}
