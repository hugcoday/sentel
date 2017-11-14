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
	Config core.Config
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
			cc := &apiContext{Context: e, Config: c}
			return h(cc)
		}
	})
	// Initialize middleware
	this.echo.Use(middleware.ApiVersion(this.version))
	// this.ech.Use(middleware.KeyAuthWithConfig(middleware.DefaultKeyAuthConfig))
	// this.ech.Use(middleware.Logger())
	// this.ech.Use(middleware.Recover())

	// Initialize api routes
	this.echo.POST("/api/v1/tenants/:id", addTenant, middleware.DefaultKeyAuth)
	this.echo.DELETE("/api/v1/tenants/:id", deleteTenant, middleware.DefaultKeyAuth)
	this.echo.GET("/api/v1/tenants/:id", getTenant, middleware.DefaultKeyAuth)

	// Product Api
	this.echo.POST("api/v1/products/:id", registerProduct, middleware.DefaultKeyAuth)
	this.echo.DELETE("/api/v1/products/:id", deleteProduct, middleware.DefaultKeyAuth)
	this.echo.GET("/api/v1/products/:id", getProduct, middleware.DefaultKeyAuth)
	this.echo.GET("/api/v1/products/:id/devices", getProductDevices, middleware.DefaultKeyAuth)

	// Device Api
	this.echo.POST("api/v1/devices/:id", registerDevice, middleware.DefaultKeyAuth)
	this.echo.GET("/devices/:id", getDevice, middleware.DefaultKeyAuth)
	this.echo.DELETE("api/v1/devices/:id", deleteDevice, middleware.DefaultKeyAuth)
	this.echo.PUT("api/v1/devices/:id", updateDevice, middleware.DefaultKeyAuth)
	this.echo.DELETE("api/v1/devices/:id/commands", purgeCommandQueue, middleware.DefaultKeyAuth)
	this.echo.GET("api/v1/devices/", getMultipleDevices, middleware.DefaultKeyAuth)
	this.echo.POST("api/v1/devices/query", queryDevices, middleware.DefaultKeyAuth)

	// Statics Api
	this.echo.GET("api/v1/statistics/devices", getRegistryStatistics, middleware.DefaultKeyAuth)
	this.echo.GET("api/v1/statistics/service", getServiceStatistics, middleware.DefaultKeyAuth)

	// Device Twin Api
	this.echo.GET("api/v1/twins/:id", getDeviceTwin, middleware.DefaultKeyAuth)
	this.echo.POST("api/v1/twins/:id/methods", invokeDeviceMethod, middleware.DefaultKeyAuth)
	this.echo.PATCH("api/v1/twins/:id", updateDeviceTwin, middleware.DefaultKeyAuth)

	// Http Runtithis. Api
	this.echo.POST("api/v1/devices/:id/messages/deviceBound/:etag/abandon",
		abandonDeviceBoundNotification, middleware.DefaultKeyAuth)
	this.echo.DELETE("api/v1/devices/:id/messages/devicesBound/:etag",
		completeDeviceBoundNotification, middleware.DefaultKeyAuth)

	this.echo.POST("api/v1/devices/:ideviceId/files",
		createFileUploadSasUri, middleware.DefaultKeyAuth)
	this.echo.GET("api/v1/devices/:id/message/deviceBound",
		receiveDeviceBoundNotification, middleware.DefaultKeyAuth)
	this.echo.POST("api/v1/devices/:deviceId/files/notifications",
		updateFileUploadStatus, middleware.DefaultKeyAuth)
	this.echo.POST("api/v1/devices/:id/messages/event", sendDeviceEvent, middleware.DefaultKeyAuth)

	// Job Api
	this.echo.POST("api/v1/jobs/:jobid/cancel", cancelJob, middleware.DefaultKeyAuth)
	this.echo.PUT("api/v1/jobs/:jobid", createJob, middleware.DefaultKeyAuth)
	this.echo.GET("api/v1/jobs/:jobid", getJob, middleware.DefaultKeyAuth)
	this.echo.GET("api/v1/jobs/query", queryJobs, middleware.DefaultKeyAuth)

	return nil
}
