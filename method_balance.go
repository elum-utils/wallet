package wallet

// Balance retrieves and returns the current balance of the wallet in NanoTON.
// It returns the balance as a uint64 or an error if retrieval fails.
func (w *Wallet) Balance() (uint64, error) {
	// Call GetBalance with the current context and block.
	// This function retrieves the balance information from the blockchain.
	balance, err := w.GetBalance(w.Context, w.Block)
	if err != nil {
		// If there is an error in retrieving the balance, return 0 as balance
		// and pass the error up the call stack for handling.
		return 0, err
	}

	// Convert the balance to NanoTON and return it as a uint64.
	// The balance is wrapped in a structure that allows it to be represented
	// with a method call `Nano()` to get the value in NanoTON.
	return balance.Nano().Uint64(), nil
}
