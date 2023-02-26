package sqlstore

import (
	"database/sql"
	"fmt"

	"github.com/jckimble/golibsignal"
)

type GroupStore struct {
	db *sql.DB
}

func NewGroupStore(db *sql.DB) *GroupStore {
	return &GroupStore{
		db: db,
	}
}

func (s GroupStore) Get(hexid string) (*libsignal.Group, error) {
	gr := &libsignal.Group{}
	row := s.db.QueryRow("SELECT id,hexid,flags,name,members FROM groups WHERE hexid=?", hexid)
	var id bytedata
	var members stringSlice
	if err := row.Scan(&id, &gr.Hexid, &gr.Flags, &gr.Name, &members); err != nil {
		return nil, fmt.Errorf("Group Doesn't Exist: %s", err)
	}
	if id != nil {
		gr.ID = []byte(id)
	}
	if members != nil {
		gr.Members = []string(members)
	}
	return gr, nil
}

func (s GroupStore) GetAll() ([]*libsignal.Group, error) {
	groups := []*libsignal.Group{}
	rows, err := s.db.Query("SELECT id,hexid,flags,name,members FROM groups")
	if err != nil {
		return nil, fmt.Errorf("Unable to get groups: %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		gr := &libsignal.Group{}
		var id bytedata
		var members stringSlice
		if err := rows.Scan(&id, &gr.Hexid, &gr.Flags, &gr.Name, &members); err != nil {
			return nil, fmt.Errorf("Unable to get Group: %s", err)
		}
		if id != nil {
			gr.ID = []byte(id)
		}
		if members != nil {
			gr.Members = []string(members)
		}
		groups = append(groups, gr)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error getting groups: %s", err)
	}
	return groups, nil
}

func (s GroupStore) Save(gr *libsignal.Group) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec("INSERT INTO groups(`id`,`hexid`,`flags`,`name`,`members`) VALUES(?,?,?,?,?);", bytedata(gr.ID), gr.Hexid, gr.Flags, gr.Name, stringSlice(gr.Members)); err == nil {
		return tx.Commit()
	}
	if _, err := tx.Exec("UPDATE groups SET flags=?,name=?,members=? WHERE hexid=?", gr.Flags, gr.Name, stringSlice(gr.Members), gr.Hexid); err != nil {
		tx.Rollback()
		return fmt.Errorf("Unable to save group: %s", err)
	}
	return tx.Commit()
}
