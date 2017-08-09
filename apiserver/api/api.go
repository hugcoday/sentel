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

package api

import (
	"strings"

	echo "github.com/labstack/echo"
)

type ApiConfig struct {
	Host        string // server host
	Port        string // server port
	LogLevel    string // Log level
	Registry    string // Registry RPC server
	ApiCategory string // Api category, aws or azure
}

var (
	apiManagers map[string]*ApiManager = make(map[string]*ApiManager)
)

type ApiDescriptor struct {
	Action  string
	Url     string
	Handler echo.HandlerFunc
}

type ApiManager struct {
	Name         string
	Config       *ApiConfig
	Handlers     []ApiDescriptor
	EchoInstance *echo.Echo
}

func (m *ApiManager) RegisterApi(action string, url string, handler echo.HandlerFunc) {
	m.Handlers = append(m.Handlers,
		ApiDescriptor{Action: action, Url: url, Handler: handler})

	// add to echo
	action = strings.ToUpper(action)
	if action == "GET" {
		m.EchoInstance.GET(url, handler)
	} else if action == "POST" {
		m.EchoInstance.POST(url, handler)
	} else if action == "PUT" {
		m.EchoInstance.PUT(url, handler)
	} else if action == "DELETE" {
		m.EchoInstance.DELETE(url, handler)
	}
}

func (m *ApiManager) Start(c *ApiConfig) error {
	address := m.Config.Host + ":" + m.Config.Port
	return m.EchoInstance.Start(address)
}

func RegisterApiManager(api *ApiManager) {
	apiManagers[api.Name] = api
}

func GetApiManager(name string) *ApiManager {
	return apiManagers[name]
}
