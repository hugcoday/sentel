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

type SerializeOptions struct {
	// Fixlengths determines whether layers should fix the values for any
	// length field that depends on the payload during serialization
	FixLengths bool
	// ComputeChecksums determines wethers layers should recompute checksums
	// based on their payloads
	ComputeChecksums bool
}

type SerializeBuffer interface {
	// Bytes return the continguous set of bytes
	Bytes() []byte
	// PrependBytes return a set of bytes which prepends the curent bytes
	// in this buffer
	PrependBytes(num int) ([]byte, error)
	// AppendBytes return a set of bytes which append the current bytes
	// in this buffer
	AppendBytes(num int) ([]byte, error)
	// Clear reset the SerializeBuffer to a new
	Clear() error
}
