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

package auth

import (
	"net"
	"sync"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/core"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type AuthService struct {
	config   core.Config
	wg       sync.WaitGroup
	listener net.Listener
	srv      *grpc.Server
}

// NewAuthService create authentication server
func NewAuthService(c core.Config) (*AuthService, error) {
	address := ":50051"
	server := &AuthService{config: c, wg: sync.WaitGroup{}}

	if addr, err := c.String("auth", "address"); err == nil && address != "" {
		address = addr
	}

	lis, err := net.Listen("tcp", address)
	if err != nil {
		glog.Fatal("Failed to listen: %v", err)
		return nil, err
	}
	server.listener = lis
	server.srv = grpc.NewServer()
	RegisterAuthServiceServer(server.srv, server)
	reflection.Register(server.srv)
	return server, nil
}

// Info
func (s *AuthService) Info() *base.ServiceInfo {
	return &base.ServiceInfo{
		ServiceName: "auth",
	}
}

// Start
func (s *AuthService) Start() {
	go func(s *AuthService) {
		s.srv.Serve(s.listener)
		s.wg.Add(1)
	}(s)
}

// Stop
func (s *AuthService) Stop() {
	s.listener.Close()
	s.wg.Wait()
}

// Wait
func (s *AuthService) Wait() {
	s.wg.Wait()
}

// Get version of Authlet service
func (s *AuthService) GetVersion(context.Context, *AuthRequest) (*AuthReply, error) {
	return nil, nil
}

// Check acl based on client id, user name and topic
func (s *AuthService) CheckAcl(context.Context, *AuthRequest) (*AuthReply, error) {
	return nil, nil
}

// Check user name and password
func (s *AuthService) CheckUserNameAndPassword(context.Context, *AuthRequest) (*AuthReply, error) {
	return nil, nil
}

// Get PSK key
func (s *AuthService) GetPskKey(context.Context, *AuthRequest) (*AuthReply, error) {
	return nil, nil
}
