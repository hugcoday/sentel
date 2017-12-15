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
	"sync"

	"github.com/cloustone/sentel/libs/sentel"
	"gopkg.in/mgo.v2"
)

type ExecutorService struct {
	config   sentel.Config
	chn      chan sentel.ServiceCommand
	wg       sync.WaitGroup
	ruleChan chan *Rule
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

func (s *ExecutorService) handleRule(r *Rule) {

}
