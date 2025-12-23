package backUpDB

import (
	"database/sql"
	"fmt"
	"passmana/model"

	_ "modernc.org/sqlite"
)

var backupDB *sql.DB

func InitBackupDB(path string) error {
	var err error
	backupDB, err = sql.Open("sqlite", path)
	if err != nil {
		return err
	}

	// Create table
	_, err = backupDB.Exec(`
        CREATE TABLE IF NOT EXISTS creds (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT NOT NULL,
            platform TEXT NOT NULL,
            nonce BLOB NOT NULL,
            cipherpass BLOB NOT NULL,
			UNIQUE(platform, username)
        );
    `)
	if err != nil {
		return err
	}
	return err
}

var Backup *BackupManager

func InitBackup(buffer int) {
	Backup = NewBackUpManager(buffer)
}
func NewBackUpManager(buffer int) *BackupManager {
	bm := &BackupManager{
		events: make(chan BackUpEvent, buffer),
		stop:   make(chan struct{}),
	}
	go bm.listen()
	return bm
}

func (bm *BackupManager) listen() {
	for {
		select {
		case evt := <-bm.events:
			bm.mu.Lock()
			switch evt.Action {
			case "add":
				backupAdd(evt.Data)
			case "update":
				backupUpdate(evt.Data)
			case "delete":
				backupDelete(evt.Data)
			}
			bm.mu.Unlock()
		case <-bm.stop:
			fmt.Println("Stop database tracking.")
			return
		}
	}
}

func backupAdd(cred model.Cred) {
	_, err := backupDB.Exec(
		`INSERT OR REPLACE INTO creds(platform, username, nonce, cipherpass) VALUES(?,?,?,?)`,
		cred.Platform, cred.Username, cred.Nonce, cred.Cipherpass,
	)

	if err != nil {
		fmt.Println("Backup add failed:", err)
	}
}

func backupUpdate(cred model.Cred) {

	_, err := backupDB.Exec(
		`UPDATE creds SET nonce=?, cipherpass=? WHERE platform=? AND username=?`,
		cred.Nonce, cred.Cipherpass, cred.Platform, cred.Username,
	)

	if err != nil {
		fmt.Println("Backup update failed:", err)
	}
}
func backupDelete(cred model.Cred) {
	_, err := backupDB.Exec(
		`DELETE FROM creds WHERE platform=? AND username=?`,
		cred.Platform, cred.Username,
	)

	if err != nil {
		fmt.Println("Backup delete failed:", err)
	}
}
