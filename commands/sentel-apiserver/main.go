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

	"github.com/cloustone/sentel/apiserver"
	"github.com/cloustone/sentel/core"

	"github.com/cloustone/sentel/apiserver/db"
	v1api "github.com/cloustone/sentel/apiserver/v1"

	"github.com/golang/glog"
)

var (
	configFileFullPath = flag.String("c", "../etc/sentel/apiserver.conf", "config file")
)

func main() {
	flag.Parse()
	glog.Info("Starting api server...")

	// Get configuration
	config, err := core.NewWithConfigFile(*configFileFullPath)
	if err != nil {
		glog.Fatal(err)
		flag.PrintDefaults()
		return
	}

	// Initialize registry
	if err := db.InitializeRegistry(config); err != nil {
		glog.Errorf("Failed to initialize registry:%v", err)
		return
	}
	glog.Infof("Registry is initialized successfuly")

	// Create api manager using configuration
	apiManager, err := apiserver.GetApiManager(config)
	if err != nil {
		glog.Error("%v", err)
		return
	}
	glog.Error(apiManager.Run())
}

func init() {
	for group, values := range defaultConfigs {
		core.RegisterConfig(group, values)
	}
	apiserver.RegisterApiManager(v1api.NewApiManager())
}
