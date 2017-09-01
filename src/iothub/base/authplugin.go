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

package base

import (
	"errors"
	"fmt"

	"github.com/golang/glog"
)

type AuthOption struct{}

const (
	AclActionNone  = 0
	AclActionRead  = 1
	AclActionWrite = 2
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
	Initialize(data interface{}, options []AuthOption) error
	Cleanup(data interface{}, options []AuthOption) error
	InitializeSecurity(data interface{}, options []AuthOption) error
	CleanupSecurity(data interface{}, options []AuthOption) error
	CheckAcl(data interface{}, clientid string, username string, topic string, access int) error
	CheckUsernameAndPasswor(data interface{}, username string, password string) error
	GetPskKey(data interface{}, hint string, identity string) (string, error)
}

// AuthPluginFactory
type AuthPluginFactory interface {
	New(options []AuthOption) (AuthPlugin, error)
}

// RegisterAuthPlugin register a auth plugin
func RegisterAuthPlugin(name string, factory AuthPluginFactory) {
	if _authPlugins[name] != nil {
		glog.Fatalf("AluthPlugin %s already registered")
		return
	}
	_authPlugins[name] = factory
}

// LoadAuthPlugin load a authPlugin
func LoadAuthPlugin(name string, options []AuthOption) (AuthPlugin, error) {
	if _authPlugins[name] == nil {
		return nil, fmt.Errorf("AuthPlugin %s does't exist", name)
	}
	return _authPlugins[name].New(options)
}

// LoadAuthPluginWithConfig load a authPlugin with config
func LoadAuthPluginWithConfig(service string, c Config) (AuthPlugin, error) {
	auth, _ := c.String(service, "authentication")
	if _authPlugins[auth] == nil {
		glog.Errorf("AuthPlugin for service %s does't exist", service)
	} else {
		auth, _ := c.String("iothub", "authentication")
		if _authPlugins[auth] == nil {
			glog.Errorf("AuthPlugin %s are not registered", service)
			return nil, fmt.Errorf("AuthPlugin %s does't exist", auth)
		}
	}
	return _authPlugins[auth].New(nil)
}
