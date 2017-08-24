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
	"github.com/satori/go.uuid"
)

const (
	maxMqttConnections = 1000000
	protocolName       = "mqtt3"
)

type mqtt struct {
	config   *base.Config
	chn      chan int
	index    int64
	sessions map[string]base.Session
	mutex    sync.Mutex // Maybe not so good
	inpacket *mqttPacket
}

// MqttFactory
type mqttFactory struct{}

// New create mqtt service factory
func (m *mqttFactory) New(c *base.Config, ch chan int) (base.Service, error) {
	t := &mqtt{config: c,
		chn:      ch,
		index:    -1,
		sessions: make(map[string]base.Session),
	}
	return t, nil
}

// MQTT Service

func (m *mqtt) NewSession(conn net.Conn) (base.Session, error) {
	id := m.CreateSessionId()
	session := newMqttSession(m, conn, id)
	return session, nil
}

// CreateSessionId create id for new session
func (m *mqtt) CreateSessionId() string {
	return uuid.NewV4().String()
}

// GetSessionTotalCount get total session count
func (m *mqtt) GetSessionTotalCount() int64 {
	return int64(len(m.sessions))
}

func (m *mqtt) RemoveSession(s base.Session) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.sessions[s.Identifier()] = nil
}
func (m *mqtt) RegisterSession(s base.Session) {
	m.mutex.Lock()
	m.sessions[s.Identifier()] = s
	m.mutex.Unlock()
}

// Run is mainloop for mqtt service
// TODO: Run is very common for each service, it should be moved to ServiceManager
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
		session, err := m.NewSession(conn)
		if err != nil {
			glog.Error("Mqtt create session failed")
			return err
		}
		glog.Infof("Mqtt new connection:%s", session.Identifier())
		m.RegisterSession(session)
		go session.Handle()
	}
	// notify main
	m.chn <- 1
	return nil
}

func init() {
	glog.Info("Registering service:%s", protocolName)
	base.RegisterServiceFactory(protocolName, &mqttFactory{})
}
