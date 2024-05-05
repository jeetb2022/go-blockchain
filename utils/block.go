package utils

import (
	"Blockchain_Project/database"
	"fmt"
)

func ValidateAddress(address [20]byte) (bool, error) {
	_, err := database.GetAccountFromDB(address)

	if err != nil {

		fmt.Println("Invalid account address")
		return false, err
	} else {
		fmt.Println("Valid account address")
		return true, nil
	}
}

