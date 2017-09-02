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
	"fmt"
	"iothub/base"
	"iothub/database"
	"net"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/golang/glog"
	"github.com/satori/go.uuid"
)

const (
	protocolName = "coap"
)

type coap struct {
	config     base.Config
	chn        chan int
	index      int64
	sessions   map[string]base.Session
	mutex      sync.Mutex // Maybe not so good
	protocol   uint8
	wg         sync.WaitGroup
	localAddrs []string
	db         database.Database
}

// CoapFactory
type CoapFactory struct{}

// New create coap service factory
func (m *CoapFactory) New(c base.Config, ch chan int) (base.Service, error) {
	var localAddrs []string = []string{}
	var db database.Database

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
	// Create database
	name := c.MustString("database", "name")
	if db, err = database.New(name, database.Option{}); err != nil {
		return nil, errors.New("Failed to create database in coap")
	}

	t := &coap{config: c,
		chn:        ch,
		index:      -1,
		sessions:   make(map[string]base.Session),
		protocol:   2,
		localAddrs: localAddrs,
		db:         db,
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

// Run is mainloop for coap service
// TODO: Run is very common for each service, it should be moved to ServiceManager
func (m *coap) Run() error {
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

// launchMqttMonitor
func (m *coap) launchCoapMonitor() error {
	glog.Info("Luanching coap monitor...")
	//sarama.Logger = glog
	khosts, _ := m.config.String("iothub", "kafka-hosts")
	consumer, err := sarama.NewConsumer(strings.Split(khosts, ","), nil)
	if err != nil {
		return fmt.Errorf("Connecting with kafka:%s failed", khosts)
	}

	partitionList, err := consumer.Partitions("iothub-coap")
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
func (m *coap) handleNotifications(topic string, value []byte) error {
	return nil
}

// handleSessionNotifications handle session notification  from kafka
func (m *coap) handleSessionNotifications(value []byte) error {
	return nil
}
