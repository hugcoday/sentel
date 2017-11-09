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

package iothub

import (
	"sync"

	"github.com/cloustone/sentel/core"
)

type Iothub struct {
	sync.Once
	tenantsMutex sync.Mutex
	tenants      map[string]*Tenant
	haproxy      *Haproxy
	brokers      map[string]*Broker
	brokersMutex sync.Mutex
	clustermgr   *clusterManager
}

var (
	_iothub *Iothub
)

// InitializeIothub create iothub global instance at startup time
func InitializeIothub(c core.Config) error {
	haproxy, err := NewHaproxy(c)
	if err != nil {
		return err
	}
	clustermgr, err := newClusterManager(c)
	if err != nil {
		return err
	}
	_iothub = &Iothub{
		tenantsMutex: sync.Mutex{},
		tenants:      make(map[string]*Tenant),
		haproxy:      haproxy,
		brokers:      make(map[string]*Broker),
		brokersMutex: sync.Mutex{},
		clustermgr:   clustermgr,
	}
	return nil
}

// getIothub return global iothub instance used in iothub packet
func getIothub() *Iothub {
	return _iothub
}

// addTenant add tenant to iothub
func (this *Iothub) addTenant(t *Tenant) error {
	return nil
}

// deleteTenant delete tenant from iothub
func (this *Iothub) deleteTenant(id string) error {
	return nil
}

// addBroker add broker to tenant
func (this *Iothub) addBroker(t *Tenant, b *Broker) error {
	return nil
}

// deleteBroker delete broker from iothub
func (this *Iothub) deleteBroker(id string) error {
	return nil
}

// startBroker start a tenant's broker
func (this *Iothub) startBroker(id string) error {
	return nil
}

// startTenantBrokers start tenant's all broker
func (this *Iothub) startTenantBrokers(t *Tenant) error {
	return nil
}

// stopBroker stop a broker
func (this *Iothub) stopBroker(id string) error {
	return nil
}

// stopTenantBrokers stop tenant's all brokers
func (this *Iothub) stopTenantBrokers(t *Tenant) error {
	return nil
}

// getTenant retrieve a tenant by id
func (this *Iothub) getTenant(id string) *Tenant {
	return nil
}

// getBroker retrieve a broker by id
func (this *Iothub) getBroker(id string) *Broker {
	return nil
}

// setBrokerStatus set broker's status
func (this *Iothub) setBrokerStatus(id string, status BrokerStatus) {
}

// getBrokerStatus retrieve broker's status
func (this *Iothub) getBrokerStatus(id string) BrokerStatus {
	return BrokerStatusStoped
}
