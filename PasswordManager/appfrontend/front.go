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
	w.ShowAndRun()
	w.SetCloseIntercept(func() {
		masterKey := config.GetMasterKey()
		// Wipe the master key
		for i := range masterKey {
			masterKey[i] = 0
		}
		masterKey = nil

		// Then close the window
		w.Close()
	})

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
		dialog.ShowInformation("Error loading", err.Error(), w)
		return
	}

	rows := len(creds)
	cols := 7 // Username, Platform, Password, Show, Copy, Delete, Edit

	maskedState := make(map[int]bool)
	for i := 0; i < rows; i++ {
		maskedState[i] = true
	}

	var table *widget.Table
	table = widget.NewTable(
		func() (int, int) {
			return rows + 1, cols
		},

		func() fyne.CanvasObject {
			return container.NewStack(
				widget.NewLabel(""),
				widget.NewButtonWithIcon("", nil, nil),
			)
		},

		func(id widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*fyne.Container).Objects[0].(*widget.Label)
			button := obj.(*fyne.Container).Objects[1].(*widget.Button)

			label.TextStyle = fyne.TextStyle{}
			label.SetText("")
			button.SetText("")
			button.SetIcon(nil)
			button.OnTapped = nil

			if id.Row == 0 {
				label.Show()
				button.Hide()
				label.TextStyle = fyne.TextStyle{Bold: true}
				switch id.Col {
				case 0:
					label.SetText("Username")
				case 1:
					label.SetText("Platform")
				case 2:
					label.SetText("Password")
				case 3:
					label.SetText("Show")
				case 4:
					label.SetText("Copy")
				case 5:
					label.SetText("Delete")
				case 6:
					label.SetText("Edit")
				}
				return
			}

			cred := creds[id.Row-1]

			switch id.Col {
			case 0: // Username
				label.Show()
				button.Hide()
				label.SetText(cred.Username)
			case 1: // Platform
				label.Show()
				button.Hide()
				label.SetText(cred.Platform)
			case 2: // Password
				label.Show()
				button.Hide()
				if maskedState[id.Row-1] {
					label.SetText("********")
				} else {
					label.SetText(string(cred.Cipherpass))
				}
			case 3: // Show/Hide
				label.Hide()
				button.Show()
				row := id.Row - 1
				if maskedState[row] {
					button.SetIcon(theme.VisibilityIcon())
				} else {
					button.SetIcon(theme.VisibilityOffIcon())
				}
				button.OnTapped = func() {
					maskedState[row] = !maskedState[row]
					table.Refresh()
				}
			case 4: // Copy
				label.Hide()
				button.Show()
				button.SetIcon(theme.ContentCopyIcon())
				button.OnTapped = func() {
					a.Clipboard().SetContent(string(cred.Cipherpass))
					dialog.ShowInformation("Copied", "Password copied to clipboard", w)
				}
			case 5: // Delete
				label.Hide()
				button.Show()
				button.SetIcon(theme.DeleteIcon())
				button.OnTapped = func() {
					confirm := dialog.NewConfirm("Confirm Delete",
						fmt.Sprintf("Delete credentials for %s on %s?", cred.Username, cred.Platform),
						func(ok bool) {
							if ok {
								dbControl.DeleteCred(cred.Platform, cred.Username)
								dialog.ShowInformation("Deleted", "Entry removed.", w)
								passwordScreen()
							}
						}, w)
					confirm.Show()
				}
			case 6: // Edit
				label.Hide()
				button.Show()
				button.SetIcon(theme.DocumentCreateIcon())
				button.OnTapped = func() {
					platform := widget.NewLabel(cred.Platform)
					username := widget.NewLabel(cred.Username)
					passwordEntry := widget.NewPasswordEntry()
					formItems := []*widget.FormItem{
						{Text: "Username", Widget: username},
						{Text: "Platform", Widget: platform},
						{Text: "Password", Widget: passwordEntry},
					}
					dialog.NewForm("Edit Password", "Save", "Cancel", formItems, func(ok bool) {
						// if ok {
						//     err := dbControl.AddCred(
						//         platformEntry.Text,
						//         usernameEntry.Text,
						//         []byte(passwordEntry.Text),
						//     )
						//     if err != nil {
						//         dialog.ShowError(err, w)
						//         return
						//     }
						//     dialog.ShowInformation("Success", "Credential added.", w)
						//     passwordScreen()
						// }
					}, w).Show()
				}
			}
		},
	)

	// Column widths
	table.SetColumnWidth(0, 120) // Username
	table.SetColumnWidth(1, 120) // Platform
	table.SetColumnWidth(2, 100) // Password
	table.SetColumnWidth(3, 50)  // Show
	table.SetColumnWidth(4, 50)  // Copy
	table.SetColumnWidth(5, 50)  // Delete
	table.SetColumnWidth(6, 50)  // Edit

	// Add button
	addButton := widget.NewButton("➕ Add New Credential", func() {
		platformEntry := widget.NewEntry()
		usernameEntry := widget.NewEntry()
		passwordEntry := widget.NewPasswordEntry()

		//widget.NewPasswordEntry()

		formItems := []*widget.FormItem{
			{Text: "Platform", Widget: platformEntry},
			{Text: "Username", Widget: usernameEntry},
			{Text: "Password", Widget: passwordEntry},
		}

		dialog.NewForm("Add New Credential", "Save", "Cancel", formItems, func(ok bool) {
			// if ok {
			//     err := dbControl.AddCred(
			//         platformEntry.Text,
			//         usernameEntry.Text,
			//         []byte(passwordEntry.Text),
			//     )
			//     if err != nil {
			//         dialog.ShowError(err, w)
			//         return
			//     }
			//     dialog.ShowInformation("Success", "Credential added.", w)
			//     passwordScreen()
			// }
		}, w).Show()
	})

	content := container.NewBorder(nil, addButton, nil, nil, table)
	w.SetContent(content)
}
