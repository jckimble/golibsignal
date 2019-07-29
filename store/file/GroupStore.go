package filestore

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/jckimble/golibsignal"
)

type GroupStore struct {
}

func (s GroupStore) Get(hexid string) (*libsignal.Group, error) {
	group := filepath.Join(".signal/groups", hexid)
	var gr libsignal.Group
	b, err := ioutil.ReadFile(group)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(b, &gr); err != nil {
		return nil, err
	}
	return &gr, nil
}

func (s GroupStore) GetAll() ([]*libsignal.Group, error) {
	groups := []*libsignal.Group{}
	err := filepath.Walk(".signal/groups", func(path string, fi os.FileInfo, err error) error {
		if !fi.IsDir() {
			if !strings.Contains(path, "avatar") {
				_, hexid := filepath.Split(path)
				gr, err := s.Get(hexid)
				if err != nil {
					return err
				}
				groups = append(groups, gr)
			}
		}
		return nil
	})
	return groups, err
}

func (s GroupStore) Save(group *libsignal.Group) error {
	b, err := yaml.Marshal(group)
	if err != nil {
		return err
	}
	os.MkdirAll(".signal/groups", 0700)
	gr := filepath.Join(".signal/groups", group.Hexid)
	return ioutil.WriteFile(gr, b, 0600)
}
