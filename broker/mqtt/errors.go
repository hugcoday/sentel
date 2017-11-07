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

import "errors"

var (
	mqttErrorInvalidProtocol = errors.New("Invalid protocol")
	mqttErrorInvalidVersion  = errors.New("Invalid protocol version")
	mqttErrorConnectPending  = errors.New("Connec pending")
	mqttErrorNoConnection    = errors.New("No connection")
	mqttErrorConnectRefused  = errors.New("Connection Refused")
	mqttErrorNotFound        = errors.New("Not found")
	mqttErrorNotSupported    = errors.New("Not supported")
	mqttErrorAutoFailed      = errors.New("Auth failed")
	mqttErrorUnkown          = errors.New("Unknown error")
)
