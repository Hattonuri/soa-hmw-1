package converters

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type NativeConverter struct{}

func (c *NativeConverter) Serialize(o *TestStruct) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(o); err != nil {
		return nil, fmt.Errorf("native marshal: %w", err)
	}
	return buf.Bytes(), nil
}

func (c *NativeConverter) Deserialize(raw []byte) (*TestStruct, error) {
	var buf bytes.Buffer
	buf.Write(raw)

	o := &TestStruct{}
	dec := gob.NewDecoder(&buf)
	if err := dec.Decode(o); err != nil {
		return nil, fmt.Errorf("native unmarshal: %w", err)
	}
	return o, nil
}
