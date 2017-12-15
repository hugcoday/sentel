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
	"fmt"

	"github.com/cloustone/sentel/libs/sentel"
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func HandleRuleNotification(cfg sentel.Config, r *Rule, action string) error {
	glog.Infof("New rule notification: ruleId=%s, ruleName=%s, action=%s", r.RuleId, r.RuleName, action)

	// Check action's validity
	switch action {
	case RuleActionNew:
	case RuleActionDelete:
	case RuleActionUpdated:
	case RuleActionStart:
	case RuleActionStop:
	default:
		return fmt.Errorf("Invalid rule action(%s) for product(%s)", action, r.ProductId)
	}
	// Get rule detail
	hosts, _ := cfg.String("conductor", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		glog.Errorf("%v", err)
		return err
	}
	defer session.Close()
	c := session.DB("registry").C("rules")
	obj := Rule{}
	if err := c.Find(bson.M{"RuleId": r.RuleId}).One(&obj); err != nil {
		glog.Errorf("Invalid rule with id(%s)", r.RuleId)
		return err
	}
	// Parse sql and target

	// Now just simply send rule to executor
	pushRule(&obj)
	return nil
}
