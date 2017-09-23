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

package collector

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/cloustone/sentel/conductor/base"
	"github.com/cloustone/sentel/libs"
	"github.com/golang/glog"
)

type CollectorService struct {
	config   libs.Config
	chn      chan base.ServiceCommand
	wg       sync.WaitGroup
	consumer sarama.Consumer
}

// CollectorServiceFactory
type CollectorServiceFactory struct{}

// New create apiService service factory
func (m *CollectorServiceFactory) New(name string, c libs.Config, ch chan base.ServiceCommand) (base.Service, error) {
	hosts, err := c.String(name, "hosts")
	if err != nil || hosts == "" {
		return nil, errors.New("Invalid kafka configuration")
	}
	consumer, err := sarama.NewConsumer(strings.Split(hosts, ","), nil)
	if err != nil {
		return nil, fmt.Errorf("Connecting with kafka:%s failed", hosts)
	}

	return &CollectorService{
		config:   c,
		wg:       sync.WaitGroup{},
		chn:      ch,
		consumer: consumer,
	}, nil

}

// Name
func (s *CollectorService) Name() string {
	return "collector"
}

// Start
func (s *CollectorService) Start() error {
	partitionList, err := s.consumer.Partitions("conductor")
	if err != nil {
		return fmt.Errorf("Failed to get list of partions:%v", err)
		return err
	}

	for partition := range partitionList {
		pc, err := s.consumer.ConsumePartition("conductor", int32(partition), sarama.OffsetNewest)
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
func (s *CollectorService) Stop() {
	s.consumer.Close()
}

// handleNotifications handle notification from kafka
func (s *CollectorService) handleNotifications(topic string, value []byte) error {
	if err := handleTopicObject(s, context.Background(), topic, value); err != nil {
		glog.Error(err)
		return err
	}
	return nil
}
