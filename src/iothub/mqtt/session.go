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
	"bytes"
	"fmt"
	"io"
	"iothub/base"
	"net"

	"github.com/golang/glog"
)

type mqttSession struct {
	mgr           *mqtt
	conn          net.Conn
	id            string
	state         uint8
	inpacket      mqttPacket
	bytesReceived int64
}

// newMqttSession create new session  for each client connection
func newMqttSession(m *mqtt, conn net.Conn, id string) *mqttSession {
	return &mqttSession{
		mgr:           m,
		conn:          conn,
		id:            id,
		bytesReceived: 0,
		state:         stateNewConnect,
		inpacket: mqttPacket{
			command:        0,
			pos:            0,
			length:         0,
			remainingCount: 0,
			payload:        []byte{},
		},
	}
}

func (s *mqttSession) Identifier() string    { return s.id }
func (s *mqttSession) Service() base.Service { return s.mgr }

// handleConnection is mainprocessor for iot device client
// Loop to read packet from conn
func (s *mqttSession) Handle() error {
	defer s.Destroy()
	for {
		if err := s.readPacket(); err != nil {
			glog.Errorf("Reading packet error occured for connection:%d", s.id)
			return err
		}
		if err := s.handlePacket(); err != nil {
			glog.Errorf("Handle packet error occured for connection:%d", s.id)
			return err
		}
	}
	return nil
}

// removeConnection remove current connection from mqttManaager if errors occured
func (s *mqttSession) Destroy() error {
	s.conn.Close()
	s.mgr.RemoveSession(s)
	return nil
}

// readPacket read a whole mqtt packet from session
// TODO: underlay's read method  should be payed attention
func (s *mqttSession) readPacket() error {
	// Assumption: Read whole packet data in one read calling
	var buf bytes.Buffer
	len, err := io.Copy(&buf, s.conn)

	if err != nil {
		return fmt.Errorf("read packet error:%s", err)
	}
	s.bytesReceived += int64(len)
	glog.Info("Reading bytes from connection:totoal(%d), current(%d)", s.bytesReceived, len)

	// Start from new packet
	if s.inpacket.command == 0 {
		cmd, err := buf.ReadByte()
		if err != nil {
			return fmt.Errorf("Reading error occured for connection:%d", s.id)
		}
		s.inpacket.command = cmd
		// Check connection state
		// Client must send CONNECT as their first command
		if s.state == stateNewConnect && (cmd&0xF0) != CONNECT {
			return fmt.Errorf("Connection state error for %d", s.id)
		}
	}

	if s.inpacket.remainingCount <= 0 {
		for {
			b, err := buf.ReadByte()
			if err != nil {
				return fmt.Errorf("Reading error occured for connection:%s", err)
			}
			s.bytesReceived++
			s.inpacket.remainingCount--
			if s.inpacket.remainingCount == -4 {
				return fmt.Errorf("Invalid protocol")
			}
			s.inpacket.remainingCount += int32(b&127) * s.inpacket.remainingMult
			s.inpacket.remainingMult *= 128
			if b&128 != 0 {
				break
			}
		}
	}
	// We have finished reading remaining length
	s.inpacket.remainingCount *= -1
	if s.inpacket.remainingCount > 0 {
		var index int32
		for index = 0; index < s.inpacket.remainingCount; index++ {
			// Assumption: whole packet had been read into buffer
			n, err := buf.ReadByte()
			if err != nil {
				return fmt.Errorf("Reading remaining packet payload error:%s", err)
			}
			s.inpacket.payload = append(s.inpacket.payload, n)
			s.bytesReceived++
		}
	}
	return nil
}

// handlePacket dispatch packet by packet type
func (s *mqttSession) handlePacket() error {
	switch s.inpacket.command & 0xF0 {
	case PINGREQ:
		return s.handlePingReq()
	case PINGRESP:
		return s.handlePingRsp()
	case PUBACK:
		return s.handlePubAck()
	case PUBCOMP:
		return s.handlePubComp()
	case PUBLISH:
		return s.handlePublish()
	case PUBREC:
		return s.handlePubRec()
	case PUBREL:
		return s.handlePubRel()
	case CONNACK:
		return s.handleConnAck()
	case SUBACK:
		return s.handleSubAck()
	case UNSUBACK:
		return s.handleUnsubAck()
	}
	return fmt.Errorf("Unrecognized protocol command:%d", int(s.inpacket.command&0xF0))
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
