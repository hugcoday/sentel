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

package base

import (
	"flag"
	"fmt"
	"os"

	"github.com/golang/glog"
	mc "github.com/koding/multiconfig"
)

type Config struct {
	Host      string          // server host:port
	LogLevel  string          // Log level
	Registry  *RegistryConfig // Registry db name
	Kafka     string
	Protocols []string
	Mqtt      *MqttConfig
	Coap      *CoapConfig
}

type MqttConfig struct {
	Host string
}

type CoapConfig struct {
	Host string
}

type RegistryConfig struct {
	Server   string
	User     string
	Password string
}

const defaultConfigFilePath = "../etc/sentel/iothub.conf"

func NewConfig() (*Config, error) {
	flag.Parse()
	config := &Config{}
	c := newLoaderWithPath(defaultConfigFilePath)
	c.mustLoad(config)
	return config, nil
}

func (c *Config) Close() {
	glog.Flush()
}

// ConfigLoader
type configLoader struct {
	mc.DefaultLoader
}

func newLoaderWithPath(path string) *configLoader {
	loader := &configLoader{}
	loader.DefaultLoader = *mc.NewWithPath(path)
	return loader
}

func mustLoadWithPath(path string, conf interface{}) {
	d := newLoaderWithPath(path)
	d.mustLoad(conf)
}

func (c *configLoader) mustLoad(conf interface{}) {
	if err := c.Load(conf); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

func (c *configLoader) mustValidate(conf interface{}) {
	c.MustValidate(conf)
}
