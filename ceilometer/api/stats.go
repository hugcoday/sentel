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

	pb "github.com/cloustone/sentel/broker/api"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/cloustone/sentel/ceilometer/collector"
	"github.com/golang/glog"
	"github.com/labstack/echo"
)

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
