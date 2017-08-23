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
	"apiserver/base"
	"apiserver/db"

	"github.com/golang/glog"
)

func main() {
	// Get configuration
	apiconfig, err := base.NewApiConfig()
	defer apiconfig.Close()

	// Initialize registry store
	if err := db.InitializeRegistry(apiconfig); err != nil {
		glog.Error("Registry initialization failed:%v", err)
		return
	}

	// Create api manager using configuration
	apiManager, err := base.CreateApiManager(apiconfig)
	if err != nil {
		glog.Error("ApiManager creation failed:%v", err)
		return
	}
	glog.Error(apiManager.Start())
}
