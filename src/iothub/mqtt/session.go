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
	"time"

	"github.com/golang/glog"
)

// Mqtt session state
const (
	mqttStateNew = 0
)

type mqttSession struct {
	mgr           *mqtt
	conn          net.Conn
	id            string
	state         uint8
	inpacket      mqttPacket
	bytesReceived int64
	pingTime      *time.Time
	keepalive     uint16
}

// newMqttSession create new session  for each client connection
func newMqttSession(m *mqtt, conn net.Conn, id string) *mqttSession {
	return &mqttSession{
		mgr:           m,
		conn:          conn,
		id:            id,
		bytesReceived: 0,
		state:         mqttStateNew,
		inpacket:      newMqttPacket(),
	}
}

func (s *mqttSession) Identifier() string    { return s.id }
func (s *mqttSession) Service() base.Service { return s.mgr }

// handle is mainprocessor for iot device client
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
	_, err := io.Copy(&buf, s.conn)
	if err != nil {
		return fmt.Errorf("read packet error:%s", err)
	}
	_, err = s.inpacket.DecodeFromBytes(buf.Bytes(), base.NilDecodeFeedback{})
	return err
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
