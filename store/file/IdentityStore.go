package filestore

import (
	"fmt"
	"gitlab.com/jckimble/golibsignal/axolotl"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"bytes"
	"os"
)

type IdentityStore struct {
}

func (s IdentityStore) GetIdentityKeyPair() (*axolotl.IdentityKeyPair, error) {
	idkeyfile := filepath.Join(".signal/identity", "identity_key")
	b, err := ioutil.ReadFile(idkeyfile)
	if err != nil {
		return nil, err
	}
	if len(b) != 64 {
		return nil, fmt.Errorf("identity key is %d not 64 bytes long", len(b))
	}
	return axolotl.NewIdentityKeyPairFromKeys(b[32:], b[:32]), nil
}
func (s IdentityStore) SetIdentityKeyPair(ikp *axolotl.IdentityKeyPair) error {
	os.MkdirAll(".signal/identity", 0700)
	idkeyfile := filepath.Join(".signal/identity", "identity_key")
	b := make([]byte, 64)
	copy(b, ikp.PublicKey.Key()[:])
	copy(b[32:], ikp.PrivateKey.Key()[:])
	return ioutil.WriteFile(idkeyfile, b, 0600)
}
func (s IdentityStore) GetLocalRegistrationID() (uint32, error) {
	regidfile := filepath.Join(".signal/identity", "regid")
	b, err := ioutil.ReadFile(regidfile)
	if err != nil {
		return 0, err
	}
	num, err := strconv.Atoi(string(b))
	if err != nil {
		return 0, err
	}
	return uint32(num), nil
}

func (s IdentityStore) exists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func (s IdentityStore) SetLocalRegistrationID(id uint32) error {
	os.MkdirAll(".signal/identity", 0700)
	regidfile := filepath.Join(".signal/identity", "regid")
	return ioutil.WriteFile(regidfile, []byte(strconv.Itoa(int(id))), 0600)
}

func (s IdentityStore) IsTrustedIdentity(id string, key *axolotl.IdentityKey) bool {
	idkeyfile := filepath.Join(".signal/identity", "remote_"+id)
	if !s.exists(idkeyfile) {
		return true
	}
	b, err := ioutil.ReadFile(idkeyfile)
	if err != nil {
		return false
	}
	return bytes.Equal(b, key.Key()[:])
}

func (s IdentityStore) SaveIdentity(id string, key *axolotl.IdentityKey) error {
	os.MkdirAll(".signal/identity", 0700)
	idkeyfile := filepath.Join(".signal/identity", "remote_"+id)
	return ioutil.WriteFile(idkeyfile, key.Key()[:], 0600)
}
