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

package main

import (
	"flag"

	"github.com/cloustone/sentel/libs/sentel"
	"github.com/golang/glog"
)

var (
	configFileFullPath = flag.String("c", "../conductor/conductor.conf", "config file")
)

func main() {
	var mgr *sentel.ServiceManager
	var config sentel.Config
	var err error

	flag.Parse()
	glog.Info("Starting condutor server...")

	// Check all registered service
	if err := sentel.CheckAllRegisteredServices(); err != nil {
		glog.Fatal(err)
		return
	}
	// Get configuration
	if config, err = sentel.NewWithConfigFile(*configFileFullPath); err != nil {
		glog.Fatal(err)
		flag.PrintDefaults()
		return
	}
	// Create service manager according to the configuration
	if mgr, err = sentel.NewServiceManager("conductor", config); err != nil {
		glog.Fatal("Failed to launch conductor ServiceManager")
		return
	}
	glog.Error(mgr.Start())
}

func init() {
	for group, values := range allDefaultConfigs {
		sentel.RegisterConfig(group, values)
	}
	//	sentel.RegisterService("api", api.Configs, &api.ApiServiceFactory{})
	//	sentel.RegisterService("collector", collector.Configs, &collector.CollectorServiceFactory{})
}
