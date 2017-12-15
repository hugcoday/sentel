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

package indicator

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/cloustone/sentel/libs/sentel"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
)

type IndicatorService struct {
	config     sentel.Config
	chn        chan sentel.ServiceCommand
	wg         sync.WaitGroup
	consumer   sarama.Consumer
	mongoHosts string // mongo hosts
}

// IndicatorServiceFactory
type IndicatorServiceFactory struct{}

// New create apiService service factory
func (m *IndicatorServiceFactory) New(name string, c sentel.Config, ch chan sentel.ServiceCommand) (sentel.Service, error) {
	// check mongo db configuration
	hosts, err := c.String("conductor", "mongo")
	if err != nil || hosts == "" {
		return nil, errors.New("Invalid mongo configuration")
	}
	// try connect with mongo db
	session, err := mgo.Dial(hosts)
	if err != nil {
		return nil, err
	}
	session.Close()

	// kafka
	khosts, err := c.String(name, "hosts")
	if err != nil || khosts == "" {
		return nil, errors.New("Invalid kafka configuration")
	}
	consumer, err := sarama.NewConsumer(strings.Split(khosts, ","), nil)
	if err != nil {
		return nil, fmt.Errorf("Connecting with kafka:%s failed", hosts)
	}

	return &IndicatorService{
		config:     c,
		wg:         sync.WaitGroup{},
		chn:        ch,
		consumer:   consumer,
		mongoHosts: hosts,
	}, nil

}

// Name
func (s *IndicatorService) Name() string {
	return "collector"
}

// Start
func (s *IndicatorService) Start() error {
	partitionList, err := s.consumer.Partitions("ceilometer")
	if err != nil {
		return fmt.Errorf("Failed to get list of partions:%v", err)
		return err
	}

	for partition := range partitionList {
		pc, err := s.consumer.ConsumePartition("ceilometer", int32(partition), sarama.OffsetNewest)
		if err != nil {
			glog.Errorf("Failed  to start consumer for partion %d:%s", partition, err)
			continue
		}
		defer pc.AsyncClose()
		s.wg.Add(1)

		go func(sarama.PartitionConsumer) {
			defer s.wg.Done()
			for msg := range pc.Messages() {
				s.handleNotifications(string(msg.Topic), msg.Value)
			}
		}(pc)
	}
	s.wg.Add(1)
	s.consumer.Close()
	return nil
}

// Stop
func (s *IndicatorService) Stop() {
	s.consumer.Close()
}

// handleNotifications handle notification from kafka
func (s *IndicatorService) handleNotifications(topic string, value []byte) error {
	//	if err := handleTopicObject(s, context.Background(), topic, value); err != nil {
	//		glog.Error(err)
	//		return err
	//	}
	return nil
}

func (s *IndicatorService) getDatabase() (*mgo.Database, error) {
	session, err := mgo.Dial(s.mongoHosts)
	if err != nil {
		glog.Fatalf("Failed to connect with mongo:%s", s.mongoHosts)
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	return session.DB("iothub"), nil
}
