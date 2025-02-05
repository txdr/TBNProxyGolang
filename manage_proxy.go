package main

import (
	"fmt"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/ZeroErrors/go-bedrockping"
	"golang.org/x/oauth2"
	"time"
)

func openIPInput(account Account) {
	ipEntry := widget.NewEntry()
	ipEntry.SetPlaceHolder("zeqa.net:19132")
	errorInput := widget.NewLabel("...")
	var playButton *widget.Button
	playButton = widget.NewButton("Play", func() {
		playButton.Disable()
		errorInput.SetText("Please wait...")
		response, err := bedrockping.Query(ipEntry.Text, 5*time.Second, 150*time.Millisecond)
		if err != nil {
			errorInput.SetText("Could not ping server. Please try again.")
			var popUp *widget.PopUp
			popUp = widget.NewModalPopUp(container.NewVBox(
				widget.NewLabel("Failed to ping, connect anyways?"),
				widget.NewButton("Yes", func() {
					popUp.Hide()
					playWindow()
					go startProxy(ipEntry.Text, 19132, &oauth2.Token{
						AccessToken: account.AccessCode,
					})
					return
				}),
				widget.NewButton("No", func() {
					popUp.Hide()
					playButton.Enable()
					return
				}),
			), fWindow.Canvas())
			popUp.Show()
			return
		}
		errorInput.SetText(fmt.Sprintf("Successful Ping! (%d/%d Players)", response.PlayerCount, response.MaxPlayers))
		time.Sleep(2 * time.Second)
		playWindow()
		go startProxy(ipEntry.Text, 19132, &oauth2.Token{
			AccessToken: account.AccessCode,
		})
	})

	fWindow.SetContent(container.NewVBox(
		widget.NewLabel("IP input"),
		widget.NewLabel("Please enter a ip address as such: address:port."),
		ipEntry,
		playButton,
		errorInput,
		createBackButton(func() {
			accountManagementUI(account)
		}),
	))
}

func playWindow() {
	killButton := widget.NewButton("Kill Relay", func() {
		killCurrentProxy(false)
		openAccounts()
	})
	killButton.Disable()
	fWindow.SetContent(container.NewVBox(
		widget.NewLabel("Proxy is running."),
		widget.NewLabel("Join at 127.0.0.1:19132, or join from the friends tab."),
		killButton,
	))
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if proxyRunning {
				killButton.Enable()
				ticker.Stop()
			}
		}
	}()
}
