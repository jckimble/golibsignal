package sqlstore

import (
	"bytes"
	"database/sql"
	"fmt"
	"gitlab.com/jckimble/golibsignal/axolotl"
	"strconv"
)

type IdentityStore struct {
	db *sql.DB
}

func (s IdentityStore) GetIdentityKeyPair() (*axolotl.IdentityKeyPair, error) {
	kv, err := getValue(s.db, "identity_key")
	if err != nil {
		return nil, err
	}
	b, err := kv.Bytes()
	if err != nil {
		return nil, err
	}
	if len(b) != 64 {
		return nil, fmt.Errorf("identity key is %d not 64 bytes long", len(b))
	}
	return axolotl.NewIdentityKeyPairFromKeys(b[32:], b[:32]), nil
}
func (s IdentityStore) SetIdentityKeyPair(ikp *axolotl.IdentityKeyPair) error {
	b := make([]byte, 64)
	copy(b, ikp.PublicKey.Key()[:])
	copy(b[32:], ikp.PrivateKey.Key()[:])
	return setValue(s.db, "identity_key", b)
}
func (s IdentityStore) GetLocalRegistrationID() (uint32, error) {
	kv, err := getValue(s.db, "regid")
	if err != nil {
		return 0, err
	}
	num, err := strconv.Atoi(kv.String())
	if err != nil {
		return 0, err
	}
	return uint32(num), nil
}

func (s IdentityStore) SetLocalRegistrationID(id uint32) error {
	return setValue(s.db, "regid", strconv.Itoa(int(id)))
}
func (s IdentityStore) IsTrustedIdentity(id string, key *axolotl.IdentityKey) bool {
	row := s.db.QueryRow("SELECT identitykey FROM contacts WHERE tel=?", "+"+id)
	var b bytedata
	if err := row.Scan(&b); err != nil {
		return true //TOFU
	}
	return bytes.Equal([]byte(b), key.Key()[:])
}

func (s IdentityStore) SaveIdentity(id string, key *axolotl.IdentityKey) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec("INSERT INTO contacts(`tel`,`identitykey`) VALUES(?,?);", "+"+id, bytedata(key.Key()[:])); err == nil {
		return tx.Commit()
	}
	if _, err := tx.Exec("UPDATE contacts SET identitykey=? WHERE tel=?", bytedata(key.Key()[:]), "+"+id); err != nil {
		tx.Rollback()
		return fmt.Errorf("Unable to save identitykey: %s", err)
	}
	return tx.Commit()
}

func NewIdentityStore(db *sql.DB) *IdentityStore {
	return &IdentityStore{
		db: db,
	}
}
