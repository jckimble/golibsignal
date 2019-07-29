package filestore

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gitlab.com/jckimble/golibsignal/axolotl"
)

type PreKeyStore struct {
}

func (s PreKeyStore) exists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}
func (s PreKeyStore) ContainsPreKey(id uint32) bool {
	prekey := filepath.Join(".signal/prekeys", fmt.Sprintf("%09d", id))
	return s.exists(prekey)
}
func (s PreKeyStore) LoadPreKey(id uint32) (*axolotl.PreKeyRecord, error) {
	prekey := filepath.Join(".signal/prekeys", fmt.Sprintf("%09d", id))
	b, err := ioutil.ReadFile(prekey)
	if err != nil {
		return nil, err
	}

	record, err := axolotl.LoadPreKeyRecord(b)
	if err != nil {
		return nil, err
	}

	return record, nil
}
func (s PreKeyStore) RemovePreKey(id uint32) {
	prekey := filepath.Join(".signal/prekeys", fmt.Sprintf("%09d", id))
	os.Remove(prekey)
}
func (s PreKeyStore) StorePreKey(id uint32, record *axolotl.PreKeyRecord) error {
	b, err := record.Serialize()
	if err != nil {
		return err
	}
	os.MkdirAll(".signal/prekeys", 0700)
	prekey := filepath.Join(".signal/prekeys", fmt.Sprintf("%09d", id))
	return ioutil.WriteFile(prekey, b, 0600)
}
