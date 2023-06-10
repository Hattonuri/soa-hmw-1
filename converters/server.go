package converters

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type Server struct {
	Port          int64
	Format        string
	MulticastAddr string
	converter     Converter
}

func NewServer(port int64, multicastAddr string, converter Converter) *Server {
	return &Server{
		Port:          port,
		MulticastAddr: multicastAddr,
		converter:     converter,
	}
}

func (s Server) SendToGroup(data []byte) error {
	resolvedAddr, err := net.ResolveUDPAddr("udp", s.MulticastAddr)
	if err != nil {
		return fmt.Errorf("resolve udp addr: %w", err)
	}
	conn, err := net.DialUDP("udp", nil, resolvedAddr)

	if err != nil {
		return fmt.Errorf("set listen conn: %w", err)
	}
	defer conn.Close()
	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("write data: %w", err)
	}
	return nil
}

func (s Server) ListenMulticastGroup() error {
	resolvedAddr, err := net.ResolveUDPAddr("udp", s.MulticastAddr)
	if err != nil {
		return fmt.Errorf("resolve udp addr: %w", err)
	}
	conn, err := net.ListenMulticastUDP("udp", nil, resolvedAddr)
	if err != nil {
		return fmt.Errorf("listen multicast udp: %w", err)
	}
	defer conn.Close()

	for {
		buf := make([]byte, 1000)
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("read data from connection: %v", err)
			continue
		}
		data, err := s.ProcessRequest(buf[:n])
		if err != nil {
			fmt.Printf("process request: %v", err)
			continue
		}
		if len(data) == 0 {
			continue
		}
		err = s.SendToGroup([]byte(data))
		if err != nil {
			fmt.Printf("failed to send to group: %v", err)
		}
	}

}

func (s Server) ProcessConverter(converter Converter) ([]byte, error) {
	person := &TestStruct{
		String_: "Hattonuri",
		Int_:    1488,
		Map_: map[string]string{
			"A": "a",
			"B": "b",
		},
		Array_: []string{
			"a", "b", "c",
		},
		Float_: 3.14,
	}

	var totalTimeSerialize int64
	var totalTimeDeserialize int64
	var totalStructSize int

	for i := 0; i < 1000; i++ {
		start := time.Now()
		bytes, err := converter.Serialize(person)
		totalTimeSerialize += time.Since(start).Microseconds()
		if err != nil {
			return nil, fmt.Errorf("serialize string: %w", err)
		}
		totalStructSize += len(bytes)

		start = time.Now()
		_, err = converter.Deserialize(bytes)
		totalTimeDeserialize += time.Since(start).Microseconds()
		if err != nil {
			return nil, fmt.Errorf("deserialize string: %w", err)
		}
	}
	return []byte(fmt.Sprintf(
		"%d - %dmcs - %dmcs\n",
		totalStructSize/1000, totalTimeSerialize/1000, totalTimeDeserialize/1000),
	), nil
}

func (s Server) ProcessRequest(buf []byte) ([]byte, error) {
	req := strings.Trim(string(buf), "\n")
	if req != "here we go" {
		return nil, nil
	}

	res, err := s.ProcessConverter(s.converter)
	if err != nil {
		return nil, fmt.Errorf("failed to process converter")
	}
	return res, nil
}
