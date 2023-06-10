package converters

import (
	"encoding/xml"

	"github.com/golang/protobuf/proto"
)

type StringMap map[string]string

type TestStruct struct {
	XMLName xml.Name `xml:"TestStruct"`

	String_ string    `xml:"string_" protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Int_    int32     `xml:"int_" protobuf:"varint,3,opt,name=age,json=age" json:"age,omitempty"`
	Map_    StringMap `protobuf:"bytes,4,rep,name=map_,proto3" json:"map_,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Array_  []string  `xml:"array_" protobuf:"bytes,5,rep,name=array_,proto3" json:"array_,omitempty"`
	Float_  float32   `xml:"float_" json:"mark,omitempty"`
}

func (m *TestStruct) Reset()         { *m = TestStruct{} }
func (m *TestStruct) String() string { return proto.CompactTextString(m) }
func (m *TestStruct) ProtoMessage()  {}

type Converter interface {
	Serialize(p *TestStruct) ([]byte, error)
	Deserialize(bytes []byte) (*TestStruct, error)
}
