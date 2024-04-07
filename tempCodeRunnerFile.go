if err := database.AddLevelDBData([]byte("CSE543"), data); err != nil {
        log.Fatal(err)
    }

    // Retrieve data from LevelDB
    if _, err := database.GetLevelDBData([]byte("CSE543")); err != nil {
        log.Fatal(err)
    }