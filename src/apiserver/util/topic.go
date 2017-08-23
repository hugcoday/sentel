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

package util

import (
	"encoding/json"
)

const (
	TopicNameTenant  = "tenant"
	TopicNameProduct = "product"
	TopicNameDevice  = "device"
)

const (
	ObjectActionRegister   = "register"
	ObjectActionUnregister = "unregister"
	ObjectActionRetrieve   = "retrieve"
	ObjectActionDelete     = "delete"
	ObjectActionUpdate     = "update"
)

type TenantTopic struct{}

// ProductTopic
type ProductTopic struct {
	ProductId   string `json:"productId"`
	ProductName string `json:"productName"`
	Action      string `json:"action"`
	TenantId    string `json:"tenantId"`

	encoded []byte
	err     error
}

func (p *ProductTopic) ensureEncoded() {
	if p.encoded == nil && p.err == nil {
		p.encoded, p.err = json.Marshal(p)
	}
}

func (p *ProductTopic) Length() int {
	p.ensureEncoded()
	return len(p.encoded)
}

func (p *ProductTopic) Encode() ([]byte, error) {
	p.ensureEncoded()
	return p.encoded, p.err
}

// DeviceTopic
type DeviceTopic struct {
	DeviceId     string `json:productId"`
	DeviceSecret string `json:productKey"`
	Action       string `json:"action"`
	ProductId    string `json:"productId"`

	encoded []byte
	err     error
}

func (p *DeviceTopic) ensureEncoded() {
	if p.encoded == nil && p.err == nil {
		p.encoded, p.err = json.Marshal(p)
	}
}

func (p *DeviceTopic) Length() int {
	p.ensureEncoded()
	return len(p.encoded)
}

func (p *DeviceTopic) Encode() ([]byte, error) {
	p.ensureEncoded()
	return p.encoded, p.err
}
