package database

import (
    "log"
    "fmt"
    "github.com/syndtr/goleveldb/leveldb"
    "github.com/syndtr/goleveldb/leveldb/util"
    "encoding/binary"
)

var (
    blockDB       *leveldb.DB
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

// ------------------------ Functions related to blockDB (Cluster 0) ------------------------

func GetCurrentHeight() (int, error) {
    iter := blockDB.NewIterator(nil, nil)
    defer iter.Release()

    height := -1
    for iter.Next() {
        key := binary.BigEndian.Uint32(iter.Key())
        if int(key) > height {
            height = int(key)
        }
    }

    if err := iter.Error(); err != nil {
        return -1, fmt.Errorf("error occurred while iterating over LevelDB: %v", err)
    }

    return height, nil
}

func AddBlockData(blockData []byte) error {
    // Get the current height
    height, err := GetCurrentHeight()
    if err != nil {
        return fmt.Errorf("error occurred while getting current height: %v", err)
    }

    // Increment the height to get the new height
    height++

    // Convert the height to a byte slice
    key := make([]byte, 4)
    binary.BigEndian.PutUint32(key, uint32(height))

    // Add the block data to the database
    err = blockDB.Put(key, blockData, nil)
    if err != nil {
        return fmt.Errorf("error occurred while adding block data to LevelDB: %v", err)
    }

    fmt.Println("Block data added successfully at height:", height)
    return nil
}


// ------------------------ Functios related to transactionDB (Cluster 2) ------------------------

func AddLevelDBData(key, value []byte) error {
    err := transactionDB.Put(key, value, nil)
    if err != nil {
        return fmt.Errorf("error occurred while adding data to LevelDB: %v", err)
    }
    fmt.Println("Data added successfully:", string(value))
    return nil
}

func GetLevelDBData(key []byte) ([]byte, error) {
    data, err := transactionDB.Get(key, nil)
    if err != nil {
        return nil, fmt.Errorf("error occurred while getting data from LevelDB: %v", err)
    }
    fmt.Println("Data fetched successfully:",string(data))
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

// Delete all data from LevelDB (for local use only)
func DeleteAllData() error {
    datArray, err := GetCompleteBlocksDBData()
    if err != nil {
        return err
    }
	fmt.Println("Deleting all data from LevelDB")
    for i := range datArray {
        if err := transactionDB.Delete([]byte(fmt.Sprintf("%d", i)), nil); err != nil {
            return fmt.Errorf("error occurred while deleting data from LevelDB: %v", err)
        }
    }
    return nil
}

// Print all blocks data from LevelDB
func PrintAllData() error {
    datArray, err := GetCompleteBlocksDBData()
    if err != nil {
        return err
    }

    for i, data := range datArray {
        fmt.Printf("Block #%d: %s\n", i, data)
    }
    return nil
}


func Close() {
    blockDB.Close()
    accountDB.Close()
    transactionDB.Close()
}