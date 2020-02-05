package evidence

import (
	amino "web3space/ethermint/components/tendermint/go-amino"
	cryptoamino "web3space/ethermint/components/tendermint/tendermint/crypto/encoding/amino"
	"web3space/ethermint/components/tendermint/tendermint/types"
)

var cdc = amino.NewCodec()

func init() {
	RegisterMessages(cdc)
	cryptoamino.RegisterAmino(cdc)
	types.RegisterEvidences(cdc)
}

// For testing purposes only
func RegisterMockEvidences() {
	types.RegisterMockEvidences(cdc)
}
