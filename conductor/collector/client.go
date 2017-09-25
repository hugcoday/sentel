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

// Client
type Client struct {
	topicBase
	ClientId        string `json:"clientId"`
	UserName        string `json:"userName"`
	IpAddress       string `json:"ipAddress"`
	Port            uint16 `json:"port"`
	CleanSession    bool   `json:"cleanSession"`
	ProtocolVersion string `json:"protocolVersion"`
	Keepalive       uint16 `json:"keepalive"`
	ConnectedAt     string `json:"connectedAt"`
}

func (p *Client) name() string { return TopicNameClient }
func (p *Client) clone() topicObject {
	return &Client{
		topicBase:       p.topicBase,
		ClientId:        p.ClientId,
		UserName:        p.UserName,
		IpAddress:       p.IpAddress,
		Port:            p.Port,
		CleanSession:    p.CleanSession,
		ProtocolVersion: p.ProtocolVersion,
		Keepalive:       p.Keepalive,
		ConnectedAt:     p.ConnectedAt,
	}
}

func (p *Client) handleTopic(service *CollectorService, ctx context.Context) error {
	db, err := service.getDatabase()
	if err != nil {
		return err
	}
	defer db.Session.Close()
	c := db.C("clients")

	result := Client{}
	if err := c.Find(bson.M{"ClientId": p.ClientId}).One(&result); err == nil {
		// Existed client found
		return c.Update(result, p)
	} else {
		return c.Insert(p)
	}
}
