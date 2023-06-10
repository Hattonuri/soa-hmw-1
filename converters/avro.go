package converters

import (
	"fmt"
	"log"

	"github.com/hamba/avro"
)

type AvroConverter struct {
	Schema avro.Schema
}

func NewAvroConverter() *AvroConverter {
	schema, err := avro.Parse(`{
		"type": "record",
		"name": "me",
		"namespace": "org.hamba.avro",
		"fields" : [
			{"name": "Int_", "type": "int"},
			{"name": "Float_", "type": "float"},
			{"name": "String_", "type": "string"},
			{"name": "Array_", "type": {"type":"array", "items": "string"}},
			{"name": "Map_", "type": {"type":"map", "values": "string"}}
		]
	}`)
	if err != nil {
		log.Fatalf("parse avro schema: %v", err)
		return nil
	}
	return &AvroConverter{Schema: schema}
}

func (c *AvroConverter) Serialize(o *TestStruct) ([]byte, error) {
	bytes, err := avro.Marshal(c.Schema, o)
	if err != nil {
		return nil, fmt.Errorf("avro marshal failed: %w", err)
	}
	return bytes, nil
}

func (c *AvroConverter) Deserialize(raw []byte) (*TestStruct, error) {
	o := &TestStruct{}
	err := avro.Unmarshal(c.Schema, raw, o)
	if err != nil {
		return nil, fmt.Errorf("avro unmarshal failed: %w", err)
	}
	return o, nil
}
