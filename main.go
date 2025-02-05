package main

import "fmt"

const API_ENDPOINT = "http://localhost:8080"

func main() {
	fmt.Println("TBNProxy v1.0.0")

	loadAccounts()

	openUI()
}
