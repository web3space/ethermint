package main

import (
	"web3space/ethermint/components/cosmos-sdk/client/flags"
	"os"
	"path"

	emintapp "web3space/ethermint/app"
	emintcrypto "web3space/ethermint/crypto"
	"web3space/ethermint/rpc"

	"web3space/ethermint/components/tendermint/go-amino"

	"web3space/ethermint/components/cosmos-sdk/client"
	clientkeys "web3space/ethermint/components/cosmos-sdk/client/keys"
	sdkrpc "web3space/ethermint/components/cosmos-sdk/client/rpc"
	cryptokeys "web3space/ethermint/components/cosmos-sdk/crypto/keys"
	sdk "web3space/ethermint/components/cosmos-sdk/types"

	authcmd "web3space/ethermint/components/cosmos-sdk/x/auth/client/cli"
	bankcmd "web3space/ethermint/components/cosmos-sdk/x/bank/client/cli"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmamino "web3space/ethermint/components/tendermint/tendermint/crypto/encoding/amino"
	"web3space/ethermint/components/tendermint/tendermint/libs/cli"
)

func main() {
	cobra.EnableCommandSorting = false

	cdc := emintapp.MakeCodec()

	tmamino.RegisterKeyType(emintcrypto.PubKeySecp256k1{}, emintcrypto.PubKeyAminoName)
	tmamino.RegisterKeyType(emintcrypto.PrivKeySecp256k1{}, emintcrypto.PrivKeyAminoName)

	cryptokeys.CryptoCdc = cdc
	clientkeys.KeysCdc = cdc

	// Read in the configuration file for the sdk
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	rootCmd := &cobra.Command{
		Use:   "emintcli",
		Short: "Ethermint Client",
	}

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(flags.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// Construct Root Command
	rootCmd.AddCommand(
		sdkrpc.StatusCommand(),
		client.ConfigCmd(emintapp.DefaultCLIHome),
		queryCmd(cdc),
		txCmd(cdc),
		rpc.EmintServeCmd(cdc),
		flags.LineBreak,
		keyCommands(),
		flags.LineBreak,
	)

	executor := cli.PrepareMainCmd(rootCmd, "EM", emintapp.DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func queryCmd(cdc *amino.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		authcmd.GetAccountCmd(cdc),
		flags.LineBreak,
		authcmd.QueryTxsByEventsCmd(cdc),
		authcmd.QueryTxCmd(cdc),
		flags.LineBreak,
	)

	// add modules' query commands
	emintapp.ModuleBasics.AddQueryCommands(queryCmd, cdc)

	return queryCmd
}

func txCmd(cdc *amino.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(
		bankcmd.SendTxCmd(cdc),
		flags.LineBreak,
		authcmd.GetSignCommand(cdc),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(cdc),
		authcmd.GetEncodeCommand(cdc),
		flags.LineBreak,
	)

	// add modules' tx commands
	emintapp.ModuleBasics.AddTxCommands(txCmd, cdc)

	return txCmd
}

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}

	cfgFile := path.Join(home, "config", "config.toml")
	if _, err := os.Stat(cfgFile); err == nil {
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}
	if err := viper.BindPFlag(flags.FlagChainID, cmd.PersistentFlags().Lookup(flags.FlagChainID)); err != nil {
		return err
	}
	if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}
