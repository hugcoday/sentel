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
	"errors"

	"github.com/golang/glog"
)

type localStorage struct {
	opt      Option
	sessions map[string]Session
}

// Open local storage
func (l *localStorage) Open() error {
	glog.Info("local storage Open")

	return nil
}

// Close local storage
func (l *localStorage) Close() {

}

// Backup serialize local storage
func (l *localStorage) Backup(shutdown bool) error {
	return nil
}

// Restore recover data from serialization
func (l *localStorage) Restore() error {
	return nil
}

// FindSession find session by id
func (l *localStorage) FindSession(c Context, id string) (*Session, error) {
	v, ok := l.sessions[id]
	if !ok {
		return nil, errors.New("Session id does not exist")
	}

	return &v, nil
}

// DeleteSession delete session by id
func (l *localStorage) DeleteSession(c Context, id string) error {
	_, ok := l.sessions[id]
	if !ok {
		return errors.New("Session id does not exist")
	}

	delete(l.sessions, id)

	return nil
}

// UpdateSession update session
func (l *localStorage) UpdateSession(c Context, s *Session) error {
	_, ok := l.sessions[s.Id]
	if !ok {
		return errors.New("Session id does not exist")
	}

	l.sessions[s.Id] = *s
	return nil
}

// RegisterSession register new session
func (l *localStorage) RegisterSession(c Context, s Session) error {
	if _, ok := l.sessions[s.Id]; ok {
		return errors.New("Session id already exists")
	} else {
		l.sessions[s.Id] = s
		return nil
	}
}

// Device
func (l *localStorage) AddDevice(c Context, d Device) error {
	return nil
}

func (l *localStorage) DeleteDevice(c Context, id string) error {
	return nil
}

func (l *localStorage) UpdateDevice(c Context, d Device) error {
	return nil
}

func (l *localStorage) GetDeviceState(c Context, id string) (int, error) {
	return 0, nil
}

func (l *localStorage) SetDeviceState(c Context, state int) error {
	return nil
}

// Topic
func (l *localStorage) TopicExist(c Context, t Topic) (bool, error) {
	return false, nil
}

func (l *localStorage) AddTopic(c Context, t Topic) error {
	return nil
}

func (l *localStorage) DeleteTopic(c Context, id string) error {
	return nil
}

func (l *localStorage) UpdateTopic(c Context, t Topic) error {
	return nil
}

func (l *localStorage) AddSubscriber(c Context, t Topic, clientid string) error {
	return nil
}

func (l *localStorage) RemoveSubscriber(c Context, t Topic, clientid string) error {
	return nil
}

func (l *localStorage) GetTopicSubscribers(c Context, t Topic) ([]string, error) {
	return nil, nil
}

// Subscription
func (l *localStorage) AddSubscription(c Context, sub string, qos uint8) error {
	return nil
}

func (l *localStorage) RetainSubscription(c Context, sub string, qos uint8) error {
	return nil
}

// Message Management
func (l *localStorage) FindMessage(clientid string, mid uint16) (bool, error) {
	return false, nil
}

func (l *localStorage) StoreMessage(clientid string, msg Message) error {
	return nil
}

func (l *localStorage) DeleteMessageWithValidator(clientid string, validator func(msg Message) bool) {

}

func (l *localStorage) DeleteMessage(clientid string, mid uint16, direction MessageDirection) error {
	return nil
}

func (l *localStorage) QueueMessage(clientid string, msg Message) error {
	return nil
}

func (l *localStorage) GetMessageTotalCount(clientid string) int {
	return 0
}

func (l *localStorage) InsertMessage(clientid string, mid uint16, direction MessageDirection, msg Message) error {
	return nil
}

func (l *localStorage) ReleaseMessage(clientid string, mid uint16, direction MessageDirection) error {
	return nil
}

func (l *localStorage) UpdateMessage(clientid string, mid uint16, direction MessageDirection, state MessageState) {

}

// localStorageFactory
type localStorageFactory struct{}

func (l *localStorageFactory) New(opt Option) (Storage, error) {
	d := &localStorage{
		opt:      opt,
		sessions: make(map[string]Session),
	}
	return d, nil
}
