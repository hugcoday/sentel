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
	"iothub/base"
	"net"
	"sync"

	"github.com/golang/glog"
)

const (
	maxMqttConnections = 1000000
	protocolName       = "mqtt3"
)

type mqtt struct {
	config   *base.Config
	chn      chan int
	index    int64
	sessions map[int64]*mqttSession
	mutex    sync.Mutex // Maybe not so good
}

// MqttFactory
type mqttFactory struct{}

// New create mqtt service factory
func (m *mqttFactory) New(c *base.Config, ch chan int) (base.Service, error) {
	t := &mqtt{config: c,
		chn:      ch,
		index:    -1,
		sessions: make(map[int64]*mqttSession),
	}
	return t, nil
}

// getSessionIndex create a index for new session
func (m *mqtt) getSessionIndex() int64 {
	m.mutex.Lock() // Not so good!
	defer m.mutex.Unlock()
	m.index++
	if m.index < maxMqttConnections && m.sessions[m.index] == nil {
		return m.index
	}
	if m.index == maxMqttConnections {
		m.index = 0
	}
	for {
		if m.sessions[m.index] == nil {
			return m.index
		}
		m.index++
	}
	return -1
}

func (m *mqtt) removeSession(s *mqttSession) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.sessions[s.id] = nil
}
func (m *mqtt) addSession(index int64, s *mqttSession) {
	m.mutex.Lock()
	m.sessions[index] = s
	m.mutex.Unlock()
}

func (m *mqtt) Run() error {
	listen, err := net.Listen("tcp", m.config.Mqtt.Host)
	if err != nil {
		glog.Errorf("Mqtt listen failed:%s", err)
		return err
	}
	glog.Info("Mqtt server is listening...")
	for {
		conn, err := listen.Accept()
		if err != nil {
			continue
		}
		index := m.getSessionIndex()
		if index < 0 {
			glog.Errorf("Mqtt max connection exceeded, drop connection occured")
			continue
		}
		glog.Infof("Mqtt new connection:%s", index)
		session := newMqttSession(m, conn, index)
		m.addSession(index, session)
		go session.handleConnection()
	}
	// notify main
	m.chn <- 1
	return nil
}

func init() {
	glog.Info("Registering service:%s", protocolName)
	base.RegisterServiceFactory(protocolName, &mqttFactory{})
}
