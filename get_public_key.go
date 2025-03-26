package wallet

import (
	"context"
	"crypto/ed25519"

	"github.com/xssnick/tonutils-go/address"
)

func (w *Wallet) GetPublicKey(wallet *address.Address) (ed25519.PublicKey, error) {

	block, err := w.Api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		return nil, err
	}

	result, err := w.Api.RunGetMethod(
		context.Background(),
		block,
		wallet,
		"get_public_key",
	)
	if err != nil {
		return nil, err
	}

	slice := result.MustInt(0)
	return slice.Bytes(), nil

}
