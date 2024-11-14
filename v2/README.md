# go-sui-sdk
Sui Golang SDK

[![Documentation (master)](https://img.shields.io/badge/docs-master-59f)](https://github.com/coming-chat/go-sui-sdk)
[![License](https://img.shields.io/badge/license-Apache-green.svg)](https://github.com/coming-chat/go-sui-sdk/blob/main/LICENSE)

The Sui Golang SDK for ComingChat. 
We welcome other developers to participate in the development and testing of sui-sdk.

## Install

```sh
go get github.com/coming-chat/go-sui
```



## Usage

### Account

```go
import "github.com/coming-chat/go-sui/account"

// Import account with mnemonic
acc, err := account.NewAccountWithMnemonic(mnemonic)

// Import account with private key
privateKey, err := hex.DecodeString("4ec5a9eefc0bb86027a6f3ba718793c813505acc25ed09447caf6a069accdd4b")
acc, err := account.NewAccount(privateKey)

// Get private key, public key, address
fmt.Printf("privateKey = %x\n", acc.PrivateKey[:32])
fmt.Printf(" publicKey = %x\n", acc.PublicKey)
fmt.Printf("   address = %v\n", acc.Address)

// Sign data
signedData := acc.Sign(data)
```



### JSON RPC Client

All data interactions on the Sui chain are implemented through the rpc client.

```go
import "github.com/coming-chat/go-sui/client"
import "github.com/coming-chat/go-sui/types"

cli, err := client.Dial(rpcUrl)

// call JSON RPC
responseObject := uint64(0) // if response is a uint64
err := cli.CallContext(ctx, &responseObject, funcName, params...)

// e.g. call get transaction
digest, err := types.NewBase64Data("/KXvTwNRHKKzAB+/Dz1O64LjVbISgIW4VUCmuuPyEfU=")
resp := types.TransactionResponse{}
err := cli.CallContext(ctx, &resp, "sui_getTransaction", digest)
print("transaction status = ", resp.Effects.Status)
print("transaction timestamp = ", resp.TimestampMs)

// And you can call some predefined methods
digest, err := types.NewBase64Data("/KXvTwNRHKKzAB+/Dz1O64LjVbISgIW4VUCmuuPyEfU=")
resp, err := cli.GetTransaction(ctx, digest)
print("transaction status = ", resp.Effects.Status)
print("transaction timestamp = ", resp.TimestampMs)

```

We currently have some rpc methods built-in, [see here](https://github.com/coming-chat/go-sui-sdk/blob/main/client/client_call.go)



### Build Transaction & Sign ( Transfer Sui )

```go
import "github.com/coming-chat/go-sui/client"
import "github.com/coming-chat/go-sui/types"
import "github.com/coming-chat/go-sui/account"

acc, err := account.NewAccountWithMnemonic(mnemonic)
signer, _ := types.NewAddressFromHex(acc.Address)

recipient, err := types.NewAddressFromHex("0x12345678.......")
suiObjectId, err := types.NewHexData("0x36d3176a796e167ffcbd823c94718e7db56b955f")
transferAmount := uint64(10000)
maxGasTransfer := 100

cli, err := client.Dial(rpcUrl)
txnBytes, err := cli.TransferSui(ctx, *signer, *recipient, suiObjectId, transferAmount, maxGasTransfer)

// Sign
signedTxn := txnBytes.SignWith(acc.PrivateKey)

```



### Send Signed Transaction

```go
txnResponse, err := cli.ExecuteTransaction(ctx, signedTxn)

print("transaction digest = ", txnResponse.Certificate.TransactionDigest)
print("transaction status = ", txnResponse.Effects.Status)
print("transaction gasFee = ", txnResponse.Effects.GasFee())
```

