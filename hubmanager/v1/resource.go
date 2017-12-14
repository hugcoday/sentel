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

import "github.com/labstack/echo"

// Check if an IoT hub name is available
func checkNameAvailability(ctx echo.Context) error {
	return nil
}

// Add a consumer group to an envent hub-compatible endpoint in an iot hub
func createEventHubConsumerGroup(ctx echo.Context) error {
	return nil
}

// Create or update the metadata of an IoT hub
func createOrUpdateMetadata(ctx echo.Context) error {
	return nil
}

// Delete IoT hub
func deleteIoTHub(ctx echo.Context) error {
	return nil
}

// Delete a consumer group from an event hub-compatible endpoint in an iot hub
func deleteEventHubConsumerGroup(ctx echo.Context) error {
	return nil
}

// Export all the device identities in the IoT hub identity registry
// to an azure storage blob container
func exportDevices(ctx echo.Context) error {
	return nil
}

// Get the non-security related metadata of an IoT hub
func getNonsecurityMetadata(ctx echo.Context) error {
	return nil
}

// Get a consumer group from the event hub-compatible device-to-cloud endpoint
// for an IoT hub
func getEventHubConsumerGroup(ctx echo.Context) error {
	return nil

}

// Get the details of a job from an IoT hub
func getJobDetail(ctx echo.Context) error {
	return nil
}

// Get a shared access policy by name from an IoT hub
func getKeyForKeyName(ctx echo.Context) error {
	return nil
}

// Get the quota metrics for an IoT hub
func getQuotaMetrics(ctx echo.Context) error {
	return nil
}

// Get the statistics from an IoT hub
func getStatics(ctx echo.Context) error {
	return nil
}

// Get the list of valid SKUs for an IoT hub
func getValidSkus(ctx echo.Context) error {
	return nil
}

// Get all IoT hubs in a resource group
func getAllHubsInGroup(ctx echo.Context) error {
	return nil
}

// Get all the IoT hubs in a subscription
func getAllHubs(ctx echo.Context) error {
	return nil
}

// Get a list of all the jobs in an IoT hub
func getAllJobsInIotHub(ctx echo.Context) error {
	return nil
}

// Get a list of the consumer groups in the event hub-compatible
// device-to-cloud endpoint in an IoT hub
func getListEventHubConsumerGroups(ctx echo.Context) error {
	return nil
}

// Get a list of all the jobs in an IoT hub
func getJobListFromHub(ctx echo.Context) error {
	return nil
}

// Get the security metadata for an IoT hub
func getSecurityMetadataFromHub(ctx echo.Context) error {
	return nil
}
