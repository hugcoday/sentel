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

// Http Runtime Api

// Abandon a cloud-to-device message
func abandonDeviceBoundNotification(ctx echo.Context) error {
	return nil
}

// Complete or rejects a cloud-to-device message
func completeDeviceBoundNotification(ctx echo.Context) error {
	return nil
}

// Retrive a storage SAS URI to upload a file
func createFileUploadSasUri(ctx echo.Context) error {
	return nil
}

// Retrive a cloud-to-device message
func receiveDeviceBoundNotification(ctx echo.Context) error {
	return nil
}

// Send a device-to-cloud message
func sendDeviceEvent(ctx echo.Context) error {
	return nil
}

// Notify an IoT hub of a complete file upload
func updateFileUploadStatus(ctx echo.Context) error {
	return nil
}
