package muni

import (
	"encoding/json"
	"testing"
)

type testJson struct {
	ID Uint64 `json:"id"`
}

func TestIDMarshalJSON(t *testing.T) {
	var tests = [...]struct {
		name string
		in   any
		out  string
	}{
		{
			"min",
			0,
			"0",
		},
		{
			"max",
			uint64(1<<64 - 1),
			"18446744073709551615",
		},
		{
			"json",
			testJson{12345},
			`{"id":"12345"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := json.Marshal(tt.in)
			if string(got) != tt.out {
				t.Errorf("marshalling type (%s): got: %s, expected: %s", tt.name, got, tt.out)
			}
		})
	}
}

type unmarshalTest[T Uint64 | testJson] struct {
	name string
	in   string
	out  T
	err  bool
}

func TestIDUnmarshalJSON(t *testing.T) {
	var testsUint64 = [...]unmarshalTest[Uint64]{
		{
			"min",
			`"0"`,
			0,
			false,
		},
		{
			"max",
			`"18446744073709551615"`,
			Uint64(1<<64 - 1),
			false,
		},
		{
			"too big",
			`"18446744073709551616"`,
			0,
			true,
		},
		{
			"negative",
			`"-1"`,
			0,
			true,
		},
	}

	var testsJson = [...]unmarshalTest[testJson]{
		{
			"normal int",
			`{"id":"12345"}`,
			testJson{ID: 12345},
			false,
		},
		{
			"json string",
			`{"id":"asd"}`,
			testJson{},
			true,
		},
		{
			"json empty",
			`{"id":""}`,
			testJson{},
			true,
		},
	}

	for _, tt := range testsUint64 {
		t.Run(tt.name, func(t *testing.T) {
			var val Uint64
			err := json.Unmarshal([]byte(tt.in), &val)

			if (err != nil) != tt.err {
				t.Errorf("%s: invalid err, got: %+v", tt.name, err)
			}

			if val != tt.out {
				t.Errorf("%s: got: %+v, expected: %+v", tt.name, val, tt.out)
			}
		})
	}

	for _, tt := range testsJson {
		t.Run(tt.name, func(t *testing.T) {
			var val testJson
			err := json.Unmarshal([]byte(tt.in), &val)

			if (err != nil) != tt.err {
				t.Errorf("%s: invalid err, got: %+v", tt.name, err)
			}

			if val.ID != tt.out.ID {
				t.Errorf("%s: got: %+v, expected: %+v", tt.name, val, tt.out)
			}
		})
	}
}
