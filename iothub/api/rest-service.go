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
	"sync"

	"github.com/cloustone/sentel/iothub/base"
	"github.com/cloustone/sentel/libs"
	"github.com/labstack/echo"
)

type RestApiService struct {
	config  libs.Config
	chn     chan base.ServiceCommand
	wg      sync.WaitGroup
	address string
	echo    *echo.Echo
}

// RestApiServiceFactory
type RestApiServiceFactory struct{}

const APIHEAD = "api/v1/"

// New create apiService service factory
func (m *RestApiServiceFactory) New(protocol string, c libs.Config, ch chan base.ServiceCommand) (base.Service, error) {
	address := "localhost:8080"
	if addr, err := c.String("authlet", "address"); err == nil && address != "" {
		address = addr
	}
	// Create echo instance and setup router
	e := echo.New()

	// Clusters
	e.GET(APIHEAD+"nodes", getAllNodes)
	e.GET(APIHEAD+"nodes/:nodeName", getNodeInfo)

	// Clients
	e.GET(APIHEAD+"nodes/:nodeName/clients", getNodeClients)
	e.GET(APIHEAD+"nodes/:nodeName/clients/:clientId", getNodeClientInfo)
	e.GET(APIHEAD+"clients/:clientId", getClusterClientInfo)

	// Session
	e.GET(APIHEAD+"nodes/:nodeName/sessions", getNodeSessions)
	e.GET(APIHEAD+"nodes/:nodeName/sessions/:clientId", getNodeSessionsClientInfo)
	e.GET(APIHEAD+"sessions/:clientId", getClusterSessionClientInfo)

	// Subscription
	e.GET(APIHEAD+"nodes/:nodeName/subscriptions", getNodeSubscriptions)
	e.GET(APIHEAD+"nodes/:nodeName/subscriptions/:clientId", getNodeSubscriptionsClientInfo)
	e.GET(APIHEAD+"subscriptions/:clientId", getClusterSubscriptionsInfo)

	// Routes
	e.GET(APIHEAD+"routes", getClusterRoutes)
	e.GET(APIHEAD+"routes/:topic", getTopicRoutes)

	// Publish & Subscribe
	e.POST(APIHEAD+"mqtt/publish", publishMqttMessage)
	e.POST(APIHEAD+"mqtt/subscribe", subscribeMqttMessage)
	e.POST(APIHEAD+"mqtt/unsubscribe", unsubscribeMqttMessage)

	// Plugins
	e.GET(APIHEAD+"nodes/:nodeName/plugins", getNodePluginsInfo)

	// Services
	e.GET(APIHEAD+"services", getClusterServicesInfo)
	e.GET(APIHEAD+"nodes/:nodeName/services", getNodeServicesInfo)

	// Metrics
	e.GET(APIHEAD+"metrics", getClusterMetricsInfo)
	e.GET(APIHEAD+"nodes/:nodeName/metrics", getNodeMetricsInfo)

	// Stats
	e.GET(APIHEAD+"stats", getClusterStats)
	e.GET(APIHEAD+"nodes/:nodeName/stats", getNodeStatsInfo)

	return &RestApiService{
		config:  c,
		wg:      sync.WaitGroup{},
		chn:     ch,
		address: address,
		echo:    e,
	}, nil

}

// Name
func (s *RestApiService) Info() *base.ServiceInfo {
	return &base.ServiceInfo{
		ServiceName: "restapi-service",
	}
}

// Start
func (s *RestApiService) Start() error {
	go func(s *RestApiService) {
		s.echo.Start(s.address)
		s.wg.Add(1)
	}(s)
	return nil
}

// Stop
func (s *RestApiService) Stop() {
	s.wg.Wait()
}

// Node

// getAllNodes return all nodes in clusters
func getAllNodes(ctx echo.Context) error {
	return nil
}

// getNodeInfo return a node's detail info
func getNodeInfo(ctx echo.Context) error {
	return nil
}

// Clients

// getNodeClients return a node's clients
func getNodeClients(ctx echo.Context) error {
	return nil
}

// getNodeClientInfo return spcicified client infor on a node
func getNodeClientInfo(ctx echo.Context) error {
	return nil
}

// getClusterClientInfo return clients info in cluster
func getClusterClientInfo(ctx echo.Context) error {
	return nil
}

// Session

// getNodeSessions return a node's session
func getNodeSessions(ctx echo.Context) error {
	return nil
}

// getNodeSessionsClient return client infor in a node's sessions
func getNodeSessionsClientInfo(ctx echo.Context) error {
	return nil
}

// getClusterSessionInfor return client info in cluster session
func getClusterSessionClientInfo(ctx echo.Context) error {
	return nil
}

// Subscription

// getNodeSubscriptions return a node's subscriptions
func getNodeSubscriptions(ctx echo.Context) error {
	return nil
}

// getNodeSubscriptionsClientInfo return client info in node's subscriptions
func getNodeSubscriptionsClientInfo(ctx echo.Context) error {
	return nil
}

// getClusterSubscriptionsInfo return client info in cluster subscriptions
func getClusterSubscriptionsInfo(ctx echo.Context) error {
	return nil
}

// Routes

// getClusterRoutes return cluster's routes table
func getClusterRoutes(ctx echo.Context) error {
	return nil
}

// getTopicRoutes return a topic's route
func getTopicRoutes(ctx echo.Context) error {
	return nil
}

// Publish & Subscribe

// publishMqttMessage will publish a mqtt message
func publishMqttMessage(ctx echo.Context) error {
	return nil
}

// subscribeMqttMessage subscribe a mqtt topic
func subscribeMqttMessage(ctx echo.Context) error {
	return nil
}

// unsubscribeMqttMessage unsubsribe mqtt topic
func unsubscribeMqttMessage(ctx echo.Context) error {
	return nil
}

// Plugins

// getNodePluginsInfo return plugins info for a node
func getNodePluginsInfo(ctx echo.Context) error {
	return nil
}

// Services

// getClusterServicesInfo return all services infor in cluster
func getClusterServicesInfo(ctx echo.Context) error {
	return nil
}

// getNodeServicesInfo return a node's service info
func getNodeServicesInfo(ctx echo.Context) error {
	return nil
}

// Metrics

// getClusterMetricsInfo return cluster metrics
func getClusterMetricsInfo(ctx echo.Context) error {
	return nil
}

// getNodeMetricsInfo return a node's metrics
func getNodeMetricsInfo(ctx echo.Context) error {
	return nil
}

// Stats

// getClusterStats return cluster stats
func getClusterStats(ctx echo.Context) error {
	return nil
}

//getNodeStatsInfo return a node's stats
func getNodeStatsInfo(ctx echo.Context) error {
	return nil
}
