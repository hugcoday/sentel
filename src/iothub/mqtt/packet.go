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
	"iothub/base"
)

const (
	// Protocol version
	PROTOCOL_NAME_V31    = "MQIsdp"
	PROTOCOL_VERSION_V31 = 3

	PROTOCOL_NAME_V311    = "MQTT"
	PROTOCOL_VERSION_V311 = 4

	// Message types
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
	mid             uint16
	pos             uint32
	toprocess       uint32
	length          uint32
	remainingCount  int32
	remainingMult   int32
	remainingLength uint32
	payload         []uint8
}

func newMqttPacket() mqttPacket {
	return mqttPacket{
		command:        0,
		pos:            0,
		length:         0,
		remainingCount: 0,
		payload:        []byte{},
	}

}

func (p *mqttPacket) PacketType() string {
	return PROTOCOL_NAME_V311
}
func (p *mqttPacket) Clear() {
	p.command = 0
	p.length = 0
}

// SerializeTo writes the serialized form of the packet into the serialize buffer
func (p *mqttPacket) SerializeTo(buf base.SerializeBuffer, opts base.SerializeOptions) error {
	return nil
}

// DecodeFromBytes decode given bytes into this protocol layer
// TODO: underlay's read method  should be payed attention, temporal implementations
func (p *mqttPacket) DecodeFromBytes(data []byte, df base.DecodeFeedback) (int, error) {
	if len(data) == 0 {
		return 0, errors.New("Invalid data packet to decode")
	}
	// Start from new packet
	p.command = data[0]
	// Compute remaining length
	var index int = 1
	if p.remainingCount <= 0 {
		for _, b := range data[1:] {
			index++
			p.remainingCount--
			if p.remainingCount == -4 {
				return 0, errors.New("Invalid protocol")
			}
			p.remainingCount += int32(b&127) * p.remainingMult
			p.remainingMult *= 128
			if b&128 != 0 {
				break
			}
		}
	}
	// We have finished reading remaining length
	p.remainingCount *= -1
	// Check wether remaining data is validity
	if int32(len(data[index:])) < p.remainingCount {
		p.Clear()
		return 0, errors.New("Packet payload is too shore")
	}
	for _, b := range data[index:] {
		p.payload = append(p.payload, b)
	}
	return index + len(p.payload), nil
}

// Length return mqtt packet length
func (p *mqttPacket) Length() uint32 {
	return p.length
}

// ReadByte read a byte from packet payload
func (p *mqttPacket) ReadByte() (uint8, error) {
	if p.pos+1 > p.remainingLength {
		return 0, errors.New("Invalid mqtt packet")
	}
	b := p.payload[p.pos]
	p.pos++
	return b, nil
}

// WriteByte  write a byte into packet payload
func (p *mqttPacket) WriteByte(b uint8) error {
	if p.pos+1 > p.length {
		return errors.New("Invalid mqtt packet")
	}
	p.payload[p.pos] = b
	p.pos++
	return nil
}

// ReadBytes read bytes from packet payload
func (p *mqttPacket) ReadBytes(count uint32) ([]uint8, error) {
	if p.pos+count > p.length {
		return nil, errors.New("Invalid mqtt packet")
	}
	p.pos += count
	return p.payload[p.pos : p.pos+count], nil
}

// WriteBytes write bytes into packet payload
func (p *mqttPacket) WriteBytes(buf []uint8) error {
	if p.pos+uint32(len(buf)) > p.length {
		return errors.New("Invalid mqtt packet")
	}
	for _, b := range buf {
		p.payload = append(p.payload[p.pos:], b)
		p.pos++
	}
	return nil
}

// ReadString read string from packet payload
func (p *mqttPacket) ReadString() (string, error) {
	len, err := p.ReadUint16()
	if err != nil {
		return "", err
	}
	if p.pos+uint32(len) > p.remainingLength {
		return "", errors.New("Invalid mqtt packet")
	}

	s := string(p.payload[p.pos : p.pos+uint32(len)])
	return s, nil
}

// WriteString write string into packet payload
func (p *mqttPacket) WriteString(data string) error {
	length := uint16(len(data))
	if err := p.WriteUint16(length); err != nil {
		return err
	}
	if err := p.WriteBytes([]uint8(data)); err != nil {
		return err
	}
	return nil
}

// ReadUint16 read word from packet payload
func (p *mqttPacket) ReadUint16() (uint16, error) {
	if p.pos+2 > p.remainingLength {
		return 0, errors.New("Invalid mqtt packet")
	}
	msb := p.payload[p.pos]
	p.pos++
	lsb := p.payload[p.pos]
	p.pos++
	w := (uint16(msb << 8)) + uint16(lsb)
	return w, nil
}

// WriteUint16 write word into packet pyload
func (p *mqttPacket) WriteUint16(data uint16) error {
	msb := uint8((data >> 8) & 0xF)
	lsb := uint8(data & 0xF)
	if err := p.WriteByte(msb); err != nil {
		return err
	}
	if err := p.WriteByte(lsb); err != nil {
		return err
	}
	return nil
}
