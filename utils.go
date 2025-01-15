package main

import "fyne.io/fyne/v2/widget"

func createBackButton(where func()) *widget.Button {
	return widget.NewButton("Back", where)
}
