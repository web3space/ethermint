package iavl

import (
	amino "web3space/ethermint/components/tendermint/go-amino"
)

var cdc = amino.NewCodec()

func init() {
	// NOTE: It's important that there be no conflicts here,
	// as that would change the canonical representations.
	RegisterWire(cdc)
}

func RegisterWire(cdc *amino.Codec) {
	// TODO
}
