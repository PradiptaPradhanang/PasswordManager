package model

type Cred struct {
	Username   string
	Platform   string
	Nonce      []byte
	Cipherpass []byte
}
