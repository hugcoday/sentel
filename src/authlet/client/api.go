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
	pb "authlet/authlet"
	"fmt"
	"strconv"

	"github.com/golang/glog"

	"google.golang.org/grpc"

	"golang.org/x/net/context"
)

type AuthOptions map[string]string

const (
	AclActionNone  = 0
	AclActionRead  = 1
	AclActionWrite = 2
)

// AuthPlugin interface for security
type AuthletApi struct {
	opts   AuthOptions
	client pb.AuthletClient
}

func (auth *AuthletApi) GetVersion(ctx context.Context) int {
	reply, _ := auth.client.GetVersion(ctx, &pb.AuthRequest{})
	version, _ := strconv.Atoi(reply.Version)
	return version
}
func (auth *AuthletApi) CheckAcl(ctx context.Context, clientid string, username string, topic string, access int) error {
	return nil
}

func (auth *AuthletApi) CheckUsernameAndPasswor(ctx context.Context, username string, password string) error {
	return nil
}
func (auth *AuthletApi) GetPskKey(ctx context.Context, hint string, identity string) (string, error) {
	return "", nil
}

func (auth *AuthletApi) Cleanup(ctx context.Context) error {
	return nil
}

func New(opts AuthOptions) (*AuthletApi, error) {
	address := ""
	api := &AuthletApi{}

	if address, ok := opts["address"]; !ok || address == "" {
		return nil, fmt.Errorf("Invalid autlet address:'%s'", address)
	}
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		glog.Fatalf("Failed to connect with authlet:%s", err)
		return nil, err
	}
	defer conn.Close()
	c := pb.NewAuthletClient(conn)
	api.client = c
	return api, nil
}
