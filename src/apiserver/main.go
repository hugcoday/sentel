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
	"log"
)

const (
	defaultConfigFilePath = "../etc/sentel/apiserver.conf"
)

var (
	apiconfig base.ApiConfig
)

func main() {
	// Get configuration
	c := base.NewLoaderWithPath(defaultConfigFilePath)
	c.MustLoad(apiconfig)

	// Initialize registry store
	if err := db.InitializeRegistryStore(apiconfig); err != nil {
		log.Fatal("Registry initialization failed:%v", err)
		return
	}

	// Create api manager using configuration
	apiManager, err := base.CreateApiManager(&apiconfig)
	if err != nil {
		log.Fatal("ApiManager creation failed:%v", err)
		return
	}
	log.Fatal(apiManager.Start())
}
