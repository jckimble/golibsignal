package sqlstore

import (
	"database/sql"
	"fmt"
	"github.com/jckimble/golibsignal/axolotl"
)

type PreKeyStore struct {
	db *sql.DB
}

func (s PreKeyStore) ContainsPreKey(id uint32) bool {
	row := s.db.QueryRow("SELECT id FROM prekeys WHERE id=?", id)
	err := row.Scan(&id)
	return err == nil
}
func (s PreKeyStore) LoadPreKey(id uint32) (*axolotl.PreKeyRecord, error) {
	row := s.db.QueryRow("SELECT `key` FROM prekeys WHERE id=?", id)
	var key bytedata
	if err := row.Scan(&key); err != nil {
		return nil, fmt.Errorf("Invalid PreKey: %s", err)
	}
	record, err := axolotl.LoadPreKeyRecord([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("Unable to load PreKey: %s", err)
	}
	return record, nil
}
func (s PreKeyStore) RemovePreKey(id uint32) {
	s.db.Exec("DELETE FROM prekeys WHERE id=?", id)
}
func (s PreKeyStore) StorePreKey(id uint32, record *axolotl.PreKeyRecord) error {
	b, err := record.Serialize()
	if err != nil {
		return err
	}
	_, err = s.db.Exec("INSERT INTO prekeys(`id`,`key`) VALUES(?,?)", id, bytedata(b))
	if err != nil {
		return fmt.Errorf("Unable to save prekey: %s", err)
	}
	return nil
}

func NewPreKeyStore(db *sql.DB) *PreKeyStore {
	return &PreKeyStore{
		db: db,
	}
}
