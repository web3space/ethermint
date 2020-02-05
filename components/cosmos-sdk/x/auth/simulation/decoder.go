package simulation

import (
	"bytes"
	"fmt"

	tmkv "web3space/ethermint/components/tendermint/tendermint/libs/kv"

	"web3space/ethermint/components/cosmos-sdk/codec"
	"web3space/ethermint/components/cosmos-sdk/x/auth/exported"
	"web3space/ethermint/components/cosmos-sdk/x/auth/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding auth type
func DecodeStore(cdc *codec.Codec, kvA, kvB tmkv.Pair) string {
	switch {
	case bytes.Equal(kvA.Key[:1], types.AddressStoreKeyPrefix):
		var accA, accB exported.Account
		cdc.MustUnmarshalBinaryBare(kvA.Value, &accA)
		cdc.MustUnmarshalBinaryBare(kvB.Value, &accB)
		return fmt.Sprintf("%v\n%v", accA, accB)
	case bytes.Equal(kvA.Key, types.GlobalAccountNumberKey):
		var globalAccNumberA, globalAccNumberB uint64
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &globalAccNumberA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &globalAccNumberB)
		return fmt.Sprintf("GlobalAccNumberA: %d\nGlobalAccNumberB: %d", globalAccNumberA, globalAccNumberB)
	default:
		panic(fmt.Sprintf("invalid account key %X", kvA.Key))
	}
}
