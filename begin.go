package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"io/ioutil"
	"net/http"
)

type APIVerifyResponse struct {
	Verified bool   `json:"verified"`
	Text     string `json:"text"`
}

var fApp fyne.App
var fWindow fyne.Window

func openUI() {
	fApp = app.New()
	fWindow = fApp.NewWindow("TBNProxy")

	fWindow.CenterOnScreen()
	fWindow.SetFixedSize(true)
	fWindow.Resize(fyne.NewSize(500, 250))

	signUpLogIn()

	fWindow.ShowAndRun()
}

func signUpLogIn() {
	token := getToken()
	response, err := http.Get(fmt.Sprintf("%s/verify/%s", API_ENDPOINT, token))
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read response body for token.")
	}
	var apiVerifyResponse APIVerifyResponse
	err = json.Unmarshal(body, &apiVerifyResponse)
	if err != nil {
		fmt.Println("Failed to parse response body for token verification.")
	}
	if apiVerifyResponse.Verified {
		openAccounts()
		return
	}

	fWindow.SetContent(container.NewVBox(
		widget.NewLabel("Welcome to TBNProxy"),
		widget.NewButton("Sign Up", signUp),
		widget.NewButton("Login", logIn),
	))
}
