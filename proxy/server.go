package proxy

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Server struct {
	Port            int64
	MulticastAddr   string
	ConvertersAddrs map[string]string
	Result          chan string
}

func (s Server) ListenMulticastGroup() error {
	resolverAddr, err := net.ResolveUDPAddr("udp", s.MulticastAddr)
	if err != nil {
		return fmt.Errorf("failed to resolve udp addr: %v", err)
	}
	conn, err := net.ListenMulticastUDP("udp", nil, resolverAddr)
	if err != nil {
		return fmt.Errorf("failed to listen multicast udp: %v", err)
	}
	defer conn.Close()

	for {
		for range s.ConvertersAddrs {
			buf := make([]byte, 4096)
			bytesRead, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				return fmt.Errorf("failed to read from UDP conn: %v", err)
			}
			s.Result <- string(buf[:bytesRead])
		}
	}
}

func (s Server) ProcessMulticast() ([]byte, error) {
	resolverAddr, err := net.ResolveUDPAddr("udp", s.MulticastAddr)
	if err != nil {
		return nil, fmt.Errorf("error with resole address")
	}

	conn, err := net.DialUDP("udp", nil, resolverAddr)
	if err != nil {
		return nil, fmt.Errorf("error in udp common addr")
	}
	defer conn.Close()

	_, err = conn.Write([]byte("here we go"))
	if err != nil {
		return nil, fmt.Errorf("failed to write to addr")
	}

	var result []string
	for i := 0; i < cap(s.Result); i++ {
		if str := <-s.Result; str != "here we go" {
			result = append(result, str)
		} else {
			result = append(result, <-s.Result)
		}
	}

	return []byte(strings.Join(result, "")), nil
}

func (s Server) ProcessRequest(request []byte) ([]byte, error) {
	req := strings.Trim(string(request), "\n")
	var addr string
	if req == "All" {
		return s.ProcessMulticast()
	}

	addr, ok := s.ConvertersAddrs[req]
	if !ok {
		return nil, fmt.Errorf("can not convert %v\n", string(request))
	}
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return nil, fmt.Errorf("bad address")
	}
	defer conn.Close()

	_, err = fmt.Fprintf(conn, "here we go")
	if err != nil {
		return nil, fmt.Errorf("failed to write data to addr %q: %v", addr, err)
	}

	buf := make([]byte, 4096)
	n, err := bufio.NewReader(conn).Read(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read data from connection: %v", err)
	}
	return buf[:n], nil
}
