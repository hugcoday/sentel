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

package azure

import (
	"sentel/apiserver/api"

	"github.com/labstack/echo"
)

func NewApi() *api.ApiManager {
	m := api.NewApiManager("azure", nil)

	// Device APIs
	m.RegisterApi("GET", "/devices/:id", getDevices)
	m.RegisterApi("GET", "/devices/", getMultiplyDevices)
	m.RegisterApi("DELETE", "/devices/id:", deleteDevices)
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

	// Job API
	return m
}

// Delete the identify of a device from the identity registry
// of an IoT Hub
func deleteDevices(c echo.Context) error {
	return nil
}

// Retrieve a device from the identify registry of an IoT hub
func getDevices(e echo.Context) error {
	return nil
}

// Get the identifies of multiple devices from The IoT hub
func getMultiplyeDevices(c echo.Context) error {
	return nil
}

// Retrieves statistics about devices identities in the IoT hub's
// identify registry
func getRegistryStatistics(c echo.Context) error {
	return nil
}

// Retrieves services statisticsfor this IoT hubs's identity registry
func getServiceStatistics(c echo.Context) error {
	return nil
}

// Delete all the pending commands for this devices from the IoT hub
func purgeCommandQueue(c echo.Context) error {
	return nil
}

// Create or update the identity of a device in the identity registry of
// an IoT Hub
func putDevices(c echo.Context) error {
	return nil
}

// Query an IoT hub to retrieve information regarding device twis
// using a SQL-like language
func queryDevices(c echo.Context) error {
	return nil
}

// Twin Api

// Get a device twin
func getDeviceTwin(c echo.Context) error {
	return nil
}

// Invoce a direct method on device
func invokeDeviceMethod(c echo.Context) error {
	return nil
}

// Updates tags and desired properties of a device twin
func updateDeviceTwin(c echo.Context) error {
	return nil
}

// Http Runtime Api

// Abandon a cloud-to-device message
func abandonDeviceBoundNotification(c echo.Context) error {
	return nil
}

// Complete or rejects a cloud-to-device message
func completeDeviceBoundNotification(c echo.Context) error {
	return nil
}

// Retrive a storage SAS URI to upload a file
func createFileUploadSasUri(c echo.Context) error {
	return nil
}

// Retrive a cloud-to-device message
func receiveDeviceBoundNotification(c echo.Notification) error {
	return nil
}

// Send a device-to-cloud message
func sendDeviceEvent(c echo.Context) error {
	return nil
}

// Notify an IoT hub of a complete file upload
func updateFileUploadStatus(c echo.Context) error {
	return nil
}
