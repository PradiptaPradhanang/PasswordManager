package dbControl

import (
	"database/sql"

	"fmt"

	_ "modernc.org/sqlite" // driver
)

type Cred struct {
	Username   string
	Platform   string
	Salt       []byte
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
            platform TEXT NOT NULL UNIQUE,
            salt BLOB NOT NULL,
            nonce BLOB NOT NULL,
            cipherpass BLOB NOT NULL
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
	rows, err := DB.Query(`SELECT username, platform,salt,nonce,cipherpass FROM creds`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var creds []Cred
	for rows.Next() {
		var c Cred

		if err := rows.Scan(&c.Username, &c.Platform, &c.Salt, &c.Nonce, &c.Cipherpass); err != nil {
			return nil, err
		}
		creds = append(creds, c)
	}
	return creds, nil
}
func AddCred(username, platform string, salt, nonce, cipherpass byte) error {
	_, err := DB.Exec(`INSERT INTO creds(platform, username, salt,nonce,cipherpass) VALUES(?,?,?,?,?)`,
		platform, username, cipherpass)
	return err
}

func UpdateCred(platform, username, cipherpass string) error {
	_, err := DB.Exec(`UPDATE creds SET cipherpass=? WHERE platform=? AND username=?`,
		cipherpass, platform, username)
	return err
}

func DeleteCred(platform, username string) error {
	_, err := DB.Exec(`DELETE FROM creds WHERE platform=? AND username=?`, platform, username)
	return err
}
