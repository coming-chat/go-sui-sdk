package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coming-chat/go-sui/types"
	"io"
	"net/http"
	"strings"
)

const (
	DevNetFaucetUrl = "https://faucet.devnet.sui.io/gas"
)

func FaucetFundAccount(address string, faucetUrl string) (string, error) {
	_, err := types.NewAddressFromHex(address)
	if err != nil {
		return "", err
	}

	var authority string
	if strings.Contains(faucetUrl, "devnet") {
		authority = "faucet.devnet.sui.io"
	} else {
		authority = "faucet.testnet.sui.io"
	}
	paramJson := fmt.Sprintf(`{"FixedAmountRequest":{"recipient":"%v"}}`, address)
	request, err := http.NewRequest(http.MethodPost, faucetUrl, bytes.NewBuffer([]byte(paramJson)))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authority", authority)
	request.Header.Set("authority", authority)
	client := http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return "", err
	}
	response := struct {
		TransferredObjects []struct {
			Amount uint64         `json:"amount"`
			Id     types.ObjectId `json:"id"`
			Digest string         `json:"transfer_tx_digest"`
		} `json:"transferred_gas_objects,omitempty"`
		Error string `json:"error,omitempty"`
	}{}
	var resByte []byte
	if res.StatusCode != 200 && res.StatusCode != 201 {
		return "", fmt.Errorf("post %v response code = %v", faucetUrl, res.Status)
	}
	defer res.Body.Close()
	resByte, err = io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(resByte, &response)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(response.Error) != "" {
		return "", errors.New(response.Error)
	}
	if len(response.TransferredObjects) <= 0 {
		return "", errors.New("transaction not found")
	}

	return response.TransferredObjects[0].Digest, nil
}
