package consensus

import (
	amino "web3space/ethermint/components/tendermint/go-amino"
	"web3space/ethermint/components/tendermint/tendermint/types"
)

var cdc = amino.NewCodec()

func init() {
	RegisterMessages(cdc)
	RegisterWALMessages(cdc)
	types.RegisterBlockAmino(cdc)
}
