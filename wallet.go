package wallet

import (
	"context"

	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

var Core *Wallet

type Wallet struct {
	*wallet.Wallet                 // Embedded wallet struct from tonutils-go
	Context        context.Context // Execution context for network operations
	Block          *ton.BlockIDExt // Current block information to validate transactions
	Api            ton.APIClientWrapped
}

// New initializes and returns a new Wallet object using the provided seed words and network configuration URL.
// It returns a pointer to a Wallet instance or an error if initialization fails.
func New(words []string, network string) (*Wallet, error) {

	client := liteclient.NewConnectionPool()

	// Retrieve configuration from the URL
	cfg, err := liteclient.GetConfigFromUrl(context.Background(), network)
	if err != nil {
		return nil, err
	}

	// Connect to mainnet lite servers
	err = client.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	// Initialize API client with proof checking and retry capabilities
	api := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	// Bind all requests to a single TON node
	ctx := client.StickyContext(context.Background())

	// Create a new wallet instance using seed words and specific configuration
	w, err := wallet.FromSeed(api, words, wallet.ConfigV5R1Final{
		NetworkGlobalID: wallet.MainnetGlobalID,
	})
	if err != nil {
		return nil, err
	}

	// Get the current masterchain block information
	block, err := api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		return nil, err
	}

	// Return the configured Wallet instance
	Core = &Wallet{
		w,
		ctx,
		block,
		api,
	}

	return Core, nil
}