package muni

import (
	"fmt"
	"strconv"
)

type Uint64 uint64

// MarshalJSON returns a json string for the ID.
func (id Uint64) MarshalJSON() ([]byte, error) {
	buff := make([]byte, 0, 22)
	buff = append(buff, '"')
	buff = strconv.AppendUint(buff, uint64(id), 10)
	buff = append(buff, '"')
	return buff, nil
}

// UnmarshalJSON converts a json (string) ID into an ID.
func (id *Uint64) UnmarshalJSON(b []byte) error {
	len := len(b)
	if len < 3 || b[0] != '"' || b[len-1] != '"' {
		return fmt.Errorf("UnmarshalJSON: parsing %s: invalid syntax", strconv.Quote(string(b)))
	}

	i, err := strconv.ParseUint(string(b[1:len-1]), 10, 64)
	if err != nil {
		// Don't here, since we are saying the same thing basically.
		return fmt.Errorf("UnmarshalJSON: parsing %s: invalid syntax", strconv.Quote(string(b)))
	}

	*id = Uint64(i)
	return nil
}
