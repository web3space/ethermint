package main

import (
	"bufio"
	"fmt"
	"web3space/ethermint/components/cosmos-sdk/types"
	"io"

	"web3space/ethermint/components/cosmos-sdk/client/flags"
	clientkeys "web3space/ethermint/components/cosmos-sdk/client/keys"
	"web3space/ethermint/components/cosmos-sdk/crypto/keys"

	emintCrypto "web3space/ethermint/crypto"
	tmcrypto "github.com/tendermint/tendermint/crypto"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagDryRun = "dry-run"
)

// keyCommands registers a sub-tree of commands to interact with
// local private key storage.
func keyCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keys",
		Short: "Add or view local private keys",
		Long: `Keys allows you to manage your local keystore for tendermint.

    These keys may be in any format supported by go-crypto and can be
    used by light-clients, full nodes, or any other application that
    needs to sign with a private key.`,
	}
	addCmd := clientkeys.AddKeyCommand()
	addCmd.RunE = runAddCmd
	cmd.AddCommand(
		clientkeys.MnemonicKeyCommand(),
		addCmd,
		clientkeys.ExportKeyCommand(),
		clientkeys.ImportKeyCommand(),
		clientkeys.ListKeysCmd(),
		clientkeys.ShowKeysCmd(),
		flags.LineBreak,
		clientkeys.DeleteKeyCommand(),
		clientkeys.UpdateKeyCommand(),
		clientkeys.ParseKeyStringCommand(),
		clientkeys.MigrateCommand(),
		flags.LineBreak,
		exportEthKeyCommand(),
	)
	return cmd
}

func getKeybase(dryrun bool, buf io.Reader) (keys.Keybase, error) {
	if dryrun {
		return keys.NewInMemory(keys.WithKeygenFunc(ethermintKeygenFunc)), nil
	}

	return keys.NewKeyring(types.DefaultKeyringServiceName,viper.GetString(flags.FlagKeyringBackend), viper.GetString(flags.FlagHome), buf, keys.WithKeygenFunc(ethermintKeygenFunc))
}

func runAddCmd(cmd *cobra.Command, args []string) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())
	kb, err := getKeybase(viper.GetBool(flagDryRun), inBuf)
	if err != nil {
		return fmt.Errorf("new keyring failed: %w", err)
	}

	return clientkeys.RunAddCmd(cmd, args, kb, inBuf)
}

//TODO check it
func ethermintKeygenFunc(bz []byte, algo keys.SigningAlgo) (tmcrypto.PrivKey, error) {
	return emintCrypto.PrivKeySecp256k1(bz[:]), nil
}
