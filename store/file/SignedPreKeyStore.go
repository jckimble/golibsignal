package filestore

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/jckimble/golibsignal/axolotl"
)

type SignedPreKeyStore struct {
}

func (s SignedPreKeyStore) exists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func (s SignedPreKeyStore) ContainsSignedPreKey(id uint32) bool {
	signed := filepath.Join(".signal/signed_prekeys", fmt.Sprintf("%09d", id))
	return s.exists(signed)
}

func (s SignedPreKeyStore) StoreSignedPreKey(id uint32, record *axolotl.SignedPreKeyRecord) error {
	b, err := record.Serialize()
	if err != nil {
		return err
	}
	os.MkdirAll(".signal/signed_prekeys", 0700)
	signed := filepath.Join(".signal/signed_prekeys", fmt.Sprintf("%09d", id))
	return ioutil.WriteFile(signed, b, 0600)
}
func (s SignedPreKeyStore) LoadSignedPreKey(id uint32) (*axolotl.SignedPreKeyRecord, error) {
	signed := filepath.Join(".signal/signed_prekeys", fmt.Sprintf("%09d", id))
	b, err := ioutil.ReadFile(signed)
	if err != nil {
		return nil, err
	}

	record, err := axolotl.LoadSignedPreKeyRecord(b)
	if err != nil {
		return nil, err
	}

	return record, nil
}
func (s SignedPreKeyStore) LoadSignedPreKeys() []axolotl.SignedPreKeyRecord {
	return []axolotl.SignedPreKeyRecord{}
}
func (s SignedPreKeyStore) RemoveSignedPreKey(id uint32) {
	signed := filepath.Join(".signal/signed_prekeys", fmt.Sprintf("%09d", id))
	os.Remove(signed)
}
