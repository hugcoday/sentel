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

package collector

import (
	"sync"

	"github.com/cloustone/sentel/conductor/base"
	"github.com/cloustone/sentel/libs"
	"github.com/labstack/echo"
)

type CollectorService struct {
	config  libs.Config
	chn     chan base.ServiceCommand
	wg      sync.WaitGroup
	address string
	echo    *echo.Echo
}

// CollectorServiceFactory
type CollectorServiceFactory struct{}

// New create apiService service factory
func (m *CollectorServiceFactory) New(protocol string, c libs.Config, ch chan base.ServiceCommand) (base.Service, error) {
	address := "localhost:8080"
	if addr, err := c.String("authlet", "address"); err == nil && address != "" {
		address = addr
	}
	// Create echo instance and setup router
	e := echo.New()

	return &CollectorService{
		config:  c,
		wg:      sync.WaitGroup{},
		chn:     ch,
		address: address,
		echo:    e,
	}, nil

}

// Name
func (s *CollectorService) Name() string {
	return "collector"
}

// Start
func (s *CollectorService) Start() error {
	go func(s *CollectorService) {
		s.echo.Start(s.address)
		s.wg.Add(1)
	}(s)
	return nil
}

// Stop
func (s *CollectorService) Stop() {
	s.wg.Wait()
}
