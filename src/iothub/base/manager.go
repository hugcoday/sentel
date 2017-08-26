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
	"iothub/db"
	"iothub/util/config"
	"strings"

	"github.com/golang/glog"
)

type ServiceManager struct {
	config   config.Config       // Global config
	services map[string]Service  // All service created by config.Protocols
	chs      map[string]chan int // Notification channel for each service
	db       db.Database
}

// NewServiceManager create ServiceManager in main context
func NewServiceManager(c config.Config) (*ServiceManager, error) {
	mgr := &ServiceManager{
		config:   c,
		chs:      make(map[string]chan int),
		services: make(map[string]Service),
	}
	// Create database instance
	if _, err := c.String("iothub", "database"); err != nil {
		return nil, errors.New("No database is specified in config file")
	}
	database, err := db.NewDatabase(c)
	if err != nil {
		return nil, err
	}
	if err := database.Open(); err != nil {
		return nil, err
	}
	mgr.db = database

	// Get supported configs
	items, err := c.String("iothub", "protocols")
	if err != nil {
		database.Close()
		return nil, fmt.Errorf("Invalid protocol declaration in config file")
	}
	protocols := strings.Split(items, ",")
	// Create service for each protocol
	for _, name := range protocols {
		ch := make(chan int)
		service, err := CreateService(name, c, ch, mgr.db)
		if err != nil {
			database.Close()
			glog.Error("Create service(%s) failed", name)
			return nil, err
		}
		glog.Info("Create service(%s) success", name)
		mgr.services[name] = service
		mgr.chs[name] = ch
	}
	return mgr, nil
}

// ServiceManger run all serice and wait to terminate
func (s *ServiceManager) Start() error {
	defer s.db.Close()
	// Run all service
	for _, service := range s.services {
		go service.Run()
	}
	// Wait all service to terminate in main context
	for name, ch := range s.chs {
		<-ch
		glog.Info("Servide(%s) is terminated", name)
	}
	return nil
}
