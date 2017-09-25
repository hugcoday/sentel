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
	"time"
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
func (p *Metric) clone() topicObject {
	return &Metric{
		topicBase:  p.topicBase,
		NodeName:   p.NodeName,
		Service:    p.Service,
		Values:     p.Values,
		UpdateTime: p.UpdateTime,
	}
}

func (p *Metric) handleTopic(service *CollectorService, ctx context.Context) error {
	db, err := service.getDatabase()
	if err != nil {
		return err
	}
	defer db.Session.Close()
	c := db.C("metrics")

	c.Insert(&Metric{
		NodeName:   p.NodeName,
		Service:    p.Service,
		Values:     p.Values,
		UpdateTime: time.Now(),
	})
	return nil
}
