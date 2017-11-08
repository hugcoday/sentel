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
	"errors"
	"sync"

	mgo "gopkg.in/mgo.v2"

	"github.com/cloustone/sentel/core"
	"github.com/labstack/echo"
)

type ApiService struct {
	config  core.Config
	chn     chan core.ServiceCommand
	wg      sync.WaitGroup
	address string
	echo    *echo.Echo
}

type apiContext struct {
	echo.Context
	config core.Config
}

type response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

// ApiServiceFactory
type ApiServiceFactory struct{}

const APIHEAD = "api/v1/"

// New create apiService service factory
func (m *ApiServiceFactory) New(protocol string, c core.Config, ch chan core.ServiceCommand) (core.Service, error) {
	// check mongo db configuration
	hosts, err := c.String("iothub", "mongo")
	if err != nil || hosts == "" {
		return nil, errors.New("Invalid mongo configuration")
	}

	// try connect with mongo db
	session, err := mgo.Dial(hosts)
	if err != nil {
		return nil, err
	}
	session.Close()

	address := "localhost:8080"
	if addr, err := c.String("iothub", "listen"); err == nil && address != "" {
		address = addr
	}
	// Create echo instance and setup router
	e := echo.New()
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			cc := &apiContext{Context: e, config: c}
			return h(cc)
		}
	})

	// e.GET(APIHEAD+"nodes", getAllNodes)
	return &ApiService{
		config:  c,
		wg:      sync.WaitGroup{},
		chn:     ch,
		address: address,
		echo:    e,
	}, nil

}

// Name
func (s *ApiService) Name() string {
	return "iothub-service"
}

// Start
func (s *ApiService) Start() error {
	go func(s *ApiService) {
		s.echo.Start(s.address)
		s.wg.Add(1)
	}(s)
	return nil
}

// Stop
func (s *ApiService) Stop() {
	s.wg.Wait()
}
