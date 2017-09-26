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
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/cloustone/sentel/iothub/base"
	"github.com/cloustone/sentel/libs"
	uuid "github.com/satori/go.uuid"

	"github.com/Shopify/sarama"
	"github.com/golang/glog"
)

const (
	maxMqttConnections = 1000000
	protocolName       = "mqtt3"
)

// MQTT service declaration
type mqtt struct {
	config     libs.Config
	chn        chan base.ServiceCommand
	index      int64
	sessions   map[string]base.Session
	mutex      sync.Mutex // Maybe not so good
	inpacket   *mqttPacket
	wg         sync.WaitGroup
	localAddrs []string
	storage    Storage
	protocol   string
	stats      *base.Stats
	metrics    *base.Metrics
}

// MqttFactory
type MqttFactory struct{}

// New create mqtt service factory
func (m *MqttFactory) New(protocol string, c libs.Config, ch chan base.ServiceCommand) (base.Service, error) {
	var localAddrs []string = []string{}
	var s Storage

	// Get all local ip address
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		glog.Errorf("Failed to get local interface:%s", err)
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
	// Create storage
	name := c.MustString("storage", "name")
	if s, err = NewStorage(name, c); err != nil {
		return nil, errors.New("Failed to create storage in mqtt")
	}

	t := &mqtt{config: c,
		chn:        ch,
		index:      -1,
		sessions:   make(map[string]base.Session),
		protocol:   protocol,
		localAddrs: localAddrs,
		storage:    s,
		stats:      base.NewStats(true),
		metrics:    base.NewMetrics(true),
	}
	return t, nil
}

// MQTT Service

// Name
func (m *mqtt) Name() string { return "matt" }

func (m *mqtt) NewSession(conn net.Conn) (base.Session, error) {
	id := m.createSessionId()
	s, err := newMqttSession(m, conn, id)
	return s, err
}

// CreateSessionId create id for new session
func (m *mqtt) createSessionId() string {
	return uuid.NewV4().String()
}

// GetSessionTotalCount get total session count
func (m *mqtt) getSessionTotalCount() int64 {
	return int64(len(m.sessions))
}

func (m *mqtt) removeSession(s base.Session) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.sessions[s.Identifier()] = nil
}
func (m *mqtt) registerSession(s base.Session) {
	m.mutex.Lock()
	m.sessions[s.Identifier()] = s
	m.mutex.Unlock()
}

// Info
func (m *mqtt) Info() *base.ServiceInfo {
	return &base.ServiceInfo{
		ServiceName: "mqtt",
	}
}

// Stats and Metrics
func (m *mqtt) GetStats() *base.Stats     { return m.stats }
func (m *mqtt) GetMetrics() *base.Metrics { return m.metrics }

// Client
func (m *mqtt) GetClients() []*base.ClientInfo       { return nil }
func (m *mqtt) GetClient(id string) *base.ClientInfo { return nil }
func (m *mqtt) KickoffClient(id string) error        { return nil }

// Session Info
func (m *mqtt) GetSessions(conditions map[string]bool) []*base.SessionInfo { return nil }
func (m *mqtt) GetSession(id string) *base.SessionInfo                     { return nil }

// Route Info
func (m *mqtt) GetRoutes() []*base.RouteInfo { return nil }
func (m *mqtt) GetRoute() *base.RouteInfo    { return nil }

// Topic info
func (m *mqtt) GetTopics() []*base.TopicInfo       { return nil }
func (m *mqtt) GetTopic(id string) *base.TopicInfo { return nil }

// SubscriptionInfo
func (m *mqtt) GetSubscriptions() []*base.SubscriptionInfo       { return nil }
func (m *mqtt) GetSubscription(id string) *base.SubscriptionInfo { return nil }

// Service Info
func (m *mqtt) GetServiceInfo() *base.ServiceInfo { return nil }

// Start is mainloop for mqtt service
func (m *mqtt) Start() error {
	host, _ := m.config.String("mqtt", "host")

	listen, err := net.Listen("tcp", host)
	if err != nil {
		glog.Errorf("Mqtt listen failed:%s", err)
		return err
	}
	// Launch montor
	// TODO:how to wait the monitor to be terminated
	if err := m.launchMqttMonitor(); err != nil {
		glog.Errorf("Mqtt monitor failed, reason:%s", err)
		//return err
	}

	glog.Infof("Mqtt server is listening on '%s'...", host)
	for {
		conn, err := listen.Accept()
		if err != nil {
			continue
		}
		session, err := m.NewSession(conn)
		if err != nil {
			glog.Errorf("Mqtt create session failed:%s", err)
			return err
		}
		m.registerSession(session)
		go func(s base.Session) {
			err := s.Handle()
			if err != nil {
				conn.Close()
				glog.Error(err)
			}
		}(session)
	}
	// notify main
	m.chn <- 1
	return nil
}

// Stop
func (m *mqtt) Stop() {
}

// launchMqttMonitor
func (m *mqtt) launchMqttMonitor() error {
	glog.Info("Luanching mqtt monitor...")
	//sarama.Logger = glog
	khosts, _ := m.config.String("iothub", "kafka-hosts")
	consumer, err := sarama.NewConsumer(strings.Split(khosts, ","), nil)
	if err != nil {
		return fmt.Errorf("Connecting with kafka:%s failed", khosts)
	}

	partitionList, err := consumer.Partitions("iothub-mqtt")
	if err != nil {
		return fmt.Errorf("Failed to get list of partions:%v", err)
		return err
	}

	for partition := range partitionList {
		pc, err := consumer.ConsumePartition("iothub", int32(partition), sarama.OffsetNewest)
		if err != nil {
			glog.Errorf("Failed  to start consumer for partion %d:%s", partition, err)
			continue
		}
		defer pc.AsyncClose()
		m.wg.Add(1)

		go func(sarama.PartitionConsumer) {
			defer m.wg.Done()
			for msg := range pc.Messages() {
				m.handleNotifications(string(msg.Topic), msg.Value)
			}
		}(pc)
	}
	m.wg.Wait()
	consumer.Close()
	return nil
}

// handleNotifications handle notification from kafka
func (m *mqtt) handleNotifications(topic string, value []byte) error {
	switch topic {
	case TopicNameSession:
		return m.handleSessionNotifications(value)
	}
	return nil
}

// handleSessionNotifications handle session notification  from kafka
func (m *mqtt) handleSessionNotifications(value []byte) error {
	// Decode value received form other mqtt node
	var topics []SessionTopic
	if err := json.Unmarshal(value, &topics); err != nil {
		glog.Errorf("Mqtt session notifications failure:%s", err)
		return err
	}
	// Get local ip address
	for _, topic := range topics {
		switch topic.Action {
		case ObjectActionUpdate:
			// Only deal with notification that is not  launched by myself
			for _, addr := range m.localAddrs {
				if addr != topic.Launcher {
					s, err := m.storage.FindSession(topic.SessionId)
					if err != nil {
						s.state = topic.State
					}
					//m.storage.UpdateSession(&StorageSession{Id: topic.SessionId, State: topic.State})
				}
			}
		case ObjectActionDelete:

		case ObjectActionRegister:
		default:
		}
	}
	return nil
}
