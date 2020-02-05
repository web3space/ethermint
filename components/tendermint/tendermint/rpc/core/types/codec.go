package coretypes

import (
	amino "web3space/ethermint/components/tendermint/go-amino"
	"web3space/ethermint/components/tendermint/tendermint/types"
)

func RegisterAmino(cdc *amino.Codec) {
	types.RegisterEventDatas(cdc)
	types.RegisterBlockAmino(cdc)
}
