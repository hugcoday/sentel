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

package db

import (
	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
)

const (
	dbNameDevices  = "devices"
	dbNameProducts = "products"
	dbNameTenants  = "tenants"
)

// Registry is wraper of mongo database about for iot object
type Registry struct {
	config  core.Config
	session *mgo.Session
	db      *mgo.Database
}

// InitializeRegistry try to connect with background database
// to confirm wether it is normal
func InitializeRegistry(c core.Config) error {
	hosts := c.MustString("registry", "hosts")
	glog.Infof("Initializing registry:%s...", hosts)
	session, err := mgo.Dial(hosts)
	if err != nil {
		return err
	}
	session.Close()
	return nil
}

// NewRegistry create registry instance
func NewRegistry(c core.Config) (*Registry, error) {
	hosts := c.MustString("registry", "hosts")
	session, err := mgo.DialWithTimeout(hosts, 5*time.Second)
	if err != nil {
		glog.Infof("Failed to initialize registry:%s", err.Error())
		return nil, err
	}
	return &Registry{session: session, db: session.DB("registry"), config: c}, nil
}

// Release release registry rources and disconnect with background database
func (r *Registry) Release() {
	r.session.Close()
}

// Tenant

// CheckTenantNamveAvailable return true if name is available
func (r *Registry) CheckTenantNameAvailable(t *Tenant) bool {
	c := r.db.C(dbNameTenants)
	err := c.Find(bson.M{"Id": t.Id}).One(nil)
	return err != nil
}

// AddTenant insert a tenant into registry
func (r *Registry) RegisterTenant(t *Tenant) error {
	c := r.db.C(dbNameTenants)
	if err := c.Find(bson.M{"Id": t.Id}).One(nil); err == nil {
		return fmt.Errorf("Tenant %s already exist", t.Id)
	}
	return c.Insert(t, nil)
}

func (r *Registry) DeleteTenant(t *Tenant) error {
	c := r.db.C(dbNameTenants)
	return c.Remove(bson.M{"Id": t.Id})
}

func (r *Registry) GetTenant(t *Tenant) error {
	c := r.db.C(dbNameTenants)
	return c.Remove(bson.M{"Id": t.Id})
}

// Product
// CheckProductNameAvailable check wethere product name is available
func (r *Registry) CheckProductNameAvailable(p *Product) bool {
	c := r.db.C(dbNameProducts)
	err := c.Find(bson.M{"Id": p.Id}).One(nil)
	return err != nil
}

// RegisterProduct register a product into registry
func (r *Registry) RegisterProduct(p *Product) error {
	c := r.db.C(dbNameProducts)
	if err := c.Find(bson.M{"Name": p.Name}).One(nil); err == nil {
		return fmt.Errorf("product %s already exist", p.Name)
	}
	return c.Insert(p, nil)

}

// DeleteProduct delete a product from registry
func (r *Registry) DeleteProduct(id string) error {
	c := r.db.C(dbNameProducts)
	return c.Remove(bson.M{"Id": id})
}

// GetProduct retrieve product detail information from registry
func (r *Registry) GetProduct(id string) (*Product, error) {
	c := r.db.C(dbNameProducts)
	product := &Product{}
	err := c.Find(bson.M{"Id": id}).One(product)
	return product, err
}

// GetProductDevices get product's device list
func (r *Registry) GetProductDevices(id string) ([]Device, error) {
	c := r.db.C(dbNameDevices)
	iter := c.Find(bson.M{"ProductId": id}).Limit(1000).Iter()
	devices := []Device{}
	var device Device

	for iter.Next(&device) {
		devices = append(devices, device)
	}
	return devices, nil
}

// UpdateProduct update product detail information in registry
func (r *Registry) UpdateProduct(p *Product) error {
	c := r.db.C(dbNameProducts)
	return c.Update(bson.M{"Id": p.Id}, p)
}

// Device

// Registerevice add a new device into registry
func (r *Registry) RegisterDevice(dev *Device) error {
	c := r.db.C(dbNameDevices)
	if err := c.Find(bson.M{"Id": dev.Id}); err == nil { // found existed device
		return fmt.Errorf("device %s already exist", dev.Id)
	}
	return c.Insert(dev)
}

// GetDevice retrieve a device information from registry/
func (r *Registry) GetDevice(id string) (*Device, error) {
	c := r.db.C(dbNameDevices)
	device := &Device{}
	err := c.Find(bson.M{"Id": id}).One(device)
	return device, err
}

// BulkRegisterDevice add a lot of devices into registry
func (r *Registry) BulkRegisterDevice(devices []Device) error {
	for _, device := range devices {
		if err := r.RegisterDevice(&device); err != nil {
			return err
		}
	}
	return nil
}

// DeleteDevice delete a device from registry
func (r *Registry) DeleteDevice(id string) error {
	c := r.db.C(dbNameDevices)
	return c.Remove(bson.M{"Id": id})
}

// BulkDeleteDevice delete a lot of devices from registry
func (r *Registry) BulkDeleteDevice(devices []string) error {
	for _, id := range devices {
		if err := r.DeleteDevice(id); err != nil {
			return err
		}
	}
	return nil
}

// UpdateDevice update device information in registry
func (r *Registry) UpdateDevice(dev *Device) error {
	c := r.db.C(dbNameDevices)
	return c.Update(bson.M{"Id": dev.Id}, dev)
}

// BulkUpdateDevice update a lot of devices in registry
func (r *Registry) BulkUpdateDevice(devices []Device) error {
	for _, device := range devices {
		if err := r.UpdateDevice(&device); err != nil {
			return err
		}
	}
	return nil
}
