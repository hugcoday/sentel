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

package sentel

import (
	"errors"
	"fmt"

	"github.com/golang/glog"
)

type Service interface {
	Start() error
	Stop()
}

type ServiceFactory interface {
	New(name string, c Config, ch chan ServiceCommand) (Service, error)
}

var (
	_serviceFactories = make(map[string]ServiceFactory)
)

// RegisterService register service with name and protocol specified
func RegisterService(name string, configs map[string]string, factory ServiceFactory) {
	if _, ok := _serviceFactories[name]; ok {
		glog.Errorf("Service '%s' is not registered", name)
	}
	RegisterConfig(name, configs)
	_serviceFactories[name] = factory
}

// CreateService create service instance according to service name
func CreateService(name string, c Config, ch chan ServiceCommand) (Service, error) {
	glog.Infof("Creating service '%s'...", name)

	if factory, ok := _serviceFactories[name]; ok && factory != nil {
		return _serviceFactories[name].New(name, c, ch)
	}
	return nil, fmt.Errorf("Invalid service '%s'", name)
}

// CheckAllRegisteredServices check all registered service simplily
func CheckAllRegisteredServices() error {
	if len(_serviceFactories) == 0 {
		return errors.New("No service registered")
	}
	for name, _ := range _serviceFactories {
		glog.Infof("Service '%s' is registered", name)
	}
	return nil
}
