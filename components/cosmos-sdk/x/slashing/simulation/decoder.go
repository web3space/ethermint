package simulation

import (
	"bytes"
	"fmt"

	"web3space/ethermint/components/tendermint/tendermint/crypto"
	tmkv "web3space/ethermint/components/tendermint/tendermint/libs/kv"

	"web3space/ethermint/components/cosmos-sdk/codec"
	sdk "web3space/ethermint/components/cosmos-sdk/types"
	"web3space/ethermint/components/cosmos-sdk/x/slashing/internal/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding slashing type
func DecodeStore(cdc *codec.Codec, kvA, kvB tmkv.Pair) string {
	switch {
	case bytes.Equal(kvA.Key[:1], types.ValidatorSigningInfoKey):
		var infoA, infoB types.ValidatorSigningInfo
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &infoA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &infoB)
		return fmt.Sprintf("%v\n%v", infoA, infoB)

	case bytes.Equal(kvA.Key[:1], types.ValidatorMissedBlockBitArrayKey):
		var missedA, missedB bool
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &missedA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &missedB)
		return fmt.Sprintf("missedA: %v\nmissedB: %v", missedA, missedB)

	case bytes.Equal(kvA.Key[:1], types.AddrPubkeyRelationKey):
		var pubKeyA, pubKeyB crypto.PubKey
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &pubKeyA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &pubKeyB)
		bechPKA := sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, pubKeyA)
		bechPKB := sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, pubKeyB)
		return fmt.Sprintf("PubKeyA: %s\nPubKeyB: %s", bechPKA, bechPKB)

	default:
		panic(fmt.Sprintf("invalid slashing key prefix %X", kvA.Key[:1]))
	}
}
