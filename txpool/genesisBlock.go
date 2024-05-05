package txpool

// import (
// 	"Blockchain_Project/blockchain"
// 	"Blockchain_Project/transaction"
// 	"Blockchain_Project/utils"
// 	"fmt"
// 	"log"
// 	"os"
// 	"strconv"
// 	"sync"
// 	"time"

// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/crypto"
// 	"github.com/ethereum/go-ethereum/rlp"
// 	"github.com/syndtr/goleveldb/leveldb"
// 	"golang.org/x/crypto/sha3"
// )

// var (
// 	blockDB       *leveldb.DB
// 	blockDBNumber *leveldb.DB
// )

// var err error

// func init() {
// 	// Initialize the block database
// 	blockDB, err = leveldb.OpenFile("./levelDB/blockDB", nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	blockDBNumber, err = leveldb.OpenFile("./levelDB/blockDBNumber", nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// // func CreateGenesisBlock() {
// // 	// Define the header for the genesis block
// // 	genesisHeader := &Header{
// // 		ParentHash:       common.Hash{},                                                     // ParentHash is empty for the genesis block
// // 		Miner:            common.HexToAddress("0x0000000000000000000000000000000000000000"), // Miner address for the genesis block
// // 		StateRoot:        common.Hash{},                                                     // StateRoot can be set to an initial state hash
// // 		TransactionsRoot: common.Hash{},                                                     // TransactionsRoot can be set to the root hash of initial transactions (if any)
// // 		Number:           uint64(0),                                                         // Block number for the genesis block
// // 		Timestamp:        uint64(time.Now().Unix()),                                         // Timestamp for the genesis block
// // 		ExtraData:        []byte("Genesis Block"),                                           // ExtraData can contain any additional information
// // 	}

// // 	// Sign the genesis header
// // 	genesisHeader.Miner = blockchain.SignHeader((*blockchain.Header)(genesisHeader))

// // 	// Create the genesis block
// // 	genesisBlock := &blockchain.Block{
// // 		Header:       (*blockchain.Header)(genesisHeader),
// // 		Transactions: []*transaction.SignedTransaction{}, // No transactions in the genesis block
// // 	}

// // 	// Add the genesis block data to the blockDB
// // 	err := AddBlockData(genesisBlock)
// // 	if err != nil {
// // 		panic(err)
// // 	}

// // 	// Store the hash of the genesis block in blockDBNumber
// // 	err = StoreBlockHash(0, genesisBlock)
// // 	if err != nil {
// // 		panic(err)
// // 	}

// //		fmt.Println("Genesis block created and added to LevelDB successfully.")
// //	}

// func CreateGenesisBlock() *Block {
// 	// Create a new block

// 	minerAddress := os.Getenv("MINER_ADDRESS")
// 	stateRoot := utils.StateRoot()
// 	header := Header{
// 		ParentHash:       common.Hash{},
// 		Miner:            common.HexToAddress(minerAddress),
// 		StateRoot:        stateRoot,
// 		TransactionsRoot: common.Hash{},
// 		Number:           0,
// 		Timestamp:        uint64(time.Now().Unix()),
// 	}
// 	header.ExtraData = blockchain.SignHeader(blockchain.Header(header))

// 	block := Block{
// 		Header:       header,
// 		Transactions: []*transaction.SignedTransaction{},
// 	}
// 	fmt.Println("Genesis Block Created")
// 	return &block
// }

// func StoreBlockHash(blockNumber uint64, block *blockchain.Block) error {
// 	// Calculate the hash of the block
// 	blockHash := RlpHash(block)

// 	// Convert the block number to a string
// 	blockNumberStr := strconv.FormatUint(blockNumber, 10)

// 	// Store the block hash in the database
// 	err := blockDBNumber.Put([]byte(blockNumberStr), blockHash[:], nil)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
// func AddBlockData(block *blockchain.Block) error {
// 	blockHash := RlpHash(block)

// 	serializedBlock, err := SerializeBlock(block)
// 	if err != nil {
// 		return err
// 	}

// 	err = blockDB.Put(blockHash[:], serializedBlock, nil)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func SerializeBlock(block *blockchain.Block) ([]byte, error) {
// 	encodedBlock, err := rlp.EncodeToBytes(block)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return encodedBlock, nil
// }

// func DeserializeBlock(encodedBlock []byte) (*Block, error) {
// 	var block Block
// 	err := rlp.DecodeBytes(encodedBlock, &block)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &block, nil
// }

// func RlpHash(x interface{}) (h common.Hash) {
// 	sha := hasherPool.Get().(crypto.KeccakState)
// 	defer hasherPool.Put(sha)
// 	sha.Reset()
// 	rlp.Encode(sha, x)
// 	sha.Read(h[:])

// 	return h
// }

// var hasherPool = sync.Pool{
// 	New: func() interface{} { return sha3.NewLegacyKeccak256() },
// }

// func Close() {
// 	blockDB.Close()
// 	blockDBNumber.Close()
// }
