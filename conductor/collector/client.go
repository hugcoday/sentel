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
	"gopkg.in/mgo.v2/bson"
)

// Client
type Client struct {
	ClientId        string `json:"clientId"`
	UserName        string `json:"userName"`
	IpAddress       string `json:"ipAddress"`
	Port            uint16 `json:"port"`
	CleanSession    bool   `json:"cleanSession"`
	ProtocolVersion string `json:"protocolVersion"`
	Keepalive       uint16 `json:"keepalive"`
	ConnectedAt     string `json:"connectedAt"`

	encoded []byte
	err     error
}

func (p *Client) ensureEncoded() {
	if p.encoded == nil && p.err == nil {
		p.encoded, p.err = json.Marshal(p)
	}
}

func (p *Client) Length() int {
	p.ensureEncoded()
	return len(p.encoded)
}

func (p *Client) Encode() ([]byte, error) {
	p.ensureEncoded()
	return p.encoded, p.err
}

func (p *Client) name() string { return TopicNameClient }

func (p *Client) handleTopic(service *CollectorService, ctx context.Context, value []byte) error {
	var client Client
	if err := json.Unmarshal(value, client); err != nil {
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
	c := session.DB("iothub").C("clients")

	result := Client{}
	if err := c.Find(bson.M{"ClientId": client.ClientId}).One(&result); err == nil {
		// Existed client found
		return c.Update(result, client)
	} else {
		return c.Insert(client)
	}
}
