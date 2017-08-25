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
				willTopic = s.observer.GetMountpoint(s) + topic
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
			default:


			}

		    }

		}
	*/
	return nil
}

// sendConnAck send CONNACK packet to client
func (s *mqttSession) sendConnAck(b uint8, reason uint8) {
}

// handleDisconnect handle disconnect packet
func (s *mqttSession) handleDisconnect() error {
	return nil
}

// handleConAck handle conack packet
func (s *mqttSession) handleConnAck() error {
	return nil
}
