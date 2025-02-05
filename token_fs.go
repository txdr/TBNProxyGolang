package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Token struct {
	Token string `json:"token"`
}

func getToken() string {
	_ = os.Mkdir("token", os.ModePerm)

	if _, err := os.Stat("token/token.json"); os.IsNotExist(err) {
		file, err := os.Create("token/token.json")
		defer file.Close()
		if err != nil {
			panic("Failed to create token file.")
		}

		packed, err := json.Marshal(Token{
			Token: "no token",
		})
		if err != nil {
			fmt.Println("Failed to pack token json (first time).")
		}
		_, err = file.Write(packed)
		if err != nil {
			fmt.Println("Failed to write token file (first time).")
		}

		return "no token"
	}

	contents, err := os.ReadFile("token/token.json")
	if err != nil {
		fmt.Println("Failed to read token file.")
	}
	var token Token
	err = json.Unmarshal(contents, &token)
	if err != nil {
		fmt.Println("Failed to parse token file.")
	}
	return token.Token
}

func updateToken(token string) {
	file, err := os.OpenFile("token/token.json", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic("Failed to open token file.")
	}
	defer file.Close()
	err = file.Truncate(0)
	if err != nil {
		fmt.Println("Failed to truncate token file.")
		return
	}
	packed, err := json.Marshal(Token{
		Token: token,
	})
	if err != nil {
		fmt.Println("Failed to pack token json.")
		return
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		fmt.Println("Failed to seek to the beginning of the file.")
		return
	}
	_, err = file.Write(packed)
	if err != nil {
		fmt.Println("Failed to write token file.")
	}
}
