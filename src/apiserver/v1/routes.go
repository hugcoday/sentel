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

import "apiserver/api"

func newApi() *api.ApiManager {
	m := api.NewApiManager("v1")

	// Tenant Api
	m.RegisterApi("POST", "/tenants/:id", addTenant)
	m.RegisterApi("DELETE", "/tenants/:id", deleteTenant)
	m.RegisterApi("GET", "/tenants/:id", getTenant)

	// Product Api
	m.RegisterApi("POST", "/products/:id", addProduct)
	m.RegisterApi("DELETE", "/products/:id", deleteProduct)
	m.RegisterApi("GET", "/products/:id", getProduct)
	m.RegisterApi("GET", "/products/:id/devices", getProductDevices)

	// Device Api
	m.RegisterApi("GET", "/devices/:id", getDevices)
	m.RegisterApi("GET", "/devices/", getMultipleDevices)
	m.RegisterApi("DELETE", "/devices/:id", deleteDevices)
	m.RegisterApi("GET", "/statistics/devices", getRegistryStatistics)
	m.RegisterApi("GET", "/statistics/service", getServiceStatistics)
	m.RegisterApi("DELETE", "/devices/:id/commands", purgeCommandQueue)
	m.RegisterApi("PUT", "/devices/:id", putDevices)
	m.RegisterApi("POST", "/devices/query", queryDevices)

	// Device Twin Api
	m.RegisterApi("GET", "/twins/:id", getDeviceTwin)
	m.RegisterApi("POST", "/twins/:id/methods", invokeDeviceMethod)
	m.RegisterApi("PATCH", "/twins/:id", updateDeviceTwin)

	// Http Runtime Api
	m.RegisterApi("POST", "/devices/:id/messages/deviceBound/:etag/abandon",
		abandonDeviceBoundNotification)
	m.RegisterApi("DELETE", "/devices/:id/messages/devicesBound/:etag",
		completeDeviceBoundNotification)
	m.RegisterApi("POST", "/devices/:ideviceId/files",
		createFileUploadSasUri)
	m.RegisterApi("GET", "/devices/:id/message/deviceBound",
		receiveDeviceBoundNotification)
	m.RegisterApi("POST", "/devices/:deviceId/files/notifications",
		updateFileUploadStatus)
	m.RegisterApi("POST", "/devices/:id/messages/event", sendDeviceEvent)

	// Resource Api
	m.RegisterApi("POST",
		"/subscriptions/:subscriptionId/providers/Microsoft.Devices/checkNameAvailibility",
		checkNameAvailability)
	m.RegisterApi("PUT",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs/:resourceName/eventHubEndPoints
		/:eventHubEndpointName/ConsumerGroups/:name`,
		createEventHubConsumerGroup)
	m.RegisterApi("PUT",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs/:resourceName`,
		createOrUpdateMetadata)
	m.RegisterApi("DELETE",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs/:resourceName`,
		deleteIoTHub)
	m.RegisterApi("GET",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs/:resourceName`,
		getNonsecurityMetadata)

	m.RegisterApi("DELETE",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs/:resourceName
		/eventHubEndpoints/:eventHubEndpointName/ConsumerGroups/:name`,
		deleteEventHubConsumerGroup)

	m.RegisterApi("GET",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs/:resourceName
		/eventHubEndpoints/:eventHubEndpointName/ConsumerGroups/:name`,
		getEventHubConsumerGroup)
	m.RegisterApi("GET",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs/:resourceName/jobs/:jobId`,
		getJobDetail)

	m.RegisterApi("GET",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs/:resourceName/jobs`,
		getAllJobsInIotHub)
	m.RegisterApi("GET",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs/:resourceName/IotHubKyes/:keyName/listkyes`,
		getKeyForKeyName)

	m.RegisterApi("GET",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs/:resourceName/quatoMetrics`,
		getQuotaMetrics)
	m.RegisterApi("GET",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs/:resourceName/IotHubStats`,
		getStatics)
	m.RegisterApi("GET",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs/:resourceName/skus`,
		getValidSkus)
	m.RegisterApi("GET",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs`,
		getAllHubsInGroup)
	m.RegisterApi("GET",
		`/subscriptions/:subscriptionId/providers/Microsoft.Devices/IotHubs`,
		getAllHubs)

	m.RegisterApi("GET",
		`/subscriptions/:subscriptionId/resourceGroups/:resourceGroupName
		/providers/Microsoft.Devices/IotHubs/:resourceName/listkyes`,
		getSecurityMetadataFromHub)

	// Job Api
	m.RegisterApi("POST", "/jobs/v2/:jobid/cancel", cancelJob)
	m.RegisterApi("PUT", "/jobs/v2/:jobid", createJob)
	m.RegisterApi("GET", "/jobs/v2/:jobid", getJob)
	m.RegisterApi("GET", "/jobs/v2/query", queryJobs)

	return m
}

func init() {
	api.RegisterApiManager(newApi())
}
