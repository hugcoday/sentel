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
	"libs"
	"strings"
	"sync"

	"github.com/golang/glog"
)

type ServiceCommand int

const (
	ServiceCommandStart = 0
	ServiceCommandStop  = 1
)

type ServiceManager struct {
	sync.Once
	config   libs.Config                    // Global config
	services map[string]Service             // All service created by config.Protocols
	chs      map[string]chan ServiceCommand // Notification channel for each service
}

const serviceManagerVersion = "0.1"

var _serviceManager *ServiceManager

// GetServiceManager create service manager and all supported service
// The function should be called in service
func GetServiceManager() *ServiceManager { return _serviceManager }

// NewServiceManager create ServiceManager only in main context
func NewServiceManager(c libs.Config) (*ServiceManager, error) {
	if _serviceManager != nil {
		return _serviceManager, errors.New("NewServiceManager had been called many times")
	}
	mgr := &ServiceManager{
		config:   c,
		chs:      make(map[string]chan ServiceCommand),
		services: make(map[string]Service),
	}
	// Get supported configs
	items := c.MustString("iothub", "services")
	services := strings.Split(items, ",")
	// Create service for each protocol
	for _, name := range services {
		ch := make(chan ServiceCommand)
		service, err := CreateService(name, c, ch)
		if err != nil {
			glog.Errorf("%s", err)
		} else {
			glog.Infof("Create service '%s' success", name)
			mgr.services[name] = service
			mgr.chs[name] = ch
		}
	}
	_serviceManager = mgr
	return _serviceManager, nil
}

// ServiceManger run all serice and wait to terminate
func (s *ServiceManager) Start() error {
	// Run all service
	for _, service := range s.services {
		go service.Start()
	}
	// Wait all service to terminate in main context
	for name, ch := range s.chs {
		<-ch
		glog.Info("Servide(%s) is terminated", name)
	}
	return nil
}

// Version
func (s *ServiceManager) GetVersion() string {
	return serviceManagerVersion
}
