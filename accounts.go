package main

import (
	"TBNProxyRewrite/authentication"
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type APIRegisterResponse struct {
	Registered bool   `json:"registered"`
	Text       string `json:"text"`
}

type APILoginResponse struct {
	LoggedIn bool   `json:"loggedIn"`
	Token    string `json:"token"`
	Text     string `json:"text"`
}

func signUp() {
	username := widget.NewEntry()
	username.SetPlaceHolder("Username")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")

	errLabel := widget.NewLabel("")

	fWindow.SetContent(container.NewVBox(
		widget.NewLabel("Sign Up"),
		username,
		password,
		widget.NewButton("Confirm", func() {
			usernameText := username.Text
			passwordText := password.Text
			response, err := http.Get(fmt.Sprintf("%s/register/%s/%s", API_ENDPOINT, usernameText, passwordText))
			if err != nil {
				panic(err)
			}
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Failed to read response body for signUp.")
			}
			var registerResponse APIRegisterResponse
			err = json.Unmarshal(body, &registerResponse)
			if err != nil {
				fmt.Println("Failed to parse response body for signUp.")
			}
			if !registerResponse.Registered {
				errLabel.SetText(registerResponse.Text)
				return
			}
			errLabel.SetText("Successfully registered. You can log in.")
		}),
		errLabel,
		createBackButton(signUpLogIn),
	))
}

func logIn() {
	username := widget.NewEntry()
	username.SetPlaceHolder("Username")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")

	errLabel := widget.NewLabel("")

	fWindow.SetContent(container.NewVBox(
		widget.NewLabel("Log In"),
		username,
		password,
		widget.NewButton("Confirm", func() {
			usernameText := username.Text
			passwordText := password.Text
			response, err := http.Get(fmt.Sprintf("%s/login/%s/%s", API_ENDPOINT, usernameText, passwordText))
			if err != nil {
				panic(err)
			}
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println("Failed to read response body for logIn.")
			}
			var logInResponse APILoginResponse
			err = json.Unmarshal(body, &logInResponse)
			if err != nil {
				fmt.Println("Failed to parse response body for logIn.")
			}
			if logInResponse.Token == "" {
				errLabel.SetText(logInResponse.Text)
				return
			}
			updateToken(logInResponse.Token)
			openAccounts()
		}),
		errLabel,
		createBackButton(signUpLogIn),
	))
}

func openAccounts() {
	var widgets []fyne.CanvasObject
	widgets = append(widgets, widget.NewLabel("Account selection & creation"))
	for name, account := range accounts {
		widgets = append(widgets, widget.NewButton(name, func() {
			accountManagementUI(account)
		}))
	}
	widgets = append(widgets, widget.NewButton("Create account", createAccountUI))

	fWindow.SetContent(container.NewVBox(widgets...))

}

func createAccountUI() {
	fWindow.SetContent(container.NewVBox(widget.NewLabel("Please wait...")))

	verifyInfo, err := authentication.StartDeviceAuth()
	if err != nil {
		fmt.Println("Failed to start device authentication.")
	}
	verifyURL, _ := url.Parse(fmt.Sprintf("%s?otc=%s", verifyInfo.VerificationURI, verifyInfo.UserCode))
	fWindow.SetContent(container.NewVBox(
		widget.NewLabel("Account creation"),
		widget.NewLabel(fmt.Sprintf("Authentication link: %s\nAuthentication code: %s.", verifyInfo.VerificationURI, verifyInfo.UserCode)),
		widget.NewHyperlink("Open Link", verifyURL),
		createBackButton(openAccounts),
	))

	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(verifyInfo.Interval))
		defer ticker.Stop()
		for range ticker.C {
			access, err := authentication.PollDeviceAuth(verifyInfo.DeviceCode)
			if err != nil {
				fmt.Println("Failed to poll device authentication.")
			}
			if access != nil {
				ticker.Stop()
				entry := widget.NewEntry()
				entry.SetPlaceHolder("Account name")
				fWindow.SetContent(container.NewVBox(
					widget.NewLabel("Account creation"),
					widget.NewLabel("Please give your account a name."),
					entry,
					widget.NewButton("Confirm", func() {
						createAccount(entry.Text, access.AccessToken)
						openAccounts()
					}),
				))
			}
		}
	}()
}

func accountManagementUI(account Account) {
	fWindow.SetContent(container.NewVBox(
		widget.NewLabel(fmt.Sprintf("Account management - \"%s\"", account.Name)),
		widget.NewButton("Play with account", func() {
			openIPInput(account)
		}),
		widget.NewButton("Delete account", func() {
			deleteAccount(account.Name)
			openAccounts()
		}),
		createBackButton(openAccounts),
	))
}
