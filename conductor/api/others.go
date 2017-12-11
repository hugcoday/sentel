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
	"github.com/golang/glog"
	"github.com/labstack/echo"
)

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
