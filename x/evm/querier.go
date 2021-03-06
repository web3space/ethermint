package evm

import (
	"strconv"

	"web3space/ethermint/components/cosmos-sdk/codec"
	sdk "web3space/ethermint/components/cosmos-sdk/types"
	sdkerrors "web3space/ethermint/components/cosmos-sdk/types/errors"
	"web3space/ethermint/utils"
	"web3space/ethermint/version"
	"web3space/ethermint/x/evm/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	abci "web3space/ethermint/components/tendermint/tendermint/abci/types"
)

// Supported endpoints
const (
	QueryProtocolVersion = "protocolVersion"
	QueryBalance         = "balance"
	QueryBlockNumber     = "blockNumber"
	QueryStorage         = "storage"
	QueryCode            = "code"
	QueryNonce           = "nonce"
	QueryHashToHeight    = "hashToHeight"
	QueryTxLogs          = "txLogs"
	QueryLogsBloom       = "logsBloom"
	QueryLogs            = "logs"
	QueryAccount         = "account"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case QueryProtocolVersion:
			return queryProtocolVersion(keeper)
		case QueryBalance:
			return queryBalance(ctx, path, keeper)
		case QueryBlockNumber:
			return queryBlockNumber(ctx, keeper)
		case QueryStorage:
			return queryStorage(ctx, path, keeper)
		case QueryCode:
			return queryCode(ctx, path, keeper)
		case QueryNonce:
			return queryNonce(ctx, path, keeper)
		case QueryHashToHeight:
			return queryHashToHeight(ctx, path, keeper)
		case QueryTxLogs:
			return queryTxLogs(ctx, path, keeper)
		case QueryLogsBloom:
			return queryBlockLogsBloom(ctx, path, keeper)
		case QueryLogs:
			return queryLogs(ctx, keeper)
		case QueryAccount:
			return queryAccount(ctx, path, keeper)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest,"unknown query endpoint")
		}
	}
}

func queryProtocolVersion(keeper Keeper) ([]byte, error) {
	vers := version.ProtocolVersion

	res, err := codec.MarshalJSONIndent(keeper.cdc, hexutil.Uint(vers))
	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}

func queryBalance(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	addr := ethcmn.HexToAddress(path[1])
	balance := keeper.GetBalance(ctx, addr)

	bRes := types.QueryResBalance{Balance: utils.MarshalBigInt(balance)}
	res, err := codec.MarshalJSONIndent(keeper.cdc, bRes)
	if err != nil {
		panic("could not marshal result to JSON: " + err.Error())
	}

	return res, nil
}

func queryBlockNumber(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	num := ctx.BlockHeight()
	bnRes := types.QueryResBlockNumber{Number: num}
	res, err := codec.MarshalJSONIndent(keeper.cdc, bnRes)
	if err != nil {
		panic("could not marshal result to JSON: " + err.Error())
	}

	return res, nil
}

func queryStorage(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	addr := ethcmn.HexToAddress(path[1])
	key := ethcmn.HexToHash(path[2])
	val := keeper.GetState(ctx, addr, key)
	bRes := types.QueryResStorage{Value: val.Bytes()}
	res, err := codec.MarshalJSONIndent(keeper.cdc, bRes)
	if err != nil {
		panic("could not marshal result to JSON: " + err.Error())
	}
	return res, nil
}

func queryCode(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	addr := ethcmn.HexToAddress(path[1])
	code := keeper.GetCode(ctx, addr)
	cRes := types.QueryResCode{Code: code}
	res, err := codec.MarshalJSONIndent(keeper.cdc, cRes)
	if err != nil {
		panic("could not marshal result to JSON: " + err.Error())
	}

	return res, nil
}

func queryNonce(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	addr := ethcmn.HexToAddress(path[1])
	nonce := keeper.GetNonce(ctx, addr)
	nRes := types.QueryResNonce{Nonce: nonce}
	res, err := codec.MarshalJSONIndent(keeper.cdc, nRes)
	if err != nil {
		panic("could not marshal result to JSON: " + err.Error())
	}

	return res, nil
}

func queryHashToHeight(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	blockHash := ethcmn.FromHex(path[1])
	blockNumber := keeper.GetBlockHashMapping(ctx, blockHash)

	bRes := types.QueryResBlockNumber{Number: blockNumber}
	res, err := codec.MarshalJSONIndent(keeper.cdc, bRes)
	if err != nil {
		panic("could not marshal result to JSON: " + err.Error())
	}

	return res, nil
}

func queryBlockLogsBloom(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	num, err := strconv.ParseInt(path[1], 10, 64)
	if err != nil {
		panic("could not unmarshall block number: " + err.Error())
	}

	bloom := keeper.GetBlockBloomMapping(ctx, num)

	bRes := types.QueryBloomFilter{Bloom: bloom}
	res, err := codec.MarshalJSONIndent(keeper.cdc, bRes)
	if err != nil {
		panic("could not marshal result to JSON: " + err.Error())
	}

	return res, nil
}

func queryTxLogs(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	txHash := ethcmn.HexToHash(path[1])
	logs := keeper.GetLogs(ctx, txHash)

	bRes := types.QueryETHLogs{Logs: logs}
	res, err := codec.MarshalJSONIndent(keeper.cdc, bRes)
	if err != nil {
		panic("could not marshal result to JSON: " + err.Error())
	}

	return res, nil
}

func queryLogs(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	logs := keeper.Logs(ctx)

	lRes := types.QueryETHLogs{Logs: logs}
	l, err := codec.MarshalJSONIndent(keeper.cdc, lRes)
	if err != nil {
		panic("could not marshal result to JSON: " + err.Error())
	}
	return l, nil
}

func queryAccount(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	addr := ethcmn.HexToAddress(path[1])
	so := keeper.GetOrNewStateObject(ctx, addr)

	lRes := types.QueryAccount{
		Balance:  utils.MarshalBigInt(so.Balance()),
		CodeHash: so.CodeHash(),
		Nonce:    so.Nonce(),
	}
	l, err := codec.MarshalJSONIndent(keeper.cdc, lRes)
	if err != nil {
		panic("could not marshal result to JSON: " + err.Error())
	}
	return l, nil
}
