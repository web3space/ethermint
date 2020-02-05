package core

import (
	ctypes "web3space/ethermint/components/tendermint/tendermint/rpc/core/types"
	rpctypes "web3space/ethermint/components/tendermint/tendermint/rpc/lib/types"
	"web3space/ethermint/components/tendermint/tendermint/types"
)

// BroadcastEvidence broadcasts evidence of the misbehavior.
// More: https://tendermint.com/rpc/#/Info/broadcast_evidence
func BroadcastEvidence(ctx *rpctypes.Context, ev types.Evidence) (*ctypes.ResultBroadcastEvidence, error) {
	err := evidencePool.AddEvidence(ev)
	if err != nil {
		return nil, err
	}
	return &ctypes.ResultBroadcastEvidence{Hash: ev.Hash()}, nil
}
