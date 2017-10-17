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
	pb "github.com/cloustone/sentel/iothub/api"

	"github.com/golang/glog"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

type SentelApi struct {
	rpcapi pb.ApiClient
	conn   *grpc.ClientConn
}

func NewSentelApi() (*SentelApi, error) {
	address := "localhost:50052"
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		glog.Fatalf("Failed to connect with iothub:%s", err)
		return nil, err
	}
	c := pb.NewApiClient(conn)
	return &SentelApi{conn: conn, rpcapi: c}, nil
}

func (s *SentelApi) Version(in *pb.VersionRequest) (*pb.VersionReply, error) {
	return s.rpcapi.Version(context.Background(), in)
}
func (s *SentelApi) Admins(in *pb.AdminsRequest) (*pb.AdminsReply, error) {
	return s.rpcapi.Admins(context.Background(), in)
}

func (s *SentelApi) Cluster(in *pb.ClusterRequest) (*pb.ClusterReply, error) {
	return s.rpcapi.Cluster(context.Background(), in)
}

func (s *SentelApi) Routes(in *pb.RoutesRequest) (*pb.RoutesReply, error) {
	return s.rpcapi.Routes(context.Background(), in)
}

func (s *SentelApi) Status(in *pb.StatusRequest) (*pb.StatusReply, error) {
	return s.rpcapi.Status(context.Background(), in)
}

func (s *SentelApi) Broker(in *pb.BrokerRequest) (*pb.BrokerReply, error) {
	return s.rpcapi.Broker(context.Background(), in)
}

func (s *SentelApi) Plugins(in *pb.PluginsRequest) (*pb.PluginsReply, error) {
	return s.rpcapi.Plugins(context.Background(), in)
}

func (s *SentelApi) Services(in *pb.ServicesRequest) (*pb.ServicesReply, error) {
	return s.rpcapi.Services(context.Background(), in)
}

func (s *SentelApi) Subscriptions(in *pb.SubscriptionsRequest) (*pb.SubscriptionsReply, error) {
	return s.rpcapi.Subscriptions(context.Background(), in)
}

func (s *SentelApi) Clients(in *pb.ClientsRequest) (*pb.ClientsReply, error) {
	return s.rpcapi.Clients(context.Background(), in)
}

func (s *SentelApi) Sessions(in *pb.SessionsRequest) (*pb.SessionsReply, error) {
	return s.rpcapi.Sessions(context.Background(), in)
}

func (s *SentelApi) Topics(in *pb.TopicsRequest) (*pb.TopicsReply, error) {
	return s.rpcapi.Topics(context.Background(), in)
}
