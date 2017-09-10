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

// Name
func (s *ApiService) Name() string { return "apiservice" }

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

// Routes delegate routes command
func (s *ApiService) Routes(ctx context.Context, req *RoutesRequest) (*RoutesReply, error) {
	mgr := base.GetServiceManager()
	reply := &RoutesReply{
		Routes: []*RouteInfo{},
		Header: &ReplyMessageHeader{Success: true},
	}

	switch req.Category {
	case "list":
		routes := mgr.GetRoutes(req.Service)
		for _, route := range routes {
			reply.Routes = append(reply.Routes, &RouteInfo{Topic: route.Topic, Route: route.Route})
		}
	case "show":
		route := mgr.GetRoute(req.Service, req.Topic)
		if route != nil {
			reply.Routes = append(reply.Routes, &RouteInfo{Topic: route.Topic, Route: route.Route})
		}
	default:
		return nil, fmt.Errorf("Invalid route command category:%s", req.Category)
	}
	return reply, nil
}

func (s *ApiService) Status(ctx context.Context, req *StatusRequest) (*StatusReply, error) {
	return nil, nil
}

// Broker delegate broker command implementation in sentel
func (s *ApiService) Broker(ctx context.Context, req *BrokerRequest) (*BrokerReply, error) {
	mgr := base.GetServiceManager()
	switch req.Category {
	case "stats":
		stats := mgr.GetStats(req.Service)
		return &BrokerReply{Stats: stats}, nil
	case "metrics":
		metrics := mgr.GetMetrics(req.Service)
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

// Clients delegate clients command implementation in sentel
func (s *ApiService) Clients(ctx context.Context, req *ClientsRequest) (*ClientsReply, error) {
	reply := &ClientsReply{
		Clients: []*ClientInfo{},
		Header:  &ReplyMessageHeader{Success: true},
	}
	mgr := base.GetServiceManager()

	switch req.Category {
	case "list":
		// Get all client information for specified service
		clients := mgr.GetClients(req.Service)
		for _, client := range clients {
			reply.Clients = append(reply.Clients,
				&ClientInfo{
					UserName:     client.UserName,
					CleanSession: client.CleanSession,
					PeerName:     client.PeerName,
					ConnectTime:  client.ConnectTime,
				})
		}
	case "show":
		// Get client information for specified client id
		if client := mgr.GetClient(req.Service, req.ClientId); client != nil {
			reply.Clients = append(reply.Clients,
				&ClientInfo{
					UserName:     client.UserName,
					CleanSession: client.CleanSession,
					PeerName:     client.PeerName,
					ConnectTime:  client.ConnectTime,
				})
		}
	case "kick":
		if err := mgr.KickoffClient(req.Service, req.ClientId); err != nil {
			reply.Header.Success = false
			reply.Header.Reason = fmt.Sprintf("%v", err)
		}
	default:
		return nil, fmt.Errorf("Invalid category:'%s' for Clients api", req.Category)
	}
	return reply, nil
}

// Sessions delegate client sessions command
func (s *ApiService) Sessions(ctx context.Context, req *SessionsRequest) (*SessionsReply, error) {
	mgr := base.GetServiceManager()
	reply := &SessionsReply{
		Header:   &ReplyMessageHeader{Success: true},
		Sessions: []*SessionInfo{},
	}
	switch req.Category {
	case "list":
		sessions := mgr.GetSessions(req.Service, req.Conditions)
		for _, session := range sessions {
			reply.Sessions = append(reply.Sessions,
				&SessionInfo{
					ClientId:           session.ClientId,
					CreatedAt:          session.CreatedAt,
					CleanSession:       session.CleanSession,
					MessageMaxInflight: session.MessageMaxInflight,
					MessageInflight:    session.MessageInflight,
					MessageInQueue:     session.MessageInQueue,
					MessageDropped:     session.MessageDropped,
					AwaitingRel:        session.AwaitingRel,
					AwaitingComp:       session.AwaitingComp,
					AwaitingAck:        session.AwaitingAck,
				})
		}
	case "show":
		session := mgr.GetSession(req.Service, req.ClientId)
		if session != nil {
			reply.Sessions = append(reply.Sessions,
				&SessionInfo{
					ClientId:           session.ClientId,
					CreatedAt:          session.CreatedAt,
					CleanSession:       session.CleanSession,
					MessageMaxInflight: session.MessageMaxInflight,
					MessageInflight:    session.MessageInflight,
					MessageInQueue:     session.MessageInQueue,
					MessageDropped:     session.MessageDropped,
					AwaitingRel:        session.AwaitingRel,
					AwaitingComp:       session.AwaitingComp,
					AwaitingAck:        session.AwaitingAck,
				})
		}
	}
	return reply, nil
}

func (s *ApiService) Topics(ctx context.Context, req *TopicsRequest) (*TopicsReply, error) {
	return nil, nil
}
