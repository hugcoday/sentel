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

package security

import (
	"authlet/authlet"
	"context"
	"libs"
)

type authletAuthPlugin struct {
	rpcapi *authlet.AuthletApi
}

func (n *authletAuthPlugin) GetVersion(ctx context.Context) int {
	return n.rpcapi.GetVersion(ctx)
}
func (n *authletAuthPlugin) Cleanup(ctx context.Context) {
	n.rpcapi.Close()
}

func (n *authletAuthPlugin) CheckAcl(ctx context.Context, clientid string, username string, topic string, access int) error {
	action := ""
	switch access {
	case AclActionRead:
		action = "r"
	case AclActionWrite:
		action = "w"
	}
	return n.rpcapi.CheckAcl(ctx, clientid, username, topic, action)
}
func (n *authletAuthPlugin) CheckUserNameAndPassword(ctx context.Context, username string, password string) error {
	return n.rpcapi.CheckUserNameAndPassword(ctx, username, password)
}
func (n *authletAuthPlugin) GetPskKey(ctx context.Context, hint string, identity string) (string, error) {
	return n.rpcapi.GetPskKey(ctx, hint, identity)
}

// AuthPluginFactory
type authletAuthPluginFactory struct{}

func (n authletAuthPluginFactory) New(ctx context.Context, c libs.Config) (AuthPlugin, error) {
	plugin := &authletAuthPlugin{}
	if rpcapi, err := authlet.New(c); err != nil {
		return nil, err
	} else {
		plugin.rpcapi = rpcapi
	}
	return plugin, nil
}
