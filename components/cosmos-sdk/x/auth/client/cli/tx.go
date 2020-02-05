package cli

import (
	"github.com/spf13/cobra"

	"web3space/ethermint/components/cosmos-sdk/client"
	"web3space/ethermint/components/cosmos-sdk/codec"
	"web3space/ethermint/components/cosmos-sdk/x/auth/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Auth transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		GetMultiSignCommand(cdc),
		GetSignCommand(cdc),
	)
	return txCmd
}
