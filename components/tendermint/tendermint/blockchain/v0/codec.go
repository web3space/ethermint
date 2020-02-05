package v0

import (
	amino "web3space/ethermint/components/tendermint/go-amino"
	"web3space/ethermint/components/tendermint/tendermint/types"
)

var cdc = amino.NewCodec()

func init() {
	RegisterBlockchainMessages(cdc)
	types.RegisterBlockAmino(cdc)
}
