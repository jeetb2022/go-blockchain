package main

import (
    "Blockchain_Project/cli"
    "os"
	"Blockchain_Project/database"
    // "Blockchain_Project/blockchain"
	"log"
)

func main() {
	defer os.Exit(0)
	cmd := cli.Client{}
	cmd.Run()

    // Add data to LevelDB with automatic key generation
    if err := database.AddDataToLevelDB([]byte("This is CSE542: Blockchain project!")); err != nil {
        log.Fatal(err)
    }

    // Print all data from LevelDB
    if err := database.PrintAllData(); err != nil {
        log.Fatal(err)
    }

	// Delete all data from LevelDB
	// if err := database.DeleteAllData(); err != nil {
    //     log.Fatal(err)
    // }

    if err := database.AddBlockData([]byte("This is CSE542: Blockchain project!")); err != nil {
        log.Fatal(err)
    }

    database.Close()
}

