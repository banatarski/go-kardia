package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetAccounts(t *testing.T) {
	accounts := GetAccounts(GenesisAddrKeys1)
	jsonAcc, _ := json.Marshal(accounts)
	fmt.Println(string(jsonAcc))

	accounts = GetAccounts(GenesisAddrKeys2)
	jsonAcc, _ = json.Marshal(accounts)
	fmt.Println(string(jsonAcc))

	accounts = GetAccounts(GenesisAddrKeys3)
	jsonAcc, _ = json.Marshal(accounts)
	fmt.Println(string(jsonAcc))
}