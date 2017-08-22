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

package util

type SentelError struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
	Err         error
}

func NewError(code int, desc string) SentelError {
	return SentelError{
		Code:        code,
		Description: desc,
	}
}

var (
	ErrorBadFormat                                  = NewError(400, "The body of the request is not valid")
	ErrorUnauthorized                               = NewError(401, "The authorization token cannot be validated")
	ErrorNotFound                                   = NewError(404, "The IoT hub instance or a device identity not exist")
	ErrorManyDevices                                = NewError(403, "The maximum number of device identities has been reached")
	ErrorPreconditionFailed                         = NewError(412, "Precondtion failed")
	ErrorTooManyRequest                             = NewError(429, "The IoT hub's identity registry operations are being throttled by the service")
	ErrorInternalError                              = NewError(500, "An internal error occured")
	ErrorInvalidErrorCode                           = NewError(600, "Invalid error code")
	ErrorInvalidApiVersion                          = NewError(1000, "Invalid api version")
	ErrorInvalidProtocolVersion                     = NewError(10001, "Invalid protocol version")
	ErrorDeviceInvalidResultCount                   = NewError(10002, "Invalid result count")
	ErrorDeviceInvalidOperation                     = NewError(10003, "Invalid operation")
	ErrorArgumentInvalid                            = NewError(10004, "Invalid argument")
	ErrorArgumentNull                               = NewError(10005, "Argument is null")
	ErrorIoTHubFormatError                          = NewError(10006, "IoT hub format error")
	ErrorDeviceStorageEntitySerializationError      = NewError(10007, "Device stoarge entity serialization error")
	ErrorBlobContainerValidationError               = NewError(10008, "Blob container validation error")
	ErrorImportWarningExistsError                   = NewError(10009, "Import warning exists error")
	ErrorInvalidSchemaVersion                       = NewError(10010, "Invalid schema version")
	ErrorDeviceDefinedMultipleTimes                 = NewError(10010, "Device defined multiple times")
	ErrorDeserializationError                       = NewError(10011, "Deserialization error")
	ErrorBulkRegistryOperationFailure               = NewError(10012, "Bulk registration operation failure")
	ErrorDefaultStorageEndpointNotConfigured        = NewError(10013, "Default storage endpoint not configured")
	ErrorInvalidFileUploadCorrelationId             = NewError(10014, "Invalid file upload correlation identifier")
	ErrorExpireFileUploadCorrelationId              = NewError(10015, "Expired file upload correlation identifier")
	ErrorInvalidStorageEndpoint                     = NewError(10016, "Invalid storage endpoint")
	ErrorInvalidMessagingEndpoing                   = NewError(10017, "Invalid message endpoint")
	ErrorInvalidFileUploadCompletionStatus          = NewError(10018, "Invalid file upload completion status")
	ErrorInvalidStorageEndpointOrBlob               = NewError(10019, "Invalid storage endpoint or blob")
	ErrorRequestCanceled                            = NewError(10020, "Request canceled")
	ErrorInvalidStorageEndpointProperty             = NewError(10021, "Invalid storage endpoint property")
	ErrorInvalidRouteTestInput                      = NewError(10022, "Invalid route test input")
	ErrorInvalidSourceOnRoute                       = NewError(10023, "Invalid source on route")
	ErrorIoTHubNotFound                             = NewError(10024, "IoT hub not found")
	ErrorIoTHubUnauthorizedAccess                   = NewError(10025, "IoT hub unauthorized access")
	ErrorIoTHubUnauthorized                         = NewError(10026, "IoT hub unauthorized")
	ErrorIoTHubSuspended                            = NewError(10027, "IoT hub suspended")
	ErrorIoTHubQuotaExceeded                        = NewError(10028, "IoT hub quota exceeded")
	ErrorIoTHubMaxCbsTokenExceeded                  = NewError(10029, "IoT hub max cbx exceeded")
	ErrorJobQuotaExceeded                           = NewError(10030, "Job quota exceeded")
	ErrorDeviceMaximumQueueDepthExceeded            = NewError(10031, "Device maximum queue depth exceeded")
	ErrorDeviceMaximumActiveFileUploadLimitExceeded = NewError(10032, "Device maximum active file upload limit excced")
	ErrorDeviceMaximumQueueSizeExceeded             = NewError(10033, "Device maximum queue size exceeded")
	ErrorDevicModelMaxPropertiesExceeded            = NewError(10034, "Device model maximum properties excedded")
	ErrorDeviceModelMaxIndexablePropertiesExceeded  = NewError(10035, "Device model maximum indexable properties exceeded")
	ErrorDeviceNotFound                             = NewError(10036, "Device not found")
	ErrorJobNotFound                                = NewError(10037, "Job not found")
	ErrorPartionNotFound                            = NewError(10038, "Partion not found")
	ErrorQuotaMetricNotFound                        = NewError(10039, "Quata metric not found")
	ErrorSystemPropertyNotFound                     = NewError(10040, "System property not found")
	ErrorAmqpAddressNotFound                        = NewError(10041, "AMQP address not found")
	ErrorDeviceNotOnline                            = NewError(10042, "Device not online")
	ErrorOperationNotAllowedInCurrentState          = NewError(10043, "Operation not allowded in current state")
	ErrorImportDeviceNotSupported                   = NewError(10045, "Import device not supported")
	ErrorBulkAddDevicesNotSupported                 = NewError(10046, "Bulk add devices not suppported")
	ErrorDeviceAlreadyExists                        = NewError(10046, "Device already exists")
	ErrorLinkCreationConflict                       = NewError(10047, "Link creation conflicti")
	ErrorModelAlreadyExists                         = NewError(10048, "Model alaready existes")
	ErrorDeviceLocked                               = NewError(10049, "Device locked")
	ErrorDeviceJobAlreadyExists                     = NewError(10050, "Device job already existes")
	ErrorDeviceMessageLockLost                      = NewError(10051, "Device message lock lost")
	ErrorJobRunPrecondtionFailed                    = NewError(10052, "Job run precondtion failed")
	ErrorDeviceMessageTooLarge                      = NewError(10053, "Device essage too large")
	ErrorTooManyDevices                             = NewError(10054, "Too many devices")
	ErrorIncompatibleDataType                       = NewError(10055, "Incompatible data type")
	ErrorThrottlingException                        = NewError(10056, "Throttling exception")
	ErrorThrottlingBacklogLimitExceeded             = NewError(10057, "Throttling backlog limit exceeded")
	ErrorThrottlingBackingTimeout                   = NewError(10058, "Throttling bacing timeout")
	ErrorThrottlingMaxActiveJobCountExceeded        = NewError(10059, "Throttling max active job count exceeded")
	ErrorServerError                                = NewError(10060, "Server error")
	ErrorJobCanceled                                = NewError(10061, "Job canceled")
	ErrorStatisticsRetrievalError                   = NewError(10062, "Statistics retrieval error")
	ErrorconnectionForcefullyClosed                 = NewError(10063, "Connection forcefully closed")
	ErrorInvalidBlobState                           = NewError(10064, "Innvalidi blob state")
	ErrorBackupTimeout                              = NewError(10065, "Backup timeout")
	ErrorStorageTimeout                             = NewError(10066, "Storage timeout")
	ErrorGenericTimeout                             = NewError(10067, "Generic timeout")
	ErrorInvalidThrottleParameter                   = NewError(10068, "Invalid throttle parameter")
	ErrorRetryAttemptsExhausted                     = NewError(10069, "Retry attempts exhausted")
	ErrorUnexpectedPropertyValue                    = NewError(10070, "Unexpected preperty value")
	ErrorServiceUnavailable                         = NewError(10071, "Service unavailable")
	ErrorIotHubFaillingOver                         = NewError(10072, "IoT hub failling over")
	ErrorConnectionUnavailable                      = NewError(10073, "Connection unavailable")
	ErrorDeviceUnavailable                          = NewError(10074, "Device unavailable")
	ErrorGatewayTimeout                             = NewError(10075, "Gateway timeout")
)

func (e *SentelError) Error() string {
	return e.Description
}
