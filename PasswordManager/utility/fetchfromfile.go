package utility

import (
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"os"
	"passmana/config"
	"passmana/crypto"
)

func readCSV() ([][]string, error) {
	file, err := os.Open(config.DBName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file) // initialize a csv reader obj

	return reader.ReadAll() // it returns two values
}

func decodeString(field string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(field)
}

func decryptBlock(entry []string, mpassword string) (p []byte, err error) {
	var salt, nonce, hashPassword []byte
	for index, field := range entry {
		switch index {
		case 2:
			{
				salt, err = decodeString(field)

				if err != nil {
					return nil, err
				}
			}
		case 3:
			{
				nonce, err = decodeString(field)
				if err != nil {
					return nil, err
				}
			}
		case 4:
			{
				hashPassword, err = decodeString(field)
				if err != nil {
					return nil, err
				}
			}

		}
	}
	cfg := crypto.MasterKeyConfig{
		Mpassword: mpassword,
		Salt:      salt,
	}
	MKey, _, err := crypto.DeriveMasterKey(cfg)
	if err != nil {
		return nil, err
	}
	password, err := crypto.Decryption(MKey, nonce, hashPassword)
	if err != nil {
		return nil, err
	}
	return password, nil

}
func Fetchfromfile(username string, platform string, mpassword string) {

	data, err := readCSV()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, slice := range data {
		//fmt.Println(index)
		// for index, field := range slice {
		// 	fmt.Println("hello")
		// 	fmt.Println(index)

		//fmt.Printf("%s,%s", field, platform)
		if slice[1] == platform {
			password, err := decryptBlock(slice, mpassword)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("Password for platform %s is %s", platform, string(password))
			break

		}

	}
}
