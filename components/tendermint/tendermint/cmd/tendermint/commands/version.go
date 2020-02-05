package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"web3space/ethermint/components/tendermint/tendermint/version"
)

// VersionCmd ...
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Version)
	},
}
