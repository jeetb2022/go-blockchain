package database

import (
	"Blockchain_Project/account"
	"Blockchain_Project/transaction"
	"encoding/binary"
	"sync"

	// "errors"
	"Blockchain_Project/blockchain"
	"fmt"
	"log"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"golang.org/x/crypto/sha3"
)

var (
	blockDB       *leveldb.DB
	blockDBNumber *leveldb.DB
	accountDB     *leveldb.DB
	transactionDB *leveldb.DB
)

func init() {
	var err error

	// Initialize the block database
	blockDB, err = leveldb.OpenFile("./levelDB/blockDB", nil)
	if err != nil {
		log.Fatal(err)
	}

	blockDBNumber, err = leveldb.OpenFile("./levelDB/blockDBNumber", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the account database
	accountDB, err = leveldb.OpenFile("./levelDB/accountDB", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the transaction database
	transactionDB, err = leveldb.OpenFile("./levelDB/transactionDB", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func SerializeBlock(block *blockchain.Block) ([]byte, error) {
	encodedBlock, err := rlp.EncodeToBytes(block)
	if err != nil {
		return nil, err
	}
	return encodedBlock, nil
}

func DeserializeBlock(encodedBlock []byte) (*blockchain.Block, error) {
	var block blockchain.Block
	err := rlp.DecodeBytes(encodedBlock, &block)
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func rlpHash(x interface{}) (h common.Hash) {
	sha := hasherPool.Get().(crypto.KeccakState)
	defer hasherPool.Put(sha)
	sha.Reset()
	rlp.Encode(sha, x)
	sha.Read(h[:])

	return h
}

var hasherPool = sync.Pool{
	New: func() interface{} { return sha3.NewLegacyKeccak256() },
}

// ------------------------ Functions related to blockDB (Cluster 0) ------------------------

func AddBlockData(block *blockchain.Block) error {
	blockHash := rlpHash(block)

	serializedBlock, err := SerializeBlock(block)
	if err != nil {
		return err
	}

	err = blockDB.Put(blockHash[:], serializedBlock, nil)
	if err != nil {
		return err
	}

	return nil
}

func GetBlockByHash(hash []byte) ([]byte, error) {
	data, err := blockDB.Get(hash, nil)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func PrintAllDataFromBlockDB() error {
	iter := blockDB.NewIterator(util.BytesPrefix(nil), nil)
	defer iter.Release()

	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		fmt.Printf("Key: %s, Value: %s\n", string(key), string(value))
	}

	if err := iter.Error(); err != nil {
		return fmt.Errorf("error occurred while iterating over LevelDB: %v", err)
	}

	return nil
}

// ------------------------ Functions related to blockDBNumber (Cluster 1) ------------------------

func StoreBlockHash(blockNumber uint64, block *blockchain.Block) error {
	// Calculate the hash of the block
	blockHash := rlpHash(block)

	// Convert the block number to a string
	blockNumberStr := strconv.FormatUint(blockNumber, 10)

	// Store the block hash in the database
	err := blockDBNumber.Put([]byte(blockNumberStr), blockHash[:], nil)
	if err != nil {
		return err
	}

	return nil
}

func GetCurrentHeight() (uint64, error) {
	iter := blockDBNumber.NewIterator(nil, nil)
	defer iter.Release()

	height := uint64(0)
	for iter.Next() {
		key := binary.BigEndian.Uint64(iter.Key())
		if int(key) > int(height) {
			height = uint64(key)
		}
	}

	if err := iter.Error(); err != nil {
		return uint64(0), fmt.Errorf("error occurred while iterating over LevelDB: %v", err)
	}

	return height, nil
}

func GetLastBlockHash() (common.Hash, error) {
	// Get the current height
	height, err := GetCurrentHeight()
	if err != nil {
		return common.Hash{}, fmt.Errorf("error occurred while getting current height: %v", err)
	}

	return RetrieveBlockHash(height)
}

func RetrieveBlockHash(blockNumber uint64) (common.Hash, error) {
	// Convert the block number to a string
	blockNumberStr := strconv.FormatUint(blockNumber, 10)

	// Retrieve the block hash from the database
	blockHashBytes, err := GetBlockByHash([]byte(blockNumberStr))
	if err != nil {
		return common.Hash{}, err
	}

	// Convert the byte slice to a common.Hash
	var blockHash common.Hash
	copy(blockHash[:], blockHashBytes)

	return blockHash, nil
}

func PrintAllDataFromBlockDBNumber() error {
	iter := blockDBNumber.NewIterator(util.BytesPrefix(nil), nil)
	defer iter.Release()

	for iter.Next() {
		blockNumber := iter.Key()
		blockHash := iter.Value()

		fmt.Printf("Block Number: %s, Block Hash: %x\n", string(blockNumber), blockHash)
	}

	if err := iter.Error(); err != nil {
		return fmt.Errorf("error occurred while iterating over LevelDB: %v", err)
	}

	return nil
}

// ------------------------ Functios related to transactionDB (Cluster 2) ------------------------

func AddAccountToDB(address [20]byte, account *account.Account) error {
	// Serialize the account object
	serializedAccount, err := rlp.EncodeToBytes(account)
	if err != nil {
		return fmt.Errorf("error serializing account: %v", err)
	}

	// Convert the address to a byte slice
	addressBytes := address[:]

	// Add the serialized account data to the account database
	err = accountDB.Put(addressBytes, serializedAccount, nil)
	if err != nil {
		return fmt.Errorf("error adding account to database: %v", err)
	}

	fmt.Println("Account added successfully:", address)
	return nil
}

func GetAccountFromDB(address [20]byte) (*account.Account, error) {
	addressBytes := address[:]

	serializedAccount, err := accountDB.Get(addressBytes, nil)
	if err != nil {
		return nil, fmt.Errorf("error retrieving account from database: %v", err)
	}

	// Deserialize the account data
	var account account.Account
	err = rlp.DecodeBytes(serializedAccount, &account)
	if err != nil {
		return nil, fmt.Errorf("error deserializing account: %v", err)
	}

	fmt.Println("Account retrieved successfully:", address)
	return &account, nil
}

func PrintAllAccountData() error {
	iter := accountDB.NewIterator(util.BytesPrefix(nil), nil)
	defer iter.Release()

	for iter.Next() {
		address := iter.Key()
		serializedAccount := iter.Value()

		// Deserialize the account data
		var account account.Account
		err := rlp.DecodeBytes(serializedAccount, &account)
		if err != nil {
			return fmt.Errorf("error deserializing account: %v", err)
		}

		fmt.Printf("Address: %x, Account: %+v\n", address, account)
	}

	if err := iter.Error(); err != nil {
		return fmt.Errorf("error occurred while iterating over LevelDB: %v", err)
	}

	return nil
}

// ------------------------ Functios related to transactionDB (Cluster 3) ------------------------

func AddLevelDBData(key, value []byte) error {
	err := transactionDB.Put(key, value, nil)
	if err != nil {
		return fmt.Errorf("error occurred while adding data to LevelDB: %v", err)
	}
	// fmt.Println("Data added successfully:", string(value))
	return nil
}

func GetLevelDBData(key []byte) ([]byte, error) {
	data, err := transactionDB.Get(key, nil)
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting data from LevelDB: %v", err)
	}
	fmt.Println("Data fetched successfully:", string(data))
	return data, nil
}

func AddDataToLevelDB(value []byte) error {
	iter := transactionDB.NewIterator(nil, nil)
	defer iter.Release()

	var i int
	for iter.Next() {
		i++
	}

	key := []byte(fmt.Sprintf("%d", i))
	fmt.Printf("Adding data to levelDB with key:%s & #%d value: %s\n", key, i, string(value))
	return AddLevelDBData(key, value)
}

// Get all blocks data from LevelDB
func GetCompleteBlocksDBData() ([][]byte, error) {
	var datArray [][]byte
	iter := transactionDB.NewIterator(util.BytesPrefix(nil), nil)
	defer iter.Release()

	for iter.Next() {
		data := make([]byte, len(iter.Value()))
		copy(data, iter.Value())
		datArray = append(datArray, data)
	}

	if err := iter.Error(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over LevelDB: %v", err)
	}

	fmt.Println("Getting all data from LevelDB:")
	fmt.Println("Blockchain Length:", len(datArray))
	return datArray, nil
}

func AddSignedTransactionToLevelDB(encodedTx []byte) error {
	return AddDataToLevelDB(encodedTx)
}

func GetSignedTransactionFromLevelDB(key []byte) (*transaction.SignedTransaction, error) {
	data, err := GetLevelDBData(key)
	if err != nil {
		return nil, err
	}

	// Deserialize the transaction
	var tx transaction.SignedTransaction
	err = rlp.DecodeBytes(data, &tx)
	if err != nil {
		return nil, fmt.Errorf("error occurred while decoding transaction: %v", err)
	}

	return &tx, nil
}

func GetCompleteSignedTransactionsDBData() ([]*transaction.SignedTransaction, error) {
	dataArray, err := GetCompleteBlocksDBData()
	if err != nil {
		return nil, err
	}

	var txs []*transaction.SignedTransaction
	for _, data := range dataArray {
		var tx transaction.SignedTransaction
		err = rlp.DecodeBytes(data, &tx)
		if err != nil {
			return nil, fmt.Errorf("error occurred while decoding transaction: %v", err)
		}
		txs = append(txs, &tx)
	}

	return txs, nil
}

func PrintAllData() error {
	dataArray, err := GetCompleteBlocksDBData()
	if err != nil {
		return err
	}

	for i, data := range dataArray {
		var tx transaction.SignedTransaction
		err = rlp.DecodeBytes(data, &tx)
		if err != nil {
			return fmt.Errorf("error occurred while decoding transaction: %v", err)
		}

		fmt.Printf("Block #%d: %+v\n", i, tx)
	}

	return nil
}

func Close() {
	blockDB.Close()
	blockDBNumber.Close()
	accountDB.Close()
	transactionDB.Close()
}
