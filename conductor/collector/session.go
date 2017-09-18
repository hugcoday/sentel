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

import "encoding/json"

// Session
type Session struct {
	ClientId     string `json:"clientId"`
	CleanSession bool   `json:"cleanSession"`
	CreatedAt    string `json:"createdAt"`

	encoded []byte
	err     error
}

func (p *Session) ensureEncoded() {
	if p.encoded == nil && p.err == nil {
		p.encoded, p.err = json.Marshal(p)
	}
}

func (p *Session) Length() int {
	p.ensureEncoded()
	return len(p.encoded)
}

func (p *Session) Encode() ([]byte, error) {
	p.ensureEncoded()
	return p.encoded, p.err
}
