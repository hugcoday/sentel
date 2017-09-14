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

package dashboard

import (
	"html/template"
	"io"
	"sync"

	"github.com/cloustone/sentel/iothub/base"
	"github.com/cloustone/sentel/libs"

	"github.com/labstack/echo"
)

type DashboardService struct {
	config libs.Config
	chn    chan base.ServiceCommand
	wg     sync.WaitGroup
	listen string
}

// DashboardServiceFactory
type DashboardServiceFactory struct{}

// New create apiService service factory
func (m *DashboardServiceFactory) New(protocol string, c libs.Config, ch chan base.ServiceCommand) (base.Service, error) {
	service := &DashboardService{
		config: c, wg: sync.WaitGroup{},
		listen: "localhost:8080",
	}

	//	service.e.Renderer = &Template{
	//		templates: template.Must(template.ParseGlob("../src/iothub/dashboard/views/*.html")),
	//	}

	if addr, err := c.String("dashboard", "listen"); err == nil && addr != "" {
		service.listen = addr
	}

	return service, nil
}

// Info
func (m *DashboardService) Info() *base.ServiceInfo {
	return &base.ServiceInfo{
		ServiceName: "dashboard",
	}
}

// Start
func (s *DashboardService) Start() error {
	go func(s *DashboardService) {
		s.wg.Add(1)
	}(s)
	return nil
}

// Stop
func (s *DashboardService) Stop() {
	// How to stop beego
	s.wg.Wait()
}

//
// Wait
func (s *DashboardService) Wait() {
	s.wg.Wait()
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
