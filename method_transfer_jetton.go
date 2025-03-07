package wallet

import (
	"context"
	"encoding/base64"
	"math/big"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

// TransferJetton creates and executes transactions to transfer Jettons
// from the current wallet to specified destination addresses.
// Each transaction may include a message.
// Parameters:
// - jetton: The address of the Jetton smart contract to interact with.
// - fromAddress: The address from which Jettons will be transferred.
// - transactions: A slice of Transaction structures containing destination addresses, amounts, and optional messages.
// Returns:
// - A base64-encoded string of the resulting transaction hash if successful.
// - An error if the transaction fails.
func (w *Wallet) TransferJetton(
	jetton string, // Address of the Jetton smart contract
	fromAddress string, // Sender's wallet address
	transactions []Transaction, // List of transactions to process
) (string, error) {

	var messages []*wallet.Message
	for _, item := range transactions {
		// Parse destination and response addresses from strings to address.Address format
		destinationAddress := address.MustParseAddr(item.Wallet)
		responseAddress := address.MustParseAddr(fromAddress)

		var comment *cell.Cell
		// If a message is provided, create a comment cell to store it
		if item.Message != "" {
			comment = cell.BeginCell().
				MustStoreUInt(0, 32).               // Store a zeroed header or type marker
				MustStoreStringSnake(item.Message). // Store the message in a snake-encoded string format
				EndCell()                           // Finalize the comment cell
		}

		// Construct an internal message for transferring Jettons
		messages = append(messages, &wallet.Message{
			Mode: wallet.PayGasSeparately + wallet.IgnoreErrors, // Set message mode to pay gas separately and allow ignoring errors
			InternalMessage: &tlb.InternalMessage{
				IHRDisabled: true,                          // Disable Instant Hypercube Routing
				Bounce:      true,                          // Allow bounced messages
				DstAddr:     address.MustParseAddr(jetton), // Set destination to the Jetton smart contract address
				Amount:      tlb.MustFromTON("0.1"),        // Amount to cover gas fees
				Body: cell.BeginCell().
					MustStoreUInt(0xf8a7ea5, 32).                           // Store the operation code for Jetton transfer
					MustStoreUInt(0, 64).                                   // Store a zeroed query ID for tracking
					MustStoreBigCoins(new(big.Int).SetUint64(item.Amount)). // Store the amount to transfer
					MustStoreAddr(destinationAddress).                      // Store the destination address for the Jettons
					MustStoreAddr(responseAddress).                         // Store the response address for transaction confirmations
					MustStoreBoolBit(false).                                // No custom payload, indicated by false
					MustStoreCoins(1).                                      // Set coins for the forward amount
					MustStoreBoolBit(true).                                 // Indicate that custom payload follows
					MustStoreRef(comment).
					EndCell(), // Finalize the message body
			},
		})
	}

	// Send the constructed messages and wait for confirmation of the transactions
	tx, _, err := w.SendManyWaitTransaction(context.Background(), messages)
	if err != nil {
		return "", err // Return an error if the transaction fails
	}

	// Encode the transaction hash to a base64 string and return it
	return base64.StdEncoding.EncodeToString(tx.Hash), nil
}
