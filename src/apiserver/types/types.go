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

package types

type CommonParameter struct {
}

// Tenant
type Tenant struct {
	CommonParameter
	id   string
	name string
}

// Product
type Product struct {
	CommonParameter
	id   string
	name string
}

// Device
type DeviceRegistryOperationError struct {
	deviceId string
	//  errorCode ErrorCode
	errorStatus string
}

type BulkRegistryOperationResult struct {
	isSuccessful bool
	errors       []DeviceRegistryOperationError
}

type Device struct {
	deviceId        string
	generationId    string
	etag            string
	connectionState int
}
