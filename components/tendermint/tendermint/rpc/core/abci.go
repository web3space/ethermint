package core

import (
	abci "web3space/ethermint/components/tendermint/tendermint/abci/types"
	"web3space/ethermint/components/tendermint/tendermint/libs/bytes"
	"web3space/ethermint/components/tendermint/tendermint/proxy"
	ctypes "web3space/ethermint/components/tendermint/tendermint/rpc/core/types"
	rpctypes "web3space/ethermint/components/tendermint/tendermint/rpc/lib/types"
)

// ABCIQuery queries the application for some information.
// More: https://tendermint.com/rpc/#/ABCI/abci_query
func ABCIQuery(
	ctx *rpctypes.Context,
	path string,
	data bytes.HexBytes,
	height int64,
	prove bool,
) (*ctypes.ResultABCIQuery, error) {
	resQuery, err := proxyAppQuery.QuerySync(abci.RequestQuery{
		Path:   path,
		Data:   data,
		Height: height,
		Prove:  prove,
	})
	if err != nil {
		return nil, err
	}
	logger.Info("ABCIQuery", "path", path, "data", data, "result", resQuery)
	return &ctypes.ResultABCIQuery{Response: *resQuery}, nil
}

// ABCIInfo gets some info about the application.
// More: https://tendermint.com/rpc/#/ABCI/abci_info
func ABCIInfo(ctx *rpctypes.Context) (*ctypes.ResultABCIInfo, error) {
	resInfo, err := proxyAppQuery.InfoSync(proxy.RequestInfo)
	if err != nil {
		return nil, err
	}
	return &ctypes.ResultABCIInfo{Response: *resInfo}, nil
}
