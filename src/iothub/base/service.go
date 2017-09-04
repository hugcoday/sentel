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

package base

import (
	"errors"
	"fmt"
	"iothub/security"
	"iothub/storage"
	"net"

	"github.com/golang/glog"
)

var (
	_serviceFactories map[string]ServiceFactory = make(map[string]ServiceFactory)
)

type Service interface {
	// NewSession create a new session
	NewSession(conn net.Conn) (Session, error)
	// Run is mainloop of current service
	Run() error
	// GetSessionTotalCount get total session count
	GetSessionTotalCount() int64
	// CreateSessionId create identifier for new session
	CreateSessionId() string
	// RegisterSession register a new session
	RegisterSession(s Session)
	// RemoveSession remove session based sessionid
	RemoveSession(s Session)
}

type ServiceFactory interface {
	New(c Config, ch chan int) (Service, error)
}

func RegisterService(name string, configs map[string]string, factory ServiceFactory) {
	RegisterConfig(name, configs)
	_serviceFactories[name] = factory
}

func CreateService(name string, c Config, ch chan int, d storage.Storage) (Service, error) {
	if _serviceFactories[name] == nil {
		return nil, fmt.Errorf("Service '%s' is not registered", name)
	}
	return _serviceFactories[name].New(c, ch)
}

func CheckAllRegisteredServices() error {
	if len(_serviceFactories) == 0 {
		return errors.New("No service registered")
	}
	for name, _ := range _serviceFactories {
		glog.Infof("Service '%s' is registered", name)
	}
	return nil
}

// LoadAuthPluginWithConfig load a authPlugin with config
func LoadAuthPluginWithConfig(service string, c Config) (security.AuthPlugin, error) {
	// Get authentication for this service, if service's authentication is not
	// specified, using iothub's authentication
	auth, err := c.String(service, "authentication")
	if err != nil {
		if auth, err = c.String("security", "authentication"); err != nil {
			return nil, fmt.Errorf("Authentication method is not specified for service '%s'", service)
		}
	}
	opts := security.AuthOptions{}
	return security.LoadAuthPlugin(auth, opts)
}
