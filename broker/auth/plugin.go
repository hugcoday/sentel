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
	"context"
	"errors"
	"fmt"

	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
)

var (
	ErrorAclDenied = errors.New("Acl denied")
)

var (
	_authPlugins = make(map[string]AuthPluginFactory)
)

// AuthPlugin interface for security
type AuthPlugin interface {
	GetVersion() int
	CheckAcl(ctx context.Context, clientid string, username string, topic string, access int) error
	CheckUsernameAndPasswor(ctx context.Context, username string, password string) error
	GetPskKey(ctx context.Context, hint string, identity string) (string, error)
	Cleanup(ctx context.Context) error
}

// AuthPluginFactory
type AuthPluginFactory interface {
	New(c core.Config) (AuthPlugin, error)
}

// RegisterAuthPlugin register a auth plugin
func RegisterAuthPlugin(name string, factory AuthPluginFactory) {
	if _authPlugins[name] != nil {
		glog.Errorf("AuthPlugin '%s' is already registered")
		return
	}
	_authPlugins[name] = factory
}

// LoadAuthPlugin load a authPlugin
func LoadAuthPlugin(name string, c core.Config) (AuthPlugin, error) {
	// Default authentication is 'none'
	if name == "" {
		glog.Warning("No authentication method is specified, using none authentication")
		name = "none"
	}
	if _authPlugins[name] == nil {
		return nil, fmt.Errorf("AuthPlugin '%s' is not registered", name)
	}
	return _authPlugins[name].New(c)
}
