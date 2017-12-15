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
	"errors"
	"net"
	"sync"

	"github.com/cloustone/sentel/iothub/base"
	"github.com/cloustone/sentel/libs/sentel"

	uuid "src/github.com/satori/go.uuid"

	"github.com/golang/glog"
)

const (
	protocolName = "coap"
)

type coap struct {
	config     sentel.Config
	chn        chan base.ServiceCommand
	index      int64
	sessions   map[string]base.Session
	mutex      sync.Mutex // Maybe not so good
	protocol   uint8
	wg         sync.WaitGroup
	localAddrs []string
}

// CoapFactory
type CoapFactory struct{}

// New create coap service factory
func (m *CoapFactory) New(protocol string, c sentel.Config, ch chan base.ServiceCommand) (base.Service, error) {
	var localAddrs []string = []string{}
	// Get all local ip address
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		glog.Errorf("Failed to get local address:%s", err)
		return nil, err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
			localAddrs = append(localAddrs, ipnet.IP.String())
		}
	}
	if len(localAddrs) == 0 {
		return nil, errors.New("Failed to get local address")
	}
	t := &coap{config: c,
		chn:        ch,
		index:      -1,
		sessions:   make(map[string]base.Session),
		protocol:   2,
		localAddrs: localAddrs,
	}
	return t, nil
}

// MQTT Service

func (m *coap) NewSession(conn net.Conn) (base.Session, error) {
	id := m.CreateSessionId()
	s, err := newCoapSession(m, conn, id)
	return s, err
}

// CreateSessionId create id for new session
func (m *coap) CreateSessionId() string {
	return uuid.NewV4().String()
}

// GetSessionTotalCount get total session count
func (m *coap) GetSessionTotalCount() int64 {
	return int64(len(m.sessions))
}

func (m *coap) RemoveSession(s base.Session) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.sessions[s.Identifier()] = nil
}
func (m *coap) RegisterSession(s base.Session) {
	m.mutex.Lock()
	m.sessions[s.Identifier()] = s
	m.mutex.Unlock()
}

// Start
func (m *coap) Start() error {
	host, _ := m.config.String("coap", "host")

	listen, err := net.Listen("tcp", host)
	if err != nil {
		glog.Errorf("Coap listen failed:%s", err)
		return err
	}
	glog.Infof("Coap server is listening on '%s'...", host)
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

func (m *coap) Stop() {}

// Name
func (m *coap) Info() *base.ServiceInfo {
	return &base.ServiceInfo{
		ServiceName: "coap",
	}
}

func (m *coap) GetMetrics() *base.Metrics            { return nil }
func (m *coap) GetStats() *base.Stats                { return nil }
func (m *coap) GetClients() []*base.ClientInfo       { return nil }
func (m *coap) GetClient(id string) *base.ClientInfo { return nil }
func (m *coap) KickoffClient(id string) error        { return nil }

func (m *coap) GetSessions(conditions map[string]bool) []*base.SessionInfo { return nil }
func (m *coap) GetSession(id string) *base.SessionInfo                     { return nil }

func (m *coap) GetRoutes() []*base.RouteInfo { return nil }
func (m *coap) GetRoute() *base.RouteInfo    { return nil }

// Topic info
func (m *coap) GetTopics() []*base.TopicInfo       { return nil }
func (m *coap) GetTopic(id string) *base.TopicInfo { return nil }

// SubscriptionInfo
func (m *coap) GetSubscriptions() []*base.SubscriptionInfo       { return nil }
func (m *coap) GetSubscription(id string) *base.SubscriptionInfo { return nil }

// Service Info
func (m *coap) GetServiceInfo() *base.ServiceInfo { return nil }
