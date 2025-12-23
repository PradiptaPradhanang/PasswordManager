package utility

import (
	"crypto/rand"
	"fmt"
	"os"
	"passmana/config"
	"passmana/encrypto"
)

// ─── YOUR EXACT SPEC ───
func CreateVault(master []byte) {

	// 1. 32-byte random challenge
	challenge := make([]byte, 32)
	if _, err := rand.Read(challenge); err != nil {
		fmt.Println(err)
		return
	}
	cfg := encrypto.MasterKeyConfig{
		Mpassword: master,
		Salt:      nil,
	}
	key, salt, err := encrypto.DeriveMasterKey(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(key)
	cipherPass, nonceMaster, err := encrypto.Encryption(key, challenge)
	if err != nil {
		fmt.Println(err)
		return
	}
	sealed := append(nonceMaster, cipherPass...) // 12+48=60

	// 5. Save proof
	os.WriteFile(config.SaltFile, salt, 0644)
	os.WriteFile(config.ChallengFile, sealed, 0600)

	// fmt.Println("Salt length:", len(salt))     // Should match Argon2 salt size
	// fmt.Println("Sealed length:", len(sealed)) // Should be nonceMaster + ciphertext
	// fmt.Println("nonceMaster:", sealed[:12])
	// fmt.Println("Ciphertext:", sealed[12:])
}

func VerifyPass(master []byte) (check bool) {
	salt, err := os.ReadFile(config.SaltFile)
	if err != nil {
		fmt.Println(err)
		return false
	}
	sealed, err := os.ReadFile(config.ChallengFile)
	if err != nil {
		fmt.Println(err)
		return false

	}
	// fmt.Println("Salt length:", len(salt))     // Should match Argon2 salt size
	// fmt.Println("Sealed length:", len(sealed)) // Should be nonceMaster + ciphertext
	// fmt.Println("nonceMaster:", sealed[:12])
	// fmt.Println("Ciphertext:", sealed[12:])
	cfg := encrypto.MasterKeyConfig{
		Mpassword: master,
		Salt:      salt,
	}
	////derive master key from the input master password
	key, _, err := encrypto.DeriveMasterKey(cfg)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println(key)
	nonceMaster := sealed[:12]
	ciphertext := sealed[12:]
	_, err = encrypto.Decryption(key, nonceMaster, ciphertext)
	if err != nil {
		fmt.Println(err)
		return false
	}
	config.SetMasterKey(key)
	return true

}
