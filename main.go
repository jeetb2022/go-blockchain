package main

import (
    "Blockchain_Project/cli"
    "os"
	"Blockchain_Project/database"
	"log"

)

func main() {
	defer os.Exit(0)
	cmd := cli.Client{}
	cmd.Run()

	data := []byte("This is CSE542: Blockchain project!")

    // Add data to LevelDB
    if err := database.AddLevelDBData([]byte("CSE542"), data); err != nil {
        log.Fatal(err)
    }

    // Retrieve data from LevelDB
    if _, err := database.GetLevelDBData([]byte("CSE542")); err != nil {
        log.Fatal(err)
    }

    // Add data to LevelDB with automatic key generation
    // if err := database.AddDataToLevelDB(data); err != nil {
    //     log.Fatal(err)
    // }

    // Print all data from LevelDB
    if err := database.PrintAllData(); err != nil {
        log.Fatal(err)
    }

	// Delete all data from LevelDB
	// if err := database.DeleteAllData(); err != nil {
    //     log.Fatal(err)
    // }
}

