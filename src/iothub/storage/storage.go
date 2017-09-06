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

package storage

import (
	"fmt"
	"time"

	"github.com/golang/glog"
)

// Session storage
type Session struct {
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

type Device struct{}
type Topic struct {
	Name string
}
type MessageState int
type Message struct {
	Id        uint
	Topic     string
	Direction MessageDirection
	State     MessageState
	Qos       uint8
	Retain    bool
	Payload   []uint8
}

type Context interface{}

// Storage option
type Option struct {
	Hosts string
}

type Storage interface {
	Open() error
	Close()
	Backup(shutdown bool) error
	Restore() error

	// Session
	FindSession(c Context, id string) (*Session, error)
	DeleteSession(c Context, id string) error
	UpdateSession(c Context, s *Session) error
	RegisterSession(c Context, s Session) error

	// Device
	AddDevice(c Context, d Device) error
	DeleteDevice(c Context, id string) error
	UpdateDevice(c Context, d Device) error
	GetDeviceState(c Context, id string) (int, error)
	SetDeviceState(c Context, state int) error

	// Topic
	TopicExist(c Context, t Topic) (bool, error)
	AddTopic(c Context, t Topic) error
	DeleteTopic(c Context, id string) error
	UpdateTopic(c Context, t Topic) error
	AddSubscriber(c Context, t Topic, clientid string) error
	RemoveSubscriber(c Context, t Topic, clientid string) error
	GetTopicSubscribers(c Context, t Topic) ([]string, error)

	// Subscription
	AddSubscription(c Context, clientid string, sub string, qos uint8) error
	RetainSubscription(c Context, clientid string, sub string, qos uint8) error
	RemoveSubscription(c Context, clientid string, sub string) error

	// Message Management
	FindMessage(clientid string, mid uint16) (bool, error)
	StoreMessage(clientid string, msg Message) error
	DeleteMessageWithValidator(clientid string, validator func(msg Message) bool)
	DeleteMessage(clientid string, mid uint16, direction MessageDirection) error

	QueueMessage(clientid string, msg Message) error
	GetMessageTotalCount(clientid string) int
	InsertMessage(clientid string, mid uint16, direction MessageDirection, msg Message) error
	ReleaseMessage(clientid string, mid uint16, direction MessageDirection) error
	UpdateMessage(clientid string, mid uint16, direction MessageDirection, state MessageState)
}

type storageFactory interface {
	New(opt Option) (Storage, error)
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
func New(name string, opt Option) (Storage, error) {
	if _allStorage[name] == nil {
		return nil, fmt.Errorf("Storage %s is not registered", name)
	}
	return _allStorage[name].New(opt)
}

func init() {
	registerStorage("local", &localStorageFactory{})
	// registerStorage("etcd", etcdStorageFactory{})
}
