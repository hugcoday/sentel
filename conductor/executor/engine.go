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

package executor

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/cloustone/sentel/libs/sentel"
	"github.com/golang/glog"
)

const publishTopicUrl = "/cluster/publish/%s"

type publishTopic struct {
	ClientId  string `json:"clientId"`
	Topic     string `json:"topic"`
	ProductId string `json:"product"`
	Content   string `json:"context"`
	encoded   []byte
	err       error
}

type ruleEngine struct {
	productId  string          // one product have one rule engine
	rules      []*Rule         // all product's rule
	consumer   sarama.Consumer // one product have one consumer
	config     sentel.Config   // configuration
	mutex      sync.Mutex      // mutex to protext rules list
	ruleChan   chan *Rule      // rule receiver
	notifyChan chan int        // stop channel
	started    bool            // indicate wether engined is started
	wg         sync.WaitGroup  // waitgroup for consumper partions
}

func newRuleEngine(c sentel.Config, productId string) (*ruleEngine, error) {
	khosts, err := c.String("conductor", "kafka")
	if err != nil || khosts == "" {
		return nil, errors.New("Invalid kafka configuration")
	}
	consumer, err := sarama.NewConsumer(strings.Split(khosts, ","), nil)
	if err != nil {
		return nil, fmt.Errorf("Connecting with kafka:%s failed", khosts)
	}

	return &ruleEngine{
		productId:  productId,
		config:     c,
		consumer:   consumer,
		rules:      []*Rule{},
		mutex:      sync.Mutex{},
		ruleChan:   make(chan *Rule),
		notifyChan: make(chan int),
		wg:         sync.WaitGroup{},
		started:    false,
	}, nil
}

func (r *ruleEngine) start() error {
	if r.started {
		return fmt.Errorf("rule engine(%s) is already started", r.productId)
	}

	// start rule receiver
	r.wg.Add(1)
	go func(r *ruleEngine) {
		for {
			select {
			case rule := <-r.ruleChan:
				r.mutex.Lock()
				r.rules = append(r.rules, rule)
				r.mutex.Unlock()
			case <-r.notifyChan:
				break
			}
		}
		r.wg.Done()
	}(r)

	// start connection with kafka
	partitionList, err := r.consumer.Partitions("conductor")
	if err != nil {
		return fmt.Errorf("Failed to get list of partions:%v", err)
	}

	topic := fmt.Sprintf(publishTopicUrl, r.productId)
	for partition := range partitionList {
		pc, err := r.consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			glog.Errorf("Failed  to start consumer for partion %d:%s", partition, err)
			continue
		}
		defer pc.AsyncClose()
		r.wg.Add(1)

		go func(sarama.PartitionConsumer) {
			defer r.wg.Done()
			for msg := range pc.Messages() {
				r.handleTopic(string(msg.Topic), msg.Value)
			}
		}(pc)
	}
	r.started = true
	return nil
}

func (r *ruleEngine) stop() {
	if r.consumer != nil {
		r.consumer.Close()
	}
	r.notifyChan <- 0
	r.wg.Wait()
}

func (r *ruleEngine) addRule(rule *Rule) {
	r.ruleChan <- rule
}

func (r *ruleEngine) handleTopic(topic string, value []byte) error {
	t := publishTopic{}
	if err := json.Unmarshal(value, &t); err != nil {
		return err
	}
	return r.execute(&t)
}

func (r *ruleEngine) execute(t *publishTopic) error {
	return nil
}
