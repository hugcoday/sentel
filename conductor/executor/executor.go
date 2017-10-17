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

//  A global executor service instance is needed because indication will send rule to it
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
	return nil
}

// Stop
func (s *ExecutorService) Stop() {
	s.chn <- 1
	s.wg.Wait()

	// stop all ruleEngine
	for _, engine := range s.engines {
		if engine != nil {
			engine.stop()
		}
	}
}

func (s *ExecutorService) handleRule(r *Rule) error {
	// Get engine instance according to product id
	if _, ok := s.engines[r.ProductId]; !ok { // not found
		engine, err := newRuleEngine(s.config, r.ProductId)
		if err != nil {
			glog.Errorf("Failed to create rule engint for product(%s)", r.ProductId)
			return err
		}
		s.engines[r.ProductId] = engine
	}
	engine := s.engines[r.ProductId]

	switch r.Action {
	case RuleActionNew:
		return engine.addRule(r)
	case RuleActionDelete:
		return engine.deleteRule(r)
	case RuleActionUpdate:
		return engine.updateRule(r)
	case RuleActionStart:
		return engine.startRule(r)
	case RuleActionStop:
		return engine.stopRule(r)
	}
	return nil
}

// HandleRuleNotification handle rule notifications recevied from kafka,
// it will check rule's validity,for example, wether rule exist in database.
func HandleRuleNotification(r *Rule) error {
	glog.Infof("New rule notification: ruleId=%s, ruleName=%s, action=%s", r.RuleId, r.RuleName, r.Action)

	// Check action's validity
	switch r.Action {
	case RuleActionNew:
	case RuleActionDelete:
	case RuleActionUpdate:
	case RuleActionStart:
	case RuleActionStop:
	default:
		return fmt.Errorf("Invalid rule action(%s) for product(%s)", r.Action, r.ProductId)
	}

	if r.RuleId == "" || r.ProductId == "" {
		return fmt.Errorf("Invalid argument")
	}
	// Now just simply send rule to executor
	executorService.ruleChan <- r
	return nil
}
