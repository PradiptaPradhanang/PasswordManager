package appfrontend

import (
	"fmt"
	"os"
	"passmana/config"
	"passmana/dbControl"
	"passmana/utility"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	_ "modernc.org/sqlite" // driver
	//"passmana/cmd
)

var a fyne.App
var w fyne.Window

func EntryPoint() {
	f, _ := os.Create("log.txt")
	fmt.Fprintln(f, "App started")
	a = app.New()
	w = a.NewWindow("PasswordManager")
	w.Resize(fyne.NewSize(460, 340))
	w.CenterOnScreen()
	if err := dbControl.OpenDatabase("vault.db"); err != nil {
		fmt.Println(err) // open or create file
	}

	if err := dbControl.CreateDatabase(); err != nil {
		fmt.Println(err) // ensure schema exists
	}
	defer dbControl.CloseDatabase()
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("Enter first master password")
	// Submit button
	_ = widget.NewButton("Unlock", func() {})
	info1, err1 := os.Stat(config.ChallengFile)
	info2, err2 := os.Stat(config.SaltFile)
	if (os.IsNotExist(err1) && os.IsNotExist(err2)) ||
		(err1 == nil && info1.Size() == 0) ||
		(err2 == nil && info2.Size() == 0) {
		firstTimeSetup()
	} else if err1 != nil {
		fmt.Println(err1)
	} else if err2 != nil {
		fmt.Println(err1)
	} else {
		unlock()
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

func unlock() {
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
		dialog.ShowInformation("Success", "Master password accepted", w)
		passwordScreen()
	})
	// Layout
	w.SetContent(container.NewVBox(
		widget.NewLabel("Welcome to PasswordManager"),
		passwordEntry,
		submit,
	))
}

func passwordScreen() {
	creds, err := dbControl.ListPassword()
	if err != nil {
		dialog.ShowInformation("error loading", err.Error(), w)
	}
	var list *widget.List
	list = widget.NewList(
		func() int { return len(creds) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("username/platform"),
				widget.NewLabel("********"),
				widget.NewButton("Show", nil),
				widget.NewButton("Copy", nil),
				widget.NewButton("Delete", nil),
				widget.NewButton("Add", nil),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			cred := creds[i]
			row := o.(*fyne.Container)
			// update labels
			row.Objects[0].(*widget.Label).SetText(fmt.Sprintf("%s / %s", cred.Username, cred.Platform))
			pwLabel := row.Objects[1].(*widget.Label)

			// update buttons
			row.Objects[2].(*widget.Button).OnTapped = func() {
				pwLabel.SetText(string(cred.Cipherpass))
			}
			row.Objects[3].(*widget.Button).OnTapped = func() {
				a.Clipboard().SetContent(string(cred.Cipherpass))
				dialog.ShowInformation("Copied", "Password copied to clipboard", w)
			}
			row.Objects[4].(*widget.Button).OnTapped = func() {
				dbControl.DeleteCred(cred.Platform, cred.Username)
				dialog.ShowInformation("Delete successful", "Username :"+cred.Username+" Platform"+cred.Platform, w)
				creds, _ = dbControl.ListPassword() // refresh data
				list.Refresh()                      // redraw list
			}
			row.Objects[5].(*widget.Button).OnTapped = func() {
				dialog.ShowInformation("Add", "Add new entry here", w)
			}
		},
	)

	w.SetContent(list)
}
