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
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/cloustone/sentel/conductor/executor"
	"github.com/cloustone/sentel/core"

	"github.com/Shopify/sarama"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
)

type IndicatorService struct {
	config     core.Config
	chn        chan core.ServiceCommand
	wg         sync.WaitGroup
	consumer   sarama.Consumer
	mongoHosts string // mongo hosts
}

// IndicatorServiceFactory
type IndicatorServiceFactory struct{}

// New create apiService service factory
func (m *IndicatorServiceFactory) New(name string, c core.Config, ch chan core.ServiceCommand) (core.Service, error) {
	// check mongo db configuration
	hosts, err := c.String("conductor", "mongo")
	if err != nil || hosts == "" {
		return nil, errors.New("Invalid mongo configuration")
	}
	// try connect with mongo db
	if session, err := mgo.Dial(hosts); err != nil {
		return nil, err
	} else {
		session.Close()
	}

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
	return "indicator"
}

// Start
func (s *IndicatorService) Start() error {
	partitionList, err := s.consumer.Partitions("conductor")
	if err != nil {
		return fmt.Errorf("Failed to get list of partions:%v", err)
	}

	for partition := range partitionList {
		pc, err := s.consumer.ConsumePartition("rule", int32(partition), sarama.OffsetNewest)
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
	s.wg.Wait()
	return nil
}

// Stop
func (s *IndicatorService) Stop() {
	s.consumer.Close()
}

type ruleTopic struct {
	RuleName  string `json:"ruleName"`
	RuleId    string `json:"ruleId"`
	ProductId string `json:"productId"`
	Action    string `json:"action"`
}

// handleNotifications handle notification from kafka
func (s *IndicatorService) handleNotifications(topic string, value []byte) error {
	rule := ruleTopic{}
	if err := json.Unmarshal(value, &topic); err != nil {
		return err
	}
	r := &executor.Rule{
		RuleName:  rule.RuleName,
		RuleId:    rule.RuleId,
		ProductId: rule.ProductId,
		Action:    rule.Action,
	}
	return executor.HandleRuleNotification(r)
}
