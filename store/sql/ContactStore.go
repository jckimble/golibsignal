package sqlstore

import (
	"database/sql"
	"fmt"
	"gitlab.com/jckimble/golibsignal"
)

type ContactStore struct {
	db *sql.DB
}

func NewContactStore(db *sql.DB) *ContactStore {
	return &ContactStore{
		db: db,
	}
}

func (s ContactStore) Get(tel string) (*libsignal.Contact, error) {
	ct := &libsignal.Contact{
		Tel:     tel,
		Devices: []uint32{1},
	}
	var devices uint32Slice
	var profileKey bytedata
	row := s.db.QueryRow("SELECT devices,name,avatar,profilekey FROM contacts WHERE tel=?", tel)
	if err := row.Scan(&devices, &ct.Name, &ct.Avatar, &profileKey); err != nil {
		return nil, fmt.Errorf("Contact Doesn't Exist: %s", err)
	}
	if profileKey != nil {
		ct.ProfileKey = []byte(profileKey)
	}
	if devices != nil {
		ct.Devices = []uint32(devices)
	}
	return ct, nil
}

func (s ContactStore) Save(ct *libsignal.Contact) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec("INSERT INTO contacts(`tel`,`devices`,`name`,`avatar`,`profilekey`) VALUES(?,?,?,?,?);", ct.Tel, uint32Slice(ct.Devices), ct.Name, ct.Avatar, bytedata(ct.ProfileKey)); err == nil {
		return tx.Commit()
	}
	if _, err := tx.Exec("UPDATE contacts SET devices=?,name=?,avatar=?,profilekey=? WHERE tel=?", uint32Slice(ct.Devices), ct.Name, ct.Avatar, bytedata(ct.ProfileKey), ct.Tel); err != nil {
		tx.Rollback()
		return fmt.Errorf("Unable to save contact: %s", err)
	}
	return tx.Commit()
}
