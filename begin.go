package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

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
	// TODO: When backend server is finished, check stored token, if verified skip this page.
	signedIn := false
	if !signedIn {
		fWindow.SetContent(container.NewVBox(
			widget.NewLabel("Welcome to TBNProxy"),
			widget.NewButton("Sign Up", signUp),
			widget.NewButton("Login", logIn),
		))
	} else {
		// Continue to profile page.
	}
}
