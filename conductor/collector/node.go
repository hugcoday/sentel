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
	"context"
	"encoding/json"
	"errors"

	"gopkg.in/mgo.v2"
)

// Node
type Node struct {
	topicBase
	NodeName  string `json:"nodeName"`
	NodeIp    string `json:"nodeIp"`
	CreatedAt string `json:"createdAt"`
}

func (p *Node) name() string { return TopicNameNode }

func (p *Node) handleTopic(service *CollectorService, ctx context.Context, value []byte) error {
	var nodes []Node
	if err := json.Unmarshal(value, nodes); err != nil {
		return err
	}

	// mongo config
	hosts, err := service.config.String("mongo", "hosts")
	if err != nil || hosts == "" {
		return errors.New("Invalid mongo configuration")
	}

	session, err := mgo.Dial(hosts)
	if err != nil {
		return err
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("iothub").C("nodes")

	for _, topic := range nodes {
		c.Insert(&Node{
			NodeName:  topic.NodeName,
			NodeIp:    topic.NodeIp,
			CreatedAt: topic.CreatedAt,
		})
	}
	return nil
}
