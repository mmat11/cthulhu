package task

import (
	"bytes"
	"encoding/gob"
)

func MarshalWMTaskContext(c *WelcomeMessageTaskContext) []byte {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(c); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func UnmarshalWMTaskContext(b []byte) (*WelcomeMessageTaskContext, error) {
	var c *WelcomeMessageTaskContext = &WelcomeMessageTaskContext{}

	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)

	if err := dec.Decode(c); err != nil {
		return nil, err
	}
	return c, nil
}
