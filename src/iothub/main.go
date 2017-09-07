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
	"iothub/base"
	"iothub/coap"
	"iothub/mqtt"
	"libs"

	"github.com/golang/glog"
)

var (
	configFileFullPath = flag.String("c", "../etc/sentel/iothub.conf", "config file")
)

func main() {
	var mgr *base.ServiceManager
	var config libs.Config
	var err error

	flag.Parse()
	glog.Info("Starting iothub server...")

	// Check all registered service
	if err := base.CheckAllRegisteredServices(); err != nil {
		glog.Fatal(err)
		return
	}
	// Get configuration
	if config, err = libs.NewWithConfigFile(*configFileFullPath); err != nil {
		glog.Fatal(err)
		flag.PrintDefaults()
		return
	}
	// Create service manager according to the configuration
	if mgr, err = base.NewServiceManager(config); err != nil {
		glog.Fatal("Failed to launch ServiceManager")
		return
	}
	glog.Error(mgr.Start())
}

func init() {
	for group, values := range allDefaultConfigs {
		libs.RegisterConfig(group, values)
	}
	base.RegisterService("mqtt", "tcp", mqtt.Configs, &mqtt.MqttFactory{})
	base.RegisterService("coap", "udp", coap.Configs, &coap.CoapFactory{})
	//	base.RegisterService("api", "rpc", coap.Configs, &api.apiFactory{})
	//	base.RegisterService("dashboard", "http", dashboard.Configs, &dashboard.dashboardFactory{})
}
