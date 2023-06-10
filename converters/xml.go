package converters

import (
	"encoding/xml"
	"fmt"
)

type XMLConverter struct {
}

func (c *XMLConverter) Serialize(o *TestStruct) ([]byte, error) {
	bytes, err := xml.MarshalIndent(o, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("xml marshal: %v", err)
	}
	return bytes, nil
}

func (c *XMLConverter) Deserialize(raw []byte) (*TestStruct, error) {
	person := &TestStruct{}
	if err := xml.Unmarshal(raw, person); err != nil {
		return nil, fmt.Errorf("failed to unmarshal xml bytes: %v", err)
	}
	return person, nil
}
