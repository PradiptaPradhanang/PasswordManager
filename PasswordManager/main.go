package main

import (
	"fmt"
	"os"
	"passmana/config"
	"passmana/utility"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	//"passmana/cmd
)

var a fyne.App
var w fyne.Window

func main() {
	//	cmd.Execute()
	f, _ := os.Create("log.txt")
	fmt.Fprintln(f, "App started")
	a = app.New()
	w = a.NewWindow("PasswordManager")
	w.Resize(fyne.NewSize(460, 340))
	w.CenterOnScreen()

	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Enter first master password")
	// Submit button
	_ = widget.NewButton("Unlock", func() {})

	info, err := os.Stat(config.ChallengFile)
	if os.IsNotExist(err) || (err == nil && info.Size() == 0) {
		firstTimeSetup()
	} else if err != nil {
		fmt.Println(err)
	} else {
		dailyUnlock()
	}

	w.ShowAndRun()
}

// -------------ONE-TIME SETUP------------------
func firstTimeSetup() {

	// Master password entry
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Enter first master password")
	// Submit button
	submit := widget.NewButton("Unlock", func() {
		masterPassword := passwordEntry.Text
		if masterPassword == "" {
			dialog.ShowInformation("Error", "Master password cannot be empty", w)
			return
		}
		_, err := os.OpenFile(config.ChallengFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
		utility.CreateVault(masterPassword)

		// TODO: Use masterPassword to derive key and unlock credentials
		dialog.ShowInformation("Success", "Master password accepted", w)
		d := dialog.NewInformation("Ready", "Vault created! Restart the app.", w)
		d.SetOnClosed(func() {
			a.Quit() // quit after OK is clicked
		})
		d.Show()
	})

	// Layout
	w.SetContent(container.NewVBox(
		widget.NewLabel("Welcome to PasswordManager"),
		passwordEntry,
		submit,
	))

}

// -------------DAILY UNLOCK--------------------
func dailyUnlock() {
	// Master password entry
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Enter master password")
	// Submit button
	submit := widget.NewButton("Unlock", func() {
		masterPassword := passwordEntry.Text
		if masterPassword == "" {
			dialog.ShowInformation("Error", "Master password cannot be empty", w)
			return
		}
		if !utility.VerifyPass(masterPassword) {
			dialog.ShowInformation("Error", "Master password is wrong", w)
			return
		}
		// TODO: Use masterPassword to derive key and unlock credentials
		dialog.ShowInformation("Success", "Master password accepted", w)
	})

	// Layout
	w.SetContent(container.NewVBox(
		widget.NewLabel("Welcome to PasswordManager"),
		passwordEntry,
		submit,
	))
}
