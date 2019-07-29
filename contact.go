package libsignal

import (
	"bytes"
	"fmt"
	"gitlab.com/jckimble/golibsignal/signalservice"

	"crypto/aes"
	"crypto/cipher"
	"errors"

	"log"
)

// ContactStore is an interface for storing contact information
type ContactStore interface {
	Get(string) (*Contact, error)
	Save(*Contact) error
}

// Contact holds contact information
type Contact struct {
	Tel        string
	Devices    []uint32
	Name       string
	Avatar     string
	ProfileKey []byte
}

// String returns human readable name for contact, returns phone number if name is empty
func (c Contact) String() string {
	if c.Name != "" {
		return c.Name
	}
	return c.Tel
}

// Profile returns Signal Profile Information
type Profile struct {
	Name   string
	Avatar string
	Key    []byte
}

// GetProfile Gets and Decrypts Signal Profile for number
func (s *Signal) GetProfile(tel string, key []byte) (*Profile, error) {
	m := map[string]interface{}{}
	if err := s.serverRequest("GET", fmt.Sprintf("/v1/profile/%s", tel), nil, "", map[int]interface{}{200: &m}); err != nil {
		return nil, err
	}
	if m["name"] == nil {
		return nil, nil
	}
	b, err := base64DecodeNonPadded(m["name"].(string))
	if err != nil {
		return nil, err
	}
	if len(b) < 12+16+1 {
		return nil, errors.New("CipherText too short!")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	de, err := aesgcm.Open(nil, b[:12], b[12:], nil)
	if err != nil {
		return nil, err
	}
	return &Profile{Name: string(bytes.Trim(de, "\x00")), Avatar: m["avatar"].(string), Key: key}, nil
}

func (s *Signal) handleContact(src string, dm *signalservice.DataMessage) (*Contact, error) {
	contact, err := s.ContactStore.Get(src)
	if err != nil {
		log.Printf("Unknown Contact: %s Error: %s", src, err)
		contact = &Contact{
			Tel:     src,
			Devices: []uint32{1},
		}
	}
	if contact.ProfileKey == nil || !bytes.Equal(contact.ProfileKey, dm.GetProfileKey()) {
		profile, err := s.GetProfile(src, dm.GetProfileKey())
		if err != nil {
			return nil, err
		}
		contact.Name = profile.Name
		contact.ProfileKey = profile.Key
		contact.Avatar = profile.Avatar
		if err := s.ContactStore.Save(contact); err != nil {
			return nil, err
		}
	}
	return contact, nil
}
