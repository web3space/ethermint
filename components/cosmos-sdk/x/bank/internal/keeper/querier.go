package keeper

import (
	abci "web3space/ethermint/components/tendermint/tendermint/abci/types"

	"web3space/ethermint/components/cosmos-sdk/codec"
	sdk "web3space/ethermint/components/cosmos-sdk/types"
	sdkerrors "web3space/ethermint/components/cosmos-sdk/types/errors"
	"web3space/ethermint/components/cosmos-sdk/x/bank/internal/types"
)

const (
	// query balance path
	QueryBalance = "balances"
)

// NewQuerier returns a new sdk.Keeper instance.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case QueryBalance:
			return queryBalance(ctx, req, k)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

// queryBalance fetch an account's balance for the supplied height.
// Height and account address are passed as first and second path components respectively.
func queryBalance(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryBalanceParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	coins := k.GetCoins(ctx, params.Address)
	if coins == nil {
		coins = sdk.NewCoins()
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, coins)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
