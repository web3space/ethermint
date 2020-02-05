package types

import (
	amino "web3space/ethermint/components/tendermint/go-amino"
	"web3space/ethermint/components/tendermint/tendermint/types"
)

var cdc = amino.NewCodec()

func init() {
	types.RegisterBlockAmino(cdc)
}
