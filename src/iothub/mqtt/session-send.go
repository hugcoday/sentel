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
	"errors"
	"fmt"

	"github.com/golang/glog"
)

// sendSimpleCommand send a simple command
func (s *mqttSession) sendSimpleCommand(cmd uint8) error {
	p := &mqttPacket{
		command:        cmd,
		remainingCount: 0,
	}
	return s.queuePacket(p)
}

// sendPingRsp send ping response to client
func (s *mqttSession) sendPingRsp() error {
	glog.Info("Sending PINGRESP to %s", s.Identifier)
	return s.sendSimpleCommand(PINGRESP)
}

// sendConnAck send connection response to client
func (s *mqttSession) sendConnAck(ack uint8, result uint8) error {
	glog.Infof("Sending CONNACK from '%s'", s.id)

	packet := &mqttPacket{
		command:         CONNACK,
		remainingLength: 2,
	}
	s.initializePacket(packet)
	packet.payload[packet.pos+0] = ack
	packet.payload[packet.pos+1] = result

	return s.queuePacket(packet)
}

// initializePacket initialize packet according to preinitialized data
func (s *mqttSession) initializePacket(p *mqttPacket) error {
	var remainingBytes = [5]uint8{}

	remainingLength := p.remainingLength
	p.remainingCount = 0
	for remainingLength > 0 && p.remainingCount < 5 {
		b := remainingLength % 128
		remainingLength = remainingLength / 128
		if remainingLength > 0 {
			b = b | 0x80
		}
		remainingBytes[p.remainingCount] = uint8(b)
		p.remainingCount++
	}
	if p.remainingCount == 5 {
		return fmt.Errorf("Invalid packet(%d) payload size", p.command)
	}
	p.length = p.remainingLength + 1 + p.remainingCount
	p.payload = make([]uint8, p.length)
	p.payload[0] = p.command
	for i := 0; i < p.remainingCount; i++ {
		p.payload[i+1] = remainingBytes[i]
	}
	p.pos = 1 + p.remainingCount
	return nil
}

// sendSubAck send subscription acknowledge to client
func (s *mqttSession) sendSubAck(mid uint16, payload []uint8) error {
	packet := new(mqttPacket)
	packet.command = SUBACK
	packet.remainingLength = 2 + int(len(payload))
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
	packet := new(mqttPacket)
	packet.command = command
	if dup {
		packet.command |= 8
	}
	packet.remainingLength = 2
	if err := s.initializePacket(packet); err != nil {
		return err
	}
	packet.payload[packet.pos+0] = uint8((mid & 0xFF00) >> 8)
	packet.payload[packet.pos+1] = uint8(mid & 0xff)
	return s.queuePacket(packet)
}

// sendPubAck
func (s *mqttSession) sendPubAck(mid uint16) error {
	glog.Info("Sending PUBACK to %s with MID:%d", s.id, mid)
	return s.sendCommandWithMid(PUBACK, mid, false)
}

// sendPubRec
func (s *mqttSession) sendPubRec(mid uint16) error {
	glog.Info("Sending PUBRREC to %s with MID:%d", s.id, mid)
	return s.sendCommandWithMid(PUBREC, mid, false)
}

func (s *mqttSession) sendPubComp(mid uint16) error {
	glog.Info("Sending PUBCOMP to %s with MID:%d", s.id, mid)
	return s.sendCommandWithMid(PUBCOMP, mid, false)
}

func (s *mqttSession) queuePacket(p *mqttPacket) error {
	p.pos = 0
	p.toprocess = 0

	s.outPacketMutex.Lock()
	s.outPackets = append(s.outPackets, p)
	s.lastOutPacket = p
	s.outPacketMutex.Unlock()
	return s.writePacket()
}

func (s *mqttSession) writePacket() error {
	s.currentOutPacketMutex.Lock()
	defer s.currentOutPacketMutex.Unlock()

	s.outPacketMutex.Lock()
	if len(s.outPackets) > 0 && s.currentOutPacket == nil {
		s.currentOutPacket = s.outPackets[0]
		s.outPackets = s.outPackets[1:]
		if len(s.outPackets) == 0 {
			s.lastOutPacket = nil
		}
	}
	s.outPacketMutex.Unlock()

	if s.state == mqttStateConnectPending {
		return errors.New("Write packet in wrong session state")
	}

	for s.currentOutPacket != nil {
		packet := s.currentOutPacket
		for packet.toprocess > 0 {
			len, err := s.netWrite(packet.payload[packet.pos:packet.toprocess])
			if err != nil {
				return nil
			}
			if len > 0 {
				packet.toprocess -= len
				packet.pos += len
			} else {
				return nil
			}
		}
		// Notify observer

		// Process net packet
		s.outPacketMutex.Lock()
		if len(s.outPackets) > 0 {
			s.currentOutPacket = s.outPackets[0]
			s.outPackets = s.outPackets[1:]
		} else {
			s.currentOutPacket = nil
			s.lastOutPacket = nil
		}
		s.outPacketMutex.Unlock()
	}
	return nil
}

func (s *mqttSession) netWrite(data []uint8) (int, error) {
	return s.conn.Write(data)
}

func (s *mqttSession) updateOutMessage(mid uint16, state int) error {
	return nil
}
