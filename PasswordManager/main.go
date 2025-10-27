package main

import (
	"flag"
	"fmt"
	"passmana/passwordgenerator"
)

func main() {
	fmt.Println("Hello, World!")
	// file, err := os.Create("file.go") // For read access.
	// if err != nil {
	// 	fmt.Println(err)
	// // 	return
	// // }
	// defer file.Close()
	Length := flag.Int("length", 12, "length of your password")
	flag.Parse()
	fmt.Println("File Created Successfully")
	output, _ := passwordgenerator.GeneratePassword(*Length)
	fmt.Println(output)

}
