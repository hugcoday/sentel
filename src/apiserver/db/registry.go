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
	"apiserver/types"
	"database/sql"
	"fmt"
)

type Registry struct {
	conn *sql.DB
	ctx  base.ApiContext
}

func InitializeRegistryStore(c base.ApiConfig) error {
	return nil
}

func NewRegistry(ctx base.ApiContext) (*Registry, error) {
	info := fmt.Sprintf("%s:%s@tcp(%s:%s)/registry",
		ctx.Config.Registry.User, ctx.Config.Registry.Password,
		ctx.Config.Registry.Server, ctx.Config.Registry.Port)

	conn, err := sql.Open("postgres", info)
	if err != nil {
		return nil, err
	}

	return &Registry{conn: conn, ctx: ctx}, nil
}

func (r *Registry) Release() {
	r.conn.Close()
}

// Tenant
func (r *Registry) CheckTenantNameAvailable(t *types.Tenant) bool {
	return false
}

func (r *Registry) AddTenant(t *types.Tenant) error {
	return nil
}

func (r *Registry) DeleteTenant(t *types.Tenant) error {
	return nil
}

func (r *Registry) GetTenant(t *types.Tenant) error {
	return nil
}

// Product
func (r *Registry) CheckProductNameAvailable(p *types.Product) bool {
	return true
}

func (r *Registry) AddProduct(p *types.Product) error {
	return nil
}
func (r *Registry) DeleteProduct(p *types.Product) error {
	return nil
}

func (r *Registry) GetProduct(p *types.Product) error {
	return nil
}

func (r *Registry) GetProductDevices(p *types.Product) error {
	return nil
}

// Device
func (r *Registry) DeleteDevice(name string) error {
	return nil
}
