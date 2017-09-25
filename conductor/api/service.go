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
	"fmt"
	"net/http"
	"sync"

	pb "github.com/cloustone/sentel/iothub/api"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/cloustone/sentel/conductor/base"
	"github.com/cloustone/sentel/conductor/collector"
	"github.com/cloustone/sentel/libs"
	"github.com/golang/glog"
	"github.com/labstack/echo"
)

type ApiService struct {
	config  libs.Config
	chn     chan base.ServiceCommand
	wg      sync.WaitGroup
	address string
	echo    *echo.Echo
}

type apiContext struct {
	echo.Context
	config libs.Config
}

// ApiServiceFactory
type ApiServiceFactory struct{}

const APIHEAD = "api/v1/"

// New create apiService service factory
func (m *ApiServiceFactory) New(protocol string, c libs.Config, ch chan base.ServiceCommand) (base.Service, error) {
	// check mongo db configuration
	hosts, err := c.String("conductor", "mongo")
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
	if addr, err := c.String("conductor", "listen"); err == nil && address != "" {
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
	return "restapi-service"
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

type responseHeader struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

// getAllNodes return all nodes in clusters
func getAllNodes(ctx echo.Context) error {
	glog.Infof("calling getAllNodes from %s", ctx.Request().RemoteAddr)

	config := ctx.(*apiContext).config
	hosts := config.MustString("condutor", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		glog.Errorf("getAllNodes:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	c := session.DB("iothub").C("nodes")
	defer session.Close()

	nodes := []collector.Node{}
	iter := c.Find(nil).Limit(100).Iter()
	err = iter.All(nodes)
	if err != nil {
		glog.Errorf("getAllNodes:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}

	return ctx.JSON(http.StatusOK, &responseHeader{
		Success: true,
		Message: "",
		Result:  nodes,
	})
}

// getNodeInfo return a node's detail info
func getNodeInfo(ctx echo.Context) error {
	glog.Infof("calling getNodeInfo from %s", ctx.Request().RemoteAddr)

	nodeName := ctx.Param("nodeName")
	if nodeName == "" {
		return ctx.JSON(http.StatusBadRequest,
			&responseHeader{
				Success: false,
				Message: "Invalid parameter",
			})
	}

	config := ctx.(*apiContext).config
	hosts := config.MustString("condutor", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		glog.Errorf("getAllNodeInfo:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	c := session.DB("iothub").C("nodes")
	defer session.Close()

	node := collector.Node{}
	if err := c.Find(bson.M{"NodeName": nodeName}).One(&node); err != nil {
		glog.Errorf("getAllNodes:%v", err)
		return ctx.JSON(http.StatusNotFound,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}

	return ctx.JSON(http.StatusOK, &responseHeader{
		Success: true,
		Result:  node,
	})

	return nil
}

// Clients

// getNodeClients return a node's all clients
func getNodeClients(ctx echo.Context) error {
	glog.Infof("calling getNodeClients from %s", ctx.Request().RemoteAddr)

	nodeName := ctx.Param("nodeName")
	if nodeName == "" {
		return ctx.JSON(http.StatusBadRequest,
			&responseHeader{
				Success: false,
				Message: "Invalid parameter",
			})
	}

	config := ctx.(*apiContext).config
	hosts := config.MustString("condutor", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		glog.Errorf("getAllNodeClients:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	c := session.DB("iothub").C("nodes")
	defer session.Close()

	node := collector.Node{}
	if err := c.Find(bson.M{"NodeName": nodeName}).One(&node); err != nil {
		glog.Errorf("getNodeClients:%v", err)
		return ctx.JSON(http.StatusNotFound,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	if node.NodeIp == "" {
		glog.Errorf("getNodeClients: cann't resolve node ip for %s", nodeName)
		return ctx.JSON(http.StatusNotFound,
			&responseHeader{
				Success: false,
				Message: fmt.Sprintf("cann't resolve node ip for %s", nodeName),
			})
	}

	sentelapi, err := newSentelApi(node.NodeIp)
	if err != nil {
		glog.Errorf("getNodeClients:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	reply, err := sentelapi.clients(&pb.ClientsRequest{Category: "list"})
	if err != nil {
		glog.Errorf("getNodeClient:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}

	return ctx.JSON(http.StatusOK, &responseHeader{
		Success: true,
		Result:  reply.Clients,
	})

	return nil
}

// getNodeClientInfo return spcicified client infor on a node
func getNodeClientInfo(ctx echo.Context) error {
	glog.Infof("calling getNodeClientInfo from %s", ctx.Request().RemoteAddr)

	nodeName := ctx.Param("nodeName")
	clientId := ctx.Param("clientId")
	if nodeName == "" || clientId == "" {
		return ctx.JSON(http.StatusBadRequest,
			&responseHeader{
				Success: false,
				Message: "Invalid parameter",
			})
	}

	config := ctx.(*apiContext).config
	hosts := config.MustString("condutor", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		glog.Errorf("getAllNodeClientInfo:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	c := session.DB("iothub").C("clients")
	defer session.Close()

	client := collector.Client{}
	if err := c.Find(bson.M{"NodeName": nodeName, "ClientId": clientId}).One(&client); err != nil {
		glog.Errorf("getNodeClientInfo:%v", err)
		return ctx.JSON(http.StatusNotFound,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	return ctx.JSON(http.StatusOK, &responseHeader{
		Success: true,
		Result:  client,
	})

}

// getClusterClientInfo return clients info in cluster
func getClusterClientInfo(ctx echo.Context) error {
	glog.Infof("calling getClusterClientInfo from %s", ctx.Request().RemoteAddr)

	clientId := ctx.Param("clientId")
	if clientId == "" {
		return ctx.JSON(http.StatusBadRequest,
			&responseHeader{
				Success: false,
				Message: "Invalid parameter",
			})
	}

	config := ctx.(*apiContext).config
	hosts := config.MustString("condutor", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		glog.Errorf("getClusterClientInfo:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	c := session.DB("iothub").C("clients")
	defer session.Close()

	clients := []collector.Client{}
	if err := c.Find(bson.M{"ClientId": clientId}).Limit(100).Iter().All(&clients); err != nil {
		glog.Errorf("getClusterClientInfo:%v", err)
		return ctx.JSON(http.StatusNotFound,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	return ctx.JSON(http.StatusOK, &responseHeader{
		Success: true,
		Result:  clients,
	})

}

// Session

// getNodeSessions return a node's session
func getNodeSessions(ctx echo.Context) error {
	glog.Infof("calling getNodeSessions from %s", ctx.Request().RemoteAddr)

	nodeName := ctx.Param("nodeName")
	if nodeName == "" {
		return ctx.JSON(http.StatusBadRequest,
			&responseHeader{
				Success: false,
				Message: "Invalid parameter",
			})
	}

	config := ctx.(*apiContext).config
	hosts := config.MustString("condutor", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		glog.Errorf("getNodeSessions:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	c := session.DB("iothub").C("sessions")
	defer session.Close()

	sessions := collector.Session{}
	if err := c.Find(bson.M{"NodeName": nodeName}).Limit(100).Iter().All(&sessions); err != nil {
		glog.Errorf("getNodeSessions:%v", err)
		return ctx.JSON(http.StatusNotFound,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	return ctx.JSON(http.StatusOK, &responseHeader{
		Success: true,
		Result:  sessions,
	})
}

// getNodeSessionsClient return client infor in a node's sessions
func getNodeSessionsClientInfo(ctx echo.Context) error {
	glog.Infof("calling getNodeSessionsClientInfo from %s", ctx.Request().RemoteAddr)

	nodeName := ctx.Param("nodeName")
	clientId := ctx.Param("clientId")
	if nodeName == "" || clientId == "" {
		return ctx.JSON(http.StatusBadRequest,
			&responseHeader{
				Success: false,
				Message: "Invalid parameter",
			})
	}

	config := ctx.(*apiContext).config
	hosts := config.MustString("condutor", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		glog.Errorf("getNodeSessionsClientInfo:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	c := session.DB("iothub").C("sessions")
	defer session.Close()

	sessions := collector.Session{}
	if err := c.Find(bson.M{"NodeName": nodeName, "ClientId": clientId}).Limit(100).Iter().All(&sessions); err != nil {
		glog.Errorf("getNodeSessionsClientInfo:%v", err)
		return ctx.JSON(http.StatusNotFound,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	return ctx.JSON(http.StatusOK, &responseHeader{
		Success: true,
		Result:  sessions,
	})
}

// getClusterSessionInfor return client info in cluster session
func getClusterSessionClientInfo(ctx echo.Context) error {
	glog.Infof("calling getClusterSessionsClientInfo from %s", ctx.Request().RemoteAddr)

	clientId := ctx.Param("clientId")
	if clientId == "" {
		return ctx.JSON(http.StatusBadRequest,
			&responseHeader{
				Success: false,
				Message: "Invalid parameter",
			})
	}

	config := ctx.(*apiContext).config
	hosts := config.MustString("condutor", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		glog.Errorf("getClusterSessionsClientInfo:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	c := session.DB("iothub").C("sessions")
	defer session.Close()

	sessions := collector.Session{}
	if err := c.Find(bson.M{"ClientId": clientId}).Limit(100).Iter().All(&sessions); err != nil {
		glog.Errorf("getClusterSessionsClientInfo:%v", err)
		return ctx.JSON(http.StatusNotFound,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	return ctx.JSON(http.StatusOK, &responseHeader{
		Success: true,
		Result:  sessions,
	})
}

// Subscription

// getNodeSubscriptions return a node's subscriptions
func getNodeSubscriptions(ctx echo.Context) error {
	glog.Infof("calling getNodeSubscriptions from %s", ctx.Request().RemoteAddr)

	nodeName := ctx.Param("nodeName")
	if nodeName == "" {
		return ctx.JSON(http.StatusBadRequest,
			&responseHeader{
				Success: false,
				Message: "Invalid parameter",
			})
	}

	config := ctx.(*apiContext).config
	hosts := config.MustString("condutor", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		glog.Errorf("getNodeSubscriptions:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	c := session.DB("iothub").C("subscriptions")
	defer session.Close()

	subs := []collector.Subscription{}
	if err := c.Find(bson.M{"NodeName": nodeName}).Limit(100).Iter().All(&subs); err != nil {
		glog.Errorf("getNodeSubscriptions:%v", err)
		return ctx.JSON(http.StatusNotFound,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	return ctx.JSON(http.StatusOK, &responseHeader{
		Success: true,
		Result:  subs,
	})
}

// getNodeSubscriptionsClientInfo return client info in node's subscriptions
func getNodeSubscriptionsClientInfo(ctx echo.Context) error {
	glog.Infof("calling getNodeSubscriptionsClientInfo from %s", ctx.Request().RemoteAddr)

	nodeName := ctx.Param("nodeName")
	clientId := ctx.Param("clientId")
	if nodeName == "" || clientId == "" {
		return ctx.JSON(http.StatusBadRequest,
			&responseHeader{
				Success: false,
				Message: "Invalid parameter",
			})
	}

	config := ctx.(*apiContext).config
	hosts := config.MustString("condutor", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		glog.Errorf("getNodeSubscriptionsClientInfo:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	c := session.DB("iothub").C("subscriptions")
	defer session.Close()

	subs := []collector.Subscription{}
	if err := c.Find(bson.M{"NodeName": nodeName, "ClientId": clientId}).Limit(100).Iter().All(&subs); err != nil {
		glog.Errorf("getNodeSubscriptionsclientInfo:%v", err)
		return ctx.JSON(http.StatusNotFound,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	return ctx.JSON(http.StatusOK, &responseHeader{
		Success: true,
		Result:  subs,
	})
}

// getClusterSubscriptionsInfo return client info in cluster subscriptions
func getClusterSubscriptionsInfo(ctx echo.Context) error {
	glog.Infof("calling getClusterSubscriptionsInfo from %s", ctx.Request().RemoteAddr)

	clientId := ctx.Param("clientId")
	if clientId == "" {
		return ctx.JSON(http.StatusBadRequest,
			&responseHeader{
				Success: false,
				Message: "Invalid parameter",
			})
	}

	config := ctx.(*apiContext).config
	hosts := config.MustString("condutor", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		glog.Errorf("getClusterSubscriptionsInfo:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	c := session.DB("iothub").C("subscriptions")
	defer session.Close()

	subs := []collector.Subscription{}
	if err := c.Find(bson.M{"ClientId": clientId}).Limit(100).Iter().All(&subs); err != nil {
		glog.Errorf("getClusterSubscriptionsInfo:%v", err)
		return ctx.JSON(http.StatusNotFound,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	return ctx.JSON(http.StatusOK, &responseHeader{
		Success: true,
		Result:  subs,
	})
}

// Routes

// getClusterRoutes return cluster's routes table
func getClusterRoutes(ctx echo.Context) error {
	glog.Infof("calling getClusterRoutes from %s", ctx.Request().RemoteAddr)
	return nil
}

// getTopicRoutes return a topic's route
func getTopicRoutes(ctx echo.Context) error {
	glog.Infof("calling getTopicRoutes from %s", ctx.Request().RemoteAddr)
	return nil
}

// Publish & Subscribe

// publishMqttMessage will publish a mqtt message
func publishMqttMessage(ctx echo.Context) error {
	glog.Infof("calling publishMqttMessage from %s", ctx.Request().RemoteAddr)
	return nil
}

// subscribeMqttMessage subscribe a mqtt topic
func subscribeMqttMessage(ctx echo.Context) error {
	glog.Infof("calling subscribeMqttMessage from %s", ctx.Request().RemoteAddr)
	return nil
}

// unsubscribeMqttMessage unsubsribe mqtt topic
func unsubscribeMqttMessage(ctx echo.Context) error {
	glog.Infof("calling unsubscribeMqttMessage from %s", ctx.Request().RemoteAddr)
	return nil
}

// Plugins

// getNodePluginsInfo return plugins info for a node
func getNodePluginsInfo(ctx echo.Context) error {
	glog.Infof("calling getNodePluginsInfo from %s", ctx.Request().RemoteAddr)
	return nil
}

// Services

// getClusterServicesInfo return all services infor in cluster
func getClusterServicesInfo(ctx echo.Context) error {
	glog.Infof("calling getClusterServicesInfo from %s", ctx.Request().RemoteAddr)
	return nil
}

// getNodeServicesInfo return a node's service info
func getNodeServicesInfo(ctx echo.Context) error {
	glog.Infof("calling getNodeServicesInfo from %s", ctx.Request().RemoteAddr)
	return nil
}

// Metrics

// getClusterMetricsInfo return cluster metrics
func getClusterMetricsInfo(ctx echo.Context) error {
	glog.Infof("calling getClusterMetricsInfo from %s", ctx.Request().RemoteAddr)
	return nil
}

// getNodeMetricsInfo return a node's metrics
func getNodeMetricsInfo(ctx echo.Context) error {
	glog.Infof("calling getNodeMetricsInfo from %s", ctx.Request().RemoteAddr)

	nodeName := ctx.Param("nodeName")
	if nodeName == "" {
		return ctx.JSON(http.StatusBadRequest,
			&responseHeader{
				Success: false,
				Message: "Invalid parameter",
			})
	}

	config := ctx.(*apiContext).config
	hosts := config.MustString("condutor", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		glog.Errorf("getNodeMetricsInfo:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	c := session.DB("iothub").C("nodes")
	defer session.Close()

	node := collector.Node{}
	if err := c.Find(bson.M{"NodeName": nodeName}).One(&node); err != nil {
		glog.Errorf("getNodeMetricsInfo:%v", err)
		return ctx.JSON(http.StatusNotFound,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	if node.NodeIp == "" {
		glog.Errorf("getNodeMetricsInfo: cann't resolve node ip for %s", nodeName)
		return ctx.JSON(http.StatusNotFound,
			&responseHeader{
				Success: false,
				Message: fmt.Sprintf("cann't resolve node ip for %s", nodeName),
			})
	}

	sentelapi, err := newSentelApi(node.NodeIp)
	if err != nil {
		glog.Errorf("getNodeMetricsInfo:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	reply, err := sentelapi.broker(&pb.BrokerRequest{Category: "metrics"})
	if err != nil {
		glog.Errorf("getNodeMetricsInfo:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}

	return ctx.JSON(http.StatusOK, &responseHeader{
		Success: true,
		Result:  reply.Metrics,
	})
}

// Stats

// getClusterStats return cluster stats
func getClusterStats(ctx echo.Context) error {
	glog.Infof("calling getClusterStats from %s", ctx.Request().RemoteAddr)
	return nil
}

//getNodeStatsInfo return a node's stats
func getNodeStatsInfo(ctx echo.Context) error {
	glog.Infof("calling getNodeStats from %s", ctx.Request().RemoteAddr)

	nodeName := ctx.Param("nodeName")
	if nodeName == "" {
		return ctx.JSON(http.StatusBadRequest,
			&responseHeader{
				Success: false,
				Message: "Invalid parameter",
			})
	}

	config := ctx.(*apiContext).config
	hosts := config.MustString("condutor", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		glog.Errorf("getNodeStatsInfo:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	c := session.DB("iothub").C("nodes")
	defer session.Close()

	node := collector.Node{}
	if err := c.Find(bson.M{"NodeName": nodeName}).One(&node); err != nil {
		glog.Errorf("getNodeStatsInfo:%v", err)
		return ctx.JSON(http.StatusNotFound,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	if node.NodeIp == "" {
		glog.Errorf("getNodeStatsInfo: cann't resolve node ip for %s", nodeName)
		return ctx.JSON(http.StatusNotFound,
			&responseHeader{
				Success: false,
				Message: fmt.Sprintf("cann't resolve node ip for %s", nodeName),
			})
	}

	sentelapi, err := newSentelApi(node.NodeIp)
	if err != nil {
		glog.Errorf("getNodeStatsInfo:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}
	reply, err := sentelapi.broker(&pb.BrokerRequest{Category: "stats"})
	if err != nil {
		glog.Errorf("getNodeStatusInfo:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&responseHeader{
				Success: false,
				Message: err.Error(),
			})
	}

	return ctx.JSON(http.StatusOK, &responseHeader{
		Success: true,
		Result:  reply.Stats,
	})
}
