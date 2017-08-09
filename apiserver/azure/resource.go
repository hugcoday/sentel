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

import "github.com/labstack/echo"

// Check if an IoT hub name is available
func checkNameAvailability(c echo.Context) error {
	return nil
}

// Add a consumer group to an envent hub-compatible endpoint in an iot hub
func createEventHubConsumerGroup(c echo.Context) error {
	return nil
}

// Create or update the metadata of an IoT hub
func createOrUpdateMetadata(c echo.Context) error {
	return nil
}

// Delete IoT hub
func deleteIoTHub(c echo.Context) error {
	return nil
}

// Delete a consumer group from an event hub-compatible endpoint in an iot hub
func deleteEventHubConsumerGroup(c echo.Context) error {
	return nil
}

// Export all the device identities in the IoT hub identity registry
// to an azure storage blob container
func exportDevices(c echo.Context) error {
	return nil
}

// Get the non-security related metadata of an IoT hub
func getNonsecurityMetadata(c echo.Context) error {
	return nil
}

// Get a consumer group from the event hub-compatible device-to-cloud endpoint
// for an IoT hub
func getEventHubConsumerGroup(c echo.Context) error {
	return nil

}

// Get the details of a job from an IoT hub
func getJob(c echo.Context) error {
	return nil
}

// Get a shared access policy by name from an IoT hub
func getKeyForKeyName(c echo.Context) error {
	return nil
}

// Get the quota metrics for an IoT hub
func getQuotaMetrics(c echo.Context) error {
	return nil
}

// Get the statistics from an IoT hub
func getStatics(c echo.Context) error {
	return nil
}

// Get the list of valid SKUs for an IoT hub
func getValidSkus(c echo.Context) error {
	return nil
}

// Get all IoT hubs in a resource group
func getAllHubsInGroup(c echo.Context) error {
	return nil
}

// Get all the IoT hubs in a subscription
func getAllHubs(c echo.Context) error {
	return nil
}

// Get a list of all the jobs in an IoT hub
func getAllJobsInIotHub(c echo.Context) error {
	return nil
}

// Get a list of the consumer groups in the event hub-compatible
// device-to-cloud endpoint in an IoT hub
func getListEventHubConsumerGroups(c echo.Context) error {
	return nil
}

// Get a list of all the jobs in an IoT hub
func getJobListFromHub(c echo.Context) error {
	return nil
}

// Get the security metadata for an IoT hub
func getSecurityMetadataFromHub(c echo.Context) error {
	return nil
}
