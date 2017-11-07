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
	"net"
	"sync"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/libs/sentel"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ApiService struct {
	config   sentel.Config
	chn      chan base.ServiceCommand
	wg       sync.WaitGroup
	listener net.Listener
	srv      *grpc.Server
}

// ApiServiceFactory
type ApiServiceFactory struct{}

// New create apiService service factory
func (m *ApiServiceFactory) New(protocol string, c sentel.Config, ch chan base.ServiceCommand) (base.Service, error) {
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
func (s *ApiService) Info() *base.ServiceInfo {
	return &base.ServiceInfo{
		ServiceName: "apiservice",
	}
}

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

// Services delegate  services command
func (s *ApiService) Services(ctx context.Context, req *ServicesRequest) (*ServicesReply, error) {
	mgr := base.GetServiceManager()
	reply := &ServicesReply{
		Header:   &ReplyMessageHeader{Success: true},
		Services: []*ServiceInfo{},
	}
	switch req.Category {
	case "list":
		services := mgr.GetAllServiceInfo()
		for _, service := range services {
			reply.Services = append(reply.Services,
				&ServiceInfo{
					ServiceName:    service.ServiceName,
					Listen:         service.Listen,
					Acceptors:      service.Acceptors,
					MaxClients:     service.MaxClients,
					CurrentClients: service.CurrentClients,
					ShutdownCount:  service.ShutdownCount,
				})
		}

	case "start":
	case "stop":
	}

	return nil, nil
}

//Subscriptions delete subscriptions command
func (s *ApiService) Subscriptions(ctx context.Context, req *SubscriptionsRequest) (*SubscriptionsReply, error) {
	mgr := base.GetServiceManager()
	reply := &SubscriptionsReply{
		Header:        &ReplyMessageHeader{Success: true},
		Subscriptions: []*SubscriptionInfo{},
	}
	switch req.Category {
	case "list":
		subs := mgr.GetSubscriptions(req.Service)
		for _, sub := range subs {
			reply.Subscriptions = append(reply.Subscriptions,
				&SubscriptionInfo{
					ClientId:  sub.ClientId,
					Topic:     sub.Topic,
					Attribute: sub.Attribute,
				})
		}
	case "show":
		sub := mgr.GetSubscription(req.Service, req.Subscription)
		if sub != nil {
			reply.Subscriptions = append(reply.Subscriptions,
				&SubscriptionInfo{
					ClientId:  sub.ClientId,
					Topic:     sub.Topic,
					Attribute: sub.Attribute,
				})
		}
	}
	return reply, nil
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
	mgr := base.GetServiceManager()
	reply := &TopicsReply{
		Header: &ReplyMessageHeader{Success: true},
		Topics: []*TopicInfo{},
	}
	switch req.Category {
	case "list":
		topics := mgr.GetTopics(req.Service)
		for _, topic := range topics {
			reply.Topics = append(reply.Topics,
				&TopicInfo{
					Topic:     topic.Topic,
					Attribute: topic.Attribute,
				})
		}
	case "show":
		topic := mgr.GetTopic(req.Service, req.Topic)
		if topic != nil {
			reply.Topics = append(reply.Topics,
				&TopicInfo{
					Topic:     topic.Topic,
					Attribute: topic.Attribute,
				})
		}
	}
	return reply, nil
}
