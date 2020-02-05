package main

import (
	"bufio"
	"fmt"
	"web3space/ethermint/components/cosmos-sdk/crypto/keys"
	"web3space/ethermint/components/cosmos-sdk/types"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"web3space/ethermint/components/cosmos-sdk/client/flags"
	"web3space/ethermint/components/cosmos-sdk/client/input"
	emintcrypto "web3space/ethermint/crypto"
)

func exportEthKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export-eth-key <name>",
		Short: "Export an Ethereum private key",
		Long:  `Export an Ethereum private key unencrypted to use in dev tooling **UNSAFE**`,
		Args:  cobra.ExactArgs(1),
		RunE:  runExportCmd,
	}
	return cmd
}

func runExportCmd(cmd *cobra.Command, args []string) error {
	kb, err := keys.NewKeyring(types.DefaultKeyringServiceName,viper.GetString(flags.FlagKeyringBackend), viper.GetString(flags.FlagHome), cmd.InOrStdin())
	if err != nil {
		return fmt.Errorf("new keyring failed: %w", err)
	}

	buf := bufio.NewReader(cmd.InOrStdin())
	decryptPassword := ""
	conf := true
	keyringBackend := viper.GetString(flags.FlagKeyringBackend)
	switch keyringBackend {
	case keys.BackendFile:
		decryptPassword, err = input.GetPassword(
			"**WARNING this is an unsafe way to export your unencrypted private key**\nEnter key password:",
			buf)
	case keys.BackendOS:
		conf, err = input.GetConfirmation(
			"**WARNING** this is an unsafe way to export your unencrypted private key, are you sure?",
			buf)
	}
	if err != nil || !conf {
		return err
	}

	// Exports private key from keybase using password
	privKey, err := kb.ExportPrivateKeyObject(args[0], decryptPassword)
	if err != nil {
		return err
	}

	// Converts key to Ethermint secp256 implementation
	emintKey, ok := privKey.(emintcrypto.PrivKeySecp256k1)
	if !ok {
		return fmt.Errorf("invalid private key type, must be Ethereum key: %T", privKey)
	}

	// Formats key for output
	privB := ethcrypto.FromECDSA(emintKey.ToECDSA())
	keyS := strings.ToUpper(hexutil.Encode(privB)[2:])

	fmt.Println(keyS)

	return nil
}
