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
	"strings"

	"github.com/cloustone/sentel/libs"

	"github.com/golang/glog"
)

type subLeaf struct {
	qos uint8
}

type subNode struct {
	level      string
	children   map[string]*subNode
	subs       map[string]*subLeaf
	retainMsg *StorageMessage
}

type localStorage struct {
	config   libs.Config
	sessions map[string]*mqttSession
	root     subNode
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
	} 

	l.sessions[s.id] = s
	return nil
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

func (l *localStorage) findNode(node *subNode, lev string) *subNode {
	for k, v := range node.children {
		if k == lev {
			return v
		}
	}

	return nil
}

func (l *localStorage) addNode(node *subNode, lev string) *subNode {
	for k, v := range node.children {
		if k == lev {
			return v
		}
	}

	tmp := &subNode{
		level:    lev,
		children: make(map[string]*subNode),
		subs:     make(map[string]*subLeaf),
	}

	node.children[lev] = tmp

	return tmp
}

// Subscription
func (l *localStorage) AddSubscription(sessionid string, topic string, qos uint8) error {
	node := &l.root
	s := strings.Split(topic, "/")
	for _, level := range s {
		node = l.addNode(node, level)
	}

	node.subs[sessionid] = &subLeaf{
		qos: qos,
	}

	return nil
}

// RetainSubscription process RETAIN flag；
func (l *localStorage) RetainSubscription(sessionid string, topic string, qos uint8) error {
	return nil
}

func (l *localStorage) RemoveSubscription(sessionid string, topic string) error {
	node := &l.root
	s := strings.Split(topic, "/")
	for _, level := range s {
		node = l.findNode(node, level)
		if node == nil {
			return nil;
		}
	}

	if _, ok := node.subs[sessionid]; ok {
		delete(node.subs, sessionid)
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

func (l *localStorage) insertMessage(s *mqttSession, mid uint, qos uint8, msg *StorageMessage) error {
	/* Check whether we've already sent this message to this client
	 * for outgoing messages only.
	 * If retain==true then this is a stale retained message and so should be
	 * sent regardless. FIXME - this does mean retained messages will received
	 * multiple times for overlapping subscriptions, although this is only the
	 * case for SUBSCRIPTION with multiple subs in so is a minor concern.
	 */
	if option, err := l.config.Bool("mqtt", "allow_duplicate_messages"); err != nil && !option {
		// TODO
	}

	outMsg := mqttMessage {
		mid:          uint16(mid),
		payload:      msg.Payload,
		qos:          qos,
		retain:       msg.Retain,
		topic:        msg.Topic,
	}

	s.QueueMessage(&outMsg)

	return nil
}

func (l *localStorage) subProcess(clientid string, msg *StorageMessage, node *subNode, setRetain bool) error {
	if msg.Retain && setRetain {
		node.retainMsg = msg
	}

	for k, v := range node.subs {
		s := l.sessions[k]
		if s.id == clientid {
			continue
		}

		/* Check for ACL topic access. */
		// TODO

		var msgQos uint8
		clientQos := v.qos
		if option, err := l.config.Bool("mqtt", "upgrade_outgoing_qos"); err != nil && option {
			msgQos = clientQos
		} else {
			if msg.Qos > clientQos {
				msgQos = clientQos
			} else {
				msgQos = msg.Qos
			}
		}

		var mid uint
		if msgQos > 0 {
			mid = s.generateMid()
		} else {
			mid = 0
		}

		l.insertMessage(s, mid, msgQos, msg)
	}
	return nil
}

func (l *localStorage) subSearch(clientid string, msg *StorageMessage, node *subNode, levels []string, setRetain bool) error {
	for k, v := range node.children {
		sr := setRetain
		if len(levels) != 0 && (k == levels[0] || k == "+") {
			if k == "+" {
				sr = false
			}
			ss := levels[1:]
			l.subSearch(clientid, msg, v, ss, sr)
			if len(ss) == 0 {
				l.subProcess(clientid, msg, v, sr)
			}
		} else if k == "#" && len(v.children) > 0 {
			l.subProcess(clientid, msg, v, sr)
		}
	}
	return nil
}

func (l *localStorage) QueueMessage(clientid string, msg StorageMessage) error {
	s := strings.Split(msg.Topic, "/")
	if msg.Retain {
		/* We have a message that needs to be retained, so ensure that the subscription
		 * tree for its topic exists.
		 */
		node := &l.root
		for _, level := range s {
			node = l.addNode(node, level)
		}
	}

	return l.subSearch(clientid, &msg, &l.root, s, true)
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
	}
	return d, nil
}
