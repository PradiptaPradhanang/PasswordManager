package database

import (
	"database/sql"
	"fmt"
	"passmana/dbControl"
	"passmana/encrypto"
)

func decryptBlock(salt []byte, nonce []byte, ciphertext []byte, mpassword string) (p []byte, err error) {
	cfg := encrypto.MasterKeyConfig{
		Mpassword: mpassword,
		Salt:      salt,
	}
	MKey, _, err := encrypto.DeriveMasterKey(cfg)
	if err != nil {
		return nil, err
	}
	password, err := encrypto.Decryption(MKey, nonce, ciphertext)
	if err != nil {
		return nil, err
	}
	return password, nil

}
func Fetch(username string, platform string, mpassword string) {

	db := dbControl.Get()

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
