package main

import (
	"TBNProxyRewrite/authentication"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"net/url"
	"time"
)

func signUp() {
	username := widget.NewEntry()
	username.SetPlaceHolder("Username")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")

	fWindow.SetContent(container.NewVBox(
		widget.NewLabel("Sign Up"),
		username,
		password,
		widget.NewButton("Confirm", func() {
			// TODO: Store token and also log in.
			openAccounts()
		}),
		createBackButton(signUpLogIn),
	))
}

func logIn() {
	username := widget.NewEntry()
	username.SetPlaceHolder("Username")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")

	fWindow.SetContent(container.NewVBox(
		widget.NewLabel("Log In"),
		username,
		password,
		widget.NewButton("Confirm", func() {
			// TODO: Confirm logged in & store token.
			openAccounts()
		}),
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
	widgets = append(widgets, widget.NewSeparator())
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
		widget.NewLabel(fmt.Sprintf("Please authenticate at %s using the code %s.", verifyInfo.VerificationURI, verifyInfo.UserCode)),
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
			// TODO: Start proxy and open base playing page.
		}),
		widget.NewButton("Delete account", func() {
			deleteAccount(account.Name)
			openAccounts()
		}),
		widget.NewSeparator(),
		createBackButton(openAccounts),
	))
}
