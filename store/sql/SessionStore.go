package sqlstore

import (
	"fmt"
	"sync"

	"database/sql"
	"github.com/jckimble/golibsignal/axolotl"
)

type SessionStore struct {
	sync.Mutex
	db *sql.DB
}

func (s *SessionStore) GetSubDeviceSessions(recipientID string) []uint32 {
	sessions := []uint32{}
	rows, err := s.db.Query("SELECT device FROM sessions WHERE recipient=?", recipientID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var device uint32
		if err := rows.Scan(&device); err != nil {
			return nil
		}
		sessions = append(sessions, device)
	}
	if err := rows.Err(); err != nil {
		return nil
	}
	return sessions
}

func (s *SessionStore) LoadSession(recipientID string, deviceID uint32) (*axolotl.SessionRecord, error) {
	row := s.db.QueryRow("SELECT data FROM sessions WHERE recipient=? AND device=?", recipientID, deviceID)
	var data bytedata
	if err := row.Scan(&data); err != nil {
		return axolotl.NewSessionRecord(), nil
	}
	record, err := axolotl.LoadSessionRecord([]byte(data))
	if err != nil {
		return nil, fmt.Errorf("Unable to load session: %s", err)
	}
	return record, nil
}

func (s *SessionStore) StoreSession(recipientID string, deviceID uint32, record *axolotl.SessionRecord) error {
	b, err := record.Serialize()
	if err != nil {
		return err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	if _, err = tx.Exec("INSERT INTO sessions(`recipient`,`device`,`data`) VALUES(?,?,?)", recipientID, deviceID, bytedata(b)); err == nil {
		return tx.Commit()
	}
	if _, err := tx.Exec("UPDATE sessions SET `data`=? WHERE recipient=? AND device=?", bytedata(b), recipientID, deviceID); err != nil {
		tx.Rollback()
		return fmt.Errorf("Unable to store session: %s", err)
	}
	return tx.Commit()
}

func (s *SessionStore) ContainsSession(recipientID string, deviceID uint32) bool {
	row := s.db.QueryRow("SELECT device FROM sessions WHERE recipient=? AND device=?", recipientID, deviceID)
	err := row.Scan(&deviceID)
	return err == nil

}

func (s *SessionStore) DeleteSession(recipientID string, deviceID uint32) {
	s.db.Exec("DELETE FROM sessions WHERE recipient=? AND device=?", recipientID, deviceID)
}

func (s *SessionStore) DeleteAllSessions(recipientID string) {
	s.db.Exec("DELETE FROM sessions WHERE recipient=?", recipientID)
}

func NewSessionStore(db *sql.DB) *SessionStore {
	return &SessionStore{
		db: db,
	}
}
