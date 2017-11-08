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
	"strings"
	"sync"

	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
)

type HubNodeInfo struct {
	NodeName  string
	NodeIp    string
	CreatedAt string
}

type ServiceCommand int

const (
	ServiceCommandStart = 0
	ServiceCommandStop  = 1
)

type ServiceManager struct {
	sync.Once
	nodeName string                         // Node name
	config   core.Config                    // Global config
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
func NewServiceManager(c core.Config) (*ServiceManager, error) {
	if _serviceManager != nil {
		return _serviceManager, errors.New("NewServiceManager had been called many times")
	}
	mgr := &ServiceManager{
		config:   c,
		chs:      make(map[string]chan ServiceCommand),
		services: make(map[string]Service),
	}
	// Get supported configs
	items := c.MustString("broker", "services")
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

// GetAllProtocolServices() return all protocol services
func (s *ServiceManager) GetAllProtocolServices() []ProtocolService {
	services := []ProtocolService{}
	for _, service := range s.services {
		if p, ok := service.(ProtocolService); ok {
			services = append(services, p)
		}
	}
	return services
}

// GetProtocolServiceByname return protocol services by name
func (s *ServiceManager) GetProtocolServices(name string) []ProtocolService {
	services := []ProtocolService{}
	for k, service := range s.services {
		if strings.IndexAny(k, name) >= 0 {
			if p, ok := service.(ProtocolService); ok {
				services = append(services, p)
			}
		}
	}
	return services
}

// Node info
func (s *ServiceManager) GetNodeInfo() *HubNodeInfo {
	return &HubNodeInfo{}
}

// Version
func (s *ServiceManager) GetVersion() string {
	return serviceManagerVersion
}

// GetStats return server's stats
func (s *ServiceManager) GetStats(serviceName string) map[string]uint64 {
	allstats := NewStats(false)
	services := s.GetProtocolServices(serviceName)

	for _, service := range services {
		stats := service.GetStats()
		allstats.AddStats(stats)
	}
	return allstats.Get()
}

// GetMetrics return server metrics
func (s *ServiceManager) GetMetrics(serviceName string) map[string]uint64 {
	allmetrics := NewMetrics(false)
	services := s.GetProtocolServices(serviceName)

	for _, service := range services {
		metrics := service.GetMetrics()
		allmetrics.AddMetrics(metrics)
	}
	return allmetrics.Get()
}

// GetClients return clients list withspecified service
func (s *ServiceManager) GetClients(serviceName string) []*ClientInfo {
	clients := []*ClientInfo{}
	services := s.GetProtocolServices(serviceName)

	for _, service := range services {
		list := service.GetClients()
		clients = append(clients, list...)
	}
	return clients
}

// GeteClient return client info with specified client id
func (s *ServiceManager) GetClient(serviceName string, id string) *ClientInfo {
	services := s.GetProtocolServices(serviceName)

	for _, service := range services {
		if client := service.GetClient(id); client != nil {
			return client
		}
	}
	return nil
}

// Kickoff Client killoff a client from specified service
func (s *ServiceManager) KickoffClient(serviceName string, id string) error {
	services := s.GetProtocolServices(serviceName)

	for _, service := range services {
		if err := service.KickoffClient(id); err == nil {
			return nil
		}
	}
	return fmt.Errorf("Failed to kick off user '%s' from service '%s'", id, serviceName)
}

// GetSessions return all sessions information for specified service
func (s *ServiceManager) GetSessions(serviceName string, conditions map[string]bool) []*SessionInfo {
	sessions := []*SessionInfo{}
	services := s.GetProtocolServices(serviceName)

	for _, service := range services {
		list := service.GetSessions(conditions)
		sessions = append(sessions, list...)
	}
	return sessions

}

// GetSession return specified session information with session id
func (s *ServiceManager) GetSession(serviceName string, id string) *SessionInfo {
	services := s.GetProtocolServices(serviceName)

	for _, service := range services {
		if info := service.GetSession(id); info != nil {
			return info
		}
	}
	return nil
}

// GetRoutes return route table information for specified service
func (s *ServiceManager) GetRoutes(serviceName string) []*RouteInfo {
	routes := []*RouteInfo{}
	services := s.GetProtocolServices(serviceName)

	for _, service := range services {
		list := service.GetRoutes()
		routes = append(routes, list...)
	}
	return routes
}

// GetRoute return route information for specified topic
func (s *ServiceManager) GetRoute(serviceName string, topic string) *RouteInfo {
	services := s.GetProtocolServices(serviceName)

	for _, service := range services {
		route := service.GetRoute(topic)
		if route != nil {
			return route
		}
	}
	return nil
}

// GetTopics return topic informaiton for specified service
func (s *ServiceManager) GetTopics(serviceName string) []*TopicInfo {
	topics := []*TopicInfo{}
	services := s.GetProtocolServices(serviceName)

	for _, service := range services {
		list := service.GetTopics()
		topics = append(topics, list...)
	}
	return topics
}

// GetTopic return topic information for specified topic
func (s *ServiceManager) GetTopic(serviceName string, topic string) *TopicInfo {
	services := s.GetProtocolServices(serviceName)

	for _, service := range services {
		info := service.GetTopic(topic)
		if info != nil {
			return info
		}
	}
	return nil
}

// GetSubscriptions return subscription informaiton for specified service
func (s *ServiceManager) GetSubscriptions(serviceName string) []*SubscriptionInfo {
	subs := []*SubscriptionInfo{}
	services := s.GetProtocolServices(serviceName)

	for _, service := range services {
		list := service.GetSubscriptions()
		subs = append(subs, list...)
	}
	return subs
}

// GetSubscription return subscription information for specified topic
func (s *ServiceManager) GetSubscription(serviceName string, sub string) *SubscriptionInfo {
	services := s.GetProtocolServices(serviceName)

	for _, service := range services {
		info := service.GetSubscription(sub)
		if info != nil {
			return info
		}
	}
	return nil
}

// GetAllServiceInfo return all service information
func (s *ServiceManager) GetAllServiceInfo() []*ServiceInfo {
	services := []*ServiceInfo{}

	for _, service := range s.services {
		services = append(services, service.Info())
	}
	return services
}
