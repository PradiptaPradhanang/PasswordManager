package dbControl

import (
	"database/sql"
	"fmt"
	"passmana/config"
	"passmana/encrypto"

	_ "modernc.org/sqlite" // driver
)

type Cred struct {
	Username   string
	Platform   string
	Nonce      []byte
	Cipherpass []byte
}

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

func ListPassword() ([]Cred, error) {
	rows, err := DB.Query(`SELECT username, platform,nonce,cipherpass FROM creds`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var creds []Cred
	for rows.Next() {
		var c Cred

		if err := rows.Scan(&c.Username, &c.Platform, &c.Nonce, &c.Cipherpass); err != nil {
			return nil, err
		}
		creds = append(creds, c)
	}
	return creds, nil
}
func AddCred(username, platform string, cipherpass []byte) error {

	masterKey := config.GetMasterKey()
	encryptedPass, nonce, err := encrypto.Encryption(masterKey, cipherpass)
	if err != nil {
		return err
	}
	_, err = DB.Exec(`INSERT INTO creds(platform, username,nonce,cipherpass) VALUES(?,?,?,?)`,
		platform, username, nonce, encryptedPass)
	return err
}

func UpdateCred(username, platform, cipherpass string) error {
	tmpCipherPass := []byte(cipherpass)
	masterKey := config.GetMasterKey()
	encryptedPass, nonce, err := encrypto.Encryption(masterKey, tmpCipherPass)
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
	return nil
}

func DeleteCred(platform, username string) error {
	_, err := DB.Exec(`DELETE FROM creds WHERE platform=? AND username=?`, platform, username)
	return err
}
