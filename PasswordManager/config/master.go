package config

type MasterKey struct {
	key []byte
}

var masterKey = &MasterKey{}

func SetMasterKey(k []byte) {
	masterKey.setMasterKey(k)
}
func (m *MasterKey) setMasterKey(k []byte) {
	m.key = make([]byte, len(k))
	copy(m.key, k)
}
func UseMasterKey(fn func([]byte)) {
	masterKey.getMasterKey(fn)
}
func (m *MasterKey) getMasterKey(fn func([]byte)) {
	fn(m.key)
}

func ClearMasterKey() {
	masterKey.clearMasterKey()
}
func (m *MasterKey) clearMasterKey() {
	for i := range m.key {
		m.key[i] = 0
	}
	m.key = nil
}
