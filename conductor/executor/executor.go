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
	"errors"
	"fmt"
	"sync"

	"github.com/cloustone/sentel/libs/sentel"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ExecutorService struct {
	config   sentel.Config
	chn      chan sentel.ServiceCommand
	wg       sync.WaitGroup
	ruleChan chan *Rule
	engines  map[string]*ruleEngine
	mutex    sync.Mutex
}

type ExecutorServiceFactory struct{}

var executorService *ExecutorService

// New create executor service factory
func (m *ExecutorServiceFactory) New(name string, c sentel.Config, ch chan sentel.ServiceCommand) (sentel.Service, error) {
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
	executorService = &ExecutorService{
		config:   c,
		wg:       sync.WaitGroup{},
		chn:      ch,
		ruleChan: make(chan *Rule),
		engines:  make(map[string]*ruleEngine),
		mutex:    sync.Mutex{},
	}
	return executorService, nil
}

// Name
func (s *ExecutorService) Name() string {
	return "executor"
}

// Start
func (s *ExecutorService) Start() error {
	// start rule channel
	go func(s *ExecutorService) {
		s.wg.Add(1)
		select {
		case r := <-s.ruleChan:
			s.handleRule(r)
		case <-s.chn:
			break
		}
	}(s)
	s.wg.Wait()
	return nil
}

// Stop
func (s *ExecutorService) Stop() {
	s.chn <- 1
}

// executeRule push a new rule to executor service
func pushRule(r *Rule) {
	executorService.ruleChan <- r
}

func (s *ExecutorService) handleRule(r *Rule) error {
	if _, ok := s.engines[r.ProductId]; !ok { // not found
		engine, err := newRuleEngine(s.config, r.ProductId)
		if err != nil {
			glog.Errorf("Failed to create rule engint for product(%s)", r.ProductId)
			return err
		}
		s.engines[r.ProductId] = engine
	}
	engine := s.engines[r.ProductId]
	engine.addRule(r)
	return nil
}

func HandleRuleNotification(cfg sentel.Config, r *Rule, action string) error {
	glog.Infof("New rule notification: ruleId=%s, ruleName=%s, action=%s", r.RuleId, r.RuleName, action)

	// Check action's validity
	switch action {
	case RuleActionNew:
	case RuleActionDelete:
	case RuleActionUpdated:
	case RuleActionStart:
	case RuleActionStop:
	default:
		return fmt.Errorf("Invalid rule action(%s) for product(%s)", action, r.ProductId)
	}
	// Get rule detail
	hosts, _ := cfg.String("conductor", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		glog.Errorf("%v", err)
		return err
	}
	defer session.Close()
	c := session.DB("registry").C("rules")
	obj := Rule{}
	if err := c.Find(bson.M{"RuleId": r.RuleId}).One(&obj); err != nil {
		glog.Errorf("Invalid rule with id(%s)", r.RuleId)
		return err
	}
	// Parse sql and target

	// Now just simply send rule to executor
	pushRule(&obj)
	return nil
}
