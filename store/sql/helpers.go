package sqlstore

import (
	"database/sql/driver"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
)

type uint32Slice []uint32

func (s uint32Slice) Value() (driver.Value, error) {
	str := []string{}
	for _, v := range s {
		str = append(str, strconv.Itoa(int(v)))
	}
	return strings.Join(str, ","), nil
}
func (s *uint32Slice) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	if sv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := sv.(string); ok {
			uintslice := []uint32{}
			spt := strings.Split(v, ",")
			for _, v := range spt {
				ui, err := strconv.Atoi(v)
				if err != nil {
					return err
				}
				uintslice = append(uintslice, uint32(ui))
			}
			// set the value of the pointer yne to YesNoEnum(v)
			*s = uint32Slice(uintslice)
			return nil
		}
	}
	return errors.New("failed to scan uint32Slice")
}

type stringSlice []string

func (s stringSlice) Value() (driver.Value, error) {
	return strings.Join([]string(s), ","), nil
}

func (s *stringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	if sv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := sv.(string); ok {
			*s = stringSlice(strings.Split(v, ","))
			return nil
		}
	}
	return errors.New("failed to scan stringSlice")
}

type bytedata []byte

func (d bytedata) Value() (driver.Value, error) {
	return base64.StdEncoding.EncodeToString([]byte(d)), nil
}

func (d *bytedata) Scan(value interface{}) error {
	if value == nil {
		*d = nil
		return nil
	}
	if sv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := sv.(string); ok {
			data, err := base64.StdEncoding.DecodeString(v)
			if err != nil {
				return err
			}
			*d = data
			return nil
		}
	}
	return errors.New("failed to scan bytedata")
}
