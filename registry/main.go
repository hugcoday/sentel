package main

import (
	"log"
	"net"
	mc "sentel/utility/config"

	grpc "github.com/grpc/grpc-go"
	//	"google.golang.org/grpc"
)

const (
	defaultConfigFilePath = "/etc/sentel/registry.toml"
)

func main() {
	// Get configuration
	loader := mc.NewWithPath(defaultConfigFilePath)
	var c RegistryConfig
	c.MustLoad(c)

	// run rpc server
	address := c.Host + ":" + c.Port
	lis, err = net.Listen("tcp", address)
	if err != nil {
		log.Fatal("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()
	pb.RegistryRegistryServer(s, &server.RegisryServer{Config: &c})
	s.Serve(lis)
}
