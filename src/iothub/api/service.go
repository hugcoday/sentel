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

package api

import (
	"fmt"
	"iothub/base"
	"libs"
	"net"
	"sync"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ApiService struct {
	config   libs.Config
	chn      chan base.ServiceCommand
	wg       sync.WaitGroup
	listener net.Listener
	srv      *grpc.Server
}

// ApiServiceFactory
type ApiServiceFactory struct{}

// New create apiService service factory
func (m *ApiServiceFactory) New(protocol string, c libs.Config, ch chan base.ServiceCommand) (base.Service, error) {
	address := ":50051"
	server := &ApiService{config: c, wg: sync.WaitGroup{}}

	if addr, err := c.String("authlet", "address"); err == nil && address != "" {
		address = addr
	}

	lis, err := net.Listen("tcp", address)
	if err != nil {
		glog.Fatal("Failed to listen: %v", err)
		return nil, err
	}
	server.listener = lis
	server.srv = grpc.NewServer()
	RegisterApiServer(server.srv, server)
	reflection.Register(server.srv)
	return server, nil

}

func (s *ApiService) GetMetrics() *base.Metrics { return nil }
func (s *ApiService) GetStats() *base.Stats     { return nil }

// Start
func (s *ApiService) Start() error {
	go func(s *ApiService) {
		s.srv.Serve(s.listener)
		s.wg.Add(1)
	}(s)
	return nil
}

// Stop
func (s *ApiService) Stop() {
	s.listener.Close()
	s.wg.Wait()
}

//
// Wait
func (s *ApiService) Wait() {
	s.wg.Wait()
}

func (s *ApiService) Version(ctx context.Context, req *VersionRequest) (*VersionReply, error) {
	version := base.GetServiceManager().GetVersion()
	return &VersionReply{Version: version}, nil
}

func (s *ApiService) Admins(ctx context.Context, req *AdminsRequest) (*AdminsReply, error) {
	return nil, nil
}

func (s *ApiService) Cluster(ctx context.Context, req *ClusterRequest) (*ClusterReply, error) {
	return nil, nil
}

func (s *ApiService) Routes(ctx context.Context, req *RoutesRequest) (*RoutesReply, error) {
	return nil, nil
}

func (s *ApiService) Status(ctx context.Context, req *StatusRequest) (*StatusReply, error) {
	return nil, nil
}

func (s *ApiService) Broker(ctx context.Context, req *BrokerRequest) (*BrokerReply, error) {
	mgr := base.GetServiceManager()
	switch req.Category {
	case "stats":
		stats := mgr.GetStats()
		return &BrokerReply{Stats: stats}, nil
	case "metrics":
		metrics := mgr.GetMetrics()
		return &BrokerReply{Metrics: metrics}, nil
	default:
	}
	return nil, fmt.Errorf("Invalid broker request with categoru:%s", req.Category)
}

func (s *ApiService) Plugins(ctx context.Context, req *PluginsRequest) (*PluginsReply, error) {
	return nil, nil
}

func (s *ApiService) Services(ctx context.Context, req *ServicesRequest) (*ServicesReply, error) {
	return nil, nil
}

func (s *ApiService) Subscriptions(ctx context.Context, req *SubscriptionsRequest) (*SubscriptionsReply, error) {
	return nil, nil
}

func (s *ApiService) Clients(ctx context.Context, req *ClientsRequest) (*ClientsReply, error) {
	return nil, nil
}

func (s *ApiService) Sessions(ctx context.Context, req *SessionsRequest) (*SessionsReply, error) {
	return nil, nil
}

func (s *ApiService) Topics(ctx context.Context, req *TopicsRequest) (*TopicsReply, error) {
	return nil, nil
}
