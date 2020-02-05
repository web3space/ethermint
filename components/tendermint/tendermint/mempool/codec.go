package mempool

import (
	amino "web3space/ethermint/components/tendermint/go-amino"
)

var cdc = amino.NewCodec()

func init() {
	RegisterMessages(cdc)
}
