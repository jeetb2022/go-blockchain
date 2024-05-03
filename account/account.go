package account

import (
	"encoding/hex"
	"fmt"
)

type Account struct {
	Address [20]byte
	Nonce   uint64
	Balance uint64
}

func ValidateAddress(mineraddr [20]byte) bool {
	myString := hex.EncodeToString(mineraddr[:])
	fmt.Println(myString)
	return true
}
