package types

import (
	sdk "web3space/ethermint/components/cosmos-sdk/types/errors"
)

// Ethermint error codes
const (
	// DefaultCodespace reserves a Codespace for Ethermint.
	DefaultCodespace  = "ethermint"

	CodeInvalidValue   uint32 = 1
	CodeInvalidChainID uint32 = 2
	CodeInvalidSender  uint32 = 3
	CodeVMExecution    uint32 = 4
	CodeInvalidNonce   uint32 = 5
	CodeInternalError  uint32 = 6
)

// CodeToDefaultMsg takes the CodeType variable and returns the error string
//func CodeToDefaultMsg(code sdk.CodeType) string {
//	switch code {
//	case CodeInvalidValue:
//		return "invalid value"
//	case CodeInvalidChainID:
//		return "invalid chain ID"
//	case CodeInvalidSender:
//		return "could not derive sender from transaction"
//	case CodeVMExecution:
//		return "error while executing evm transaction"
//	case CodeInvalidNonce:
//		return "invalid nonce"
//	default:
//		return sdk.CodeToDefaultMsg(code)
//	}
//}

var (
	ErrInvalidValue = sdk.Register(DefaultCodespace, CodeInvalidValue, "invalid value")
	ErrInvalidChainID = sdk.Register(DefaultCodespace, CodeInvalidChainID, "invalid chainid")
	ErrInvalidSender = sdk.Register(DefaultCodespace, CodeInvalidSender, "invalid sender")
	ErrVMExecution = sdk.Register(DefaultCodespace, CodeVMExecution, "vme execution failed")
	ErrInvalidNonce = sdk.Register(DefaultCodespace, CodeInvalidNonce, "invalid nonce")
	ErrInternalError = sdk.Register(DefaultCodespace, CodeInternalError, "internal errot")
)

// ErrInvalidValue returns a standardized SDK error resulting from an invalid value.
func WrapErrInvalidValue(msg string) error {
	return sdk.Wrap(ErrInvalidValue, msg)
}

// ErrInvalidChainID returns a standardized SDK error resulting from an invalid chain ID.
func WrapErrInvalidChainID(msg string) error {
	return sdk.Wrap(ErrInvalidChainID, msg)
}

// ErrInvalidSender returns a standardized SDK error resulting from an invalid transaction sender.
func WrapErrInvalidSender(msg string) error {
	return sdk.Wrap(ErrInvalidSender, msg)
}

// ErrVMExecution returns a standardized SDK error resulting from an error in EVM execution.
func WrapErrVMExecution(msg string) error {
	return sdk.Wrap(ErrVMExecution, msg)
}

// ErrVMExecution returns a standardized SDK error resulting from an error in EVM execution.
func WrapErrInvalidNonce(msg string) error {
	return sdk.Wrap(ErrInvalidNonce, msg)
}

func WrapErrInternalError(msg string) error {
	return sdk.Wrap(ErrInternalError, msg)
}