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
	"time"

	"gopkg.in/mgo.v2"
)

// Metric
type Metric struct {
	topicBase
	NodeName   string            `json:"nodeName"`
	Service    string            `json:"service"`
	Values     map[string]uint64 `json:"values"`
	UpdateTime time.Time         `json:"updateTime"`
}

func (p *Metric) name() string { return TopicNameStats }

func (p *Metric) handleTopic(service *CollectorService, ctx context.Context, value []byte) error {
	var metric Metric
	if err := json.Unmarshal(value, &metric); err != nil {
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
	c := session.DB("iothub").C("metrics")

	c.Insert(&Metric{
		NodeName:   metric.NodeName,
		Service:    metric.Service,
		Values:     metric.Values,
		UpdateTime: time.Now(),
	})
	return nil
}
