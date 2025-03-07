package wallet

import (
	"context"
	"encoding/base64"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

// TransferNFT performs the transfer of one or more NFTs from the current owner's wallet
// to new owners specified in the transactions. It also allows for an optional message
// to be attached to each transaction.
// Parameters:
// - responseAddress: The address to receive responses or notifications about the transaction.
// - transactions: A variadic list of transactions containing NFT addresses, target addresses, and optional messages.
// Returns:
// - A base64-encoded string of the transaction hash if successful.
// - An error if the transaction fails.
func (w *Wallet) TransferNFT(
	responseAddress string, // Address for receiving transaction responses
	transactions ...TransactionNFT, // Variadic list of NFT transfer transactions
) (string, error) {

	// Parse the response address from string to address.Address format
	resAddress := address.MustParseAddr(responseAddress)

	var messages []*wallet.Message
	for _, item := range transactions {
		// Parse the destination NFT address from string to address.Address format
		destinationNFTAddress := address.MustParseAddr(item.AddressNFT)
		// Parse the new owner address from string to address.Address format
		newOwnerAddr := address.MustParseAddr(item.AddressTarget)

		// Construct an internal message for the NFT transfer
		messages = append(messages, &wallet.Message{
			Mode: wallet.PayGasSeparately + wallet.IgnoreErrors, // Set message mode to pay gas separately and ignore errors
			InternalMessage: &tlb.InternalMessage{
				IHRDisabled: true,                   // Disable Instant Hypercube Routing
				Bounce:      true,                   // Allow message to be bounced
				DstAddr:     destinationNFTAddress,  // Set the destination address to the NFT address
				Amount:      tlb.MustFromTON("0.1"), // Set the minimum amount to cover gas fees
				// Construct the message body for ownership transfer
				Body: cell.BeginCell().
					MustStoreUInt(0x5fcc3d14, 32). // Store the operation code for NFT transfer
					MustStoreUInt(0, 64).          // Store a zeroed query ID
					MustStoreAddr(newOwnerAddr).   // Store the new owner's address
					MustStoreAddr(resAddress).     // Store the response address
					MustStoreBoolBit(false).       // No custom payload
					MustStoreCoins(1).             // Set forwarding coins
					MustStoreBoolBit(true).        // Custom payload follows
					MustStoreRef(cell.BeginCell().
							MustStoreUInt(0, 32).               // Store a zeroed header or type marker
							MustStoreStringSnake(item.Message). // Store the message in a snake-encoded string format
							EndCell()).                         // Store the comment cell reference
					EndCell(), // Finalize the message body
			},
		})
	}

	// Send the transactions and wait for confirmation
	tx, _, err := w.SendManyWaitTransaction(context.Background(), messages)
	if err != nil {
		return "", err // Return an error if the transaction fails
	}

	// Return the hash of the transaction encoded as a base64 string
	return base64.StdEncoding.EncodeToString(tx.Hash), nil
}
