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

import "context"

type authletAuthPlugin struct{}

func (n *authletAuthPlugin) GetVersion() int {
	return 0
}
func (n *authletAuthPlugin) Cleanup(ctx context.Context) error {
	return nil
}
func (n *authletAuthPlugin) CheckAcl(ctx context.Context, clientid string, username string, topic string, access int) error {
	return nil
}
func (n *authletAuthPlugin) CheckUsernameAndPasswor(ctx context.Context, username string, password string) error {
	return nil
}
func (n *authletAuthPlugin) GetPskKey(ctx context.Context, hint string, identity string) (string, error) {
	return "", nil
}

// AuthPluginFactory
type authletAuthPluginFactory struct{}

func (n authletAuthPluginFactory) New(ctx context.Context, opts AuthOptions) (AuthPlugin, error) {
	return &authletAuthPlugin{}, nil
}
