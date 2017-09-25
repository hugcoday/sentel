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

	"gopkg.in/mgo.v2/bson"
)

// Session
type Session struct {
	topicBase
	Action             string `json:"action"`
	ClientId           string `json:"clientId"`
	CleanSession       bool   `json:"cleanSession"`
	MessageMaxInflight uint64 `json:"messageMaxInflight"`
	MessageInflight    uint64 `json:"messageInflight"`
	MessageInQueue     uint64 `json:"messageInQueue"`
	MessageDropped     uint64 `json:"messageDropped"`
	AwaitingRel        uint64 `json:"awaitingRel"`
	AwaitingComp       uint64 `json:"awaitingComp"`
	AwaitingAck        uint64 `json:"awaitingAck"`
	CreatedAt          string `json:"createdAt"`
}

func (p *Session) name() string { return TopicNameSubscription }
func (p *Session) clone() topicObject {
	return &Session{
		topicBase:          p.topicBase,
		Action:             p.Action,
		ClientId:           p.ClientId,
		CleanSession:       p.CleanSession,
		MessageMaxInflight: p.MessageMaxInflight,
		MessageInflight:    p.MessageInflight,
		MessageInQueue:     p.MessageInQueue,
		MessageDropped:     p.MessageDropped,
		AwaitingRel:        p.AwaitingRel,
		AwaitingComp:       p.AwaitingComp,
		AwaitingAck:        p.AwaitingAck,
		CreatedAt:          p.CreatedAt,
	}
}

func (p *Session) handleTopic(service *CollectorService, ctx context.Context) error {
	db, err := service.getDatabase()
	if err != nil {
		return err
	}
	defer db.Session.Close()
	c := db.C("subscriptions")

	switch p.Action {
	case ObjectActionUpdate:
		result := Session{}
		if err := c.Find(bson.M{"ClientId": p.ClientId}).One(&result); err == nil {
			return c.Update(result, p)
		} else {
			c.Insert(p)
		}
	case ObjectActionDelete:
	case ObjectActionRegister:
	default:
	}
	return nil
}
