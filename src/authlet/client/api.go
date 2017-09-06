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
	"errors"
	"fmt"
	"strconv"

	"github.com/golang/glog"

	"google.golang.org/grpc"

	"golang.org/x/net/context"
)

type AuthOptions map[string]string

const (
	AclActionNameNone  = ""
	AclActionNameRead  = "r"
	AclActionNameWrite = "w"
)

// AuthPlugin interface for security
type AuthletApi struct {
	opts   AuthOptions
	client pb.AuthletClient
	conn   *grpc.ClientConn
}

func (auth *AuthletApi) GetVersion(ctx context.Context) int {
	reply, _ := auth.client.GetVersion(ctx, &pb.AuthRequest{})
	version, _ := strconv.Atoi(reply.Version)
	return version
}

func (auth *AuthletApi) CheckAcl(ctx context.Context, clientid string, username string, topic string, access string) error {
	reply, err := auth.client.CheckAcl(ctx, &pb.AuthRequest{
		Clientid: clientid,
		Username: username,
		Topic:    topic,
		Access:   access,
	})
	if err != nil || reply.Result != true {
		return errors.New("Acl denied")
	}
	return nil
}

func (auth *AuthletApi) CheckUserNameAndPasswor(ctx context.Context, username string, password string) error {
	reply, err := auth.client.CheckUserNameAndPassword(ctx, &pb.AuthRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return err
	}
	if reply.Result != true {
		return errors.New("Acl denied")
	}
	return nil
}

func (auth *AuthletApi) GetPskKey(ctx context.Context, hint string, identity string) (string, error) {
	reply, err := auth.client.GetPskKey(ctx, &pb.AuthRequest{
		Hint:     hint,
		Username: identity,
	})
	return reply.Key, err
}

func (auth *AuthletApi) Close() {
	auth.conn.Close()
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
	api.client = pb.NewAuthletClient(conn)
	api.conn = conn
	return api, nil
}
