package converters

import (
	"fmt"

	"github.com/golang/protobuf/proto"
)

type ProtoConverter struct{}

func (c *ProtoConverter) Serialize(o *TestStruct) ([]byte, error) {
	bytes, err := proto.Marshal(o)
	if err != nil {
		return nil, fmt.Errorf("proto marshal: %w", err)
	}
	return bytes, nil
}

func (c *ProtoConverter) Deserialize(raw []byte) (*TestStruct, error) {
	o := &TestStruct{}
	if err := proto.Unmarshal(raw, o); err != nil {
		return nil, fmt.Errorf("proto unmarshal: %w", err)
	}
	return o, nil
}
