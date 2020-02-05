package merkle

import (
	amino "web3space/ethermint/components/tendermint/go-amino"
)

var cdc *amino.Codec

func init() {
	cdc = amino.NewCodec()
	cdc.Seal()
}
