package app

import (
	"fmt"
	"math/big"
	"time"

	"web3space/ethermint/components/cosmos-sdk/codec"
	"web3space/ethermint/components/cosmos-sdk/store"
	sdk "web3space/ethermint/components/cosmos-sdk/types"
	"web3space/ethermint/components/cosmos-sdk/x/auth"
	"web3space/ethermint/components/cosmos-sdk/x/auth/types"
	"web3space/ethermint/components/cosmos-sdk/x/mock"
	"web3space/ethermint/components/cosmos-sdk/x/params"

	"web3space/ethermint/crypto"
	emint "web3space/ethermint/types"
	evmtypes "web3space/ethermint/x/evm/types"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	abci "web3space/ethermint/components/tendermint/tendermint/abci/types"
	tmcrypto "web3space/ethermint/components/tendermint/tendermint/crypto"
	"os"
	"web3space/ethermint/components/tendermint/tendermint/libs/log"
	dbm "web3space/ethermint/components/tendermint/tm-db"
)

type testSetup struct {
	ctx          sdk.Context
	cdc          *codec.Codec
	accKeeper    auth.AccountKeeper
	supplyKeeper types.SupplyKeeper
	anteHandler  sdk.AnteHandler
}

func newTestSetup() testSetup {
	db := dbm.NewMemDB()
	authCapKey := sdk.NewKVStoreKey("authCapKey")
	keySupply := sdk.NewKVStoreKey("keySupply")
	keyParams := sdk.NewKVStoreKey("params")
	tkeyParams := sdk.NewTransientStoreKey("transient_params")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeIAVL, db)

	if err := ms.LoadLatestVersion(); err != nil {
		fmt.Printf(err.Error() + "\n")
		os.Exit(1)
	}
	cdc := MakeCodec()
	cdc.RegisterConcrete(&sdk.TestMsg{}, "test/TestMsg", nil)

	// Set params keeper and subspaces
	paramsKeeper := params.NewKeeper(cdc, keyParams, tkeyParams)
	authSubspace := paramsKeeper.Subspace(auth.DefaultParamspace)

	ctx := sdk.NewContext(
		ms,
		abci.Header{ChainID: "3", Time: time.Now().UTC()},
		true,
		log.NewNopLogger(),
	)

	// Add keepers
	accKeeper := auth.NewAccountKeeper(cdc, authCapKey, authSubspace, auth.ProtoBaseAccount)
	accKeeper.SetParams(ctx, types.DefaultParams())
	supplyKeeper := mock.NewDummySupplyKeeper(accKeeper)
	anteHandler := NewAnteHandler(accKeeper, supplyKeeper)

	return testSetup{
		ctx:          ctx,
		cdc:          cdc,
		accKeeper:    accKeeper,
		supplyKeeper: supplyKeeper,
		anteHandler:  anteHandler,
	}
}

func newTestMsg(addrs ...sdk.AccAddress) *sdk.TestMsg {
	return sdk.NewTestMsg(addrs...)
}

func newTestCoins() sdk.Coins {
	return sdk.Coins{sdk.NewInt64Coin(emint.DenomDefault, 500000000)}
}

func newTestStdFee() auth.StdFee {
	return auth.NewStdFee(220000, sdk.NewCoins(sdk.NewInt64Coin(emint.DenomDefault, 150)))
}

// GenerateAddress generates an Ethereum address.
func newTestAddrKey() (sdk.AccAddress, tmcrypto.PrivKey) {
	privkey, _ := crypto.GenerateKey()
	addr := ethcrypto.PubkeyToAddress(privkey.ToECDSA().PublicKey)

	return sdk.AccAddress(addr.Bytes()), privkey
}

func newTestSDKTx(
	ctx sdk.Context, msgs []sdk.Msg, privs []tmcrypto.PrivKey,
	accNums []uint64, seqs []uint64, fee auth.StdFee,
) sdk.Tx {

	sigs := make([]auth.StdSignature, len(privs))
	for i, priv := range privs {
		signBytes := auth.StdSignBytes(ctx.ChainID(), accNums[i], seqs[i], fee, msgs, "")

		sig, err := priv.Sign(signBytes)
		if err != nil {
			panic(err)
		}

		sigs[i] = auth.StdSignature{
			PubKey:    priv.PubKey(),
			Signature: sig,
		}
	}

	return auth.NewStdTx(msgs, fee, sigs, "")
}

func newTestEthTx(ctx sdk.Context, msg *evmtypes.EthereumTxMsg, priv tmcrypto.PrivKey) sdk.Tx {
	chainID, ok := new(big.Int).SetString(ctx.ChainID(), 10)
	if !ok {
		panic(fmt.Sprintf("invalid chainID: %s", ctx.ChainID()))
	}

	privkey, ok := priv.(crypto.PrivKeySecp256k1)
	if !ok {
		panic(fmt.Sprintf("invalid private key type: %T", priv))
	}

	msg.Sign(chainID, privkey.ToECDSA())
	return msg
}
