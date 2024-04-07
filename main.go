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


    // Add data to LevelDB with automatic key generation
    if err := database.AddDataToLevelDB([]byte("This is CSE542: Blockchain project!")); err != nil {
        log.Fatal(err)
    }

    if err := database.AddDataToLevelDB([]byte("This is CSE543: Blockchain project!")); err != nil {
        log.Fatal(err)
    }

    if err := database.AddDataToLevelDB([]byte("This is CSE544: Blockchain project!")); err != nil {
        log.Fatal(err)
    }

    /*  
        // CUSTOM KEY ADD & DELETION OF DATA
        
        // Add data to LevelDB with custom key
        if err := database.AddLevelDBData([]byte("CSE543"), []byte("This is CSE544: Blockchain project!")); err != nil {
            log.Fatal(err)
        }

        // Retrieve data from LevelDB
        if _, err := database.GetLevelDBData([]byte("CSE543")); err != nil {
            log.Fatal(err)
        }
    */

    // Print all data from LevelDB
    if err := database.PrintAllData(); err != nil {
        log.Fatal(err)
    }

	// Delete all data from LevelDB
	if err := database.DeleteAllData(); err != nil {
        log.Fatal(err)
    }
}

