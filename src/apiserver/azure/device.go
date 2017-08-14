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
	"apiserver/api"
	"net/http"

	"github.com/labstack/echo"
)

// Retrieve a device from the identify registry of an IoT hub
func getDevices(c echo.Context) error {
	//	id := c.Param("id")
	//	cc := c.(*api.ApiContext)
	//	_, err := cc.Registry.GetDevice(c, id)
	return nil
}

// Delete the identify of a device from the identity registry
// of an IoT Hub
func deleteDevices(c echo.Context) error {
	id := c.Param("id")
	cc := c.(*api.ApiContext)
	cc.Registry.DeleteDevice(c, id)
	return c.NoContent(http.StatusNoContent)
}

// Get the identifies of multiple devices from The IoT hub
func getMultipleDevices(c echo.Context) error {
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
