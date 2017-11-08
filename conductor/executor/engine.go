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

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Shopify/sarama"
	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
)

// one product has one scpecific topic name
const publishTopicUrl = "/cluster/publish/%s"

// publishTopic is topic object received from kafka
type publishTopic struct {
	ClientId  string `json:"clientId"`
	Topic     string `json:"topic"`
	ProductId string `json:"product"`
	Content   string `json:"content"`
	encoded   []byte
	err       error
}

// ruleEngine manage product's rules, add, start and stop rule
type ruleEngine struct {
	productId  string           // one product have one rule engine
	rules      map[string]*Rule // all product's rule
	consumer   sarama.Consumer  // one product have one consumer
	config     core.Config      // configuration
	mutex      sync.Mutex       // mutex to protext rules list
	notifyChan chan int         // stop channel
	started    bool             // indicate wether engined is started
	wg         sync.WaitGroup   // waitgroup for consumper partions
}

// newRuleEngine create a engine according to product id and configuration
func newRuleEngine(c core.Config, productId string) (*ruleEngine, error) {
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
		rules:      make(map[string]*Rule),
		mutex:      sync.Mutex{},
		notifyChan: make(chan int),
		wg:         sync.WaitGroup{},
		started:    false,
	}, nil
}

// start will start the rule engine, receiving topic and rule
func (p *ruleEngine) start() error {
	if p.started {
		return fmt.Errorf("rule engine(%s) is already started", p.productId)
	}

	// start connection with kafka
	partitionList, err := p.consumer.Partitions("conductor")
	if err != nil {
		return fmt.Errorf("Failed to get list of partions:%v", err)
	}

	topic := fmt.Sprintf(publishTopicUrl, p.productId)
	for partition := range partitionList {
		pc, err := p.consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			glog.Errorf("Failed  to start consumer for partion %d:%s", partition, err)
			continue
		}
		defer pc.AsyncClose()
		p.wg.Add(1)

		go func(sarama.PartitionConsumer) {
			defer p.wg.Done()
			for msg := range pc.Messages() {
				t := publishTopic{}
				if err := json.Unmarshal(msg.Value, &t); err != nil {
					glog.Errorf("Failed to handle topic:%v", err)
					continue
				}
				if err := p.execute(&t); err != nil {
					glog.Errorf("Failed to handle topic:%v", err)
					continue
				}
			}
		}(pc)
	}
	p.started = true
	return nil
}

// stop will stop the engine
func (p *ruleEngine) stop() {
	if p.consumer != nil {
		p.consumer.Close()
	}
	p.notifyChan <- 0
	p.wg.Wait()
}

// getRuleObject get all rule's information from backend database
func (p *ruleEngine) getRuleObject(r *Rule) (*Rule, error) {
	hosts, _ := p.config.String("conductor", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		glog.Errorf("%v", err)
		return nil, err
	}
	defer session.Close()
	c := session.DB("registry").C("rules")
	obj := Rule{}
	if err := c.Find(bson.M{"RuleId": r.RuleId}).One(&obj); err != nil {
		glog.Errorf("Invalid rule with id(%s)", r.RuleId)
		return nil, err
	}
	return &obj, nil
}

// addRule add a rule received from apiserver to this engine
func (p *ruleEngine) addRule(r *Rule) error {
	glog.Infof("ruld:%s is added", r.RuleId)

	obj, err := p.getRuleObject(r)
	if err != nil {
		return err
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.rules[r.RuleId]; ok {
		return fmt.Errorf("rule:%s already exist", r.RuleId)
	}
	p.rules[r.RuleId] = obj
	return nil
}

// delteRule remove a rule from current rule engine
func (p *ruleEngine) deleteRule(r *Rule) error {
	glog.Infof("Rule:%s is deleted", r.RuleId)

	// Get rule detail
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.rules[r.RuleId]; ok {
		delete(p.rules, r.RuleId)
		return nil
	}
	return fmt.Errorf("rule:%s doesn't exist", r.RuleId)
}

// updateRule update rule in engine
func (p *ruleEngine) updateRule(r *Rule) error {
	glog.Infof("Rule:%s is updated", r.RuleId)

	obj, err := p.getRuleObject(r)
	if err != nil {
		return err
	}

	// Get rule detail
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.rules[r.RuleId]; ok {
		p.rules[r.RuleId] = obj
		return nil
	}
	return fmt.Errorf("rule:%s doesn't exist", r.RuleId)
}

// startRule start rule in engine
func (p *ruleEngine) startRule(r *Rule) error {
	glog.Infof("rule:%s is started", r.RuleId)

	// Check wether the rule engine is started
	if p.started == false {
		if err := p.start(); err != nil {
			glog.Errorf("%v", err)
			return err
		}
		p.started = true
	}

	// Start the rule
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.rules[r.RuleId]; ok {
		p.rules[r.RuleId].Status = RuleStatusStarted
		return nil
	}
	return fmt.Errorf("rule:%s doesn't exist", r.RuleId)
}

// stopRule stop rule in engine
func (p *ruleEngine) stopRule(r *Rule) error {
	glog.Infof("rule:%s is stoped", r.RuleId)

	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.rules[r.RuleId]; !ok { // not found
		return fmt.Errorf("Invalid rule:%s", r.RuleId)
	}
	p.rules[r.RuleId].Status = RuleStatusStoped
	// Stop current engine if all rules are stoped
	for _, rule := range p.rules {
		// If one of rule is not stoped, don't stop current engine
		if rule.Status != RuleStatusStoped {
			return nil
		}
	}
	p.stop()
	return nil
}

// execute rule to process published topic
// Data recevied from iothub will be processed here and transformed into database
func (p *ruleEngine) execute(t *publishTopic) error {
	return nil
}
