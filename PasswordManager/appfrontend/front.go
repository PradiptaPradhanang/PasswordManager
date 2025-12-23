package appfrontend

import (
	"fmt"
	"os"
	"passmana/config"
	"passmana/dbControl"
	"passmana/encrypto"
	"passmana/utility"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
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
	w.Resize(fyne.NewSize(500, 500))
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

	w.SetCloseIntercept(func() {
		dialog.NewConfirm("Exit", "Are you sure you want to quit?", func(ok bool) {
			if ok {
				config.ClearMasterKey()
				fmt.Println("clear the masterkey")
				// Then close the window
				w.Close()
			}
		}, w).Show()
	})
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
		utility.CreateVault([]byte(masterPassword))
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
		if !utility.VerifyPass([]byte(masterPassword)) {
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
		dialog.ShowError(err, w)
		return
	}

	var rows []fyne.CanvasObject

	// Header row
	header := container.NewGridWithColumns(7,
		widget.NewLabelWithStyle("Username", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Platform", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Password", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel("Show/Hide"),
		widget.NewLabel("Copy"),
		widget.NewLabel("Delete"),
		widget.NewLabel("Edit"),
	)
	rows = append(rows, header)

	for _, cred := range creds {
		passwordLabel := widget.NewLabel("********")
		showBtn := widget.NewButtonWithIcon("", theme.VisibilityIcon(), nil)
		copyBtn := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), nil)
		deleteBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
		editBtn := widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), nil)

		// Show/hide logic
		showing := false
		showBtn.OnTapped = func() {
			if showing {
				passwordLabel.SetText("********")
				showBtn.SetIcon(theme.VisibilityIcon())
			} else {
				config.UseMasterKey(func(masterKey []byte) {
					password, _ := encrypto.Decryption(masterKey, cred.Nonce, cred.Cipherpass)
					passwordLabel.SetText(string(password))
				})
				showBtn.SetIcon(theme.VisibilityOffIcon())
			}
			showing = !showing
		}

		// Copy logic
		copyBtn.OnTapped = func() {
			a.Clipboard().SetContent(string(cred.Cipherpass))
			dialog.ShowInformation("Copied", "Password copied to clipboard", w)
			// Optional: auto-clear clipboard
			go func() {
				time.Sleep(15 * time.Second)
				a.Clipboard().SetContent("")
				fmt.Println("clipboard cleared")
				//showClipboardWarningModal()
			}()

		}

		// Delete logic
		deleteBtn.OnTapped = func() {
			dialog.NewConfirm("Delete", "Are you sure?", func(ok bool) {
				if ok {
					dbControl.DeleteCred(cred.Platform, cred.Username)
					passwordScreen()
				}
			}, w).Show()
		}

		// Edit logic
		editBtn.OnTapped = func() {
			passwordEntry := widget.NewPasswordEntry()
			//modal design
			formItems := []*widget.FormItem{
				{Text: "Username", Widget: widget.NewLabel(cred.Username)},
				{Text: "Platform", Widget: widget.NewLabel(cred.Platform)},
				{Text: "New Password", Widget: passwordEntry},
			}
			formDialog := dialog.NewForm("Edit Password", "Save", "Cancel", formItems, func(ok bool) {
				if ok {
					password := passwordEntry.Text
					if len(password) < 8 {
						dialog.ShowError(fmt.Errorf("password is less than 8 characters,credentials not added"), w)
						return
					}
					err := dbControl.UpdateCred(cred.Username, cred.Platform, password)
					if err != nil {
						dialog.ShowError(err, w)
						return
					}
					dialog.ShowInformation("Success", "Credential added.", w)
					passwordScreen()
				}
			}, w)

			formDialog.Resize(fyne.NewSize(400, 300))
			formDialog.Show()
		}

		row := container.NewGridWithColumns(7,
			widget.NewLabel(cred.Username),
			widget.NewLabel(cred.Platform),
			passwordLabel,
			showBtn,
			copyBtn,
			deleteBtn,
			editBtn,
		)
		rows = append(rows, row)
	}

	list := container.NewVScroll(container.NewVBox(rows...))
	//scroller
	addButton := widget.NewButton("➕ Add New Credential", func() {
		platformEntry := widget.NewEntry()
		usernameEntry := widget.NewEntry()
		passwordEntry := widget.NewPasswordEntry()

		formItems := []*widget.FormItem{
			{Text: "Platform", Widget: platformEntry},
			{Text: "Username", Widget: usernameEntry},
			{Text: "Password", Widget: passwordEntry},
		}
		formDialog := dialog.NewForm("Add New Credential", "Save", "Cancel", formItems, func(ok bool) {
			if ok {
				platform := strings.TrimSpace(platformEntry.Text)
				username := strings.TrimSpace(usernameEntry.Text)
				password := passwordEntry.Text

				if platform == "" || username == "" {
					dialog.ShowError(fmt.Errorf("please fill all fields correctly,credentials not added"), w)
					return
				}
				if len(password) < 8 {
					dialog.ShowError(fmt.Errorf("password is less than 8 characters,credentials not added"), w)
					return
				}
				err := dbControl.AddCred(username, platform, []byte(password))
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				dialog.ShowInformation("Success", "Credential added.", w)
				passwordScreen()
			}
		}, w)

		formDialog.Resize(fyne.NewSize(400, 300))
		formDialog.Show()
	})

	content := container.NewBorder(nil, addButton, nil, nil, list)
	w.SetContent(content)
}

/*
func showClipboardWarningModal() {
	dialog.ShowConfirm("Security Notice",
		"Clipboard history may still contain your password. Do you want to clear it?",
		func(confirm bool) {
			if confirm {
				go func() {
					err := exec.Command("cmd", "/c", `echo off | clip`).Run()
					if err != nil {
						fmt.Println("Failed to clear clipboard:", err)
					}
				}()
			}
		}, w)
}
*/
