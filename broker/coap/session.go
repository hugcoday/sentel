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

package coap

import (
	"net"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/core"

	"github.com/golang/glog"
)

type coapSession struct {
	mgr      *coap
	config   core.Config
	conn     net.Conn
	id       string
	state    uint8
	observer base.SessionObserver
}

// newCoapSession create new session  for each client connection
func newCoapSession(m *coap, conn net.Conn, id string) (*coapSession, error) {
	s := &coapSession{
		mgr:      m,
		config:   m.config,
		conn:     conn,
		id:       id,
		observer: nil,
	}
	return s, nil
}

func (s *coapSession) RegisterObserver(o base.SessionObserver) {
	if s.observer != nil {
		glog.Error("MqttSession register multiple observer")
	}
	s.observer = o
}

// handle is mainprocessor for iot device client
// Loop to read packet from conn
func (s *coapSession) Handle() error {
	return nil
}

// removeConnection remove current connection from coapManaager if errors occured
func (s *coapSession) Destroy() error            { return nil }
func (s *coapSession) Identifier() string        { return "" }
func (s *coapSession) Service() base.Service     { return nil }
func (s *coapSession) GetStats() *base.Stats     { return nil }
func (s *coapSession) GetMetrics() *base.Metrics { return nil }
func (s *coapSession) Info() *base.SessionInfo   { return nil }
