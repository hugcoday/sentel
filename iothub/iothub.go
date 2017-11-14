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
)

type Iothub struct {
	sync.Once
	config     core.Config
	mutex      sync.Mutex
	tenants    map[string]*Tenant
	brokers    map[string]*Broker
	clustermgr *clusterManager
}

type Tenant struct {
	id           string    `json:"tenantId"`
	createdAt    time.Time `json:"createdAt"`
	brokersCount int32     `json:"brokersCount"`
	brokers      map[string]*Broker
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
	session, err := mgo.DialWithTimeout(hosts, 5*time.Second)
	if err != nil {
		return err
	}
	session.Close()

	clustermgr, err := newClusterManager(c)
	if err != nil {
		return err
	}
	_iothub = &Iothub{
		config:     c,
		mutex:      sync.Mutex{},
		tenants:    make(map[string]*Tenant),
		brokers:    make(map[string]*Broker),
		clustermgr: clustermgr,
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
	this.mutex.Lock()
	defer this.mutex.Unlock()

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
	brokers, err := this.clustermgr.createBrokers(tid, tenant.brokersCount)
	if err != nil {
		// TODO should we update database status
		glog.Fatalf("Failed to created brokers for tenant(%s)", tid)
	}
	for _, broker := range brokers {
		this.brokers[broker.bid] = broker
	}

	return nil
}

// deleteTenant delete tenant from iothub
func (this *Iothub) deleteTenant(tid string) error {
	// check wether the tenant already exist
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if _, ok := this.tenants[tid]; ok {
		return fmt.Errorf("tenant(%s) already exist", tid)
	}

	if err := this.clustermgr.deleteBrokers(tid); err != nil {
		return err
	}

	// delete tenant from iothub
	tenant := this.tenants[tid]
	for bid, _ := range tenant.brokers {
		delete(this.brokers, bid)
	}
	delete(this.tenants, tid)

	// remove tenant from database
	hosts, err := this.config.String("iothub", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		return err
	}
	defer session.Close()

	c := session.DB("iothub").C("tenants")
	if err := c.Remove(bson.M{"tenantId": tid}); err != nil {
		return err
	}
	return nil
}

// updateTenant update a local tenant with database
func (this *Iothub) updateTenant(tid string) error {
	// check wether the tenant already exist
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if _, ok := this.tenants[tid]; ok {
		return fmt.Errorf("tenant(%s) already exist", tid)
	}

	lt := this.tenants[tid]

	// retrieve tenant from database
	tenant, err := this.getTenantFromDatabase(tid)
	if err != nil {
		return err
	}

	// rollback brokers if local tenant and tenant in database is not same
	if lt.brokersCount != tenant.brokersCount {
		return this.clustermgr.rollbackTenantBrokers(tenant)
	}

	return nil
}

// addBroker add broker to tenant
func (this *Iothub) addBroker(tid string, b *Broker) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if _, ok := this.tenants[tid]; !ok {
		return fmt.Errorf("invalid tenant(%s)", tid)
	}
	tenant := this.tenants[tid]
	tenant.brokers[b.bid] = b
	this.brokers[b.bid] = b

	return nil
}

// deleteBroker delete broker from iothub
func (this *Iothub) deleteBroker(bid string) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if _, ok := this.brokers[bid]; !ok {
		return fmt.Errorf("invalid broker(%s)", bid)
	}
	broker := this.brokers[bid]
	tenant := this.tenants[broker.tid]

	// stop broker when deleted
	if err := this.clustermgr.deleteBroker(broker); err != nil {
		return err
	}
	delete(this.brokers, bid)
	delete(tenant.brokers, bid)

	return nil
}

// startBroker start a tenant's broker
func (this *Iothub) startBroker(bid string) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if _, ok := this.brokers[bid]; !ok {
		return fmt.Errorf("invalid broker(%s)", bid)
	}

	broker := this.brokers[bid]
	if err := this.clustermgr.startBroker(broker); err != nil {
		return err
	}
	broker.status = BrokerStatusStarted
	return nil
}

// stopBroker stop a broker
func (this *Iothub) stopBroker(bid string) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if _, ok := this.brokers[bid]; !ok {
		return fmt.Errorf("invalid broker(%s)", bid)
	}

	broker := this.brokers[bid]
	if err := this.clustermgr.stopBroker(broker); err != nil {
		return err
	}
	broker.status = BrokerStatusStoped

	return nil
}

// startTenantBrokers start tenant's all broker
func (this *Iothub) startTenantBrokers(tid string) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if _, ok := this.tenants[tid]; !ok {
		return fmt.Errorf("invalid tenant(%s)", tid)
	}

	tenant := this.tenants[tid]
	for bid, broker := range tenant.brokers {
		if broker.status != BrokerStatusStarted {
			if err := this.clustermgr.startBroker(broker); err != nil {
				glog.Errorf("Failed to start broker(%s) for tenant(%s)", bid, tid)
				continue
			}
		}
	}
	return nil
}

// stopTenantBrokers stop tenant's all brokers
func (this *Iothub) stopTenantBrokers(tid string) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if _, ok := this.tenants[tid]; !ok {
		return fmt.Errorf("invalid tenant(%s)", tid)
	}

	tenant := this.tenants[tid]
	for bid, broker := range tenant.brokers {
		if broker.status != BrokerStatusStoped {
			if err := this.clustermgr.stopBroker(broker); err != nil {
				glog.Errorf("Failed to stop broker(%s) for tenant(%s)", bid, tid)
				continue
			}
		}
	}
	return nil
}

// getTenant retrieve a tenant by id
func (this *Iothub) getTenant(id string) *Tenant {
	return this.tenants[id]
}

// getBroker retrieve a broker by id
func (this *Iothub) getBroker(id string) *Broker {
	return this.brokers[id]
}

// setBrokerStatus set broker's status
func (this *Iothub) setBrokerStatus(bid string, status BrokerStatus) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if _, ok := this.brokers[bid]; !ok {
		return fmt.Errorf("Invalid broker(%s)", bid)
	}
	broker := this.brokers[bid]
	if broker.status != status {
		var err error
		switch status {
		case BrokerStatusStarted:
			err = this.clustermgr.startBroker(broker)
		case BrokerStatusStoped:
			err = this.clustermgr.stopBroker(broker)
		default:
			err = fmt.Errorf("Invalid broker status to set for broker(%s)", bid)
		}
		if err != nil {
			return err
		}
	}
	broker.status = status
	return nil
}

// getBrokerStatus retrieve broker's status
func (this *Iothub) getBrokerStatus(bid string) BrokerStatus {
	if _, ok := this.brokers[bid]; !ok {
		return BrokerStatusInvalid
	}
	return this.brokers[bid].status
}
