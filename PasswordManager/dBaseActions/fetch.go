package database

import (
	"database/sql"
	"fmt"
	"passmana/crypto"
	"passmana/database"
)

func decryptBlock(salt []byte, nonce []byte, ciphertext []byte, mpassword string) (p []byte, err error) {
	cfg := crypto.MasterKeyConfig{
		Mpassword: mpassword,
		Salt:      salt,
	}
	MKey, _, err := crypto.DeriveMasterKey(cfg)
	if err != nil {
		return nil, err
	}
	password, err := crypto.Decryption(MKey, nonce, ciphertext)
	if err != nil {
		return nil, err
	}
	return password, nil

}
func Fetch(username string, platform string, mpassword string) {

	db := database.Get()

	row := db.QueryRow(`SELECT salt, nonce, cipherpass FROM creds WHERE platform = ? AND username = ?`, platform, username)

	var salt, nonce, cipherpass []byte
	err := row.Scan(&salt, &nonce, &cipherpass)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No credentials found for this platform and username.")
		} else {
			fmt.Println("Query error:", err)
		}
		return
	}
	password, err := decryptBlock(salt, nonce, cipherpass, mpassword)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Password is %s", password)

}
