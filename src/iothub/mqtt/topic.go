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

package mqtt

import (
	"encoding/json"
)

const (
	TopicNameSession = "mqtt-session"
)

const (
	ObjectActionRegister   = "register"
	ObjectActionUnregister = "unregister"
	ObjectActionRetrieve   = "retrieve"
	ObjectActionDelete     = "delete"
	ObjectActionUpdate     = "update"
)

type TenantTopic struct{}

type SessionTopic struct {
	Launcher  string `json:"launcher"`
	SessionId string `json:"sessionId"`
	Action    string `json:"action"`
	State     uint8  `json:"oldState"`

	encoded []byte
	err     error
}

func (p *SessionTopic) ensureEncoded() {
	if p.encoded == nil && p.err == nil {
		p.encoded, p.err = json.Marshal(p)
	}
}

func (p *SessionTopic) Length() int {
	p.ensureEncoded()
	return len(p.encoded)
}

func (p *SessionTopic) Encode() ([]byte, error) {
	p.ensureEncoded()
	return p.encoded, p.err
}
