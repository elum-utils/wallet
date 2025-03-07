package wallet

type Transaction struct {
	Wallet  string `json:"wallet" binding:"required"`      // The recipient wallet address
	Amount  uint64 `json:"amount" binding:"required,gt=0"` // The amount of tokens to withdraw; must be greater than zero
	Message string `json:"message" binding:"required"`     // An optional message or comment for the transaction
}

type TransactionNFT struct {
	AddressNFT    string `json:"address_nft" binding:"require"`
	AddressTarget string `json:"address_target" binding:"require"`
	Message       string `json:"message" binding:"omitempty"`
}
