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

package v1

import (
	"apiserver/base"
	"apiserver/middleware"
)

func newApi() *base.ApiManager {
	m := base.NewApiManager("v1")

	// Tenant Api
	m.RegisterApi("POST", "/tenants/:id", addTenant)
	m.RegisterApi("DELETE", "/tenants/:id", deleteTenant, middleware.DefaultKeyAuth)
	m.RegisterApi("GET", "/tenants/:id", getTenant, middleware.DefaultKeyAuth)

	// Product Api
	m.RegisterApi("POST", "/products/:id", addProduct, middleware.DefaultKeyAuth)
	m.RegisterApi("DELETE", "/products/:id", deleteProduct, middleware.DefaultKeyAuth)
	m.RegisterApi("GET", "/products/:id", getProduct, middleware.DefaultKeyAuth)
	m.RegisterApi("GET", "/products/:id/devices", getProductDevices, middleware.DefaultKeyAuth)

	// Device Api
	m.RegisterApi("GET", "/devices/:id", getDevices, middleware.DefaultKeyAuth)
	m.RegisterApi("GET", "/devices/", getMultipleDevices, middleware.DefaultKeyAuth)
	m.RegisterApi("DELETE", "/devices/:id", deleteDevices, middleware.DefaultKeyAuth)
	m.RegisterApi("GET", "/statistics/devices", getRegistryStatistics, middleware.DefaultKeyAuth)
	m.RegisterApi("GET", "/statistics/service", getServiceStatistics, middleware.DefaultKeyAuth)
	m.RegisterApi("DELETE", "/devices/:id/commands", purgeCommandQueue, middleware.DefaultKeyAuth)
	m.RegisterApi("PUT", "/devices/:id", putDevices, middleware.DefaultKeyAuth)
	m.RegisterApi("POST", "/devices/query", queryDevices, middleware.DefaultKeyAuth)

	// Device Twin Api
	m.RegisterApi("GET", "/twins/:id", getDeviceTwin, middleware.DefaultKeyAuth)
	m.RegisterApi("POST", "/twins/:id/methods", invokeDeviceMethod, middleware.DefaultKeyAuth)
	m.RegisterApi("PATCH", "/twins/:id", updateDeviceTwin, middleware.DefaultKeyAuth)

	// Http Runtime Api
	m.RegisterApi("POST", "/devices/:id/messages/deviceBound/:etag/abandon",
		abandonDeviceBoundNotification, middleware.DefaultKeyAuth)
	m.RegisterApi("DELETE", "/devices/:id/messages/devicesBound/:etag",
		completeDeviceBoundNotification, middleware.DefaultKeyAuth)
	m.RegisterApi("POST", "/devices/:ideviceId/files",
		createFileUploadSasUri, middleware.DefaultKeyAuth)
	m.RegisterApi("GET", "/devices/:id/message/deviceBound",
		receiveDeviceBoundNotification, middleware.DefaultKeyAuth)
	m.RegisterApi("POST", "/devices/:deviceId/files/notifications",
		updateFileUploadStatus, middleware.DefaultKeyAuth)
	m.RegisterApi("POST", "/devices/:id/messages/event", sendDeviceEvent, middleware.DefaultKeyAuth)

	// Resource Api
	m.RegisterApi("POST",
		"/subscriptions/:subscriptionId/providers/Microsoft.Devices/checkNameAvailibility",
		checkNameAvailability, middleware.DefaultKeyAuth)
	m.RegisterApi("PUT",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs/:resourceName/eventHubEndPoints
		/:eventHubEndpointName/ConsumerGroups/:name`,
		createEventHubConsumerGroup, middleware.DefaultKeyAuth)
	m.RegisterApi("PUT",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs/:resourceName`,
		createOrUpdateMetadata, middleware.DefaultKeyAuth)
	m.RegisterApi("DELETE",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs/:resourceName`,
		deleteIoTHub, middleware.DefaultKeyAuth)
	m.RegisterApi("GET",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs/:resourceName`,
		getNonsecurityMetadata, middleware.DefaultKeyAuth)

	m.RegisterApi("DELETE",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs/:resourceName
		/eventHubEndpoints/:eventHubEndpointName/ConsumerGroups/:name`,
		deleteEventHubConsumerGroup, middleware.DefaultKeyAuth)

	m.RegisterApi("GET",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs/:resourceName
		/eventHubEndpoints/:eventHubEndpointName/ConsumerGroups/:name`,
		getEventHubConsumerGroup, middleware.DefaultKeyAuth)
	m.RegisterApi("GET",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs/:resourceName/jobs/:jobId`,
		getJobDetail, middleware.DefaultKeyAuth)

	m.RegisterApi("GET",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs/:resourceName/jobs`,
		getAllJobsInIotHub, middleware.DefaultKeyAuth)
	m.RegisterApi("GET",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs/:resourceName/IotHubKyes/:keyName/listkyes`,
		getKeyForKeyName, middleware.DefaultKeyAuth)

	m.RegisterApi("GET",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs/:resourceName/quatoMetrics`,
		getQuotaMetrics, middleware.DefaultKeyAuth)
	m.RegisterApi("GET",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs/:resourceName/IotHubStats`,
		getStatics, middleware.DefaultKeyAuth)
	m.RegisterApi("GET",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs/:resourceName/skus`,
		getValidSkus, middleware.DefaultKeyAuth)
	m.RegisterApi("GET",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs`,
		getAllHubsInGroup, middleware.DefaultKeyAuth)
	m.RegisterApi("GET",
		`/subscriptions/:subscriptionId/providers/Microsoft.Devices/IotHubs`,
		getAllHubs, middleware.DefaultKeyAuth)

	m.RegisterApi("GET",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs/:resourceName/listkyes`,
		getSecurityMetadataFromHub, middleware.DefaultKeyAuth)

	// Job Api
	m.RegisterApi("POST", "/jobs/v2/:jobid/cancel", cancelJob, middleware.DefaultKeyAuth)
	m.RegisterApi("PUT", "/jobs/v2/:jobid", createJob, middleware.DefaultKeyAuth)
	m.RegisterApi("GET", "/jobs/v2/:jobid", getJob, middleware.DefaultKeyAuth)
	m.RegisterApi("GET", "/jobs/v2/query", queryJobs, middleware.DefaultKeyAuth)

	return m
}

func init() {
	base.RegisterApiManager(newApi())
}
