package node

import (
	amino "web3space/ethermint/components/tendermint/go-amino"
	cryptoamino "web3space/ethermint/components/tendermint/tendermint/crypto/encoding/amino"
)

var cdc = amino.NewCodec()

func init() {
	cryptoamino.RegisterAmino(cdc)
}
