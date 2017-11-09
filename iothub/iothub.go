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
	"errors"
	"fmt"
	"sync"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
	uuid "github.com/satori/go.uuid"
)

type Iothub struct {
	sync.Once
	config       core.Config
	tenantsMutex sync.Mutex
	tenants      map[string]*Tenant
	brokers      map[string]*Broker
	brokersMutex sync.Mutex
	clustermgr   *clusterManager
}

type Tenant struct {
	id           string    `json:"tenantId"`
	createdAt    time.Time `json:"createdAt"`
	brokersCount int       `json:"brokersCount"`
	brokers      []*Broker
}

var (
	_iothub *Iothub
)

// InitializeIothub create iothub global instance at startup time
func InitializeIothub(c core.Config) error {
	// check mongo db configuration
	hosts, err := c.String("iothub", "mongo")
	if err != nil || hosts == "" {
		return errors.New("Invalid mongo configuration")
	}
	// try connect with mongo db
	session, err := mgo.Dial(hosts)
	if err != nil {
		return err
	}
	session.Close()

	clustermgr, err := newClusterManager(c)
	if err != nil {
		return err
	}
	_iothub = &Iothub{
		config:       c,
		tenantsMutex: sync.Mutex{},
		tenants:      make(map[string]*Tenant),
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

// getTenantFromDatabase retrieve tenant from database
func (this *Iothub) getTenantFromDatabase(tid string) (*Tenant, error) {
	hosts, err := this.config.String("iothub", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	tenant := Tenant{}
	c := session.DB("iothub").C("tenants")
	if err := c.Find(bson.M{"tenantId": tid}).One(&tenant); err != nil {
		return nil, err
	}
	return &tenant, nil
}

// addTenant add tenant to iothub
func (this *Iothub) addTenant(tid string) error {
	// check wether the tenant already exist
	this.tenantsMutex.Lock()
	defer this.tenantsMutex.Unlock()

	if _, ok := this.tenants[tid]; ok {
		return fmt.Errorf("tenant(%s) already exist", tid)
	}

	// retrieve tenant from database
	tenant, err := this.getTenantFromDatabase(tid)
	if err != nil {
		return err
	}
	// save new tenant into iothub
	this.tenants[tid] = tenant

	// create brokers accroding to tenant's request
	for i := 0; i < tenant.brokersCount; i++ {
		bid := uuid.NewV4().String()
		broker, err := this.clustermgr.startBroker(bid)
		if err != nil {
			// TODO should we update database status
			glog.Fatalf("Failed to created broker for tenant(%s)", tid)
			continue
		}
		this.brokersMutex.Lock()
		this.brokers[bid] = broker
		this.brokersMutex.Unlock()
	}

	return nil
}

// deleteTenant delete tenant from iothub
func (this *Iothub) deleteTenant(tid string) error {
	return nil
}

// updateTenant update a tenant infromation
func (this *Iothub) updateTenant(tid string) error {
	return nil
}

// addBroker add broker to tenant
func (this *Iothub) addBroker(tid string, b *Broker) error {
	return nil
}

// deleteBroker delete broker from iothub
func (this *Iothub) deleteBroker(bid string) error {
	return nil
}

// startBroker start a tenant's broker
func (this *Iothub) startBroker(bid string) error {
	return nil
}

// startTenantBrokers start tenant's all broker
func (this *Iothub) startTenantBrokers(tid string) error {
	return nil
}

// stopBroker stop a broker
func (this *Iothub) stopBroker(bid string) error {
	return nil
}

// stopTenantBrokers stop tenant's all brokers
func (this *Iothub) stopTenantBrokers(tid string) error {
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
