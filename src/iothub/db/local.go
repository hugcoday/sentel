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

package db

import (
	"iothub/base"
	"iothub/util/config"

	"github.com/golang/glog"
)

type localDatabase struct {
	config   config.Config
}

// MqttFactory
type localDatabaseFactory struct{}


func (m *localDatabaseFactory) New(c config.Config) (Database, error) {
	d := &localDatabase{config: c}
	return d, nil
}