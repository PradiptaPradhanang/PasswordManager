package config

var masterKey []byte

func SetMasterKey(k []byte) {
	masterKey = k
}

func GetMasterKey() []byte {
	return masterKey
}

func ClearMasterKey() {
	for i := range masterKey {
		masterKey[i] = 0
	}
	masterKey = nil
}
