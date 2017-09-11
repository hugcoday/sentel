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

// MQTT specified session infor
const (
	infoCleanSession       = "clean_session"
	infoMessageMaxInflight = "inflight_max"
	infoMessageInflight    = "message_in_flight"
	infoMessageInQueue     = "message_in_queue"
	infoMessageDropped     = "message_dropped"
	infoAwaitingRel        = "message_awaiting_rel"
	infoAwaitingComp       = "message_awaitng_comp"
	infoAwaitingAck        = "message_aiwaiting _ack"
	infoCreatedAt          = "created_at"
)

// Stats declarations
const (
	statClientsMax         = "clients/max"
	statClientsCount       = "client/count"
	statQueuesMax          = "queues/max"
	statQueuesCount        = "queues/count"
	statRetainedMax        = "retained/max"
	statRetainedCount      = "retained/count"
	statSessionsMax        = "sessions/max"
	statSessionsCount      = "sessions/count"
	statSubscriptionsMax   = "subscriptions/max"
	statSubscriptionsCount = "subscriptions/count"
	statTopicsMax          = "topics/max"
	statTopicsCount        = "topic/count"
)

// Metrics declarations
const (
	metricBytesReceived         = "bytes/recevied"
	metricBytesSent             = "bytes/sent"
	metricMessageDroped         = "messages/droped"
	metricMessageQos0Recevied   = "messages/qos0/received"
	metricMessageQos0Sent       = "messages/qos0/sent"
	metricMessageQos1Received   = "messages/qos1/recevied"
	metricMessageQos1Sent       = "messages/qos1/sent"
	metricMessageOos2Recevied   = "messages/qos2/received"
	metricMessageOos2Sent       = "messages/qos2/sent"
	metricMessageRetained       = "messages/retained"
	metricMessageReceived       = "messages/received"
	metricMessageSent           = "messages/sent"
	metricPacketConnack         = "packets/connack"
	metricPacketConnect         = "packets/connect"
	metricPacketDisconnect      = "packets/disconnect"
	metricPacketPingreq         = "packets/pingreq"
	metricPacketPingresp        = "packets/pingresp"
	metricPacketPubackRecevied  = "packets/puback/received"
	metricPacketPubackSent      = "packets/puback/sent"
	metricPacketPubcompReceived = "packets/pubcomp/received"
	metricPacketPubcompSent     = "packets/pubcomp/sent"
	metricPacketPublishReceived = "packets/publish/received"
	metricPacketPublishSent     = "packets/publish/sent"
	metricPacketPubrecReceived  = "packets/pubrec/received"
	metricPacketPubrecSent      = "packets/pubrec/sent"
	metricPacketPubrelReceived  = "packets/pubrel/received"
	metricPacketPubrelSent      = "packets/pubrel/sent"
	metricPacketReceived        = "packes/received"
	metricPacketSent            = "packets/sent"
	metricPacketSuback          = "packets/subback"
	metricPacketSubscribe       = "packets/subscribe"
	metricPacketUnsuback        = "packets/unsuback"
	metricPacketUnsubscribe     = "packets/unsubscribe"
)
