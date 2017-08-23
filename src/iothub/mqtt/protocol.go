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

const (
	// Protocol version
	PROTOCOL_NAME_V31    = "MQIsdp"
	PROTOCOL_VERSION_V31 = 3

	PROTOCOL_NAME_V311    = "MQTT"
	PROTOCOL_VERSION_V311 = 4

	// Message types
	CONNECT     = 0x10
	CNNNACK     = 0x20
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
	command        uint8
	remainingCount uint16
	mid            uint16
	pos            uint32
	toprocess      uint32
	length         uint32
	remainingMult  uint32
	payload        []byte
}
