package encrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/argon2"
)

type MasterKeyConfig struct {
	Mpassword string
	Salt      []byte
}

func GenerateRandomSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	return salt, err
}

func DeriveMasterKey(cfg MasterKeyConfig) (keyout []byte, saltout []byte, err error) {

	salt := cfg.Salt
	if salt == nil {
		salt, err = GenerateRandomSalt(16) //With 128 bits of randomness,
		// 	the chance of two salts colliding is astronomically low — even across millions of entries.
		if err != nil {
			fmt.Println(err)
			return nil, nil, err
		}
	}
	key := argon2.IDKey([]byte(cfg.Mpassword), salt, 3, 64*1024, 4, 32)
	//// Derive a 256-bit/32byte encryption key from the master password using Argon2id.
	// This uses 3 iteration, 64MB of memory, and 4 threads for resistance against brute-force and GPU attacks.
	// The salt ensures that identical passwords produce unique keys.
	return key, salt, nil
}

func Encryption(mKey []byte, pass []byte) (enPass []byte, nonceout []byte, err error) {

	cblock, err := aes.NewCipher(mKey) // create cipher block
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	aesgcm, err := cipher.NewGCM(cblock) //wrap aes in gcm mode
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	nonce, err := GenerateRandomSalt(12) // requires 12 byte size for aesgcm
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	ciphertext := aesgcm.Seal(nil, nonce, pass, nil) // encrypt plain text ,; dont set dst as nonce
	return ciphertext, nonce, nil
}

func Decryption(mKey []byte, nonce []byte, pass []byte) (p []byte, err error) {

	cblock, err := aes.NewCipher(mKey) // create cipher block
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(cblock) //wrap aes in gcm mode
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if len(nonce) != aesgcm.NonceSize() {
		fmt.Printf("invalid nonce size: expected %d, got %d", aesgcm.NonceSize(), len(nonce))
		return nil, nil
	}
	password, err := aesgcm.Open(nil, nonce, pass, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return password, nil
}
