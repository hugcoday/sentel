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

package base

import (
	mw "apiserver/middleware"
	"strings"

	echo "github.com/labstack/echo"
)

type apiDescriptor struct {
	action  string
	url     string
	handler echo.HandlerFunc
}

type ApiManager struct {
	Version  string
	Config   *ApiConfig
	handlers []apiDescriptor
	ech      *echo.Echo
}

var (
	apiManagers map[string]*ApiManager = make(map[string]*ApiManager)
)

func NewApiManager(version string) *ApiManager {
	m := &ApiManager{
		Version:  version,
		Config:   nil,
		handlers: []apiDescriptor{},
		ech:      echo.New(),
	}

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

func (m *ApiManager) Start() error {
	address := m.Config.Host + ":" + m.Config.Port
	return m.ech.Start(address)
}

func RegisterApiManager(api *ApiManager) {
	apiManagers[api.Version] = api
}

func CreateApiManager(c *ApiConfig) (*ApiManager, error) {
	m := apiManagers[c.Version]
	m.Config = c

	m.ech.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			cc := &ApiContext{Context: e, Config: m.Config}
			return h(cc)
		}
	})
	m.ech.Use(mw.ApiVersion(c.Version))
	//	m.ech.Use(middleware.Logger())
	//	m.ech.Use(middleware.Recover())

	return m, nil
}
