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

import "io"

// ClientInfo
type ClientInfo struct {
	UserName     string
	CleanSession bool
	PeerName     string
	ConnectTime  string
}

// SessionInfo
type SessionInfo struct {
	ClientId           string
	CleanSession       bool
	MessageMaxInflight uint64
	MessageInflight    uint64
	MessageInQueue     uint64
	MessageDropped     uint64
	AwaitingRel        uint64
	AwaitingComp       uint64
	AwaitingAck        uint64
	CreatedAt          string
}

// RouteInfo
type RouteInfo struct {
	Topic string
	Route []string
}

// TopicInfo
type TopicInfo struct {
	Topic     string
	Attribute string
}

// SubscriptionInfo
type SubscriptionInfo struct {
	ClientId  string
	Topic     string
	Attribute string
}

type ProtocolService interface {
	// Service Infomation
	GetServiceInfo() *ServiceInfo

	// Stats and Metrics
	GetStats() *Stats
	GetMetrics() *Metrics

	// Client
	GetClients() []*ClientInfo
	GetClient(id string) *ClientInfo
	KickoffClient(id string) error

	// Session
	GetSessions(conditions map[string]bool) []*SessionInfo
	GetSession(id string) *SessionInfo

	// Route info
	GetRoutes() []*RouteInfo
	GetRoute(id string) *RouteInfo

	// Topic info
	GetTopics() []*TopicInfo
	GetTopic(id string) *TopicInfo

	// SubscriptionInfo
	GetSubscriptions() []*SubscriptionInfo
	GetSubscription(id string) *SubscriptionInfo
}

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

type Session interface {
	// Identifier get session identifier
	Identifier() string
	// Info return session information
	Info() *SessionInfo
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

type Packet interface {
	// PacketType return type name of packet
	PacketType() string

	// DecodeFromReader decode packet from given reader
	DecodeFromReader(r io.Reader, df DecodeFeedback) error

	// DecodeFromBytes decode packet from given
	DecodeFromBytes(data []uint8, df DecodeFeedback) error

	// SerializeTo writes the serialized form of the packet into the serialize buffer
	SerializeTo(buf SerializeBuffer, opts SerializeOptions) error

	// Clear clear packet content and payload
	Clear()

	// Length return length of the packet
	Length() int
}
