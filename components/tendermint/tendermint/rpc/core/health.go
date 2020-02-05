package core

import (
	ctypes "web3space/ethermint/components/tendermint/tendermint/rpc/core/types"
	rpctypes "web3space/ethermint/components/tendermint/tendermint/rpc/lib/types"
)

// Health gets node health. Returns empty result (200 OK) on success, no
// response - in case of an error.
// More: https://tendermint.com/rpc/#/Info/health
func Health(ctx *rpctypes.Context) (*ctypes.ResultHealth, error) {
	return &ctypes.ResultHealth{}, nil
}
