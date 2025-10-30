package main

import (
	"passmana/cmd"
	// "flag"
	// "fmt"
	// "os"
	// "passmana/database"
	// "passmana/passwordgenerator"
	// "bufio"
	// "strings"
	// "passmana/masterDBname"
)

func main() {

	cmd.Execute()
	// var args []string = os.Args
	// if args[1] == "put" {
	// 	fmt.Print("Enter master password:")
	// 	mpassword,_:=bufio.NewReader(os.Stdin).ReadString('\n')
	// 	mpassword=strings.TrimSpace(mpassword)
	// 	database.Insert(args[2], args[3], args[4],mpassword)
	// } else if args[1] == "get" {
	// 	//retrieve(args[2], args[3])
	// } else if args[1] == "create" {
	// 	Length := 12 //flag.Int("length", 12, "length of your password")
	// 	flag.Parse()
	// 	output, _ := passwordgenerator.GeneratePassword(Length)
	// 	fmt.Println(output)
	// 	//createPassword(args[2], args[3])
	// } else {
	// 	fmt.Println("Operation isnot supported", args[1])
	// }

}
