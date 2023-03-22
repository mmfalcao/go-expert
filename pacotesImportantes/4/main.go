package main

import (
	"encoding/json"
	"os"
)

type Account struct {
	Number  int `json:"number"`
	Balance int `json:"balance"`
}

func main() {
	account := Account{Number: 1, Balance: 100}
	res, err := json.Marshal(account)
	if err != nil {
		panic(err)
	}
	println(string(res))

	err = json.NewEncoder(os.Stdout).Encode(account)
	if err != nil {
		panic(err)
	}

	jsonDummy := []byte(`{"number":1,"balance":100}`)
	var accountX Account
	err = json.Unmarshal(jsonDummy, &accountX)
	if err != nil {
		panic(err)
	}
	println(accountX.Balance)
}
