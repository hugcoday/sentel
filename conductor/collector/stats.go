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

// Stat
type Stats struct {
	topicBase
	NodeName   string `json:"nodeName"`
	Service    string `json:"service"`
	Action     string `json:"action"`
	UpdateTime time.Time
	Values     map[string]uint64 `json:"values"`
}

func (p *Stats) name() string { return TopicNameStats }

func (p *Stats) handleTopic(service *CollectorService, ctx context.Context) error {
	db, err := service.getDatabase()
	if err != nil {
		return err
	}
	defer db.Session.Close()
	c := db.C("stats")

	switch p.Action {
	case ObjectActionUpdate:
		c.Insert(&Stats{
			NodeName:   p.NodeName,
			Service:    p.Service,
			Values:     p.Values,
			UpdateTime: time.Now(),
		})
	case ObjectActionDelete:
	case ObjectActionRegister:
	default:
	}
	return nil
}
