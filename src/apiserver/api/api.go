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

type apiDescriptor struct {
	action  string
	url     string
	handler echo.HandlerFunc
}

type ApiManager struct {
	Name     string
	Config   *ApiConfig
	handlers []apiDescriptor
	ech      *echo.Echo
}

func NewApiManager(name string, c *ApiConfig) *ApiManager {
	m := &ApiManager{
		Name:     name,
		Config:   c,
		handlers: []apiDescriptor{},
		ech:      echo.New(),
	}
	//	m.ech.Use(middleware.Logger())
	//	m.ech.Use(middleware.Recover())
	return m
}

func (m *ApiManager) RegisterApi(action string, url string, handler echo.HandlerFunc) {
	m.handlers = append(m.handlers,
		apiDescriptor{action: action, url: url, handler: handler})

	// add to echo
	action = strings.ToUpper(action)
	if action == "GET" {
		m.ech.GET(url, handler)
	} else if action == "POST" {
		m.ech.POST(url, handler)
	} else if action == "PUT" {
		m.ech.PUT(url, handler)
	} else if action == "DELETE" {
		m.ech.DELETE(url, handler)
	}
}

func (m *ApiManager) Start(c *ApiConfig) error {
	address := m.Config.Host + ":" + m.Config.Port
	return m.ech.Start(address)
}

func RegisterApiManager(api *ApiManager) {
	apiManagers[api.Name] = api
}

func GetApiManager(name string) *ApiManager {
	return apiManagers[name]
}
