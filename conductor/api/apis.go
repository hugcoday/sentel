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
	"fmt"
	"net/http"

	pb "github.com/cloustone/sentel/iothub/api"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/cloustone/sentel/conductor/collector"
	"github.com/golang/glog"
	"github.com/labstack/echo"
)

type response struct {
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
			&response{
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
			&response{
				Success: false,
				Message: err.Error(),
			})
	}

	return ctx.JSON(http.StatusOK, &response{
		Success: true,
		Message: "",
		Result:  nodes,
	})
}

// getNodesUsersInfo return users statics for each node with specificed condition
func getNodesUsersInfo(ctx echo.Context) error {
	glog.Infof("calling getNodesUserInfo from %s", ctx.Request().RemoteAddr)

	config := ctx.(*apiContext).config
	hosts := config.MustString("condutor", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		glog.Errorf("getAllNodes:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&response{
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
			&response{
				Success: false,
				Message: err.Error(),
			})
	}

	return ctx.JSON(http.StatusOK, &response{
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
			&response{
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
			&response{
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
			&response{
				Success: false,
				Message: err.Error(),
			})
	}

	return ctx.JSON(http.StatusOK, &response{
		Success: true,
		Result:  node,
	})
}

// Clients

// getNodeClients return a node's all clients
func getNodeClients(ctx echo.Context) error {
	glog.Infof("calling getNodeClients from %s", ctx.Request().RemoteAddr)

	nodeName := ctx.Param("nodeName")
	if nodeName == "" {
		return ctx.JSON(http.StatusBadRequest,
			&response{
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
			&response{
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
			&response{
				Success: false,
				Message: err.Error(),
			})
	}
	if node.NodeIp == "" {
		glog.Errorf("getNodeClients: cann't resolve node ip for %s", nodeName)
		return ctx.JSON(http.StatusNotFound,
			&response{
				Success: false,
				Message: fmt.Sprintf("cann't resolve node ip for %s", nodeName),
			})
	}

	sentelapi, err := newSentelApi(node.NodeIp)
	if err != nil {
		glog.Errorf("getNodeClients:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&response{
				Success: false,
				Message: err.Error(),
			})
	}
	reply, err := sentelapi.clients(&pb.ClientsRequest{Category: "list"})
	if err != nil {
		glog.Errorf("getNodeClient:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&response{
				Success: false,
				Message: err.Error(),
			})
	}

	return ctx.JSON(http.StatusOK, &response{
		Success: true,
		Result:  reply.Clients,
	})
}

// getNodeClientInfo return spcicified client infor on a node
func getNodeClientInfo(ctx echo.Context) error {
	glog.Infof("calling getNodeClientInfo from %s", ctx.Request().RemoteAddr)

	nodeName := ctx.Param("nodeName")
	clientId := ctx.Param("clientId")
	if nodeName == "" || clientId == "" {
		return ctx.JSON(http.StatusBadRequest,
			&response{
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
			&response{
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
			&response{
				Success: false,
				Message: err.Error(),
			})
	}
	return ctx.JSON(http.StatusOK, &response{
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
			&response{
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
			&response{
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
			&response{
				Success: false,
				Message: err.Error(),
			})
	}
	return ctx.JSON(http.StatusOK, &response{
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
			&response{
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
			&response{
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
			&response{
				Success: false,
				Message: err.Error(),
			})
	}
	return ctx.JSON(http.StatusOK, &response{
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
			&response{
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
			&response{
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
			&response{
				Success: false,
				Message: err.Error(),
			})
	}
	return ctx.JSON(http.StatusOK, &response{
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
			&response{
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
			&response{
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
			&response{
				Success: false,
				Message: err.Error(),
			})
	}
	return ctx.JSON(http.StatusOK, &response{
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
			&response{
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
			&response{
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
			&response{
				Success: false,
				Message: err.Error(),
			})
	}
	return ctx.JSON(http.StatusOK, &response{
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
			&response{
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
			&response{
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
			&response{
				Success: false,
				Message: err.Error(),
			})
	}
	return ctx.JSON(http.StatusOK, &response{
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
			&response{
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
			&response{
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
			&response{
				Success: false,
				Message: err.Error(),
			})
	}
	return ctx.JSON(http.StatusOK, &response{
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

	config := ctx.(*apiContext).config
	hosts := config.MustString("condutor", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		glog.Errorf("getNodeStatsInfo:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&response{Success: false, Message: err.Error()})
	}
	c := session.DB("iothub").C("stats")
	defer session.Close()

	metrics := []collector.Metric{}
	if err := c.Find(nil).Iter().All(&metrics); err != nil {
		glog.Errorf("getClusterStats:%v", err)
		return ctx.JSON(http.StatusNotFound, &response{Success: false, Message: err.Error()})
	}
	services := map[string]map[string]uint64{}
	for _, metric := range metrics {
		if service, ok := services[metric.Service]; !ok { // not found
			services[metric.Service] = metric.Values
		} else {
			for key, val := range metric.Values {
				if _, ok := service[key]; !ok {
					service[key] = val
				} else {
					service[key] += val
				}
			}
		}
	}
	return ctx.JSON(http.StatusOK, &response{Success: true, Result: services})
}

// getNodeMetricsInfo return a node's metrics
func getNodeMetricsInfo(ctx echo.Context) error {
	glog.Infof("calling getNodeMetricsInfo from %s", ctx.Request().RemoteAddr)

	nodeName := ctx.Param("nodeName")
	if nodeName == "" {
		return ctx.JSON(http.StatusBadRequest,
			&response{
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
			&response{
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
			&response{
				Success: false,
				Message: err.Error(),
			})
	}
	if node.NodeIp == "" {
		glog.Errorf("getNodeMetricsInfo: cann't resolve node ip for %s", nodeName)
		return ctx.JSON(http.StatusNotFound,
			&response{
				Success: false,
				Message: fmt.Sprintf("cann't resolve node ip for %s", nodeName),
			})
	}

	sentelapi, err := newSentelApi(node.NodeIp)
	if err != nil {
		glog.Errorf("getNodeMetricsInfo:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&response{
				Success: false,
				Message: err.Error(),
			})
	}
	reply, err := sentelapi.broker(&pb.BrokerRequest{Category: "metrics"})
	if err != nil {
		glog.Errorf("getNodeMetricsInfo:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&response{
				Success: false,
				Message: err.Error(),
			})
	}

	return ctx.JSON(http.StatusOK, &response{
		Success: true,
		Result:  reply.Metrics,
	})
}

// Stats

// getClusterStats return cluster stats
func getClusterStats(ctx echo.Context) error {
	glog.Infof("calling getClusterStats from %s", ctx.Request().RemoteAddr)

	config := ctx.(*apiContext).config
	hosts := config.MustString("condutor", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		glog.Errorf("getNodeStatsInfo:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&response{Success: false, Message: err.Error()})
	}
	c := session.DB("iothub").C("stats")
	defer session.Close()

	stats := []collector.Stats{}
	if err := c.Find(nil).Iter().All(&stats); err != nil {
		glog.Errorf("getClusterStats:%v", err)
		return ctx.JSON(http.StatusNotFound, &response{Success: false, Message: err.Error()})
	}
	services := map[string]map[string]uint64{}
	for _, stat := range stats {
		if service, ok := services[stat.Service]; !ok { // not found
			services[stat.Service] = stat.Values
		} else {
			for key, val := range stat.Values {
				if _, ok := service[key]; !ok {
					service[key] = val
				} else {
					service[key] += val
				}
			}
		}
	}
	return ctx.JSON(http.StatusOK, &response{Success: true, Result: services})
}

//getNodeStatsInfo return a node's stats
func getNodeStatsInfo(ctx echo.Context) error {
	glog.Infof("calling getNodeStats from %s", ctx.Request().RemoteAddr)

	nodeName := ctx.Param("nodeName")
	if nodeName == "" {
		return ctx.JSON(http.StatusBadRequest,
			&response{
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
			&response{
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
			&response{
				Success: false,
				Message: err.Error(),
			})
	}
	if node.NodeIp == "" {
		glog.Errorf("getNodeStatsInfo: cann't resolve node ip for %s", nodeName)
		return ctx.JSON(http.StatusNotFound,
			&response{
				Success: false,
				Message: fmt.Sprintf("cann't resolve node ip for %s", nodeName),
			})
	}

	sentelapi, err := newSentelApi(node.NodeIp)
	if err != nil {
		glog.Errorf("getNodeStatsInfo:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&response{
				Success: false,
				Message: err.Error(),
			})
	}
	reply, err := sentelapi.broker(&pb.BrokerRequest{Category: "stats"})
	if err != nil {
		glog.Errorf("getNodeStatusInfo:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&response{
				Success: false,
				Message: err.Error(),
			})
	}

	return ctx.JSON(http.StatusOK, &response{
		Success: true,
		Result:  reply.Stats,
	})
}
