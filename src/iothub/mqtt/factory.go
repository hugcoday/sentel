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
	"iothub/base"

	"github.com/golang/glog"
)

const protocolName = "mqtt3"

type mqttServiceFactory struct{}

// New create mqtt service factory
func (m *mqttServiceFactory) New(c *base.Config, ch chan struct{}) (base.Service, error) {
	return &mqtt{config: c, chn: ch}, nil
}

func init() {
	glog.Info("Registering service:%s", protocolName)
	base.RegisterServiceFactory(protocolName, &mqttServiceFactory{})
}
