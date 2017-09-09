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

package base

type SessionObserver interface {
	OnGetMountPoint() string
	OnConnect(s Session, userdata interface{}) error
	OnDisconnect(s Session, userdta interface{}) error
	OnPublish(s Session, userdata interface{}) error
	OnMessage(s Session, userdata interface{}) error
	OnSubscribe(s Session, userdata interface{}) error
	OnUnsubscribe(s Session, userdata interface{}) error
	OnAuthenticate(s Session, username string, password string) error
}

type SessionInfo map[string]string

type Session interface {
	// Identifier get session identifier
	Identifier() string
	// Info return session information
	Info() SessionInfo
	// GetService get the service ower for current session
	Service() Service
	// Handle indicate service to handle the packet
	Handle() error
	// Destroy will release current session
	Destroy() error
	// RegisterObserver register observer on session
	RegisterObserver(SessionObserver)
	// Get Stats
	GetStats() *Stats
	// Get Metrics
	GetMetrics() *Metrics
}
