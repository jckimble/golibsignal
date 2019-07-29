package filestore

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
	"sync"
)

type Config struct {
	Tel    string `yaml:"tel"`
	Server string `yaml:"server"`

	mu sync.Mutex
}

func (c *Config) loadConfig() error {
	f, err := os.Open(".signal/config.yml")
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewDecoder(f).Decode(c)
}

func (c *Config) GetTel() (string, error) {
	if c.Tel == "" {
		c.mu.Lock()
		defer c.mu.Unlock()
		if err := c.loadConfig(); err != nil {
			return "", err
		}
	}
	return c.Tel, nil
}

func (c *Config) GetServer() (string, error) {
	if c.Server == "" {
		c.mu.Lock()
		defer c.mu.Unlock()
		if err := c.loadConfig(); err != nil {
			return "", err
		}
	}
	return c.Server, nil
}
func (c *Config) GetHTTPPassword() (string, error) {
	passFile := filepath.Join(".signal/identity", "http_password")
	b, err := ioutil.ReadFile(passFile)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
func (c *Config) SetHTTPPassword(p string) error {
	os.MkdirAll(".signal/identity", 0700)
	passFile := filepath.Join(".signal/identity", "http_password")
	return ioutil.WriteFile(passFile, []byte(p), 0600)
}
func (c *Config) GetHTTPSignalingKey() ([]byte, error) {
	keyFile := filepath.Join(".signal/identity", "http_signaling_key")
	b, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	return b, nil
}
func (c *Config) SetHTTPSignalingKey(b []byte) error {
	os.MkdirAll(".signal/identity", 0700)
	keyFile := filepath.Join(".signal/identity", "http_signaling_key")
	return ioutil.WriteFile(keyFile, b, 0600)
}
