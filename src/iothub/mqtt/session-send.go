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
	"fmt"

	"github.com/golang/glog"
)

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
func (s *mqttSession) sendConnAck(ack uint8, result uint8) error {
	glog.Info("Sending CONNACK from %s", s.id)

	packet := &mqttPacket{
		command:         CONNACK,
		remainingLength: 2,
	}
	packet.payload[packet.pos+0] = ack
	packet.payload[packet.pos+1] = result
	if err := s.initializePacket(packet); err != nil {
		return nil
	}

	return s.queuePacket(packet)
}

// initializePacket initialize packet according to preinitialized data
func (s *mqttSession) initializePacket(p *mqttPacket) error {
	var remainingBytes = [5]uint8{}

	remainingLength := p.remainingLength
	p.remainingCount = 0
	for {
		b := remainingLength % 128
		remainingLength = remainingLength / 128
		if remainingLength > 0 {
			b = b | 0x80
		}
		remainingBytes[p.remainingLength] = uint8(b)
		p.remainingCount++
		if remainingLength < 0 || p.remainingCount >= 5 {
			break
		}
		p.length = p.remainingLength + 1 + uint32(p.remainingCount)
	}
	if p.remainingCount == 5 {
		return fmt.Errorf("Invalid packet(%d) payload size", p.command)
	}
	p.payload = make([]uint8, p.length)
	p.payload[0] = p.command
	for index, b := range remainingBytes {
		p.payload[index+1] = b
	}
	p.pos = 1 + uint32(p.remainingCount)
	return nil
}

func (s *mqttSession) queuePacket(p *mqttPacket) error {
	return nil

}

// sendSubAck send subscription acknowledge to client
func (s *mqttSession) sendSubAck(mid uint16, payload []uint8) error {
	packet := new(mqttPacket)
	packet.command = SUBACK
	packet.remainingLength = 2 + uint32(len(payload))
	if err := s.initializePacket(packet); err != nil {
		return err
	}

	packet.WriteUint16(mid)
	if len(payload) > 0 {
		packet.WriteBytes(payload)
	}
	return s.queuePacket(packet)
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

func (s *mqttSession) sendPubComp(mid uint16) error {
	return nil
}
