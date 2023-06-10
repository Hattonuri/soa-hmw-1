package converters

import (
	"fmt"

	"github.com/shamaton/msgpack/v2"
)

type MsgPackConverter struct {
}

func (c *MsgPackConverter) Serialize(o *TestStruct) ([]byte, error) {
	bytes, err := msgpack.Marshal(o)
	if err != nil {
		return nil, fmt.Errorf("msgpack marshal: %w", err)
	}
	return bytes, nil
}

func (c *MsgPackConverter) Deserialize(raw []byte) (*TestStruct, error) {
	o := &TestStruct{}
	if err := msgpack.Unmarshal(raw, o); err != nil {
		return nil, fmt.Errorf("msgpack unmarshal: %w", err)
	}
	return o, nil
}
