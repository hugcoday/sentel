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

package api

import "errors"

type SentelError struct {
	Code        int
	Description string
	Err         error
}

func NewError(code int, desc string) SentelError {
	return SentelError{
		Code:        code,
		Description: desc,
		Err:         errors.New(desc),
	}
}

var (
	ErrorBadFormat          = NewError(400, "The body of the request is not valid")
	ErrorUnauthorized       = NewError(401, "The authorization token cannot be validated")
	ErrorNotFound           = NewError(404, "The IoT hub instance or a device identity not exist")
	ErrorManyDevices        = NewError(403, "The maximum number of device identities has been reached")
	ErrorPreconditionFailed = NewError(412, "The etag in the request does not match the etag of the existing resource")
	ErrorTooManyRequest     = NewError(429, "The IoT hub's identity registry operations are being throttled by the service")
	ErrorInternalError      = NewError(500, "An internal error occured")

	ErrorInvalidErrorCode                      = NewError(1000, "Invalid error code")
	ErrorInvalidProtocolVersion                = NewError(10001, "Invalid protocol version")
	ErrorDeviceInvalidResultCount              = NewError(10002, "Invalid result count")
	ErrorDeviceInvalidOperation                = NewError(10003, "Invalid operation")
	ErrorArgumentInvalid                       = NewError(10004, "Invalid argument")
	ErrorArgumentNull                          = NewError(10005, "Argument is null")
	ErrorIoTHubFormatError                     = NewError(10006, "IoT hub format error")
	ErrorDeviceStorageEntitySerializationError = NewError(10007, "Device stoarge entity serialization error")
	ErrorBlobContainerValidationError          = NewError(10008, "Blob container validation error")
	ErrorImportWarningExistsError              = NewError(10009, "Import warning exists error")
	ErrorInvalidSchemaVersion                  = NewError(10010, "Invalid schema version")
	ErrorDeviceDefinedMultipleTimes            = NewError(10010, "Device defined multiple times")
	ErrorDeserializationError                  = NewError(10011, "Deserialization error")
	ErrorBulkRegistryOperationFailure          = NewError(10012, "Bulk registration operation failure")
	ErrorDefaultStorageEndpointNotConfigured   = NewError(10013, "Default storage endpoint not configured")
	ErrorInvalidFileUploadCorrelationId        = NewError(10014, "Invalid file upload correlation identifier")
	ErrorExpireFileUploadCorrelationId         = NewError(10015, "Expired file upload correlation identifier")
	ErrorInvalidStorageEndpoint                = NewError(10016, "Invalid storage endpoint")
	ErrorInvalidMessagingEndpoing              = NewError(10017, "Invalid message endpoint")
)
