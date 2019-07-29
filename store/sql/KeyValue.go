package sqlstore

import (
	"database/sql"
	"encoding/base64"
	"fmt"
)

type KeyValue struct {
	Key   string
	Value string
}

func (kv *KeyValue) SetBytes(b []byte) {
	kv.Value = base64.StdEncoding.EncodeToString(b)
}

func (kv KeyValue) Bytes() ([]byte, error) {
	return base64.StdEncoding.DecodeString(kv.Value)
}
func (kv KeyValue) String() string {
	return kv.Value
}

func setValue(db *sql.DB, key string, val interface{}) error {
	kv := &KeyValue{
		Key: key,
	}
	if v, ok := val.(string); ok {
		kv.Value = v
	} else if b, ok := val.([]byte); ok {
		kv.SetBytes(b)
	} else {
		return fmt.Errorf("Unsupported Type")
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec("INSERT INTO config(`key`,`value`) VALUES(?,?);", kv.Key, kv.Value); err == nil {
		return tx.Commit()
	}
	if _, err := tx.Exec("UPDATE config SET value=? WHERE `key`=?", kv.Value, kv.Key); err != nil {
		tx.Rollback()
		return fmt.Errorf("Unable to set %s: %s", key, err)
	}
	return tx.Commit()
}

func getValue(db *sql.DB, key string) (*KeyValue, error) {
	kv := &KeyValue{}
	row := db.QueryRow("SELECT `key`,value FROM config WHERE `key`=?", key)
	if err := row.Scan(&kv.Key, &kv.Value); err != nil {
		return nil, fmt.Errorf("%s is not set: %s", key, err)
	}
	return kv, nil
}
