package telegram

import (
	"bytes"
	"encoding/gob"
)

func MarshalUpdate(u *Update) []byte {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(u); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func UnmarshalUpdate(b []byte) (*Update, error) {
	var u *Update = &Update{}

	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)

	if err := dec.Decode(u); err != nil {
		return nil, err
	}
	return u, nil
}
