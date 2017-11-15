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

import (
	"github.com/cloustone/sentel/apiserver"
	"github.com/cloustone/sentel/apiserver/middleware"
	"github.com/cloustone/sentel/core"

	echo "github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

const APIHEAD = "api/v1/"

// v1apiManager mananage version 1 apis
type v1apiManager struct {
	version string
	config  core.Config
	echo    *echo.Echo
}

type apiContext struct {
	echo.Context
	config core.Config
}

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

// NewApiManager create api manager instance
func NewApiManager() apiserver.ApiManager {
	return &v1apiManager{
		echo:    echo.New(),
		version: "v1",
	}
}

// GetVersion return api's version
func (this *v1apiManager) GetVersion() string { return this.version }

// Run loop to wait api server to terminate
func (this *v1apiManager) Run() error {
	address := this.config.MustString("apiserver", "listen")
	return this.echo.Start(address)
}

// Initialize initialize api manager with configuration
func (this *v1apiManager) Initialize(c core.Config) error {
	this.config = c
	this.echo.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			cc := &apiContext{Context: e, config: c}
			return h(cc)
		}
	})
	// Initialize middleware
	this.echo.Use(middleware.ApiVersion(this.version))
	this.echo.Use(mw.KeyAuthWithConfig(middleware.DefaultKeyAuthConfig))
	this.echo.Use(mw.Logger())

	// Initialize api routes
	this.echo.POST("/api/v1/tenants/:id", addTenant)
	this.echo.DELETE("/api/v1/tenants/:id", deleteTenant)
	this.echo.GET("/api/v1/tenants/:id", getTenant)

	// Product Api
	this.echo.POST("api/v1/products/:id", registerProduct)
	this.echo.DELETE("/api/v1/products/:id", deleteProduct)
	this.echo.GET("/api/v1/products/:id", getProduct)
	this.echo.GET("/api/v1/products/:id/devices", getProductDevices)

	// Rule
	this.echo.POST("api/v1/rules/", addRule)
	this.echo.DELETE("/api/v1/rules/:id", deleteRule)
	this.echo.GET("/api/v1/rules/:id", getRule)
	this.echo.PATCH("/api/v1/rules/:id", updateRule)

	// Device Api
	this.echo.POST("api/v1/devices/:id", registerDevice)
	this.echo.GET("/devices/:id", getDevice)
	this.echo.DELETE("api/v1/devices/:id", deleteDevice)
	this.echo.PUT("api/v1/devices/:id", updateDevice)
	this.echo.DELETE("api/v1/devices/:id/commands", purgeCommandQueue)
	this.echo.GET("api/v1/devices/", getMultipleDevices)
	this.echo.POST("api/v1/devices/query", queryDevices)

	// Statics Api
	this.echo.GET("api/v1/statistics/devices", getRegistryStatistics)
	this.echo.GET("api/v1/statistics/service", getServiceStatistics)

	// Device Twin Api
	this.echo.GET("api/v1/twins/:id", getDeviceTwin)
	this.echo.POST("api/v1/twins/:id/methods", invokeDeviceMethod)
	this.echo.PATCH("api/v1/twins/:id", updateDeviceTwin)

	// Http Runtithis. Api
	this.echo.POST("api/v1/devices/:id/messages/deviceBound/:etag/abandon", abandonDeviceBoundNotification)
	this.echo.DELETE("api/v1/devices/:id/messages/devicesBound/:etag", completeDeviceBoundNotification)

	this.echo.POST("api/v1/devices/:ideviceId/files", createFileUploadSasUri)
	this.echo.GET("api/v1/devices/:id/message/deviceBound", receiveDeviceBoundNotification)
	this.echo.POST("api/v1/devices/:deviceId/files/notifications", updateFileUploadStatus)
	this.echo.POST("api/v1/devices/:id/messages/event", sendDeviceEvent)

	// Job Api
	this.echo.POST("api/v1/jobs/:jobid/cancel", cancelJob)
	this.echo.PUT("api/v1/jobs/:jobid", createJob)
	this.echo.GET("api/v1/jobs/:jobid", getJob)
	this.echo.GET("api/v1/jobs/query", queryJobs)

	return nil
}
