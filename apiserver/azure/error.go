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

import "errors"

type AzureError struct {
	Code        int
	Description string
	Err         error
}

func NewError(code int, desc string) *AzureError {
	return &AzureError{
		Code:        code,
		Description: desc,
		Err:         errors.New(desc),
	}
}

var (
	AzureErrorBadFormat          = NewError(400, "The body of the request is not valid")
	AzureErrorUnauthorized       = NewError(401, "The authorization token cannot be validated")
	AzureErrorNotFound           = NewError(404, "The IoT hub instance or a device identity not exist")
	AzureErrorManyDevices        = NewError(403, "The maximum number of device identities has been reached")
	AzureErrorPreconditionFailed = NewError(412, "The etag in the request does not match the etag of the existing resource")
	AzureErrorTooManyRequest     = NewError(429, "The IoT hub's identity registry operations are being throttled by the service")
	AzureErrorInternalError      = NewError(500, "An internal error occured")
)
