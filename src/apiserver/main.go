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
	"apiserver/api"
	"apiserver/aws"
	"apiserver/azure"
	config "utility/config"
)

const (
	defaultConfigFilePath = "/etc/sentel-apiserver.toml"
)

func main() {
	// Get configuration
	c := config.NewWithPath(defaultConfigFilePath)
	var apiConfig api.ApiConfig
	c.MustLoad(apiConfig)

	// Create api manager using configuration
	apiManager := api.GetApiManager(apiConfig.ApiCategory)
	apiManager.Start(&apiConfig)
}

func init() {
	api.RegisterApiManager(azure.NewApi())
	api.RegisterApiManager(aws.NewApi())
}
