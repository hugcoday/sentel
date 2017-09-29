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
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/cloustone/sentel/libs"

	"github.com/cloustone/sentel/iothub/base"

	auth "github.com/cloustone/sentel/iothub/auth"
	uuid "github.com/satori/go.uuid"

	"github.com/golang/glog"
)

// Mqtt session state
const (
	mqttStateInvalid        = 0
	mqttStateNew            = 1
	mqttStateConnected      = 2
	mqttStateDisconnecting  = 3
	mqttStateConnectAsync   = 4
	mqttStateConnectPending = 5
	mqttStateConnectSrv     = 6
	mqttStateDisconnectWs   = 7
	mqttStateDisconnected   = 8
	mqttStateExpiring       = 9
)

// mqtt protocol
const (
	mqttProtocolInvalid = 0
	mqttProtocol31      = 1
	mqttProtocol311     = 2
	mqttProtocolS       = 3
)

type mqttSession struct {
	mgr               *mqtt
	config            libs.Config
	storage           Storage
	authapi           auth.IAuthAPI
	conn              net.Conn
	id                string
	clientID          string
	state             uint8
	inpacket          mqttPacket
	bytesReceived     int64
	pingTime          *time.Time
	address           string
	keepalive         uint16
	protocol          uint8
	observer          base.SessionObserver
	username          string
	password          string
	lastMessageIn     time.Time
	lastMessageOut    time.Time
	cleanSession      uint8
	isDroping         bool
	willMsg           *mqttMessage
	stateMutex        sync.Mutex
	sendStopChannel   chan int
	sendPacketChannel chan *mqttPacket
	sendMsgChannel    chan *mqttMessage
	waitgroup         sync.WaitGroup
	stats             *base.Stats
	metrics           *base.Metrics

	// resume field
	msgs              []*mqttMessage
	storedMsgs        []*mqttMessage
}

// newMqttSession create new session  for each client connection
func newMqttSession(m *mqtt, conn net.Conn, id string) (*mqttSession, error) {
	// Get session message queue size, if it is not set, default is 10
	qsize, err := m.config.Int(m.protocol, "session_queue_size")
	if err != nil {
		qsize = 10
	}
	authapi, err := auth.NewAuthApi(m.config)
	if err != nil {
		return nil, err
	}
	var msgqsize int
	msgqsize, err = m.config.Int(m.protocol, "session_msg_queue_size")
	if err != nil {
		msgqsize = 10
	}

	s := &mqttSession{
		mgr:               m,
		config:            m.config,
		storage:           m.storage,
		conn:              conn,
		id:                id,
		bytesReceived:     0,
		state:             mqttStateNew,
		inpacket:          newMqttPacket(),
		protocol:          mqttProtocolInvalid,
		observer:          nil,
		sendStopChannel:   make(chan int),
		sendPacketChannel: make(chan *mqttPacket, qsize),
		authapi:           authapi,
		stats:             base.NewStats(true),
		metrics:           base.NewMetrics(true),
		sendMsgChannel:    make(chan *mqttMessage, msgqsize),
		msgs:              make([]*mqttMessage, msgqsize),
		storedMsgs:        make([]*mqttMessage, msgqsize),
	}

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
func (s *mqttSession) GetStats() *base.Stats     { return s.stats }
func (s *mqttSession) GetMetrics() *base.Metrics { return s.metrics }
func (s *mqttSession) Info() *base.SessionInfo   { return nil }

// launchPacketSendHandler launch goroutine to send packet queued for client
func (s *mqttSession) launchPacketSendHandler() {
	go func(stopChannel chan int, packetChannel chan *mqttPacket, msgChannel chan *mqttMessage) {
		defer s.waitgroup.Add(1)

		for {
			select {
			case <-stopChannel:
				return
			case p := <-packetChannel:
				for p.toprocess > 0 {
					len, err := s.conn.Write(p.payload[p.pos:p.toprocess])
					if err != nil {
						glog.Fatal("Failed to send packet to '%s:%s'", s.id, err)
						return
					}
					if len > 0 {
						p.toprocess -= len
						p.pos += len
					} else {
						glog.Fatal("Failed to send packet to '%s'", s.id)
						return
					}
				}
			case msg := <-msgChannel:
				s.processMessage(msg)
			case <-time.After(1 * time.Second):
				s.processTimeout()
			}
		}
	}(s.sendStopChannel, s.sendPacketChannel, s.sendMsgChannel)
}

// processMessage proceess messages
func (s *mqttSession) processMessage(msg *mqttMessage) error {

	return nil
}

// processTimeout proceess timeout
func (s *mqttSession) processTimeout() error {

	return nil
}

// Handle is mainprocessor for iot device client
func (s *mqttSession) Handle() error {

	glog.Infof("Handling session:%s", s.id)
	defer s.Destroy()

	s.launchPacketSendHandler()
	for {
		var err error
		if err = s.inpacket.DecodeFromReader(s.conn, base.NilDecodeFeedback{}); err != nil {
			glog.Error(err)
			return err
		}
		switch s.inpacket.command & 0xF0 {
		case PINGREQ:
			err = s.handlePingReq()
		case CONNECT:
			err = s.handleConnect()
		case DISCONNECT:
			err = s.handleDisconnect()
		case PUBLISH:
			err = s.handlePublish()
		case PUBREL:
			err = s.handlePubRel()
		case SUBSCRIBE:
			err = s.handleSubscribe()
		case UNSUBSCRIBE:
			err = s.handleUnsubscribe()
		default:
			err = fmt.Errorf("Unrecognized protocol command:%d", int(s.inpacket.command&0xF0))
		}
		if err != nil {
			glog.Error(err)
			return err
		}
		// Check sesstion state
		if s.state == mqttStateDisconnected {
			break
		}
		s.inpacket.Clear()
	}
	return nil
}

// Destroy will destory the current session
func (s *mqttSession) Destroy() error {
	// Stop packet sender goroutine
	s.sendStopChannel <- 1
	s.waitgroup.Wait()
	if s.conn != nil {
		s.conn.Close()
	}
	s.mgr.removeSession(s)
	return nil
}

// generateId generate id fro session or client
func (s *mqttSession) generateId() string {
	return uuid.NewV4().String()
}

// handlePingReq handle ping request packet
func (s *mqttSession) handlePingReq() error {
	glog.Infof("Received PINGREQ from %s", s.Identifier())
	return s.sendPingRsp()
}

// handleConnect handle connect packet
func (s *mqttSession) handleConnect() error {
	glog.Infof("Handling CONNECT packet from %s", s.id)

	if s.state != mqttStateNew {
		return errors.New("Invalid session state")
	}
	// Check protocol name and version
	protocolName, err := s.inpacket.readString()
	if err != nil {
		return err
	}
	protocolVersion, err := s.inpacket.readByte()
	if err != nil {
		return err
	}
	switch protocolName {
	case PROTOCOL_NAME_V31:
		if protocolVersion&0x7F != PROTOCOL_VERSION_V31 {
			s.sendConnAck(0, CONNACK_REFUSED_PROTOCOL_VERSION)
			return fmt.Errorf("Invalid protocol version '%d' in CONNECT packet", protocolVersion)
		}
		s.protocol = mqttProtocol311

	case PROTOCOL_NAME_V311:
		if protocolVersion&0x7F != PROTOCOL_VERSION_V311 {
			s.sendConnAck(0, CONNACK_REFUSED_PROTOCOL_VERSION)
			return fmt.Errorf("Invalid protocol version '%d' in CONNECT packet", protocolVersion)
		}
		// Reserved flags is not set to 0, must disconnect
		if s.inpacket.command&0x0F != 0x00 {
			return fmt.Errorf("Invalid protocol version '%d' in CONNECT packet", protocolVersion)
		}
		s.protocol = mqttProtocol311
	default:
		return fmt.Errorf("Invalid protocol name '%s' in CONNECT packet", protocolName)
	}

	// Check connect flags
	cflags, err := s.inpacket.readByte()
	if err != nil {
		return nil
	}
	/*
		if s.mgr.protocol == mqttProtocol311 {
			if cflags&0x01 != 0x00 {
				return errors.New("Invalid protocol version in connect flags")
			}
		}
	*/
	cleanSession := (cflags & 0x02) >> 1
	will := cflags & 0x04
	willQos := (cflags & 0x18) >> 3
	if willQos >= 3 { // qos level3 is not supported
		return fmt.Errorf("Invalid Will Qos in CONNECT from %s", s.id)
	}

	willRetain := (cflags & 0x20) == 0x20
	passwordFlag := cflags & 0x40
	usernameFlag := cflags & 0x80
	keepalive, err := s.inpacket.readUint16()
	if err != nil {
		return err
	}
	s.keepalive = keepalive

	// Deal with client identifier
	clientid, err := s.inpacket.readString()
	if err != nil {
		return err
	}
	if clientid == "" {
		if s.protocol == mqttProtocol31 {
			s.sendConnAck(0, CONNACK_REFUSED_IDENTIFIER_REJECTED)
		} else {
			if option, err := s.config.Bool(s.mgr.protocol, "allow_zero_length_clientid"); err != nil && (!option || cleanSession == 0) {
				s.sendConnAck(0, CONNACK_REFUSED_IDENTIFIER_REJECTED)
				return errors.New("Invalid mqtt packet with client id")
			}
			
			clientid = s.generateId()
		}
	}

	// TODO: clientid_prefixes check

	
	// Deal with topc
	var willTopic string
	var willMsg *mqttMessage
	var payload []uint8

	if will > 0 {
		willMsg = new(mqttMessage)
		// Get topic
		topic, err := s.inpacket.readString()
		if err != nil || topic == "" {
			return nil
		}
		willTopic = topic
		if s.observer != nil {
			willTopic = s.observer.OnGetMountPoint() + topic
		}
		if err := checkTopicValidity(willTopic); err != nil {
			return err
		}
		// Get willtopic's payload
		willPayloadLength, err := s.inpacket.readUint16()
		if err != nil {
			return err
		}
		if willPayloadLength > 0 {
			payload, err = s.inpacket.readBytes(int(willPayloadLength))
			if err != nil {
				return err
			}
		}
	} else {
		if s.protocol == mqttProtocol311 {
			if willQos != 0 || willRetain {
				return mqttErrorInvalidProtocol
			}
		}
	} // else will

	var username string
	var password string
	if usernameFlag > 0 {
		username, err = s.inpacket.readString()
		if err == nil {
			if passwordFlag > 0 {
				password, err = s.inpacket.readString()
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
			if passwordFlag > 0 {
				return mqttErrorInvalidProtocol
			}
		}
	}

	if usernameFlag > 0 {
		if s.observer != nil {
			err := s.observer.OnAuthenticate(s, username, password)
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
		allowAnonymous, _ := s.config.Bool(s.mgr.protocol, "allow_anonymous")
		if usernameFlag > 0 && allowAnonymous == false {
			// Dont allow anonymous client connection
			s.sendConnAck(0, CONNACK_REFUSED_NOT_AUTHORIZED)
			return mqttErrorInvalidProtocol
		}
	}
	// Check wether username will be used as client id,
	// The connection request will be refused if the option is set
	if option, err := s.config.Bool(s.mgr.protocol, "user_name_as_client_id"); err != nil && option {
		if s.username != "" {
			clientid = s.username
		} else {
			s.sendConnAck(0, CONNACK_REFUSED_NOT_AUTHORIZED)
			return mqttErrorInvalidProtocol
		}
	}
	conack := 0
	// Find if the client already has an entry, this must be done after any security check
	if found, _ := s.storage.FindSession(clientid); found != nil {
		// Found old session
		if found.state == mqttStateInvalid {
			glog.Errorf("Invalid session(%s) in store", found.id)
		}
		if s.protocol == mqttProtocol311 {
			if cleanSession == 0 {
				conack |= 0x01
			}
		}
		s.cleanSession = cleanSession

		if s.cleanSession == 0 && found.cleanSession == 0 {
			// Resume last session   // fix me ssddn
			s.storage.UpdateSession(s)
			// Notify other mqtt node to release resource
			base.AsyncProduceMessage(s.config,
				TopicNameSession,
				&SessionTopic{
					Launcher:  s.conn.LocalAddr().String(),
					SessionId: clientid,
					Action:    ObjectActionUpdate,
					State:     mqttStateDisconnecting,
				})
		}

	} else {
		// Register the session in storage
		s.storage.RegisterSession(s)
	}

	if willMsg != nil {
		s.willMsg = willMsg
		s.willMsg.topic = willTopic
		if len(payload) > 0 {
			s.willMsg.payload = payload
		} else {
			s.willMsg.payload = nil
		}
		s.willMsg.qos = willQos
		s.willMsg.retain = willRetain
	}
	s.clientID = clientid
	s.cleanSession = cleanSession
	s.pingTime = nil
	s.isDroping = false

	// Remove any queued messages that are no longer allowd through ACL
	// Assuming a possible change of username
	s.storage.DeleteMessageWithValidator(
		clientid,
		func(msg StorageMessage) bool {
			err := s.authapi.CheckAcl(context.Background(), clientid, username, willTopic, auth.AclActionRead)
			if err != nil {
				return false
			}
			return true
		})

	s.state = mqttStateConnected
	err = s.sendConnAck(uint8(conack), CONNACK_ACCEPTED)
	return err
}

// handleDisconnect handle disconnect packet
func (s *mqttSession) handleDisconnect() error {
	glog.Infof("Received DISCONNECT from %s", s.id)

	if s.inpacket.remainingLength != 0 {
		return mqttErrorInvalidProtocol
	}
	if s.protocol == mqttProtocol311 && (s.inpacket.command&0x0F) != 0x00 {
		s.disconnect()
		return mqttErrorInvalidProtocol
	}
	s.state = mqttStateDisconnecting
	s.disconnect()
	return nil
}

// disconnect will disconnect current connection because of protocol error
func (s *mqttSession) disconnect() {
	if s.state == mqttStateDisconnected {
		return
	}
	if s.cleanSession > 0 {
		s.storage.DeleteSession(s.id)
		s.id = ""
	}
	s.state = mqttStateDisconnected
	s.conn.Close()
	s.conn = nil
}

// handleSubscribe handle subscribe packet
func (s *mqttSession) handleSubscribe() error {
	payload := make([]uint8, 0)

	glog.Infof("Received SUBSCRIBE from %s", s.id)
	if s.protocol == mqttProtocol311 {
		if (s.inpacket.command & 0x0F) != 0x02 {
			return mqttErrorInvalidProtocol
		}
	}
	// Get message identifier
	mid, err := s.inpacket.readUint16()
	if err != nil {
		return err
	}
	// Deal each subscription
	for s.inpacket.pos < s.inpacket.remainingLength {
		sub := ""
		qos := uint8(0)
		if sub, err = s.inpacket.readString(); err != nil {
			return err
		}
		if checkTopicValidity(sub) != nil {
			glog.Errorf("Invalid subscription topic %s from %s, disconnecting", sub, s.id)
			return mqttErrorInvalidProtocol
		}
		if qos, err = s.inpacket.readByte(); err != nil {
			return err
		}

		if qos > 2 {
			glog.Errorf("Invalid Qos in subscription %s from %s", sub, s.id)
			return mqttErrorInvalidProtocol
		}

		if s.observer != nil {
			mp := s.observer.OnGetMountPoint()
			sub = mp + sub
		}
		if qos != 0x80 {
			if err := s.storage.AddSubscription(s.id, sub, qos); err != nil {
				return err
			}
			if err := s.storage.RetainSubscription(s.id, sub, qos); err != nil {
				return err
			}
		}
		payload = append(payload, qos)
	}

	if s.protocol == mqttProtocol311 && len(payload) == 0 {
		return mqttErrorInvalidProtocol
	}
	return s.sendSubAck(mid, payload)
}

// handleUnsubscribe handle unsubscribe packet
func (s *mqttSession) handleUnsubscribe() error {
	glog.Infof("Received UNSUBSCRIBE from %s", s.id)

	if s.protocol == mqttProtocol311 && (s.inpacket.command&0x0f) != 0x02 {
		return mqttErrorInvalidProtocol
	}
	mid, err := s.inpacket.readUint16()
	if err != nil {
		return err
	}
	// Iterate all subscription
	for s.inpacket.pos < s.inpacket.remainingLength {
		sub, err := s.inpacket.readString()
		if err != nil {
			return mqttErrorInvalidProtocol
		}
		if err := checkTopicValidity(sub); err != nil {
			return fmt.Errorf("Invalid unsubscription string from %s, disconnecting", s.id)
		}
		s.storage.RemoveSubscription(s.id, sub)
	}

	return s.sendCommandWithMid(UNSUBACK, mid, false)
}

// handlePublish handle publish packet
func (s *mqttSession) handlePublish() error {
	glog.Infof("Received PUBLISH from %s", s.id)

	var topic string
	var mid uint16
	var err error
	var payload []uint8


	dup := (s.inpacket.command & 0x08) >> 3
	qos := (s.inpacket.command & 0x06) >> 1
	if qos == 3 {
		return fmt.Errorf("Invalid Qos in PUBLISH from %s, disconnectiing", s.id)
	}
	retain := (s.inpacket.command & 0x01)

	// Topic
	if topic, err = s.inpacket.readString(); err != nil {
		return fmt.Errorf("Invalid topic in PUBLISH from %s", s.id)
	}
	if checkTopicValidity(topic) != nil {
		return fmt.Errorf("Invalid topic in PUBLISH(%s) from %s", topic, s.id)
	}
	if s.observer != nil && s.observer.OnGetMountPoint() != "" {
		topic = s.observer.OnGetMountPoint() + topic
	}

	if qos > 0 {
		mid, err = s.inpacket.readUint16()
		if err != nil {
			return err
		}
	}

	// Payload
	payloadlen := s.inpacket.remainingLength - s.inpacket.pos
	if payloadlen > 0 {
		limitSize, _ := s.config.Int(s.mgr.protocol, "message_size_limit")
		if payloadlen > limitSize {
			return mqttErrorInvalidProtocol
		}
		payload, err = s.inpacket.readBytes(payloadlen)
		if err != nil {
			return err
		}
	}
	// Check for topic access
	if s.observer != nil {
		err := s.authapi.CheckAcl(context.Background(), s.id, s.username, topic, auth.AclActionWrite)
		switch err {
		case auth.ErrorAclDenied:
			return mqttErrorInvalidProtocol
		default:
			return err
		}
	}
	glog.Infof("Received PUBLISH from %s(d:%d, q:%d r:%d, m:%d, '%s',..(%d)bytes",
		s.id, dup, qos, retain, mid, topic, payloadlen)

	// Check wether the message has been stored
	dup = 0
	var storedMsg *mqttMessage
	if qos > 0 {
		for _, storedMsg = range s.storedMsgs {
			if storedMsg.mid == mid && storedMsg.direction == mqttMessageDirectionIn {
				dup = 1
				break
			}
		}
	}
	msg := StorageMessage{
		ID:        uint(mid),
		SourceID:  s.id,
		Topic:     topic,
		Direction: MessageDirectionIn,
		State:     0,
		Qos:       qos,
		Retain:    (retain > 0),
		Payload:   payload,
	}

	switch qos {
	case 0:
		err = s.storage.QueueMessage(s.id, msg)
	case 1:
		err = s.storage.QueueMessage(s.id, msg)
		err = s.sendPubAck(mid)
	case 2:
		err = nil
		if dup > 0 {
			err = s.storage.InsertMessage(s.id, mid, MessageDirectionIn, msg)
		}
		if err == nil {
			err = s.sendPubRec(mid)
		}
	default:
		err = mqttErrorInvalidProtocol
	}

	return err
}

// handlePubRel handle pubrel packet
func (s *mqttSession) handlePubRel() error {
	// Check protocol specifal requriement
	if s.protocol == mqttProtocol311 {
		if (s.inpacket.command & 0x0F) != 0x02 {
			return mqttErrorInvalidProtocol
		}
	}
	// Get message identifier
	mid, err := s.inpacket.readUint16()
	if err != nil {
		return err
	}

	s.storage.DeleteMessage(s.id, mid, MessageDirectionIn)
	return s.sendPubComp(mid)
}

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
	glog.Infof("Sending PINGRESP to %s", s.id)
	return s.sendSimpleCommand(PINGRESP)
}

// sendConnAck send connection response to client
func (s *mqttSession) sendConnAck(ack uint8, result uint8) error {
	glog.Infof("Sending CONNACK from %s", s.id)

	packet := &mqttPacket{
		command:         CONNACK,
		remainingLength: 2,
	}
	packet.initializePacket()
	packet.payload[packet.pos+0] = ack
	packet.payload[packet.pos+1] = result

	return s.queuePacket(packet)
}

// sendSubAck send subscription acknowledge to client
func (s *mqttSession) sendSubAck(mid uint16, payload []uint8) error {
	glog.Infof("Sending SUBACK on %s", s.id)
	packet := &mqttPacket{
		command:         SUBACK,
		remainingLength: 2 + int(len(payload)),
	}

	packet.initializePacket()
	packet.writeUint16(mid)
	if len(payload) > 0 {
		packet.writeBytes(payload)
	}
	return s.queuePacket(packet)
}

// sendCommandWithMid send command with message identifier
func (s *mqttSession) sendCommandWithMid(command uint8, mid uint16, dup bool) error {
	packet := &mqttPacket{
		command:         command,
		remainingLength: 2,
	}
	if dup {
		packet.command |= 8
	}
	packet.initializePacket()
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
	p.toprocess = p.length
	s.sendPacketChannel <- p
	return nil
}

func (s *mqttSession) QueueMessage(msg *mqttMessage) error {
	s.sendMsgChannel <- msg
	return nil
}

func (s *mqttSession) updateOutMessage(mid uint16, state int) error {
	return nil
}

func (s *mqttSession) generateMid() uint16 {
	return 0
}

func (s *mqttSession) sendPublish(subQos uint8, srcQos uint8, topic string, payload []uint8) error {
	/* Check for ACL topic access. */
	// TODO

	var qos uint8
	if option, err := s.config.Bool(s.mgr.protocol, "upgrade_outgoing_qos"); err != nil && option {
		qos = subQos
	} else {
		if srcQos > subQos {
			qos = subQos
		} else {
			qos = srcQos
		}
	}

	var mid uint16
	if qos > 0 {
		mid = s.generateMid()
	} else {
		mid = 0
	}

	packet := newMqttPacket()
	packet.command = PUBLISH
	packet.remainingLength = 2 + len(topic) + len(payload)
	if qos > 0 {
		packet.remainingLength += 2
	}
	packet.initializePacket()
	packet.writeString(topic)
	if qos > 0 {
		packet.writeUint16(mid)
	}
	packet.writeBytes(payload)

	return s.queuePacket(&packet)
}
