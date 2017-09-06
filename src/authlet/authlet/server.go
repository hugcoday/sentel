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
	config libs.Config
}

func LaunchAuthServer(c libs.Config, wg sync.WaitGroup) error {
	address := ":50051"
	if addr, err := c.String("authlet", "address"); err == nil && address != "" {
		address = addr
	}

	lis, err := net.Listen("tcp", address)
	if err != nil {
		glog.Fatal("Failed to listen: %v", err)
		return err
	}
	s := grpc.NewServer()
	RegisterAuthletServer(s, &AuthServer{})
	reflection.Register(s)
	return s.Serve(lis)
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
