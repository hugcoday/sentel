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

import "net"

type mqttSession struct {
	mgr  *mqtt
	conn net.Conn
	id   int64
}

func newMqttSession(m *mqtt, conn net.Conn, id int64) *mqttSession {
	return &mqttSession{mgr: m, conn: conn, id: id}
}

func (s *mqttSession) handleConnection() {
}
