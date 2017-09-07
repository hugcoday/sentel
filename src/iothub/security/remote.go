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
	"context"
	"libs"
)

type remoteAuthPlugin struct{}

func (n *remoteAuthPlugin) GetVersion(ctx context.Context) int {
	return 0
}
func (n *remoteAuthPlugin) Cleanup(ctx context.Context) {}
func (n *remoteAuthPlugin) CheckAcl(ctx context.Context, clientid string, username string, topic string, access int) error {
	return nil
}
func (n *remoteAuthPlugin) CheckUserNameAndPassword(ctx context.Context, username string, password string) error {
	return nil
}
func (n *remoteAuthPlugin) GetPskKey(ctx context.Context, hint string, identity string) (string, error) {
	return "", nil
}

// AuthPluginFactory
type remoteAuthPluginFactory struct{}

func (n remoteAuthPluginFactory) New(ctx context.Context, c libs.Config) (AuthPlugin, error) {
	return &remoteAuthPlugin{}, nil
}
