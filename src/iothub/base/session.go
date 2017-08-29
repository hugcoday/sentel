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
	OnConnect(Session, userdata interface{}) error
	OnDisconnect(Session, userdta interface{}) error
	OnPublish(Session, userdata interface{}) error
	OnMessage(Session, userdata interface{}) error
	OnSubscribe(Session, userdata interface{}) error
	OnUnsubscribe(Session, userdata interface{}) error
}

type Session interface {
	// Identifier get session identifier
	Identifier() string
	// GetService get the service ower for current session
	Service() Service
	// Handle indicate service to handle the packet
	Handle() error
	// Destroy will release current session
	Destroy() error
	// RegisterObserver register observer on session
	RegisterObserver(SessionObserver)
}
