package filestore

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/jckimble/golibsignal/axolotl"
)

type SessionStore struct {
	sync.Mutex
}

func (s *SessionStore) exists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func (s *SessionStore) sessionFilePath(recipientID string, deviceID uint32) string {
	return filepath.Join(".signal/sessions", fmt.Sprintf("%s_%d", recipientID, deviceID))
}

func (s *SessionStore) GetSubDeviceSessions(recipientID string) []uint32 {
	sessions := []uint32{}

	filepath.Walk(".signal/sessions", func(path string, fi os.FileInfo, err error) error {
		if !fi.IsDir() {
			i := strings.LastIndex(path, "_")
			id, _ := strconv.Atoi(path[i+1:])
			sessions = append(sessions, uint32(id))
		}
		return nil
	})
	return sessions
}

func (s *SessionStore) LoadSession(recipientID string, deviceID uint32) (*axolotl.SessionRecord, error) {
	sfile := s.sessionFilePath(recipientID, deviceID)
	b, err := ioutil.ReadFile(sfile)
	if err != nil {
		return axolotl.NewSessionRecord(), nil
	}
	record, err := axolotl.LoadSessionRecord(b)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (s *SessionStore) StoreSession(recipientID string, deviceID uint32, record *axolotl.SessionRecord) error {
	sfile := s.sessionFilePath(recipientID, deviceID)
	b, err := record.Serialize()
	if err != nil {
		return err
	}
	os.MkdirAll(".signal/sessions", 0700)
	return ioutil.WriteFile(sfile, b, 0600)
}

func (s *SessionStore) ContainsSession(recipientID string, deviceID uint32) bool {
	sfile := s.sessionFilePath(recipientID, deviceID)
	return s.exists(sfile)
}

func (s *SessionStore) DeleteSession(recipientID string, deviceID uint32) {
	sfile := s.sessionFilePath(recipientID, deviceID)
	_ = os.Remove(sfile)
}

func (s *SessionStore) DeleteAllSessions(recipientID string) {
	sessions := s.GetSubDeviceSessions(recipientID)
	for _, dev := range sessions {
		_ = os.Remove(s.sessionFilePath(recipientID, dev))
	}
}
