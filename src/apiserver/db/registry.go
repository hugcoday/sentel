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
	"apiserver/base"
	"fmt"

	"github.com/go-xorm/xorm"
	"github.com/golang/glog"
)

type Registry struct {
	ctx base.ApiContext
	orm *xorm.Engine
}

func InitializeRegistry(c *base.ApiConfig) error {
	info := fmt.Sprintf("%s:%s@tcp(%s:%s)/registry",
		c.Registry.User, c.Registry.Password,
		c.Registry.Server, c.Registry.Port)
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

func NewRegistry(ctx base.ApiContext) (*Registry, error) {
	info := fmt.Sprintf("%s:%s@tcp(%s:%s)/registry",
		ctx.Config.Registry.User, ctx.Config.Registry.Password,
		ctx.Config.Registry.Server, ctx.Config.Registry.Port)
	orm, err := xorm.NewEngine("postgres", info)
	if err != nil {
		glog.Error("Create xorm engine failed:%s", err)
		return nil, err
	}
	return &Registry{orm: orm, ctx: ctx}, nil
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
	return true
}

// RegisterProduct register a product into registry
func (r *Registry) RegisterProduct(p *Product) error {
	return nil
}

// DeleteProduct delete a product from registry
func (r *Registry) DeleteProduct(id string) error {
	return nil
}

// GetProduct retrieve product detail information from registry
func (r *Registry) GetProduct(id string) (*Product, error) {
	return nil, nil
}

// GetProductDevices get product's device list
func (r *Registry) GetProductDevices(id string) ([]Device, error) {
	devices := []Device{}
	return devices, nil
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
