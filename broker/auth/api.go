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
	"errors"
	"fmt"
	"strconv"

	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"

	"google.golang.org/grpc"

	"golang.org/x/net/context"
)

const (
	AclActionNone  = ""
	AclActionRead  = "r"
	AclActionWrite = "w"
)

// IAuthAPI ...
type IAuthAPI interface {
	GetVersion(ctx context.Context) int
	CheckAcl(ctx context.Context, clientid string, username string, topic string, access string) error
	CheckUserNameAndPassword(ctx context.Context, username string, password string) error
	GetPskKey(ctx context.Context, hint string, identity string) (string, error)
}

// AuthPlugin interface for security
type AuthApi struct {
	config core.Config
	client AuthServiceClient
	conn   *grpc.ClientConn
}

func (auth *AuthApi) GetVersion(ctx context.Context) int {
	reply, _ := auth.client.GetVersion(ctx, &AuthRequest{})
	version, _ := strconv.Atoi(reply.Version)
	return version
}

func (auth *AuthApi) CheckAcl(ctx context.Context, clientid string, username string, topic string, access string) error {
	reply, err := auth.client.CheckAcl(ctx, &AuthRequest{
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

func (auth *AuthApi) CheckUserNameAndPassword(ctx context.Context, username string, password string) error {
	reply, err := auth.client.CheckUserNameAndPassword(ctx, &AuthRequest{
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

func (auth *AuthApi) GetPskKey(ctx context.Context, hint string, identity string) (string, error) {
	reply, err := auth.client.GetPskKey(ctx, &AuthRequest{
		Hint:     hint,
		Username: identity,
	})
	return reply.Key, err
}

func (auth *AuthApi) Close() {
	auth.conn.Close()
}

func NewAuthApi(c core.Config) (IAuthAPI, error) {
	address := ""
	if address, err := c.String("auth", "address"); err != nil || address == "" {
		return nil, fmt.Errorf("Invalid autlet address:'%s'", address)
	}

	if address == "dummy" {
		return NewDummyAuthService(), nil
	}

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		glog.Fatalf("Failed to connect with authagent:%s", err)
		return nil, err
	}

	api := &AuthApi{config: c}
	api.client = NewAuthServiceClient(conn)
	api.conn = conn
	return api, nil
}
