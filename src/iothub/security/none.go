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

type noneAuthPlugin struct{}

func (n *noneAuthPlugin) GetVersion() int {
	return 0
}
func (n *noneAuthPlugin) Initialize(data interface{}, opts AuthOptions) error {
	return nil
}
func (n *noneAuthPlugin) Cleanup(data interface{}, opts AuthOptions) error {
	return nil
}
func (n *noneAuthPlugin) InitializeSecurity(data interface{}, opts AuthOptions) error {
	return nil
}
func (n *noneAuthPlugin) CleanupSecurity(data interface{}, opts AuthOptions) error {
	return nil
}
func (n *noneAuthPlugin) CheckAcl(data interface{}, clientid string, username string, topic string, access int) error {
	return nil
}
func (n *noneAuthPlugin) CheckUsernameAndPasswor(data interface{}, username string, password string) error {
	return nil
}
func (n *noneAuthPlugin) GetPskKey(data interface{}, hint string, identity string) (string, error) {
	return "", nil
}

// AuthPluginFactory
type noneAuthPluginFactory struct{}

func (n noneAuthPluginFactory) New(opts AuthOptions) (AuthPlugin, error) {
	return &noneAuthPlugin{}, nil
}
