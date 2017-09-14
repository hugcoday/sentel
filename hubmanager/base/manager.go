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
	"fmt"

	"github.com/cloustone/sentel/libs"

	"github.com/cloustone/sentel/hubmanager/middleware"

	echo "github.com/labstack/echo"
)

type apiDescriptor struct {
	action  int
	url     string
	handler echo.HandlerFunc
}

type ApiManager struct {
	Version  string
	config   libs.Config
	handlers []apiDescriptor
	ech      *echo.Echo
}

const (
	GET    = 0
	POST   = 1
	DELETE = 2
	PATCH  = 3
	PUT    = 4
)

var (
	apiManagers map[string]*ApiManager = make(map[string]*ApiManager)
)

func NewApiManager(version string, c libs.Config) *ApiManager {
	m := &ApiManager{
		Version:  version,
		config:   c,
		handlers: []apiDescriptor{},
		ech:      echo.New(),
	}

	return m
}

func (m *ApiManager) RegisterApi(action int, url string, handler echo.HandlerFunc, h ...echo.MiddlewareFunc) {
	m.handlers = append(m.handlers,
		apiDescriptor{action: action, url: url, handler: handler})
	switch action {
	case GET:
		m.ech.GET(url, handler, h...)
	case POST:
		m.ech.POST(url, handler, h...)
	case PUT:
		m.ech.PUT(url, handler, h...)
	case DELETE:
		m.ech.DELETE(url, handler, h...)
	case PATCH:
		m.ech.PATCH(url, handler, h...)
	}
}

func (m *ApiManager) Start() error {
	address := ":1385"
	if addr, err := m.config.String("apiserver", "address"); err == nil && addr != "" {
		address = addr
	}
	return m.ech.Start(address)
}

func RegisterApiManager(api *ApiManager) {
	apiManagers[api.Version] = api
}

func CreateApiManager(c libs.Config) (*ApiManager, error) {
	version, err := c.String("apiserver", "version")
	if err != nil {
		return nil, err
	}
	if _, ok := apiManagers[version]; !ok {
		return nil, fmt.Errorf("There isn't manager for '%s'", version)
	}

	m := apiManagers[version]
	m.config = c

	m.ech.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			cc := &ApiContext{Context: e, Config: c}
			return h(cc)
		}
	})
	m.ech.Use(middleware.ApiVersion(version))
	//m.ech.Use(mw.KeyAuthWithConfig(middleware.DefaultKeyAuthConfig))
	//	m.ech.Use(middleware.Logger())
	//	m.ech.Use(middleware.Recover())

	return m, nil
}
