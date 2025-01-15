package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Account struct {
	Name       string `json:"name"`
	AccessCode string `json:"access-code"`
}

var accounts = make(map[string]Account)

func loadAccounts() {
	accounts = map[string]Account{}

	err := os.MkdirAll("accounts", os.ModePerm)
	if err != nil {
		fmt.Println("Failed to make accounts folder.")
	}

	files, err := os.ReadDir("accounts")
	if err != nil {
		fmt.Println("Failed to read accounts folder.")
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		contents, err := os.ReadFile(fmt.Sprintf("accounts/%s", file.Name()))
		if err != nil {
			fmt.Println(fmt.Sprintf("Failed to read account data for %s.", file.Name()))
		}
		var account Account
		err = json.Unmarshal(contents, &account)
		accounts[account.Name] = account
	}
}

func createAccount(name, authCode string) {
	account := Account{
		Name:       name,
		AccessCode: authCode,
	}

	file, err := os.Create(fmt.Sprintf("accounts/%s.json", name))
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to create account file for %s.", name))
	}
	defer file.Close()

	packed, err := json.Marshal(account)
	_, err = file.Write(packed)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to write account file for %s.", name))
	}

	accounts[name] = account
}

func deleteAccount(name string) {
	err := os.Remove(fmt.Sprintf("accounts/%s.json", name))
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to delete account file for %s.", name))
	}
	delete(accounts, name)
}
