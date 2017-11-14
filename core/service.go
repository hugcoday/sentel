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

package core

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/golang/glog"
)

type ServiceCommand int

const (
	ServiceCommandStart = 0
	ServiceCommandStop  = 1
)

type Service interface {
	Name() string
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

type ServiceManager struct {
	sync.Once
	nodeName string                         // Node name
	config   Config                         // Global config
	services map[string]Service             // All service created by config.Protocols
	chs      map[string]chan ServiceCommand // Notification channel for each service
}

const serviceManagerVersion = "0.1"

var (
	_serviceManager *ServiceManager
)

// GetServiceManager create service manager and all supported service
// The function should be called in service
func GetServiceManager() *ServiceManager { return _serviceManager }

// NewServiceManager create ServiceManager only in main context
func NewServiceManager(name string, c Config) (*ServiceManager, error) {
	if _serviceManager != nil {
		return _serviceManager, errors.New("NewServiceManager had been called many times")
	}
	mgr := &ServiceManager{
		config:   c,
		chs:      make(map[string]chan ServiceCommand),
		services: make(map[string]Service),
	}
	// Get supported configs
	items := c.MustString(name, "services")
	services := strings.Split(items, ",")
	// Create service for each protocol
	for _, name := range services {
		// Format service name
		name = strings.Trim(name, " ")
		ch := make(chan ServiceCommand)
		service, err := CreateService(name, c, ch)
		if err != nil {
			glog.Errorf("%s", err)
		} else {
			glog.Infof("Create service '%s' successfully", name)
			mgr.services[name] = service
			mgr.chs[name] = ch
		}
	}
	_serviceManager = mgr
	return _serviceManager, nil
}

// Run launch all serices and wait to terminate
func (s *ServiceManager) Run() error {
	if err := CheckAllRegisteredServices(); err != nil {
		return err
	}
	// Run all service
	glog.Infof("There are %d service in iothub", len(s.services))
	for _, service := range s.services {
		glog.Infof("Starting service:'%s'...", service.Name())
		go service.Start()
	}
	// Wait all service to terminate in main context
	for name, ch := range s.chs {
		<-ch
		glog.Info("Servide(%s) is terminated", name)
	}
	return nil
}

// StartService launch specified service
func (s *ServiceManager) StartService(name string) error {
	// Return error if service has already been started
	for id, service := range s.services {
		if strings.IndexAny(id, name) >= 0 && service != nil {
			return fmt.Errorf("The service '%s' has already been started", name)
		}
	}
	ch := make(chan ServiceCommand)
	service, err := CreateService(name, s.config, ch)
	if err != nil {
		glog.Errorf("%s", err)
	} else {
		glog.Infof("Create service '%s' success", name)
		s.services[name] = service
		s.chs[name] = ch
	}
	return nil
}

// StopService stop specified service
func (s *ServiceManager) StopService(id string) error {
	for name, service := range s.services {
		if name == id && service != nil {
			service.Stop()
			s.services[name] = nil
			close(s.chs[name])
		}
	}
	return nil
}

// GetServicesByName return service instance by name, or matched by part of name
func (s *ServiceManager) GetServicesByName(name string) []Service {
	services := []Service{}
	for k, service := range s.services {
		if strings.IndexAny(k, name) >= 0 {
			services = append(services, service)
		}
	}
	return services
}
