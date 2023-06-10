package converters

import (
	"encoding/json"
	"fmt"
)

type JsonConverter struct {
}

func (c *JsonConverter) Serialize(o *TestStruct) ([]byte, error) {
	bytes, err := json.Marshal(o)
	if err != nil {
		return nil, fmt.Errorf("marshal json: %w", err)
	}
	return bytes, nil
}

func (c *JsonConverter) Deserialize(raw []byte) (*TestStruct, error) {
	o := &TestStruct{}
	err := json.Unmarshal(raw, o)
	if err != nil {
		return nil, fmt.Errorf("unmarshal json: %w", err)
	}
	return o, nil
}
