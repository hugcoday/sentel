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
	"errors"
	"fmt"
	"libs"
	"strconv"

	"github.com/golang/glog"

	"google.golang.org/grpc"

	"golang.org/x/net/context"
)

const (
	AclActionNone  = ""
	AclActionead   = "r"
	AclActionWrite = "w"
)

// AuthPlugin interface for security
type AuthletApi struct {
	config libs.Config
	client AuthletClient
	conn   *grpc.ClientConn
}

func (auth *AuthletApi) GetVersion(ctx context.Context) int {
	reply, _ := auth.client.GetVersion(ctx, &AuthRequest{})
	version, _ := strconv.Atoi(reply.Version)
	return version
}

func (auth *AuthletApi) CheckAcl(ctx context.Context, clientid string, username string, topic string, access string) error {
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

func (auth *AuthletApi) CheckUserNameAndPassword(ctx context.Context, username string, password string) error {
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

func (auth *AuthletApi) GetPskKey(ctx context.Context, hint string, identity string) (string, error) {
	reply, err := auth.client.GetPskKey(ctx, &AuthRequest{
		Hint:     hint,
		Username: identity,
	})
	return reply.Key, err
}

func (auth *AuthletApi) Close() {
	auth.conn.Close()
}

func New(c libs.Config) (*AuthletApi, error) {
	address := ""
	api := &AuthletApi{config: c}

	if address, err := c.String("authlet", "address"); err != nil || address == "" {
		return nil, fmt.Errorf("Invalid autlet address:'%s'", address)
	}
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		glog.Fatalf("Failed to connect with authlet:%s", err)
		return nil, err
	}
	api.client = NewAuthletClient(conn)
	api.conn = conn
	return api, nil
}
