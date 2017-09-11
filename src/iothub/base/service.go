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
	"libs"

	"github.com/golang/glog"
)

type Service interface {
	Name() string // mqtt:tcp, mqtt:ssl
	Start() error
	Stop()
}

type ServiceFactory interface {
	New(name string, c libs.Config, ch chan ServiceCommand) (Service, error)
}

var (
	_serviceFactories = make(map[string]ServiceFactory)
)

// RegisterService register service with name and protocol specified
func RegisterService(name string, configs map[string]string, factory ServiceFactory) {
	libs.RegisterConfig(name, configs)
	_serviceFactories[name] = factory
}

// CreateService create service instance according to service name
func CreateService(name string, c libs.Config, ch chan ServiceCommand) (Service, error) {
	if _, ok := _serviceFactories[name]; ok {
		return nil, fmt.Errorf("Service '%s' is not registered", name)
	}
	return _serviceFactories[name].New(name, c, ch)
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
