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
}

func NewRegistry(ctx api.ApiContext) (*Registry, error) {
	info := fmt.Sprintf("%s:%s@tcp(%s:%s)/registry",
		ctx.Config.Registry.User, ctx.Config.Registry.Password,
		ctx.Config.Registry.Server, ctx.Config.Registry.Port)

	conn, err := sql.Open("postgres", info)
	if err != nil {
		return nil, err
	}

	return &Registry{conn: conn}, nil
}

func (r *Registry) Release(ctx api.ApiContext) {
	r.conn.Close()
}

func (r *Registry) DeleteDevice(ctx api.ApiContext, name string) error {
	return nil
}
