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
	"net"

	"github.com/golang/glog"
)

const (
	mqttStateNewConnect = 0
)

type mqttSession struct {
	mgr           *mqtt
	conn          net.Conn
	id            int64
	state         uint8
	inpacket      mqttPacket
	bytesReceived int64
}

// newMqttSession create new session  for each client connection
func newMqttSession(m *mqtt, conn net.Conn, id int64) *mqttSession {
	return &mqttSession{
		mgr:           m,
		conn:          conn,
		id:            id,
		inpacket:      mqttPacket{},
		bytesReceived: 0,
		state:         mqttStateNewConnect,
	}
}

// handleConnection is mainprocessor for iot device client
// Loop to read packet from conn
func (s *mqttSession) handleConnection() {
	defer s.removeConnection()
	for {
		if err := s.readPacket(); err != nil {
			glog.Errorf("Reading packet error occured for connection:%d", s.id)
			return
		}
		if err := s.handlePacket(); err != nil {
			glog.Errorf("Handle packet error occured for connection:%d", s.id)
			return
		}
	}
}

// removeConnection remove current connection from mqttManaager if errors occured
func (s *mqttSession) removeConnection() {
	s.conn.Close()
	s.mgr.removeSession(s)
}

func (s *mqttSession) readPacket() error {
	var bytes []byte
	for {
		// Read data from client
		n, err := s.conn.Read(bytes)
		// Check wether reading error occured, exit mainloop if error occured
		if err != nil {
			return err
		}
		s.bytesReceived += int64(n)
		// Start from new packet
		if s.inpacket.command == 0 {
			if n > 0 {
				s.inpacket.command = bytes[0]
				// Check connection state
				// Client must send CONNECT as their first command
				if s.state == mqttStateNewConnect && (bytes[0]&0xF0) != CONNECT {
					return fmt.Errorf("Connection state error for %d", s.id)
				}

			} else {
				return fmt.Errorf("Reading error occured for connection:%d", s.id)
			}
		}

		if s.inpacket.remainingCount <= 0 {

		}
	}
	return nil
}

// handlePacket dispatch packet by packet type
func (s *mqttSession) handlePacket() error {
	return nil
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
