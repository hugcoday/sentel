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

package metric

import (
	"time"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/ceilometer/collector"
	"github.com/cloustone/sentel/libs/sentel"
)

type MetricService struct {
	config    sentel.Config
	chn       chan base.ServiceCommand
	keepalive *time.Ticker
	stat      *time.Ticker
	name      string
	createdAt string
	ip        string
}

// MetricServiceFactory
type MetricServiceFactory struct{}

// New create apiService service factory
func (m *MetricServiceFactory) New(protocol string, c sentel.Config, ch chan base.ServiceCommand) (base.Service, error) {
	// Get node ip, name and created time
	return &MetricService{
		config: c,
	}, nil
}

// Name
func (s *MetricService) Info() *base.ServiceInfo {
	return &base.ServiceInfo{
		ServiceName: "report-service",
	}
}

// Start
func (s *MetricService) Start() error {
	// Launch timer scheduler
	duration, err := s.config.Int("mqttbroker", "report_duration")
	if err != nil {
		duration = 2
	}
	s.keepalive = time.NewTicker(1 * time.Second)
	s.stat = time.NewTicker(time.Duration(duration) * time.Second)
	go func(*MetricService) {
		for {
			select {
			case <-s.keepalive.C:
				s.reportKeepalive()
			case <-s.stat.C:
				s.reportHubStats()
			}
		}
	}(s)
	return nil
}

// Stop
func (s *MetricService) Stop() {
	s.keepalive.Stop()
	s.stat.Stop()
}

// reportHubStats report current iothub stats
func (s *MetricService) reportHubStats() {
	mgr := base.GetServiceManager()
	// Stats
	stats := mgr.GetStats("mqtt")
	collector.AsyncReport(s.config, collector.TopicNameStats,
		&collector.Stats{
			NodeName: s.name,
			Service:  "mqtt",
			Values:   stats,
		})

	// Metrics
	metrics := mgr.GetMetrics("mqtt")
	collector.AsyncReport(s.config, collector.TopicNameMetric,
		&collector.Metric{
			NodeName: s.name,
			Service:  "mqtt",
			Values:   metrics,
		})
}

// reportKeepalive report node information to cluster manager
func (s *MetricService) reportKeepalive() {
	mgr := base.GetServiceManager()
	// Node
	node := mgr.GetNodeInfo()
	collector.AsyncReport(s.config, collector.TopicNameNode,
		&collector.Node{
			NodeName:  node.NodeName,
			NodeIp:    node.NodeIp,
			CreatedAt: node.CreatedAt,
		})
}
