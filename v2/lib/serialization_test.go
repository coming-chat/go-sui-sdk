package lib

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerialization(t *testing.T) {
	hexStr := "0x12333aabcc"

	hexdata, err := NewHexData(hexStr)
	assert.Nil(t, err)
	assert.Equal(t, hexStr, hexdata.String())

	base64data := Bytes(hexdata.Data()).GetBase64Data()
	base64Str := base64data.String()

	t.Log(base64Str)
	t.Log(hexStr)
}

func TestJson(t *testing.T) {
	hexdata, err := NewHexData("0x12333aabcc")
	assert.Nil(t, err)

	dataJson, err := json.Marshal(hexdata)
	assert.Nil(t, err)

	hexdata2 := HexData{}
	err = json.Unmarshal(dataJson, &hexdata2)
	assert.Nil(t, err)
	assert.Equal(t, hexdata.Data(), hexdata2.Data())

	base64data := Bytes(hexdata.Data()).GetBase64Data()
	dataJsonb, err := json.Marshal(base64data)
	assert.Nil(t, err)

	base64data2 := Base64Data{}
	err = json.Unmarshal(dataJsonb, &base64data2)
	assert.Nil(t, err)
	assert.Equal(t, base64data.Data(), base64data2.Data())
	assert.Equal(t, hexdata.Data(), base64data2.Data())
}
