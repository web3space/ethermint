package pex

import (
	amino "web3space/ethermint/components/tendermint/go-amino"
)

var cdc *amino.Codec = amino.NewCodec()

func init() {
	RegisterMessages(cdc)
}
