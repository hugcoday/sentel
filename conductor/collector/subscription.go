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

	mgo "gopkg.in/mgo.v2"
)

// Subscription
type Subscription struct {
	topicBase
	ClientId  string    `json:"clientId"`
	Toopic    string    `json:"topic"`
	Qos       int       `json:"qos"`
	CreatedAt time.Time `json:"createdAt"`
	Action    string    `json:"action"`
}

func (p *Subscription) name() string { return TopicNameSubscription }

func (p *Subscription) handleTopic(service *CollectorService, ctx context.Context) error {
	session, err := mgo.Dial(service.mongoHosts)
	if err != nil {
		return err
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("iothub").C("subscriptions")

	switch p.Action {
	case ObjectActionUpdate:
		c.Insert(p)
	case ObjectActionDelete:
	case ObjectActionRegister:
	default:
	}
	return nil
}
