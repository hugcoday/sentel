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

// Message state
const (
	mqttMessageStateInvalid        = 0
	mqttMessageStatePublishQos0    = 1
	mqttMessageStatePublishQos1    = 2
	mqttMessageStateWaitForPubAck  = 3
	mqttMessageStatePublishQos2    = 4
	mqttMessageStateWaitForPubRec  = 5
	mqttMessageStateResendPubRel   = 6
	mqttMessageStateWaitForPubRel  = 7
	mqttMessageStateResendPubComp  = 8
	mqttMessageStateWaitForPubComp = 9
	mqttMessateStateSendPubRec     = 10
	mqttMessateStateQueued         = 11
)

// Message direction
const (
	mqttMessageDirectionIn  = 0
	mqttMessageDirectionOut = 1
)

type mqttMessage struct {
	mid       uint16
	direction int
	topic     string
	payload   []uint8
	qos       uint8
	retain    bool
}
