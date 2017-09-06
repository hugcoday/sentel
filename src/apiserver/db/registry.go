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
	"errors"
	"fmt"
	"libs"

	"github.com/go-xorm/xorm"
	"github.com/golang/glog"
)

type Registry struct {
	orm    *xorm.Engine
	config libs.Config
}

func InitializeRegistry(c libs.Config) error {
	var user, pwd, server, port string
	var err error

	if user, err = c.String("registry", "user"); err != nil {
		return err
	}
	if pwd, err = c.String("registry", "password"); err != nil {
		return err
	}
	if server, err = c.String("registry", "server"); err != nil {
		return err
	}
	if port, err = c.String("registry", "port"); err != nil {
		return err
	}
	info := fmt.Sprintf("%s:%s@tcp(%s:%s)/registry", user, pwd, server, port)
	orm, err := xorm.NewEngine("postgres", info)
	if err != nil {
		glog.Error("Create xorm engine failed:%s", err)
		return err
	}
	orm.ShowSQL(true)
	err = orm.CreateTables(&Tenant{}, &Product{}, &Device{})
	if err != nil {
		glog.Error("Created Registry tables failed:%s", err)
		return err
	}
	return nil
}

func NewRegistry(c libs.Config) (*Registry, error) {
	var user, pwd, server, port string
	var err error

	if user, err = c.String("registry", "user"); err != nil {
		return nil, err
	}
	if pwd, err = c.String("registry", "password"); err != nil {
		return nil, err
	}
	if server, err = c.String("registry", "server"); err != nil {
		return nil, err
	}
	if port, err = c.String("registry", "port"); err != nil {
		return nil, err
	}
	info := fmt.Sprintf("%s:%s@tcp(%s:%s)/registry", user, pwd, server, port)
	orm, err := xorm.NewEngine("postgres", info)
	if err != nil {
		glog.Error("Create xorm engine failed:%s", err)
		return nil, err
	}
	return &Registry{orm: orm, config: c}, nil
}

func (r *Registry) Release() {
}

// Tenant
func (r *Registry) CheckTenantNameAvailable(t *Tenant) bool {
	return false
}

func (r *Registry) AddTenant(t *Tenant) error {
	return nil
}

func (r *Registry) DeleteTenant(t *Tenant) error {
	return nil
}

func (r *Registry) GetTenant(t *Tenant) error {
	return nil
}

// Product
// CheckProductNameAvailable check wethere product name is available
func (r *Registry) CheckProductNameAvailable(p *Product) bool {
	has, _ := r.orm.Exist(p)
	return has
}

// RegisterProduct register a product into registry
func (r *Registry) RegisterProduct(p *Product) error {
	if has, err := r.orm.Exist(p); err == nil && has == true {
		return errors.New("product already exist")
	}
	_, err := r.orm.Insert(p)
	return err
}

// DeleteProduct delete a product from registry
func (r *Registry) DeleteProduct(id string) error {
	_, err := r.orm.Delete(&Product{Id: id})
	return err
}

// GetProduct retrieve product detail information from registry
func (r *Registry) GetProduct(id string) (*Product, error) {
	p := new(Product)
	p.Id = id
	_, err := r.orm.Get(p)
	return p, err
}

// GetProductDevices get product's device list
func (r *Registry) GetProductDevices(id string) ([]Device, error) {
	devices := []Device{}
	err := r.orm.Iterate(&Device{Id: id}, func(idx int, bean interface{}) error {
		dev := bean.(*Device)
		devices = append(devices, *dev)
		return nil
	})
	return devices, err
}

// UpdateProduct update product detail information in registry
func (r *Registry) UpdateProduct(p *Product) error {
	return nil
}

// Device

// Registerevice add a new device into registry
func (r *Registry) RegisterDevice(dev *Device) error {
	return nil
}

// GetDevice retrieve a device information from registry/
func (r *Registry) GetDevice(id string) (*Device, error) {
	return nil, nil
}

// BulkRegisterDevice add a lot of devices into registry
func (r *Registry) BulkRegisterDevice(devices []Device) error {
	return nil
}

// DeleteDevice delete a device from registry
func (r *Registry) DeleteDevice(id string) error {
	return nil
}

// BulkDeleteDevice delete a lot of devices from registry
func (r *Registry) BulkDeleteDevice(devices []string) error {
	return nil
}

// UpdateDevice update device information in registry
func (r *Registry) UpdateDevice(dev *Device) error {
	return nil
}

// BulkUpdateDevice update a lot of devices in registry
func (r *Registry) BulkUpdateDevice(devices []Device) error {
	return nil
}
