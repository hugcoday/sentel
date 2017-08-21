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
	"apiserver/api"
	"database/sql"
	"fmt"
)

type Registry struct {
	conn *sql.DB
	ctx  api.ApiContext
}

func NewRegistry(ctx api.ApiContext) (*Registry, error) {
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
func (r *Registry) CheckTenantNameAvailable(t *api.Tenant) bool {
	return false
}

func (r *Registry) AddTenant(t *api.Tenant) error {
	return nil
}

func (r *Registry) DeleteTenant(t *api.Tenant) error {
	return nil
}

func (r *Registry) GetTenant(t *api.Tenant) error {
	return nil
}

// Product
func (r *Registry) CheckProductNameAvailable(p *api.Product) bool {
	return true
}

func (r *Registry) AddProduct(p *api.Product) error {
	return nil
}
func (r *Registry) DeleteProduct(p *api.Product) error {
	return nil
}

func (r *Registry) GetProduct(p *api.Product) error {
	return nil
}

func (r *Registry) GetProductDevices(p *api.Product) error {
	return nil
}

// Device
func (r *Registry) DeleteDevice(name string) error {
	return nil
}
