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
	mgr      *mqtt
	conn     net.Conn
	id       int64
	inpacket *mqttPacket
}

// newMqttSession create new session  for each client connection
func newMqttSession(m *mqtt, conn net.Conn, id int64) *mqttSession {
	return &mqttSession{mgr: m, conn: conn, id: id}
}

// handleConnection is mainprocessor for iot device client
func (s *mqttSession) handleConnection() {
}

// handlePingReq handle ping request packet
func (s *mqttSession) handlePingReq() error {
	return nil
}

// handlePingRsp handle ping response packet
func (s *mqttSession) handlePingRsp() error {
	return nil
}

// handlePubAck handle publish ack packet
func (s *mqttSession) handlePubAck() error {
	return nil
}

// handlePubComp handle publish comp packet
func (s *mqttSession) handlePubComp() error {
	return nil
}

// handlePublish handle publish packet
func (s *mqttSession) handlePublish() error {
	return nil
}

// handlePubRec handle pubrec packet
func (s *mqttSession) handlePubRec() error {
	return nil
}

// handlePubRel handle pubrel packet
func (s *mqttSession) handlePubRel() error {
	return nil
}

// handleConnect handle connect packet
func (s *mqttSession) handleConnect() error {
	return nil
}

// handleDisconnect handle disconnect packet
func (s *mqttSession) handleDisconnect() error {
	return nil
}

// handleSubscribe handle subscribe packet
func (s *mqttSession) handleSubscribe() error {
	return nil
}

// handleUnsubscribe handle unsubscribe packet
func (s *mqttSession) handleUnsubscribe() error {
	return nil
}

// handleConAck handle conack packet
func (s *mqttSession) handleConnAck() error {
	return nil
}

// handleSubAck handle suback packet
func (s *mqttSession) handleSubAck() error {
	return nil
}

// handleUnsubAck handle unsuback packet
func (s *mqttSession) handleUnsubAck() error {
	return nil
}
