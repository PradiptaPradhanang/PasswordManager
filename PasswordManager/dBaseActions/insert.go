package database

import (
	"fmt"
	"os"
	"passmana/config"
	"passmana/dbControl"
	"passmana/encrypto"
	//"golang.org/x/crypto/bcrypt"
)

func Insert(username string, platform string, password string, masterPassword string) {

	//converts string into byte array because bcrypt only accept byte array
	cfg := encrypto.MasterKeyConfig{
		Mpassword: masterPassword,
		Salt:      nil,
	}
	mKey, salt, err := encrypto.DeriveMasterKey(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	bytePassword := []byte(password)
	enPassword, nonce, err := encrypto.Encryption(bytePassword, mKey)
	//hashValue, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return
	}
	//data := username + "," + platform + ","
	//open the file in WRITE mode, if file is not present then it will create
	input, err := os.OpenFile(config.DBName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	//adding platform and username
	//l1 is size of the string passed
	// metadata, err := input.WriteString(data)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// tmpsalt := base64.StdEncoding.EncodeToString(salt)
	// tmpnonce := base64.StdEncoding.EncodeToString(nonce)
	// enHashPassowrd := base64.StdEncoding.EncodeToString(enPassword)
	fmt.Print("database inserting..")
	db := dbControl.Get()
	_, err = db.Exec(`INSERT INTO creds (username, platform, salt, nonce, cipherpass)
                   VALUES (?, ?, ?, ?, ?)`, username, platform, salt, nonce, enPassword)
	//total := tmpsalt + "," + tmpnonce + ","
	// caldata, err := input.WriteString(total)
	if err != nil {
		// if strings.Contains(err.Error(), "UNIQUE constraint failed") {
		// 	fmt.Println("Platform already exists. Use update instead.")
		// } else {
		fmt.Println("Insert error:", err)
		//}
		///// to do
		return
	}
	fmt.Print("Credentials Saved")
	// //adding password
	// hash, err := input.WriteString(base64.StdEncoding.EncodeToString(enPassword)) //string(hash) function save the hash int the string format
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// //adding new line
	// l3, err := input.WriteString("\n")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// //check the data was stored or not
	// if metadata != 0 && hash != 0 && l3 != 0 && caldata != 0 {
	//
	// }

	//close the file
	err = input.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
