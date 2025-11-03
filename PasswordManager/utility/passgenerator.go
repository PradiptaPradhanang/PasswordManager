package utility

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GeneratePassword(length int) (string, error) {

	const (
		uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
		digits           = "0123456789"
		specialChars     = "!@#$%^&*()-_=+[]{}|;:'\",.<>/?"
	)

	totalChars := uppercaseLetters + lowercaseLetters + digits + specialChars
	fmt.Println(totalChars)
	var password string

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(totalChars))))
		if err != nil {
			return "", err
		}
		password += string(totalChars[randomIndex.Int64()])
	}
	return password, nil

}
