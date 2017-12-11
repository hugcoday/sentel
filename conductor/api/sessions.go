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
	"net/http"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/cloustone/sentel/conductor/collector"
	"github.com/golang/glog"
	"github.com/labstack/echo"
)

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
