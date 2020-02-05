package evm

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	sdk "web3space/ethermint/components/cosmos-sdk/types"
	sdkerrors "web3space/ethermint/components/cosmos-sdk/types/errors"
	authutils "web3space/ethermint/components/cosmos-sdk/x/auth/client/utils"
	emint "web3space/ethermint/types"
	"web3space/ethermint/x/evm/types"

	tm "web3space/ethermint/components/tendermint/tendermint/types"
)

// NewHandler returns a handler for Ethermint type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.EthereumTxMsg:
			return handleETHTxMsg(ctx, keeper, msg)
		case *types.EmintMsg:
			return handleEmintMsg(ctx, keeper, *msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized ethermint Msg type: %v", msg.Type())
			return &sdk.Result{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// Handle an Ethereum specific tx
func handleETHTxMsg(ctx sdk.Context, keeper Keeper, msg types.EthereumTxMsg) (*sdk.Result, error) {
	if err := msg.ValidateBasic(); err != nil {
		return &sdk.Result{}, err
	}

	// parse the chainID from a string to a base-10 integer
	intChainID, ok := new(big.Int).SetString(ctx.ChainID(), 10)
	if !ok {
		return &sdk.Result{}, emint.WrapErrInvalidChainID(fmt.Sprintf("invalid chainID: %s", ctx.ChainID()))
	}

	// Verify signature and retrieve sender address
	sender, err := msg.VerifySig(intChainID)
	if err != nil {
		return &sdk.Result{}, emint.WrapErrInvalidSender(err.Error())
	}

	// Encode transaction by default Tx encoder
	txEncoder := authutils.GetTxEncoder(types.ModuleCdc)
	txBytes, err := txEncoder(msg)
	if err != nil {
		return &sdk.Result{}, emint.WrapErrInternalError(err.Error())
	}
	txHash := tm.Tx(txBytes).Hash()
	ethHash := common.BytesToHash(txHash)

	st := types.StateTransition{
		Sender:       sender,
		AccountNonce: msg.Data.AccountNonce,
		Price:        msg.Data.Price,
		GasLimit:     msg.Data.GasLimit,
		Recipient:    msg.Data.Recipient,
		Amount:       msg.Data.Amount,
		Payload:      msg.Data.Payload,
		Csdb:         keeper.csdb.WithContext(ctx),
		ChainID:      intChainID,
		THash:        &ethHash,
		Simulate:     ctx.IsCheckTx(),
	}
	// Prepare db for logs
	keeper.csdb.Prepare(ethHash, common.Hash{}, keeper.txCount.get())
	keeper.txCount.increment()

	bloom, res, err := st.TransitionCSDB(ctx)
	if err != nil {
		return &res, err
	}
	keeper.bloom.Or(keeper.bloom, bloom)
	return &res, nil
}

func handleEmintMsg(ctx sdk.Context, keeper Keeper, msg types.EmintMsg) (*sdk.Result, error) {
	if err := msg.ValidateBasic(); err != nil {
		return &sdk.Result{}, err
	}

	// parse the chainID from a string to a base-10 integer
	intChainID, ok := new(big.Int).SetString(ctx.ChainID(), 10)
	if !ok {
		return &sdk.Result{}, emint.WrapErrInvalidChainID(fmt.Sprintf("invalid chainID: %s", ctx.ChainID()))
	}

	st := types.StateTransition{
		Sender:       common.BytesToAddress(msg.From.Bytes()),
		AccountNonce: msg.AccountNonce,
		Price:        msg.Price.BigInt(),
		GasLimit:     msg.GasLimit,
		Amount:       msg.Amount.BigInt(),
		Payload:      msg.Payload,
		Csdb:         keeper.csdb.WithContext(ctx),
		ChainID:      intChainID,
		Simulate:     ctx.IsCheckTx(),
	}

	if msg.Recipient != nil {
		to := common.BytesToAddress(msg.Recipient.Bytes())
		st.Recipient = &to
	}

	// Prepare db for logs
	keeper.csdb.Prepare(common.Hash{}, common.Hash{}, keeper.txCount.get()) // Cannot provide tx hash
	keeper.txCount.increment()

	_, res, err := st.TransitionCSDB(ctx)
	return &res, err
}
