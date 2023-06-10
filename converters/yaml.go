package converters

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type YAMLConverter struct {
}

func (c *YAMLConverter) Serialize(o *TestStruct) ([]byte, error) {
	bytes, err := yaml.Marshal(o)
	if err != nil {
		return nil, fmt.Errorf("yaml marshal: %w", err)
	}
	return bytes, nil
}

func (c *YAMLConverter) Deserialize(raw []byte) (*TestStruct, error) {
	person := &TestStruct{}
	if err := yaml.Unmarshal(raw, person); err != nil {
		return nil, fmt.Errorf("yaml unmarshal: %w", err)
	}
	return person, nil
}
