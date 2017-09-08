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

package mqtt

import (
	"context"
	"fmt"
	"libs"
	"time"

	"github.com/golang/glog"
)

// Session storage
type StorageSession struct {
	Id                 string
	Username           string
	Password           string
	Keepalive          uint16
	LastMid            uint16
	State              uint8
	LastMessageInTime  time.Time
	LastMessageOutTime time.Time
	Ping               time.Time
	CleanSession       uint8
	SubscribeCount     uint32
	Protocol           uint8
	RefCount           uint8
}

type MessageDirection int

const (
	MessageDirectionIn  MessageDirection = 0
	MessageDirectionOut MessageDirection = 1
)

type StorageDevice struct{}
type StorageTopic struct {
	Name string
}
type MessageState int
type StorageMessage struct {
	Id        uint
	Topic     string
	Direction MessageDirection
	State     MessageState
	Qos       uint8
	Retain    bool
	Payload   []uint8
}

type Storage interface {
	Open() error
	Close()
	Backup(shutdown bool) error
	Restore() error

	// Session
	FindSession(c context.Context, id string) (*StorageSession, error)
	DeleteSession(c context.Context, id string) error
	UpdateSession(c context.Context, s *StorageSession) error
	RegisterSession(c context.Context, s StorageSession) error

	// Device
	AddDevice(c context.Context, d StorageDevice) error
	DeleteDevice(c context.Context, id string) error
	UpdateDevice(c context.Context, d StorageDevice) error
	GetDeviceState(c context.Context, id string) (int, error)
	SetDeviceState(c context.Context, state int) error

	// Topic
	TopicExist(c context.Context, t StorageTopic) (bool, error)
	AddTopic(c context.Context, t StorageTopic) error
	DeleteTopic(c context.Context, id string) error
	UpdateTopic(c context.Context, t StorageTopic) error
	AddSubscriber(c context.Context, t StorageTopic, clientid string) error
	RemoveSubscriber(c context.Context, t StorageTopic, clientid string) error
	GetTopicSubscribers(c context.Context, t StorageTopic) ([]string, error)

	// Subscription
	AddSubscription(c context.Context, clientid string, sub string, qos uint8) error
	RetainSubscription(c context.Context, clientid string, sub string, qos uint8) error
	RemoveSubscription(c context.Context, clientid string, sub string) error

	// Message Management
	FindMessage(clientid string, mid uint16) (bool, error)
	StoreMessage(clientid string, msg StorageMessage) error
	DeleteMessageWithValidator(clientid string, validator func(StorageMessage) bool)
	DeleteMessage(clientid string, mid uint16, direction MessageDirection) error

	QueueMessage(clientid string, msg StorageMessage) error
	GetMessageTotalCount(clientid string) int
	InsertMessage(clientid string, mid uint16, direction MessageDirection, msg StorageMessage) error
	ReleaseMessage(clientid string, mid uint16, direction MessageDirection) error
	UpdateMessage(clientid string, mid uint16, direction MessageDirection, state MessageState)
}

type storageFactory interface {
	New(c libs.Config) (Storage, error)
}

var _allStorage = make(map[string]storageFactory)

func registerStorage(name string, s storageFactory) {
	if _allStorage[name] != nil {
		glog.Fatalf("Storage %s already registered", name)
		return
	}
	_allStorage[name] = s
}

// New storage lookup registered storage list, create a new storage instance
func NewStorage(name string, c libs.Config) (Storage, error) {
	if _allStorage[name] == nil {
		return nil, fmt.Errorf("Storage %s is not registered", name)
	}
	return _allStorage[name].New(c)
}

func init() {
	//registerStorage("local", &localStorageFactory{})
}
