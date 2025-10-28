package database

import (
	"fmt"
	"os"
	"passmana/config"

	"golang.org/x/crypto/bcrypt"
)

func Insert(username string, sitename string, password string) {
	//converts string into byte array because bcrypt only accept byte array
	bytePassword := []byte(password)
	hashValue, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := username + "," + sitename + ","
	//open the file in WRITE mode, if file is not present then it will create
	input, err := os.OpenFile(config.DBName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	//adding platform and username
	//l1 is size of the string passed
	metadata, err := input.WriteString(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	//adding password
	hash, err := input.WriteString(string(hashValue)) //string(hash) function save the hash int the string format
	if err != nil {
		fmt.Println(err)
		return
	}
	//adding new line
	l3, err := input.WriteString("\n")
	if err != nil {
		fmt.Println(err)
		return
	}
	//check the data was stored or not
	if metadata != 0 && hash != 0 && l3 != 0 {
		fmt.Print("Credentials Saved")
	}

	//close the file
	err = input.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
