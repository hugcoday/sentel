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

package base

type Packet interface {
	// PacketType return type name of packet
	PacketType() string
	// DecodeFromBytes decode given bytes into this protocol layer
	DecodeFromBytes(data []byte, df DecodeFeedback) (int, error)
	// SerializeTo writes the serialized form of the packet into the serialize buffer
	SerializeTo(buf SerializeBuffer, opts SerializeOptions) error
	// Clear clear packet content and payload
	Clear()
}
