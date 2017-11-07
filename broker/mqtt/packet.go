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
	"io"

	"github.com/cloustone/sentel/broker/base"
)

const (
	// Protocol version
	PROTOCOL_NAME_V31    = "MQIsdp"
	PROTOCOL_VERSION_V31 = 3

	PROTOCOL_NAME_V311    = "MQTT"
	PROTOCOL_VERSION_V311 = 4

	// Message types
	INVALID     = 0x00
	CONNECT     = 0x10
	CONNACK     = 0x20
	PUBLISH     = 0x30
	PUBACK      = 0x40
	PUBREC      = 0x50
	PUBREL      = 0x60
	PUBCOMP     = 0x70
	SUBSCRIBE   = 0x80
	SUBACK      = 0x90
	UNSUBSCRIBE = 0xA0
	UNSUBACK    = 0xB0
	PINGREQ     = 0xC0
	PINGRESP    = 0xD0
	DISCONNECT  = 0xE0

	// CONNACK result
	CONNACK_ACCEPTED                      = 0
	CONNACK_REFUSED_PROTOCOL_VERSION      = 1
	CONNACK_REFUSED_IDENTIFIER_REJECTED   = 2
	CONNACK_REFUSED_SERVER_UNAVAILABLE    = 3
	CONNACK_REFUSED_BAD_USERNAME_PASSWORD = 4
	CONNACK_REFUSED_NOT_AUTHORIZED        = 5

	MQTT_MAX_PAYLOAD = 268435455
)

type mqttPacket struct {
	command         uint8
	dup             bool
	qos             uint8
	retain          bool
	mid             uint16
	pos             int
	toprocess       int
	length          int
	remainingCount  int
	remainingMult   int
	remainingLength int
	payload         []uint8
	buf             []uint8
}

func newMqttPacket() mqttPacket {
	return mqttPacket{
		command:        0,
		dup:            false,
		qos:            0,
		retain:         false,
		pos:            0,
		length:         0,
		remainingCount: 0,
		remainingMult:  1,
		payload:        []uint8{},
		buf:            []uint8{},
	}
}

func (p *mqttPacket) PacketType() string {
	return PROTOCOL_NAME_V311
}
func (p *mqttPacket) Clear() {
	p.command = 0
	p.dup = false
	p.qos = 0
	p.retain = false
	p.pos = 0
	p.length = 0
	p.toprocess = 0
	p.remainingCount = 0
	p.remainingLength = 0
	p.remainingMult = 1
	p.payload = []uint8{}
	p.buf = []uint8{}
}

func (p *mqttPacket) initializePacket() error {
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
		return errors.New("Invalid packet payload size")
	}
	p.length = p.remainingLength + 1 + p.remainingCount
	p.payload = make([]uint8, p.length)

	// assemble fixed header
	p.payload[0] = p.command
	if p.dup {
		p.payload[0] |= 8
	}
	if p.qos > 0 {
		p.payload[0] |= p.qos << 1
	}
	if p.retain {
		p.payload[0] |= 1
	}
	for i := 0; i < p.remainingCount; i++ {
		p.payload[i+1] = remainingBytes[i]
	}
	p.pos = 1 + p.remainingCount
	return nil
}

// SerializeTo writes the serialized form of the packet into the serialize buffer
func (p *mqttPacket) SerializeTo(buf base.SerializeBuffer, opts base.SerializeOptions) error {
	return nil
}

// Write implement Writer interface
func (p *mqttPacket) Write(buf []byte) (int, error) {
	p.buf = buf
	return len(p.buf), nil
}

// DecodeFromBytes decode given bytes into this protocol layer
func (p *mqttPacket) DecodeFromBytes(r io.Reader, df base.DecodeFeedback) (int, error) {
	return 0, errors.New("Not Implementented")
}

// DecodeFromReader decode packet from given reader
func (p *mqttPacket) DecodeFromReader(r io.Reader, df base.DecodeFeedback) error {
	if p.command == 0 {
		// Read command
		if length, err := io.CopyN(p, r, 1); err != nil || length != 1 {
			return mqttErrorInvalidProtocol
		}
		p.command = p.buf[0]
	}
	// Compute remaining length
	if p.remainingCount <= 0 {
		for {
			// Read length byte
			if length, err := io.CopyN(p, r, 1); err != nil || length != 1 {
				return mqttErrorInvalidProtocol
			}
			p.remainingCount--
			if p.remainingCount < -4 {
				return mqttErrorInvalidProtocol
			}
			p.remainingLength += int(p.buf[0]&127) * p.remainingMult
			p.remainingMult *= 128
			if p.buf[0]&128 == 0 {
				break
			}
		}
	}
	// We have finished reading remaining length
	p.remainingCount *= -1
	if p.remainingLength > 0 {
		p.toprocess = p.remainingLength
		p.payload = []uint8{}
	}

	if p.toprocess > 0 {
		length, err := io.CopyN(p, r, int64(p.toprocess))
		if err != nil || length != int64(p.toprocess) {
			p.Clear()
			return mqttErrorInvalidProtocol
		}
		p.payload = p.buf
	}
	p.pos = 0
	return nil
}

// Length return mqtt packet length
func (p *mqttPacket) Length() int {
	return p.length
}

// ReadByte read a byte from packet payload
func (p *mqttPacket) readByte() (uint8, error) {
	if p.pos+1 > p.remainingLength {
		return 0, mqttErrorInvalidProtocol
	}
	b := p.payload[p.pos]
	p.pos++
	return b, nil
}

// WriteByte  write a byte into packet payload
func (p *mqttPacket) writeByte(b uint8) error {
	if p.pos+1 > p.length {
		return mqttErrorInvalidProtocol
	}
	p.payload[p.pos] = b
	p.pos++
	return nil
}

// ReadBytes read bytes from packet payload
func (p *mqttPacket) readBytes(count int) ([]uint8, error) {
	if p.pos+count > p.remainingLength {
		return nil, mqttErrorInvalidProtocol
	}
	r := p.payload[p.pos : p.pos+count]
	p.pos += count
	return r, nil
}

// WriteBytes write bytes into packet payload
func (p *mqttPacket) writeBytes(buf []uint8) error {
	if p.pos+len(buf) > p.length {
		return mqttErrorInvalidProtocol
	}

	for _, b := range buf {
		p.payload[p.pos] = b
		p.pos++
	}
	return nil
}

// ReadString read string from packet payload
func (p *mqttPacket) readString() (string, error) {
	len, err := p.readUint16()
	if err != nil {
		return "", err
	}
	if p.pos+int(len) > p.remainingLength {
		return "", mqttErrorInvalidProtocol
	}

	s := string(p.payload[p.pos : p.pos+int(len)])
	p.pos += int(len)
	return s, nil
}

// WriteString write string into packet payload
func (p *mqttPacket) writeString(data string) error {
	length := uint16(len(data))
	if err := p.writeUint16(length); err != nil {
		return err
	}
	if err := p.writeBytes([]uint8(data)); err != nil {
		return err
	}
	return nil
}

// ReadUint16 read word from packet payload
func (p *mqttPacket) readUint16() (uint16, error) {
	if p.pos+2 > p.remainingLength {
		return 0, mqttErrorInvalidProtocol
	}
	msb := p.payload[p.pos]
	p.pos++
	lsb := p.payload[p.pos]
	p.pos++
	w := (uint16(msb << 8)) + uint16(lsb)
	return w, nil
}

// WriteUint16 write word into packet pyload
func (p *mqttPacket) writeUint16(data uint16) error {
	msb := uint8((data >> 8) & 0xF)
	lsb := uint8(data & 0xF)
	if err := p.writeByte(msb); err != nil {
		return err
	}
	if err := p.writeByte(lsb); err != nil {
		return err
	}
	return nil
}
