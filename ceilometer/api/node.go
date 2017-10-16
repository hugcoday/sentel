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
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/cloustone/sentel/ceilometer/collector"
	"github.com/golang/glog"
	"github.com/labstack/echo"
)

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
		glog.Errorf("getNodeInfo:%v", err)
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
		glog.Errorf("getNodeInfo:%v", err)
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

// getNodesClientInfoWithinTimeScope return each node's client info in
// specified time scope
func getNodesClientInfoWithinTimeScope(ctx echo.Context) error {

	// Check parameter's validity
	from, err1 := time.Parse("yyyy-mm-dd hh:mm:ss", ctx.Param("from"))
	to, err2 := time.Parse("yyyy-mm-dd hh:mm:ss", ctx.Param("to"))
	duration, err3 := time.ParseDuration(ctx.Param("unit"))
	glog.Infof("getNodesClientInfoWithinTimeScope(from=%v, to=%v, unit=%v", from, to, duration)

	if err1 != nil || err2 != nil || err3 != nil {
		return ctx.JSON(http.StatusBadRequest,
			&response{Success: false, Message: "time format is wrong"})
	}

	if to.Sub(from) < duration {
		return ctx.JSON(http.StatusBadRequest,
			&response{Success: false, Message: "time format is wrong"})
	}

	config := ctx.(*apiContext).config
	hosts := config.MustString("condutor", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			&response{Success: false, Message: err.Error()})
	}

	c := session.DB("iothub").C("nodes")
	defer session.Close()

	// Get all nodes
	nodes := []collector.Node{}
	if err := c.Find(nil).Limit(100).Iter().All(&nodes); err != nil {
		glog.Errorf("getNodesClientInfo:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&response{
				Success: false,
				Message: err.Error(),
			})
	}
	// For each node, query clients collection to get client's count
	c = session.DB("iothub").C("clients")
	results := map[string][]int{}
	for _, node := range nodes {
		f := from
		result := []int{}
		for {
			t := f.Add(duration)
			if to.Sub(t) <= duration {
				break
			}
			query := bson.M{"nodeId": node.NodeName, "updateTime": bson.M{"$gte": f, "$lt": t}}
			count, err := c.Find(query).Count()
			if err != nil {
				result = append(result, count)
			} else {
				result = append(result, 0)
			}
			f = f.Add(duration)
		}
		results[node.NodeName] = result
	}
	return ctx.JSON(http.StatusOK, &response{
		Success: true,
		Result:  results,
	})

}

//getNodesClientInfo return clients static infor for each node
func getNodesClientInfo(ctx echo.Context) error {
	glog.Infof("calling getNodesClientInfo from %s", ctx.Request().RemoteAddr)

	// Deal specifully if timescope is specified
	from := ctx.Param("from")
	if from != "" {
		return getNodesClientInfoWithinTimeScope(ctx)
	}

	// Retrun last statics for each node
	config := ctx.(*apiContext).config
	hosts := config.MustString("condutor", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		glog.Errorf("getNodeInfo:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&response{
				Success: false,
				Message: err.Error(),
			})
	}

	c := session.DB("iothub").C("nodes")
	defer session.Close()

	// Get all nodes
	nodes := []collector.Node{}
	if err := c.Find(nil).Limit(100).Iter().All(&nodes); err != nil {
		glog.Errorf("getNodesClientInfo:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&response{
				Success: false,
				Message: err.Error(),
			})
	}

	// For each node, query clients collection to get client's count
	result := map[string]int{}
	c = session.DB("iothub").C("clients")
	for _, node := range nodes {
		count, err := c.Find(bson.M{"nodeId": node.NodeName}).Limit(100).Count()
		if err != nil {
			result[node.NodeName] = count
		} else {
			result[node.NodeName] = 0
		}
	}
	return ctx.JSON(http.StatusOK, &response{
		Success: true,
		Result:  result,
	})
}

// getNodeClientsWithinTimeScope return a node's clients statics within
// timescope
func getNodeClientsWithinTimeScope(ctx echo.Context) error {
	// Check parameter's validity
	from, err1 := time.Parse("yyyy-mm-dd hh:mm:ss", ctx.Param("from"))
	to, err2 := time.Parse("yyyy-mm-dd hh:mm:ss", ctx.Param("to"))
	duration, err3 := time.ParseDuration(ctx.Param("unit"))
	nodeName := ctx.Param("nodeName")

	glog.Infof("getNodeClientsWithinTimeScope(node:%s, from=%v, to=%v, unit=%v", nodeName, from, to, duration)

	if err1 != nil || err2 != nil || err3 != nil || nodeName == "" {
		return ctx.JSON(http.StatusBadRequest,
			&response{Success: false, Message: "time format is wrong"})
	}

	if to.Sub(from) < duration {
		return ctx.JSON(http.StatusBadRequest,
			&response{Success: false, Message: "time format is wrong"})
	}

	config := ctx.(*apiContext).config
	hosts := config.MustString("condutor", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			&response{Success: false, Message: err.Error()})
	}

	defer session.Close()

	c := session.DB("iothub").C("clients")
	result := []int{}
	for {
		t := from.Add(duration)
		if to.Sub(t) <= duration {
			break
		}
		query := bson.M{"nodeName": nodeName, "updateTime": bson.M{"$gte": from, "$lt": t}}
		count, err := c.Find(query).Count()
		if err != nil {
			result = append(result, count)
		} else {
			result = append(result, 0)
		}
	}
	return ctx.JSON(http.StatusOK, &response{
		Success: true,
		Result:  result,
	})
}

// getNodeClients return a node's all clients
func getNodeClients(ctx echo.Context) error {
	glog.Infof("calling getNodeClients from %s", ctx.Request().RemoteAddr)

	// Deal specifully if timescope is specified
	from := ctx.Param("from")
	if from != "" {
		return getNodeClientsWithinTimeScope(ctx)
	}

	// Retrun last statics for this node
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
		glog.Errorf("getNodeInfo:%v", err)
		return ctx.JSON(http.StatusInternalServerError,
			&response{
				Success: false,
				Message: err.Error(),
			})
	}
	defer session.Close()
	c := session.DB("iothub").C("nodes")

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
	result := map[string]int{}
	c = session.DB("iothub").C("clients")
	count, err := c.Find(bson.M{"nodeId": node.NodeName}).Limit(100).Count()
	if err != nil {
		result[nodeName] = count
	} else {
		result[node.NodeName] = 0
	}
	return ctx.JSON(http.StatusOK, &response{
		Success: true,
		Result:  result,
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
