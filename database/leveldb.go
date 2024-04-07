package database

import (
    "fmt"
    "log"

    "github.com/syndtr/goleveldb/leveldb"
    "github.com/syndtr/goleveldb/leveldb/util"
)

var db *leveldb.DB

func init() {
    var err error

    db, err = leveldb.OpenFile("./levelDBData", nil)
    if err != nil {
        log.Fatal(err)
    }
}

func AddLevelDBData(key, value []byte) error {
    err := db.Put(key, value, nil)
    if err != nil {
        return fmt.Errorf("error occurred while adding data to LevelDB: %v", err)
    }
    fmt.Println("Data added successfully:", string(value))
    return nil
}

func GetLevelDBData(key []byte) ([]byte, error) {
    data, err := db.Get(key, nil)
    if err != nil {
        return nil, fmt.Errorf("error occurred while getting data from LevelDB: %v", err)
    }
    fmt.Println("Data fetched successfully:",string(data))
    return data, nil
}

func AddDataToLevelDB(value []byte) error {
    iter := db.NewIterator(nil, nil)
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
    iter := db.NewIterator(util.BytesPrefix(nil), nil)
    defer iter.Release()

    for iter.Next() {
        datArray = append(datArray, iter.Value())
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
        if err := db.Delete([]byte(fmt.Sprintf("%d", i)), nil); err != nil {
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
