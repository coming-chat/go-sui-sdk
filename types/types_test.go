package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAddressFromHex(t *testing.T) {
	addr, err := NewAddressFromHex("0x2")
	assert.Nil(t, err)

	t.Log(addr)
}
