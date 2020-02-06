package types

import (
	"web3space/ethermint/components/cosmos-sdk/codec"
)

// ModuleCdc defines the codec to be used by evm module
var ModuleCdc *codec.Codec// = codec.New()

//func init() {
//	cdc := codec.New()
//
//	codec.RegisterCrypto(cdc)
//
//	ModuleCdc = cdc.Seal()
//}

// RegisterCodec registers concrete types and interfaces on the given codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&EthereumTxMsg{}, "ethermint/MsgEthereumTx", nil)
	cdc.RegisterConcrete(&EmintMsg{}, "ethermint/MsgEmint", nil)

	ModuleCdc = cdc
}
