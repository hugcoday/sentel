package main

import (
  "registry/server"
  "registry/config"
  "registry/protocol"

	"log"
	"net"

	mc "utility/config"
	grpc	"google.golang.org/grpc"
)

const (
	defaultConfigFilePath = "/etc/sentel/registry.toml"
)

func main() {
	// Get configuration
	loader := mc.NewWithPath(defaultConfigFilePath)
	var c config.RegistryConfig
	loader.MustLoad(c)

	// run rpc server
	address := c.Host + ":" + c.Port
  lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()
	protocol.RegisterRegistryServer(s, &server.RegistryServer{Config: &c})
	s.Serve(lis)
}
