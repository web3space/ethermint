package client

import (
	"github.com/gorilla/mux"

	"web3space/ethermint/components/cosmos-sdk/client/context"
	"web3space/ethermint/components/cosmos-sdk/client/rpc"
)

// Register routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	rpc.RegisterRPCRoutes(cliCtx, r)
}
