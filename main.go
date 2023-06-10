package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/hattonuri/soa-hmw-1/config"
	"github.com/hattonuri/soa-hmw-1/converters"
	"github.com/hattonuri/soa-hmw-1/models"
	"github.com/hattonuri/soa-hmw-1/proxy"
)

type Server interface {
	ProcessRequest([]byte) ([]byte, error)
}

type MulticastScheduler interface {
	ListenMulticastGroup() error
}

type Scheduler struct {
	GroupAddr string
	Address   string
	Server    Server
}

func MakeScheduler(groupAddr, serverType string, serviceInfos map[string]models.ServiceInfo) (*Scheduler, error) {
	convertersAddrs := make(map[string]string)
	port := serviceInfos[serverType].Port
	for name, info := range serviceInfos {
		if name == serverType {
			port = int64(info.Port)
		} else {
			convertersAddrs[name] = name + ":" + strconv.Itoa(int(info.Port))
		}
	}

	scheduler := &Scheduler{
		Address:   fmt.Sprintf("%s:%d", serverType, port),
		GroupAddr: groupAddr,
	}

	switch serverType {
	case "proxy":
		scheduler.Server = proxy.Server{
			Port:            port,
			ConvertersAddrs: convertersAddrs,
			MulticastAddr:   groupAddr,
			Result:          make(chan string, len(serviceInfos)-1),
		}
	case "native":
		scheduler.Server = converters.NewServer(port, groupAddr, &converters.NativeConverter{})
	case "xml":
		scheduler.Server = converters.NewServer(port, groupAddr, &converters.XMLConverter{})
	case "json":
		scheduler.Server = converters.NewServer(port, groupAddr, &converters.JsonConverter{})
	case "proto":
		scheduler.Server = converters.NewServer(port, groupAddr, &converters.ProtoConverter{})
	case "avro":
		scheduler.Server = converters.NewServer(port, groupAddr, converters.NewAvroConverter())
	case "yaml":
		scheduler.Server = converters.NewServer(port, groupAddr, &converters.YAMLConverter{})
	case "msgpack":
		scheduler.Server = converters.NewServer(port, groupAddr, &converters.MsgPackConverter{})
	}

	return scheduler, nil
}

func (s *Scheduler) Listen() error {
	conn, err := net.ListenPacket("udp", s.Address)
	if err != nil {
		return fmt.Errorf("listen packet: %w", err)
	}
	defer conn.Close()

	for {
		buf := make([]byte, 4096)
		bytesRead, addr, err := conn.ReadFrom(buf)
		if err != nil {
			return fmt.Errorf("reading from connection: %w", err)
		}

		res, err := s.Server.ProcessRequest(buf[:bytesRead])
		if err != nil {
			fmt.Printf("error processing request: %v", err)
			conn.WriteTo([]byte(err.Error()+"\n"), addr)
			continue
		}
		_, err = conn.WriteTo([]byte(res), addr)
		if err != nil {
			return fmt.Errorf("write to connection: %w", err)
		}
	}
}

func main() {
	content, err := ioutil.ReadFile("services.json")
	if err != nil {
		log.Fatalf("bad services.json: %v", err)
		return
	}

	var serviceInfos map[string]models.ServiceInfo
	if err := json.Unmarshal(content, &serviceInfos); err != nil {
		log.Fatalf("services.json unmarshal: %v", err)
		return
	}
	cfg := &config.Config{}
	config.Init(cfg)
	if err := config.Init(cfg); err != nil {
		log.Fatalf("failed to parse config: %v", err)
		return
	}

	scheduler, err := MakeScheduler(cfg.GroupAddr, os.Args[1], serviceInfos)
	if err != nil {
		log.Fatalf("make scheduler %v", err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		obj, ok := scheduler.Server.(MulticastScheduler)
		if !ok {
			return
		}
		obj.ListenMulticastGroup()
	}()

	go func() {
		defer wg.Done()
		scheduler.Listen()
	}()

	wg.Wait()
}
