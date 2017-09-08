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
	"errors"
	"libs"

	"github.com/golang/glog"
)

type subNode struct {
	id  string
	qos uint8
}

type localStorage struct {
	config   libs.Config
	sessions map[string]*mqttSession
	subs     map[string][]subNode
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
func (l *localStorage) FindSession(id string) (*mqttSession, error) {
	v, ok := l.sessions[id]
	if !ok {
		return nil, errors.New("Session id does not exist")
	}

	return v, nil
}

// DeleteSession delete session by id
func (l *localStorage) DeleteSession(id string) error {
	_, ok := l.sessions[id]
	if !ok {
		return errors.New("Session id does not exist")
	}

	delete(l.sessions, id)

	return nil
}

// UpdateSession update session
func (l *localStorage) UpdateSession(s *mqttSession) error {
	_, ok := l.sessions[s.id]
	if !ok {
		return errors.New("Session id does not exist")
	}

	l.sessions[s.id] = s
	return nil
}

// RegisterSession register new session
func (l *localStorage) RegisterSession(s *mqttSession) error {
	if _, ok := l.sessions[s.id]; ok {
		return errors.New("Session id already exists")
	} else {
		l.sessions[s.id] = s
		return nil
	}
}

// Device
// AddDevice
// func (l *localStorage) AddDevice(d Device) error {
// 	return nil
// }

// func (l *localStorage) DeleteDevice(id string) error {
// 	return nil
// }

// func (l *localStorage) UpdateDevice(d Device) error {
// 	return nil
// }

// func (l *localStorage) GetDeviceState(id string) (int, error) {
// 	return 0, nil
// }

// func (l *localStorage) SetDeviceState(state int) error {
// 	return nil
// }

// // Topic
// func (l *localStorage) TopicExist(t Topic) (bool, error) {
// 	return false, nil
// }

// func (l *localStorage) AddTopic(t Topic) error {
// 	return nil
// }

// func (l *localStorage) DeleteTopic(id string) error {
// 	return nil
// }

// func (l *localStorage) UpdateTopic(t Topic) error {
// 	return nil
// }

// func (l *localStorage) AddSubscriber(t Topic, clientid string) error {
// 	return nil
// }

// func (l *localStorage) RemoveSubscriber(t Topic, clientid string) error {
// 	return nil
// }

// func (l *localStorage) GetTopicSubscribers(t Topic) ([]string, error) {
// 	return nil, nil
// }

// Subscription
func (l *localStorage) AddSubscription(sessionid string, topic string, qos uint8) error {
	var node subNode
	node.id = sessionid
	node.qos = qos
	if _, ok := l.subs[topic]; ok {
		l.subs[topic] = append(l.subs[topic], node)
	} else {
		l.subs[topic] = make([]subNode, 1)
		l.subs[topic][0] = node
	}
	return nil
}

func (l *localStorage) RetainSubscription(sessionid string, topic string, qos uint8) error {
	return nil
}

func (l *localStorage) RemoveSubscription(sessionid string, topic string) error {
	if _, ok := l.subs[topic]; ok {
		var index int
		var value subNode
		for index, value = range l.subs[topic] {
			if value.id == sessionid {
				break
			}
		}

		copy(l.subs[topic][index:], l.subs[topic][index+1:])
		l.subs[topic] = l.subs[topic][:len(l.subs[topic])-1]
	} else {
		return errors.New("Topic name is not exists")
	}

	return nil
}

// Message Management
func (l *localStorage) FindMessage(clientid string, mid uint16) (bool, error) {
	return false, nil
}

func (l *localStorage) StoreMessage(clientid string, msg StorageMessage) error {
	return nil
}

func (l *localStorage) DeleteMessageWithValidator(clientid string, validator func(msg StorageMessage) bool) {

}

func (l *localStorage) DeleteMessage(clientid string, mid uint16, direction MessageDirection) error {
	return nil
}

func (l *localStorage) QueueMessage(clientid string, msg StorageMessage) error {
	return nil
}

func (l *localStorage) GetMessageTotalCount(clientid string) int {
	return 0
}

func (l *localStorage) InsertMessage(clientid string, mid uint16, direction MessageDirection, msg StorageMessage) error {
	return nil
}

func (l *localStorage) ReleaseMessage(clientid string, mid uint16, direction MessageDirection) error {
	return nil
}

func (l *localStorage) UpdateMessage(clientid string, mid uint16, direction MessageDirection, state MessageState) {

}

// localStorageFactory
type localStorageFactory struct{}

func (l *localStorageFactory) New(c libs.Config) (Storage, error) {
	d := &localStorage{
		config:   c,
		sessions: make(map[string]*mqttSession),
		subs:     make(map[string][]subNode),
	}
	return d, nil
}


