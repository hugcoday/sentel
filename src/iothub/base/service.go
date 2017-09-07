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
	"strings"

	"github.com/golang/glog"
)

var (
	_serviceFactories map[string]ServiceFactory = make(map[string]ServiceFactory)
)

type Service interface {
	Run() error
}

type ServiceFactory interface {
	New(protocol string, c libs.Config, ch chan ServiceCommand) (Service, error)
}

func RegisterService(name string, protocol string, configs map[string]string, factory ServiceFactory) {
	s := name + ":" + protocol
	libs.RegisterConfig(s, configs)
	_serviceFactories[s] = factory
}

func CreateService(name string, c libs.Config, ch chan ServiceCommand) (Service, error) {
	if _serviceFactories[name] == nil {
		return nil, fmt.Errorf("Service '%s' is not registered", name)
	}
	ps := strings.Split(name, ":")
	if len(ps) != 2 {
		return nil, fmt.Errorf("Service '%s' is not rightly configured", name)
	}
	return _serviceFactories[name].New(ps[1], c, ch)
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
