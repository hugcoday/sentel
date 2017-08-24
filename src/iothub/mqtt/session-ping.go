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
