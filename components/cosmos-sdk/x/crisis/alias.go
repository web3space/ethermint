package crisis

// nolint

import (
	"web3space/ethermint/components/cosmos-sdk/x/crisis/internal/keeper"
	"web3space/ethermint/components/cosmos-sdk/x/crisis/internal/types"
)

const (
	ModuleName           = types.ModuleName
	DefaultParamspace    = types.DefaultParamspace
	EventTypeInvariant   = types.EventTypeInvariant
	AttributeValueCrisis = types.AttributeValueCrisis
	AttributeKeyRoute    = types.AttributeKeyRoute
)

var (
	RegisterCodec            = types.RegisterCodec
	ErrNoSender              = types.ErrNoSender
	ErrUnknownInvariant      = types.ErrUnknownInvariant
	NewGenesisState          = types.NewGenesisState
	DefaultGenesisState      = types.DefaultGenesisState
	NewMsgVerifyInvariant    = types.NewMsgVerifyInvariant
	ParamKeyTable            = types.ParamKeyTable
	NewInvarRoute            = types.NewInvarRoute
	NewKeeper                = keeper.NewKeeper
	ModuleCdc                = types.ModuleCdc
	ParamStoreKeyConstantFee = types.ParamStoreKeyConstantFee
)

type (
	GenesisState       = types.GenesisState
	MsgVerifyInvariant = types.MsgVerifyInvariant
	InvarRoute         = types.InvarRoute
	Keeper             = keeper.Keeper
)
