package conn

import (
	amino "web3space/ethermint/components/tendermint/go-amino"
	cryptoamino "web3space/ethermint/components/tendermint/tendermint/crypto/encoding/amino"
)

var cdc *amino.Codec = amino.NewCodec()

func init() {
	cryptoamino.RegisterAmino(cdc)
	RegisterPacket(cdc)
}
