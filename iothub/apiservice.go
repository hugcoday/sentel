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

package iothub

import (
	"errors"
	"net/http"
	"sync"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
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

const APIHEAD = "iothub/api/v1/"

// New create apiService service factory
func (this *ApiServiceFactory) New(protocol string, c core.Config, ch chan core.ServiceCommand) (core.Service, error) {
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
	if addr, err := c.String("api", "listen"); err == nil && address != "" {
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

	// Tenant
	e.POST(APIHEAD+"tenant/:id", addTenant)
	e.DELETE(APIHEAD+"tenant/:id", deleteTenant)
	e.GET(APIHEAD+"tenant/:id", getTenantInfo)

	// Broker
	e.GET(APIHEAD+"brokers", getAllBrokerStatus)
	e.GET(APIHEAD+"broker/:id", getBrokerStatus)
	e.PUT(APIHEAD+"broker/:id", updateBrokerStatus)

	return &ApiService{
		config:  c,
		wg:      sync.WaitGroup{},
		chn:     ch,
		address: address,
		echo:    e,
	}, nil

}

// Name
func (this *ApiService) Name() string {
	return "api"
}

// Start
func (this *ApiService) Start() error {
	go func(s *ApiService) {
		this.echo.Start(this.address)
		this.wg.Add(1)
	}(this)
	return nil
}

// Stop
func (this *ApiService) Stop() {
	this.wg.Wait()
}

// Tenant

// addTenant add a tenant to iothub
func addTenant(ctx echo.Context) error {
	glog.Infof("iothub: add tenant(%s)", ctx.Param("id"))

	id := ctx.Param("id")
	hub := getIothub()
	if err := hub.addTenant(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, &response{Success: true})
}

// deleteTenant delete a tenant from iothub
func deleteTenant(ctx echo.Context) error {
	glog.Infof("iothub: delete tenant(%s)", ctx.Param("id"))

	id := ctx.Param("id")
	hub := getIothub()
	if err := hub.deleteTenant(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, &response{Success: true})
}

type tenantResponse struct {
	id           string            `json:"tenantId"`
	createdAt    time.Time         `json:"createdAt"`
	brokersCount int32             `json:"brokersCount"`
	brokers      map[string]string `json:"brokers"`
}

// getTenantInfo retrieve a tenant information from iothub
func getTenantInfo(ctx echo.Context) error {
	glog.Infof("iothub: get tenant(%s) info", ctx.Param("id"))

	tid := ctx.Param("id")
	hub := getIothub()
	tenant := hub.getTenant(tid)
	if tenant == nil {
		return ctx.JSON(http.StatusInternalServerError,
			&response{Success: false, Message: "Invalid tenant id"})
	}

	rsp := &tenantResponse{
		id:           tid,
		createdAt:    tenant.createdAt,
		brokersCount: tenant.brokersCount,
		brokers:      make(map[string]string),
	}
	for bid, broker := range tenant.brokers {
		switch broker.status {
		case BrokerStatusInvalid:
			rsp.brokers[bid] = "invalid"
		case BrokerStatusStarted:
			rsp.brokers[bid] = "started"
		case BrokerStatusStoped:
			rsp.brokers[bid] = "stoped"
		}
	}
	return ctx.JSON(http.StatusOK, &response{Success: true, Result: rsp})
}

// Broker

// getBrokerStatus retrieve a broker status
func getBrokerStatus(ctx echo.Context) error {
	return nil
}

// updateBrokerStatus update a broker status
func updateBrokerStatus(ctx echo.Context) error {
	return nil
}

// getAllBrokerStatus retrieve all broker status
func getAllBrokerStatus(ctx echo.Context) error {
	return nil
}
