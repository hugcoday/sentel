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

import "github.com/golang/glog"

// sendSimpleCommand send a simple command
func (s *mqttSession) sendSimpleCommand(cmd uint8) error {
	p := mqttPacket{
		command:        cmd,
		remainingCount: 0,
	}
	return s.sendPacket(&p)
}

// sendPacket send packet to client
// TODO:
func (s *mqttSession) sendPacket(p *mqttPacket) error {
	return nil
}

func (s *mqttSession) updateOutMessage(mid uint16, state int) error {
	return nil
}

// sendPingRsp send ping response to client
func (s *mqttSession) sendPingRsp() error {
	glog.Info("Sending PINGRESP to %s", s.Identifier)
	return s.sendSimpleCommand(PINGRESP)
}

// sendConnAck send connection response to client
func (s *mqttSession) sendConnAck(mid uint16, reason uint16) error {
	return nil
}

// sendSubAck send subscription acknowledge to client
func (s *mqttSession) sendSubAck(mid uint16, payload []uint8) error {
	return nil
}

// sendCommandWithMid send command with message identifier
func (s *mqttSession) sendCommandWithMid(command uint8, mid uint16, dup bool) error {
	return nil
}

// sendPubAck
func (s *mqttSession) sendPubAck(mid uint16) error {
	return nil
}

// sendPubRec
func (s *mqttSession) sendPubRec(mid uint16) error {
	return nil
}

func (s *mqttSession) sendPubRel(mid uint16) error {
	return nil
}
