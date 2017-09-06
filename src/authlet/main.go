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

package main

import (
	"authlet/authlet"
	"flag"
	"libs"

	"github.com/golang/glog"
)

var (
	configFileFullPath = flag.String("c", "../etc/sentel/authlet.conf", "config file")
)

func main() {
	var config libs.Config
	var err error

	flag.Parse()
	glog.Info("Starting authlet rpc server...")

	// Get configuration
	if config, err = libs.NewWithConfigFile(*configFileFullPath); err != nil {
		glog.Fatal(err)
		flag.PrintDefaults()
		return
	}

	if server, err := authlet.NewAuthServer(config); err != nil {
		glog.Fatal("Failed to launch Authlet Server")
		return
	} else {
		server.Start()
		//server.Stop()
		server.Wait()
	}
}
