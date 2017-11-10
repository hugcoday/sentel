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

package iothub

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
)

const (
	notifyActionCreate = "create"
	notifyActionDelete = "delete"
	notifyActionUpdate = "update"
)

// TenantNofiy is notification object from api server by kafka
type tenantNotify struct {
	action      string `json:"action"`
	id          string `json:"tenantid"`
	brokerCount string `json:"brokerCount"`
}

type NotifyService struct {
	config   core.Config
	chn      chan core.ServiceCommand
	wg       sync.WaitGroup
	consumer sarama.Consumer
}

// NotifyServiceFactory
type NotifyServiceFactory struct{}

// New create apiService service factory
func (m *NotifyServiceFactory) New(protocol string, c core.Config, ch chan core.ServiceCommand) (core.Service, error) {
	// kafka
	khosts, err := c.String("iothub", "hosts")
	if err != nil || khosts == "" {
		return nil, errors.New("Invalid kafka configuration")
	}
	consumer, err := sarama.NewConsumer(strings.Split(khosts, ","), nil)
	if err != nil {
		return nil, fmt.Errorf("Connecting with kafka:%s failed", khosts)
	}

	return &NotifyService{
		config:   c,
		wg:       sync.WaitGroup{},
		chn:      ch,
		consumer: consumer,
	}, nil
}

// Name
func (this *NotifyService) Name() string {
	return "notify-service"
}

// Start
func (this *NotifyService) Start() error {
	partitionList, err := this.consumer.Partitions("apiserver")
	if err != nil {
		return fmt.Errorf("Failed to get list of partions:%v", err)
		return err
	}

	for partition := range partitionList {
		pc, err := this.consumer.ConsumePartition("apiserver", int32(partition), sarama.OffsetNewest)
		if err != nil {
			glog.Errorf("Failed  to start consumer for partion %d:%s", partition, err)
			continue
		}
		defer pc.AsyncClose()
		this.wg.Add(1)

		go func(sarama.PartitionConsumer) {
			defer this.wg.Done()
			for msg := range pc.Messages() {
				this.handleNotifications(string(msg.Topic), msg.Value)
			}
		}(pc)
	}
	this.wg.Wait()
	return nil
}

// Stop
func (this *NotifyService) Stop() {
	this.consumer.Close()
	this.wg.Wait()
}

// handleNotifications handle notification from kafka
func (this *NotifyService) handleNotifications(topic string, value []byte) error {
	switch topic {
	case "tenant":
		obj := &tenantNotify{}
		if err := json.Unmarshal(value, obj); err != nil {
			return err
		}
		return this.handleTenantNotify(obj)
	}

	return nil
}

// handleTenantNotify handle notification about tenant from api server
func (this *NotifyService) handleTenantNotify(tf *tenantNotify) error {
	glog.Infof("iothub-notifyservice: tenant(%s) notification received", tf.id)

	hub := getIothub()

	switch tf.action {
	case notifyActionCreate:
		return hub.addTenant(tf.id)
	case notifyActionDelete:
		return hub.deleteTenant(tf.id)
	case notifyActionUpdate:
		return hub.updateTenant(tf.id)
	}
	return nil
}
