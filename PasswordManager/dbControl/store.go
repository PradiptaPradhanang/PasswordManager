package dbControl

import (
	"database/sql"
	"fmt"
	"passmana/backUpDB"
	"passmana/config"
	"passmana/encrypto"
	"passmana/model"

	_ "modernc.org/sqlite" // driver
)

//var Backup = backUpDB.Backup

var DB *sql.DB

func Get() *sql.DB {
	return DB
}
func CreateDatabase() error {
	if DB == nil {
		return fmt.Errorf("DB not initialized")
	}

	// Create table
	_, err := DB.Exec(`
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

	// Create index
	_, err = DB.Exec(`
        CREATE INDEX IF NOT EXISTS idx_platform_username 
        ON creds(platform, username);
    `)
	return err
}

func OpenDatabase(path string) error {
	var err error
	DB, err = sql.Open("sqlite", path)
	return err
}

func CloseDatabase() {
	if DB != nil {
		DB.Close()
	}
}

func ListPassword() ([]model.Cred, error) {
	rows, err := DB.Query(`SELECT username, platform,nonce,cipherpass FROM creds`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var creds []model.Cred
	for rows.Next() {
		var c model.Cred

		if err := rows.Scan(&c.Username, &c.Platform, &c.Nonce, &c.Cipherpass); err != nil {
			return nil, err
		}
		creds = append(creds, c)
	}
	return creds, nil
}
func AddCred(username, platform string, cipherpass []byte) error {
	var (
		encryptedPass []byte
		nonce         []byte
		err           error
	)

	config.UseMasterKey(func(masterKey []byte) {
		encryptedPass, nonce, err = encrypto.Encryption(masterKey, cipherpass)
	})
	if err != nil {
		return err
	}

	_, err = DB.Exec(`INSERT INTO creds(platform, username,nonce,cipherpass) VALUES(?,?,?,?)`,
		platform, username, nonce, encryptedPass)
	cred := model.Cred{
		Username:   username,
		Platform:   platform,
		Nonce:      nonce,
		Cipherpass: encryptedPass,
	}
	// Send backup event
	backUpDB.Backup.Send(backUpDB.BackUpEvent{
		Action: "add",
		Data:   cred,
	})

	return err

}

func UpdateCred(username, platform, cipherpass string) error {
	var (
		encryptedPass []byte
		nonce         []byte
		err           error
	)
	tmpCipherPass := []byte(cipherpass)
	config.UseMasterKey(func(masterKey []byte) {
		encryptedPass, nonce, err = encrypto.Encryption(masterKey, tmpCipherPass)
	})
	if err != nil {
		return err
	}
	res, err := DB.Exec(`UPDATE creds SET nonce=?, cipherpass=? WHERE platform=? AND username=?`,
		nonce, encryptedPass, platform, username)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no matching credential found to update")
	}
	cred := model.Cred{
		Username:   username,
		Platform:   platform,
		Nonce:      nonce,
		Cipherpass: encryptedPass,
	}
	backUpDB.Backup.Send(backUpDB.BackUpEvent{
		Action: "update",
		Data:   cred,
	})

	return nil
}

func DeleteCred(platform, username string) error {
	_, err := DB.Exec(`DELETE FROM creds WHERE platform=? AND username=?`, platform, username)
	cred := model.Cred{
		Username:   username,
		Platform:   platform,
		Nonce:      nil,
		Cipherpass: nil,
	}
	backUpDB.Backup.Send(backUpDB.BackUpEvent{
		Action: "delete",
		Data:   cred,
	})
	return err
}
