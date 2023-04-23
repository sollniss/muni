package muni

import (
	"errors"
	"strconv"
)

type Uint64 uint64

// MarshalJSON returns a json string for the ID.
func (id Uint64) MarshalJSON() ([]byte, error) {
	buff := make([]byte, 0, 22)
	buff = append(buff, '"')
	buff = strconv.AppendInt(buff, int64(id), 10)
	buff = append(buff, '"')
	return buff, nil
}

// UnmarshalJSON converts a json (string) ID into an ID.
func (id *Uint64) UnmarshalJSON(b []byte) error {
	len := len(b)
	if len < 3 || b[0] != '"' || b[len-1] != '"' {
		return errors.New("UnmarshalJSON: invalid ID")
	}

	i, err := strconv.ParseInt(string(b[1:len-1]), 10, 64)
	if err != nil {
		return err
	}

	*id = Uint64(i)
	return nil
}
