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

type requestBase struct {
	Format           string `json:"format"`
	AccessKeyId      string `json:"accessKeyID"`
	Signature        string `json:"signature"`
	Timestamp        string `json:"timestamp"`
	SignatureVersion string `json:"signatureVersion"`
	SignatueNonce    string `json:"signatureNonce"`
	RegionId         string `json:"regiionID"`
}

type response struct {
	RequestId string      `json:"requestID"`
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Result    interface{} `json:"result"`
}

func logInfo(ctx echo.Context, fmt string, args ...interface{}) {
}

func logDebug(ctx echo.Context, fmt string, args ...interface{}) {
}

func logFatal(ctx echo.Context, fmt string, args ...interface{}) {
}

func logError(ctx echo.Context, fmt string, args ...interface{}) {
}
