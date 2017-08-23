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
	"fmt"

	"github.com/golang/glog"
)

var (
	serviceFactories map[string]ServiceFactory = make(map[string]ServiceFactory)
)

type Service interface {
	Run() error
}

type ServiceFactory interface {
	New(c *Config, ch chan int) (Service, error)
}

func RegisterServiceFactory(name string, factory ServiceFactory) {
	serviceFactories[name] = factory
}

func CreateService(name string, c *Config, ch chan int) (Service, error) {
	if serviceFactories[name] == nil {
		glog.Error("Service %s is not registered", name)
		return nil, fmt.Errorf("Service %s is not registered", name)
	}
	return serviceFactories[name].New(c, ch)
}
