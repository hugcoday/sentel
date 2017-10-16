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

package collector

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	TopicNameNode         = "/cluster/nodes"
	TopicNameClient       = "/cluster/clients"
	TopicNameSession      = "/cluster/sessions"
	TopicNameSubscription = "/cluster/subscriptions"
	TopicNamePublish      = "/cluster/publish"
	TopicNameMetric       = "/cluster/metrics"
	TopicNameStats        = "/cluster/stats"
)

const (
	ObjectActionRegister   = "register"
	ObjectActionUnregister = "unregister"
	ObjectActionRetrieve   = "retrieve"
	ObjectActionDelete     = "delete"
	ObjectActionUpdate     = "update"
)

type topicBase struct {
	encoded []byte
	err     error
}

func (p *topicBase) ensureEncoded() {
	if p.encoded == nil && p.err == nil {
		p.encoded, p.err = json.Marshal(p)
	}
}

func (p *topicBase) Length() int {
	p.ensureEncoded()
	return len(p.encoded)
}

func (p *topicBase) Encode() ([]byte, error) {
	p.ensureEncoded()
	return p.encoded, p.err
}

type topicObject interface {
	name() string
	clone() topicObject // Not realy clone, just construct a new object
	handleTopic(s *CollectorService, ctx context.Context) error
}

var _topicObjects map[string]topicObject = make(map[string]topicObject)

func registerTopicObject(t topicObject) {
	if _, ok := _topicObjects[t.name()]; !ok {
		_topicObjects[t.name()] = t
	}
}

func handleTopicObject(s *CollectorService, ctx context.Context, topic string, value []byte) error {
	if obj, ok := _topicObjects[topic]; !ok || obj == nil {
		return fmt.Errorf("No valid handler for topic:%s", topic)
	}

	obj := _topicObjects[topic].clone()
	if err := json.Unmarshal(value, &obj); err != nil {
		return err
	}

	return obj.handleTopic(s, ctx)
}
