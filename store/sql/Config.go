package sqlstore

import (
	"database/sql"
)

type Config struct {
	db *sql.DB
}

func (c Config) GetTel() (string, error) {
	kv, err := getValue(c.db, "phone")
	if err != nil {
		return "", err
	}
	return kv.String(), nil
}

func (c Config) GetServer() (string, error) {
	kv, err := getValue(c.db, "server")
	if err != nil {
		return "", err
	}
	return kv.String(), nil
}

func (c Config) GetHTTPPassword() (string, error) {
	kv, err := getValue(c.db, "http_password")
	if err != nil {
		return "", err
	}
	return kv.String(), nil
}

func (c Config) SetHTTPPassword(p string) error {
	return setValue(c.db, "http_password", p)
}

func (c Config) GetHTTPSignalingKey() ([]byte, error) {
	kv, err := getValue(c.db, "http_signaling_key")
	if err != nil {
		return nil, err
	}
	return kv.Bytes()
}

func (c Config) SetHTTPSignalingKey(b []byte) error {
	return setValue(c.db, "http_signaling_key", b)
}

func NewConfig(db *sql.DB) *Config {
	return &Config{
		db: db,
	}
}
