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

package authlet

import (
	"libs"
	"net"
	"sync"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type AuthServer struct {
	config   libs.Config
	wg       sync.WaitGroup
	listener net.Listener
	srv      *grpc.Server
}

// NewAuthServer create authentication server
func NewAuthServer(c libs.Config) (*AuthServer, error) {
	address := ":50051"
	server := &AuthServer{config: c, wg: sync.WaitGroup{}}

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
	RegisterAuthletServer(server.srv, server)
	reflection.Register(server.srv)
	return server, nil
}

// Start
func (s *AuthServer) Start() {
	go func(s *AuthServer) {
		s.srv.Serve(s.listener)
		s.wg.Add(1)
	}(s)
}

// Stop
func (s *AuthServer) Stop() {
	s.listener.Close()
	s.wg.Wait()
}

// Wait
func (s *AuthServer) Wait() {
	s.wg.Wait()
}

// Get version of Authlet service
func (s *AuthServer) GetVersion(context.Context, *AuthRequest) (*AuthReply, error) {
	return nil, nil
}

// Check acl based on client id, user name and topic
func (s *AuthServer) CheckAcl(context.Context, *AuthRequest) (*AuthReply, error) {
	return nil, nil
}

// Check user name and password
func (s *AuthServer) CheckUserNameAndPassword(context.Context, *AuthRequest) (*AuthReply, error) {
	return nil, nil
}

// Get PSK key
func (s *AuthServer) GetPskKey(context.Context, *AuthRequest) (*AuthReply, error) {
	return nil, nil
}
