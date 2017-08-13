//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use this file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.

package client

import (
  pb "registry/protocol"
  util "lib/config"

	"context"
	"log"

	grpc "google.golang.org/grpc"
)

type RegistryApi struct {
	conn     *grpc.ClientConn
	registry pb.RegistryClient
}

func New(c util.Config) (*RegistryApi, error) {
	conn, err := grpc.Dial(c.GetKey("registry"), grpc.WithInsecure())
	if err != nil {
		log.Fatal("did not connect: %v", err)
		return nil, err
	}
	r := pb.NewRegistryClient(conn)
	return &RegistryApi{conn: conn, registry: r}, nil
}

func (r *RegistryApi) AddDevice(ctx context.Context, name string) error {
  _, err := r.registry.AddDevice(ctx, &pb.DeviceAddRequest{Name: "hello"})
	if err != nil {
		log.Fatal("AddDevice rpc fail failed: %v", err)
	}
  return nil
}
