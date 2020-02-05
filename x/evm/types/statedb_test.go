package types

import (
	"math/big"
	"testing"

	"web3space/ethermint/components/cosmos-sdk/codec"
	"web3space/ethermint/components/cosmos-sdk/store"
	sdkstore "web3space/ethermint/components/cosmos-sdk/store/types"
	sdk "web3space/ethermint/components/cosmos-sdk/types"
	"web3space/ethermint/components/cosmos-sdk/x/auth"
	"web3space/ethermint/components/cosmos-sdk/x/params"
	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"web3space/ethermint/types"

	abci "web3space/ethermint/components/tendermint/tendermint/abci/types"
	tmlog "web3space/ethermint/components/tendermint/tendermint/libs/log"
	dbm "web3space/ethermint/components/tendermint/tm-db"
)

func newTestCodec() *codec.Codec {
	cdc := codec.New()

	RegisterCodec(cdc)
	types.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

func setupStateDB() (*CommitStateDB, error) {
	accKey := sdk.NewKVStoreKey("acc")
	storageKey := sdk.NewKVStoreKey(EvmStoreKey)
	codeKey := sdk.NewKVStoreKey(EvmCodeKey)
	logger := tmlog.NewNopLogger()

	db := dbm.NewMemDB()

	// create logger, codec and root multi-store
	cdc := newTestCodec()
	cms := store.NewCommitMultiStore(db)

	// The ParamsKeeper handles parameter storage for the application
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	paramsKeeper := params.NewKeeper(cdc, keyParams, tkeyParams)
	// Set specific supspaces
	authSubspace := paramsKeeper.Subspace(auth.DefaultParamspace)
	ak := auth.NewAccountKeeper(cdc, accKey, authSubspace, types.ProtoBaseAccount)

	// mount stores
	keys := []*sdk.KVStoreKey{accKey, storageKey, codeKey}
	for _, key := range keys {
		cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, nil)
	}

	cms.SetPruning(sdkstore.PruneNothing)

	// load latest version (root)
	if err := cms.LoadLatestVersion(); err != nil {
		return nil, err
	}

	ms := cms.CacheMultiStore()
	ctx := sdk.NewContext(ms, abci.Header{}, false, logger)
	return NewCommitStateDB(ctx, ak, storageKey, codeKey), nil
}

func TestBloomFilter(t *testing.T) {
	stateDB, err := setupStateDB()
	require.NoError(t, err)

	// Prepare db for logs
	tHash := ethcmn.BytesToHash([]byte{0x1})
	stateDB.Prepare(tHash, common.Hash{}, 0)

	contractAddress := ethcmn.BigToAddress(big.NewInt(1))

	// Generate and add a log to test
	log := ethtypes.Log{Address: contractAddress}
	stateDB.AddLog(&log)

	// Get log from db
	logs := stateDB.GetLogs(tHash)
	require.Equal(t, len(logs), 1)

	// get logs bloom from the log
	bloomInt := ethtypes.LogsBloom(logs)
	bloomFilter := ethtypes.BytesToBloom(bloomInt.Bytes())

	// Check to make sure bloom filter will succeed on
	require.True(t, ethtypes.BloomLookup(bloomFilter, contractAddress))
	require.False(t, ethtypes.BloomLookup(bloomFilter, ethcmn.BigToAddress(big.NewInt(2))))
}
