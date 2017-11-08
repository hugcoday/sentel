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

	"github.com/cloustone/sentel/ceilometer/api"
	"github.com/cloustone/sentel/ceilometer/collector"
	"github.com/cloustone/sentel/core"

	"github.com/golang/glog"
)

var (
	configFileFullPath = flag.String("c", "../ceilometer/ceilometer.conf", "config file")
)

func main() {
	var mgr *core.ServiceManager
	var config core.Config
	var err error

	flag.Parse()
	glog.Info("Starting ceilometer server...")

	// Check all registered service
	if err := core.CheckAllRegisteredServices(); err != nil {
		glog.Fatal(err)
		return
	}
	// Get configuration
	if config, err = core.NewWithConfigFile(*configFileFullPath); err != nil {
		glog.Fatal(err)
		flag.PrintDefaults()
		return
	}
	// Create service manager according to the configuration
	if mgr, err = core.NewServiceManager("ceilometer", config); err != nil {
		glog.Fatal("Failed to launch ServiceManager")
		return
	}
	glog.Error(mgr.Start())
}

func init() {
	for group, values := range defaultConfigs {
		core.RegisterConfig(group, values)
	}
	core.RegisterService("api", api.Configs, &api.ApiServiceFactory{})
	core.RegisterService("collector", collector.Configs, &collector.CollectorServiceFactory{})
}
