package filestore

import (
	"gitlab.com/jckimble/golibsignal"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ContactStore struct {
}

func (ContactStore) Get(tel string) (*libsignal.Contact, error) {
	contact := filepath.Join(".signal/contacts", tel)
	var ct libsignal.Contact
	b, err := ioutil.ReadFile(contact)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(b, &ct); err != nil {
		return nil, err
	}
	return &ct, nil
}

func (ContactStore) Save(ct *libsignal.Contact) error {
	b, err := yaml.Marshal(ct)
	if err != nil {
		return err
	}
	os.MkdirAll(".signal/contacts", 0700)
	contact := filepath.Join(".signal/contacts", ct.Tel)
	return ioutil.WriteFile(contact, b, 0600)
}
