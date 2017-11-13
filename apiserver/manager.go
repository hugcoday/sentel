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

package apiserver

import (
	"fmt"

	"github.com/cloustone/sentel/core"
)

// ApiManager represent api manager interface for each version
type ApiManager interface {
	// Initialize api manager with configuration
	Initialize(c core.Config) error
	// Get apimanager's version
	GetVersion() string
	// Mainloop method
	Run() error
}

var (
	// apiManagers hold all version's api manager
	apiManagers map[string]ApiManager = make(map[string]ApiManager)
)

// RegisterApiManager will allow each version module to register itself
func RegisterApiManager(api ApiManager) {
	version := api.GetVersion()
	if _, ok := apiManagers[version]; ok {
		panic("Same api manager is registered")
	}
	apiManagers[api.GetVersion()] = api
}

// GetApiManger return api manager specified by configuration
func GetApiManager(c core.Config) (ApiManager, error) {
	version := c.MustString("apiserver", "version")
	if _, ok := apiManagers[version]; !ok {
		return nil, fmt.Errorf("Manager %s doesn't exist", version)
	}

	apimgr := apiManagers[version]
	apimgr.Initialize(c)
	return apimgr, nil
}
