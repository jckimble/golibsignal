package sqlstore

import (
	"database/sql"
	"fmt"
	"github.com/jckimble/golibsignal/axolotl"
)

type SignedPreKeyStore struct {
	db *sql.DB
}

func (s SignedPreKeyStore) ContainsSignedPreKey(id uint32) bool {
	row := s.db.QueryRow("SELECT id FROM signedprekeys WHERE id=?", id)
	err := row.Scan(&id)
	return err == nil
}

func (s SignedPreKeyStore) StoreSignedPreKey(id uint32, record *axolotl.SignedPreKeyRecord) error {
	b, err := record.Serialize()
	if err != nil {
		return err
	}
	_, err = s.db.Exec("INSERT INTO signedprekeys(`id`,`key`) VALUES(?,?)", id, bytedata(b))
	if err != nil {
		return fmt.Errorf("Unable to store SignedPreKey: %s", err)
	}
	return nil
}
func (s SignedPreKeyStore) LoadSignedPreKey(id uint32) (*axolotl.SignedPreKeyRecord, error) {
	row := s.db.QueryRow("SELECT `key` FROM signedprekeys WHERE id=?", id)
	var key bytedata
	if err := row.Scan(&key); err != nil {
		return nil, fmt.Errorf("Invalid Signed PreKey: %s", err)
	}
	record, err := axolotl.LoadSignedPreKeyRecord([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("Unable to load SingnedPreKey: %s", err)
	}

	return record, nil
}
func (s SignedPreKeyStore) RemoveSignedPreKey(id uint32) {
	s.db.Exec("DELETE FROM signedprekeys WHERE id=?", id)
}

func NewSignedPreKeyStore(db *sql.DB) *SignedPreKeyStore {
	return &SignedPreKeyStore{
		db: db,
	}
}
