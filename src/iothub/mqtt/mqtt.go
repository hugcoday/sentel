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
	"iothub/base"
	"iothub/storage"
	"libs"
	"net"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/golang/glog"
	"github.com/satori/go.uuid"
)

const (
	maxMqttConnections = 1000000
	protocolName       = "mqtt3"
)

type mqtt struct {
	config     libs.Config
	chn        chan base.ServiceCommand
	index      int64
	sessions   map[string]base.Session
	mutex      sync.Mutex // Maybe not so good
	inpacket   *mqttPacket
	wg         sync.WaitGroup
	localAddrs []string
	storage    storage.Storage
	protocol   string
}

// MqttFactory
type MqttFactory struct{}

// New create mqtt service factory
func (m *MqttFactory) New(protocol string, c libs.Config, ch chan base.ServiceCommand) (base.Service, error) {
	var localAddrs []string = []string{}
	var s storage.Storage

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
	if s, err = storage.New(name, storage.Option{}); err != nil {
		return nil, errors.New("Failed to create storage in mqtt")
	}

	t := &mqtt{config: c,
		chn:        ch,
		index:      -1,
		sessions:   make(map[string]base.Session),
		protocol:   protocol,
		localAddrs: localAddrs,
		storage:    s,
	}
	return t, nil
}

// MQTT Service

func (m *mqtt) NewSession(conn net.Conn) (base.Session, error) {
	id := m.CreateSessionId()
	s, err := newMqttSession(m, conn, id)
	return s, err
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
		m.RegisterSession(session)
		go func(s base.Session) {
			err := s.Handle()
			if err != nil {
				glog.Error(err)
			}
		}(session)
	}
	// notify main
	m.chn <- 1
	return nil
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
					m.storage.UpdateSession(nil,
						&storage.Session{Id: topic.SessionId, State: topic.State})
				}
			}
		case ObjectActionDelete:

		case ObjectActionRegister:
		default:
		}
	}
	return nil
}
