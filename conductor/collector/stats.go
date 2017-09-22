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

// Stat
type Stats struct {
	NodeName   string `json:"nodeName"`
	Service    string `json:"service"`
	Action     string `json:"action"`
	UpdateTime time.Time
	Values     map[string]uint64 `json:"values"`
	encoded    []byte
	err        error
}

func (p *Stats) ensureEncoded() {
	if p.encoded == nil && p.err == nil {
		p.encoded, p.err = json.Marshal(p)
	}
}

func (p *Stats) Length() int {
	p.ensureEncoded()
	return len(p.encoded)
}

func (p *Stats) Encode() ([]byte, error) {
	p.ensureEncoded()
	return p.encoded, p.err
}

func (p *Stats) name() string { return TopicNameStats }

func (p *Stats) handleTopic(service *CollectorService, ctx context.Context, value []byte) error {
	var stats []Stats
	if err := json.Unmarshal(value, stats); err != nil {
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
	c := session.DB("iothub.nodes").C("stats")

	for _, topic := range stats {
		switch topic.Action {
		case ObjectActionUpdate:
			c.Insert(&Stats{
				NodeName:   topic.NodeName,
				Service:    topic.Service,
				Values:     topic.Values,
				UpdateTime: time.Now(),
			})
		case ObjectActionDelete:
		case ObjectActionRegister:
		default:
		}
	}
	return nil
}
