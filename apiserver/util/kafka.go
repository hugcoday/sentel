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
package util

import (
	"errors"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
)

func SyncProduceMessage(cfg core.Config, topic string, key string, value sarama.Encoder) error {
	// Get kafka server
	kafka, err := cfg.String("kafka", "hosts")
	if err != nil || kafka == "" {
		return errors.New("Invalid kafka configuration")
	}

	//	sarama.Logger = c.Logger()

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: value,
	}

	producer, err := sarama.NewSyncProducer(strings.Split(kafka, ","), config)
	if err != nil {
		glog.Errorf("Failed to produce message:%s", err.Error())
		return err
	}
	defer producer.Close()

	if _, _, err := producer.SendMessage(msg); err != nil {
		glog.Errorf("Failed to send producer message:%s", err.Error())
	}
	return err
}

func AsyncProduceMessage(cfg core.Config, key string, topic string, value sarama.Encoder) error {
	// Get kafka server
	kafka, err := cfg.String("kafka", "hosts")
	if err != nil || kafka == "" {
		return errors.New("Invalid kafka configuration")
	}

	//	sarama.Logger = c.Logger()

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: value,
	}

	producer, err := sarama.NewAsyncProducer(strings.Split(kafka, ","), config)
	if err != nil {
		glog.Errorf("Failed to produce message:%s", err.Error())
		return err
	}
	defer producer.Close()

	go func(p sarama.AsyncProducer) {
		errors := p.Errors()
		success := p.Successes()
		for {
			select {
			case err := <-errors:
				if err != nil {
					glog.Error(err)
				}
			case <-success:
			}
		}
	}(producer)

	producer.Input() <- msg
	return err
}
