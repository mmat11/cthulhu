package bot

import (
	"bytes"
	"encoding/gob"
)

func MarshalCountersMap(c CountersMap) []byte {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(c); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func UnmarshalCountersMap(b []byte) (CountersMap, error) {
	var c CountersMap

	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)

	if err := dec.Decode(&c); err != nil {
		return nil, err
	}
	return c, nil
}
