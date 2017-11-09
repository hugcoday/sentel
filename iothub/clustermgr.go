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

import "github.com/cloustone/sentel/core"

type clusterManager struct {
	ip   string
	port string
}

// newClusterManager retrieve clustermanager instance connected with clustermgr
func newClusterManager(c core.Config) (*clusterManager, error) {
	return nil, nil
}

// createBroker start specified node
func (this *clusterManager) createBroker(id string) (*Broker, error) {
	return nil, nil
}

// startBroker start specified node
func (this *clusterManager) startBroker(id string) error {
	return nil
}

// stopBroker stop specified node
func (this *clusterManager) stopBroker(id string) error {
	return nil
}

// deleteBroker stop and delete specified node
func (this *clusterManager) deleteBroker(id string) error {
	return nil
}

// rollbackBrokers rollback tenant's brokers
func (this *clusterManager) rollbackTenantBrokers(oldTenant *Tenant, newTenant *Tenant) error {
	return nil
}
