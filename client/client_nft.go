package client

import (
	"context"

	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
)

// MintNFT
// Create an unsigned transaction to mint a nft at devnet
func (c *Client) MintNFT(
	ctx context.Context,
	signer suiAddress,
	nftName, nftDescription, nftUri string,
	gas *suiObjectID,
	gasBudget uint64,
) (*types.TransactionBytes, error) {
	packageId, _ := sui_types.NewAddressFromHex("0x2")
	args := []any{
		nftName, nftDescription, nftUri,
	}
	return c.MoveCall(
		ctx,
		signer,
		*packageId,
		"devnet_nft",
		"mint",
		[]string{},
		args,
		gas,
		types.NewSafeSuiBigInt(gasBudget),
	)
}

func (c *Client) GetNFTsOwnedByAddress(ctx context.Context, address suiAddress) ([]types.SuiObjectResponse, error) {
	return c.BatchGetObjectsOwnedByAddress(
		ctx, address, types.SuiObjectDataOptions{
			ShowType:    true,
			ShowContent: true,
			ShowOwner:   true,
		}, "0x2::devnet_nft::DevNetNFT",
	)
}
