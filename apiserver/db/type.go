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

import "time"

// Tenant
type Tenant struct {
	Id   string
	Name string
}

// Product
type Product struct {
	Id           string
	Name         string
	Description  string
	TimeCreated  time.Time
	TimeModified time.Time
	CategoryId   string
}

// Device
type Device struct {
	Id           string
	Name         string
	ProductId    string
	ProductKey   string
	DeviceStatus string
	DeviceSecret string
	TimeCreated  time.Time
	TimeModified time.Time
}

// Rule
type Rule struct {
	Id           string
	Name         string
	ProductId    string
	TimeCreated  time.Time
	TimeModified time.Time
	Status       string
	Method       string
	Target       string
}
