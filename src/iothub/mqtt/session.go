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
	"iothub/db"
	"iothub/util/config"
	"net"
	"time"

	"github.com/golang/glog"
	"github.com/satori/go.uuid"
)

// Mqtt session state
const (
	mqttStateNew            = 0
	mqttStateConnected      = 1
	mqttStateDisconnecting  = 2
	mqttStateConnectAsync   = 3
	mqttStateConnectPending = 4
	mqttStateConnectSrv     = 5
	mqttStateDisconnectWs   = 6
	mqttStateDisconnected   = 7
	mqttStateExpiring       = 8
)

// mqtt protocol
const (
	mqttProtocolInvalid = 0
	mqttProtocol31      = 1
	mqttProtocol311     = 2
	mqttProtocolS       = 3
)

type mqttSession struct {
	mgr            *mqtt
	config         config.Config
	db             db.Database
	authplugin     base.AuthPlugin
	conn           net.Conn
	id             string
	state          uint8
	inpacket       mqttPacket
	bytesReceived  int64
	pingTime       *time.Time
	address        string
	keepalive      uint16
	protocol       uint8
	observer       base.SessionObserver
	username       string
	password       string
	lastMessageIn  time.Time
	lastMessageOut time.Time
}

// newMqttSession create new session  for each client connection
func newMqttSession(m *mqtt, conn net.Conn, id string) (*mqttSession, error) {
	s := &mqttSession{
		mgr:           m,
		config:        m.config,
		conn:          conn,
		id:            id,
		bytesReceived: 0,
		state:         mqttStateNew,
		inpacket:      newMqttPacket(),
		protocol:      mqttProtocolInvalid,
		observer:      nil,
	}
	// Load database and plugin for each session
	db, err := db.NewDatabase(m.config)
	if err != nil {
		return nil, err
	}
	plugin, err := base.LoadAuthPluginWithConfig("mqtt", m.config)
	if err != nil {
		return nil, err
	}
	s.db = db
	s.authplugin = plugin

	return s, nil
}

func (s *mqttSession) Identifier() string    { return s.id }
func (s *mqttSession) Service() base.Service { return s.mgr }
func (s *mqttSession) RegisterObserver(o base.SessionObserver) {
	if s.observer != nil {
		glog.Error("MqttSession register multiple observer")
	}
	s.observer = o
}

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

// generateId generate id fro session or client
func (s *mqttSession) generateId() string {
	return uuid.NewV4().String()
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
	case CONNECT:
		return s.handleConnect()
	case DISCONNECT:
		return s.handleDisconnect()
	case SUBSCRIBE:
		return s.handleSubscribe()
	case UNSUBSCRIBE:
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

func (s *mqttSession) updateOutMessage(mid uint16, state int) error {
	return nil
}

// handlePingReq handle ping request packet
func (s *mqttSession) handlePingReq() error {
	glog.Info("Received PINGREQ from %s", s.Identifier())
	return s.sendPingRsp()
}

// sendPingRsp send ping response to client
func (s *mqttSession) sendPingRsp() error {
	glog.Info("Sending PINGRESP to %s", s.Identifier)
	return s.sendSimpleCommand(PINGRESP)
}

// handlePingRsp handle ping response packet
func (s *mqttSession) handlePingRsp() error {
	glog.Info("Received PINGRSP form %s", s.Identifier())
	s.pingTime = nil
	return nil
}

// handleConnect handle connect packet
func (s *mqttSession) handleConnect() error {
	glog.Info("Handling CONNECT packet...")
	/*
		if s.state != mqttStateNew {
			return errors.New("Invalid session state")
		}
		// Check protocol name and version
		protocolName, err := s.inpacket.ReadString()
		if err != nil {
			return err
		}
		protocolVersion, err := s.inpacket.ReadByte()
		if err != nil {
			return err
		}
		if protocolName != PROTOCOL_NAME_V31 {
			if protocolVersion&0x7F != PROTOCOL_VERSION_V31 {
				glog.Errorf("Invalid protocol version %d in CONNECT packet", protocolVersion)
				s.sendConnAck(0, CONNACK_REFUSED_PROTOCOL_VERSION)
				return errors.New("Invalid protocol version %d in CONNECT packet", protocolVersion)
			}
			s.protocol = mqttProtocol311

		} else if protocolName != PROTOCOL_NAME_V311 {
			if protocolVersion&0x7F != PROTOCOL_VERSION_V311 {
				s.sendConnAck(0, CONNACK_REFUSED_PROTOCOL_VERSION)
				return errors.New("Invalid protocol version %d in CONNECT packet", protocolVersion)
			}
			// Reserved flags is not set to 0, must disconnect
			if s.inpacket.command&0x0F != 0x00 {
				return errors.New("Invalid protocol version %d in CONNECT packet", protocolVersion)
			}
			s.protocol = mqttProtocol311
		} else {
			return errors.New("Invalid protocol version %d in CONNECT packet", protocolVersion)
		}

		// Check connect flags
		cflags, err := s.inpacket.ReadByte()
		if err != nil {
			return nil
		}
		if s.mgr.protocol == mqttProtocol311 {
			if cflags&0x01 != 0x00 {
				return errors.New("Invalid protocol version in connect flags")
			}
		}
		cleanSesion := (cflags & 0x02) >> 1
		will = clfags & 0x04
		willQos = (cflags & 0x18) >> 3
		if willQos == 3 { // qos level3 is not supported
			return fmt.Errorf("Invalid Will Qos in CONNECT from %s", s.conn.Address())
		}

		willRetain := (cflags & 0x20) == 0x20
		passwordFlag := cflags & 0x40
		usernameFlag := cflags & 0x80
		keepalive, err := s.inpacket.ReadUint16()
		if err != nil {
			return err
		}
		s.keepalive = keepalive

		// Deal with client identifier
		clientid, err := s.inpacket.ReadString()
		if err != nil {
			return err
		}
		if len(clientid) == 0 {
			if s.protocol == mqttProtocol31 {
				s.sendConnAck(0, CONNACK_REFUSED_IDENTIFIER_REJECTED)
			} else {
				if cleanSession == 0 {
					s.sendConnAck(0, CONNACk_REFUSED_IDENTIFIER_REJECTED)
					return errors.New("Invalid mqtt packet with client id")
				} else {
					clientid = s.generateId()
				}
			}
		}
		// Deal with topc
		var willTopic string = ""
		var msg *mqttMessage = nil
		var willTopicPayload []uint8

		if will {
			// Get topic
			topic, err := s.inpacket.readString()
			if err != nil || topic == "" {
				return nil
			}
			willTopic = topic
			if s.ovserver != nil {
				willTopic = s.observer.OnGetMountPoint(s) + topic
			}
			if err := checkTopicValidity(willTopic); err != nil {
				return err
			}
			// Get willtopic's payload
			willPayloadLength, err := s.inpacket.ReadUint16()
			if err != nil {
				return err
			}
			willToicPayload, err = s.inpacket.ReadBytes(willPayloadLength)
			if err != nil {
				return err
			}
		} else {
			if s.protocol == mqttProtocol311 {
				if willQos != 0 || willRetain != 0 {
					return mqttErrorInvalidPacket
				}
			}
		} // else will

		var username string
		if usernameFlag {
			username, err = s.inpacket.ReadString()
			if err == nil {
				if passwordFlag {
					password, err = s.inpacket.ReadString()
					if err == mqttErrorInvalidProtocol {
						if s.protocol == mqttProtocol31 {
							passwordFlag = 0
						} else if s.protocol == mqttProtocol311 {
							return err
						} else {
							return err
						}
					}
				}
			} else {
				if s.protocol == mqttProtocol31 {
					usernameFlag = 0
				} else {
					return err
				}
			}
		} else { // username flag
			if s.protocol == mqttProtocol311 {
				if passwordFlag {
					return mqttErrorInvalidProtocol
				}
			}
		}

		if usernameFlag {
			if s.observer != nil {
				err := s.observer.Authenticate(s, username, password)
				switch err {
				case nil:
				case base.IotErrorAuthFailed:
					s.sendConnAck(0, CONNACK_REFUSED_NOT_AUTHORIZED)
					s.disconnect()
					return err
				default:
					s.disconnect()
					return err

				}
				// Get username and passowrd sucessfuly
				s.username = username
				s.password = password
			}
			// Get anonymous allow configuration
			allowAnonymous, _ := s.config.Bool("mqtt", "allow_anonymous")
			if !usernameFlag && allowAnonymous == false {
				// Dont allow anonymous client connection
				s.sendConnAck(0, CONNACK_REFUSED_NOT_AUTHORIZED)
				return mqttErrorInvalidProtocol
			}
		}
		// Check wether username will be used as client id,
		// The connection request will be refused if the option is set
		if s.observer != nil && s.observer.UserUserNameAsClientId() {
			if s.username != "" {
				clientid = s.username
			} else {
				s.sendConnAck(0, CONNACK_REFUSED_NOT_AUTHORIZED)
				return mqttErrorInvalidProtocol
			}
		}
		// Find if the client already has an entry, this must be done after any security check
	*/
	return nil
}

// disconnect will disconnect current connection because of protocol error
func (s *mqttSession) disconnect() {
}

// sendConnAck send CONNACK packet to client
func (s *mqttSession) sendConnAck(b uint8, reason uint8) {
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

// handleSubAck handle suback packet
func (s *mqttSession) handleSubAck() error {
	return nil
}

// handleUnsubAck handle unsuback packet
func (s *mqttSession) handleUnsubAck() error {
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
	var dup, qos, retain uint8
	var topic string
	var mid uint16
	var err error
	var payload []uint8
	var stored bool

	dup = (s.inpacket.command & 0x08) >> 3
	qos = (s.inpacket.command & 0x06) >> 1
	if qos == 3 {
		return fmt.Errorf("Invalid Qos in PUBLISH from %s, disconnectiing.", s.id)
	}

	retain = (s.inpacket.command & 0x01)

	// Topic
	if topic, err = s.inpacket.ReadString(); err != nil || topic == "" {
		return fmt.Errorf("Invalid topic in PUBLISH from %s", s.id)
	}
	if checkTopicValidity(topic) != nil {
		return fmt.Errorf("Invalid topic in PUBLISH(%s) from %s", topic, s.id)
	}
	if s.observer != nil && s.observer.OnGetMountPoint() != "" {
		topic = s.observer.OnGetMountPoint() + topic
	}

	if qos > 0 {
		mid, err = s.inpacket.ReadUint16()
		if err != nil {
			return err
		}
	}
	// Payload
	payloadlen := s.inpacket.remainingLength - s.inpacket.pos
	if payloadlen > 0 {
		limitSize, _ := s.config.Int("mqtt", "message_size_limit")
		if payloadlen > uint32(limitSize) {
			return mqttErrorInvalidProtocol
		}
		payload, err = s.inpacket.ReadBytes(payloadlen)
		if err != nil {
			return err
		}
	}
	// Check for topic access
	if s.observer != nil {
		err := s.authplugin.CheckAcl(s, s.id, s.username, topic, base.AclActionWrite)
		switch err {
		case base.ErrorAclDenied:
			return mqttErrorInvalidProtocol
		default:
			return err
		}
	}
	glog.Info("MQTT received PUBLISH from %s(d%d, q%d r%, m%d, '%s',..(%d)bytes",
		s.id, dup, qos, retain, mid, topic, payloadlen)

	// Check wether the message has been stored
	if qos > 0 {
		if found, err := s.db.FindMessage(s.id, uint(mid)); err != nil {
			return err
		} else {
			stored = found
		}
	}
	msg := db.Message{
		Id:        uint(mid),
		Direction: db.MessageDirectionIn,
		State:     0,
		Qos:       qos,
		Retain:    (retain > 0),
		Payload:   payload,
	}

	if !stored {
		dup = 0
		if err := s.db.StoreMessage(s.id, msg); err != nil {
			return err
		}
	} else {
		dup = 1
	}

	switch qos {
	case 0:
		err = s.db.QueueMessage(s.id, msg)
	case 1:
		err = s.db.QueueMessage(s.id, msg)
		err = s.sendPubAck(mid)
	case 2:
		err = nil
		if dup > 0 {
			err = s.db.InsertMessage(s.id, int(mid), db.MessageDirectionIn, msg)
		}
		if err == nil {
			err = s.sendPubRec(mid)
		}
	default:
		err = mqttErrorInvalidProtocol
	}

	return err
}

// sendPubAck
func (s *mqttSession) sendPubAck(mid uint16) error {
	return nil
}

// sendPubRec
func (s *mqttSession) sendPubRec(mid uint16) error {
	return nil
}

// handlePubRec handle pubrec packet
func (s *mqttSession) handlePubRec() error {
	mid, err := s.inpacket.ReadUint16()
	if err != nil {
		return err
	}
	glog.Info("Client %s received PUBRED mid:%d", s.id, mid)
	err = s.updateOutMessage(mid, mqttMessageStateWaitForPubComp)
	if err == base.IotErrorNotFound {
		glog.Errorf("Received  PUBREC from %s for an unknown packet identifier", s.id)
	} else if err != nil {
		return err
	}
	return s.sendPubRel(mid)
}

// handlePubRel handle pubrel packet
func (s *mqttSession) handlePubRel() error {
	return nil
}

func (s *mqttSession) sendPubRel(mid uint16) error {
	return nil
}
