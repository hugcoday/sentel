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
	"apiserver/base"
	"apiserver/db"
	"net/http"
	"time"

	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

// Device internal definition
type registerDeviceRequest struct {
	RequestCommonParameter
	ProductKey string `json:"productKey"`
	DeviceName string `json:"productName"`
}

type registerDeviceResponse struct {
	ResponseCommonParameter
	DeviceId     string `json:"deviceId"`
	DeviceName   string `json:"deviceName"`
	DeviceSecret string `json:deviceSecret"`
	DeviceStatus string `json:deviceStatus"`
	ProductKey   string `json:"productKey"`
	TimeCreated  string `json:"timeCreated"`
}

// RegisterDevice register a new device in IoT hub
func registerDevice(c echo.Context) error {
	logInfo(c, "registerDevice(%s) called", c.Param("id"))
	// Get product
	req := new(registerDeviceRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	// Connect with registry
	r, err := db.NewRegistry(*c.(*base.ApiContext))
	if err != nil {
		logFatal(c, "Registry connection failed")
		return err
	}
	defer r.Release()

	// Insert device into registry, the created product
	// will be modified to retrieve specific information sucha as
	// product.id and creation time
	dp := db.Device{
		Name:        req.DeviceName,
		ProductKey:  req.ProductKey,
		TimeCreated: time.Now().String(),
	}
	rcp := ResponseCommonParameter{
		RequestId:    uuid.NewV4().String(),
		Success:      true,
		ErrorMessage: "",
	}
	err = r.RegisterDevice(&dp)
	if err != nil {
		rcp.Success = false
		rcp.ErrorMessage = err.Error()
		return c.JSON(http.StatusOK, rcp)
	}
	rsp := &registerDeviceResponse{
		ResponseCommonParameter: rcp,
		DeviceId:                dp.Id,
		DeviceName:              dp.Name,
		ProductKey:              dp.ProductKey,
		DeviceSecret:            dp.DeviceSecret,
		TimeCreated:             dp.TimeCreated,
	}
	return c.JSON(http.StatusOK, rsp)

}

// Retrieve a device from the identify registry of an IoT hub
func getDevice(c echo.Context) error {
	logInfo(c, "getDevice(%s) called", c.Param("id"))
	// Get product
	req := new(registerDeviceRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	// Connect with registry
	registry, err := db.NewRegistry(*c.(*base.ApiContext))
	if err != nil {
		logFatal(c, "Registry connection failed")
		return err
	}
	defer registry.Release()

	rcp := ResponseCommonParameter{
		RequestId:    uuid.NewV4().String(),
		Success:      true,
		ErrorMessage: "",
	}
	// Get device into registry, the created product
	dev, err := registry.GetDevice(c.Param("id"))
	if err != nil {
		logError(c, "Registry.GetDevice(%s) failed:%v", c.Param("id"), err)
		rcp.Success = false
		rcp.ErrorMessage = err.Error()
		return c.JSON(http.StatusOK, rcp)
	}
	rsp := &registerDeviceResponse{
		ResponseCommonParameter: rcp,
		DeviceId:                dev.Id,
		DeviceName:              dev.Name,
		ProductKey:              dev.ProductKey,
		DeviceSecret:            dev.DeviceSecret,
		TimeCreated:             dev.TimeCreated,
		DeviceStatus:            dev.DeviceStatus,
	}
	return c.JSON(http.StatusOK, rsp)
}

// Delete the identify of a device from the identity registry
// of an IoT Hub
func deleteDevice(c echo.Context) error {
	logInfo(c, "deleteDevice(%s) called", c.Param("id"))
	// Get product
	req := new(registerDeviceRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	// Connect with registry
	registry, err := db.NewRegistry(*c.(*base.ApiContext))
	if err != nil {
		logFatal(c, "Registry connection failed")
		return err
	}
	defer registry.Release()

	rcp := ResponseCommonParameter{
		RequestId:    uuid.NewV4().String(),
		Success:      true,
		ErrorMessage: "",
	}
	// Get device into registry, the created product
	if err := registry.DeleteDevice(c.Param("id")); err != nil {
		logError(c, "Registry.DeleteDevice(%s) failed:%v", c.Param("id"), err)
		rcp.Success = false
		rcp.ErrorMessage = err.Error()
	}
	return c.JSON(http.StatusOK, rcp)
}

type updateDeviceRequest struct {
	RequestCommonParameter
	DeviceId     string `json:"deviceId"`
	DeviceName   string `json:"productName"`
	ProductKey   string `json:"productKey"`
	DeviceSecret string `json:"deviceSecrt"`
	DeviceStatus string `json:"deviceStatus"`
}
type updateDeviceResponse struct {
	ResponseCommonParameter
	DeviceId     string `json:"deviceId"`
	DeviceName   string `json:"deviceName"`
	DeviceSecret string `json:deviceSecret"`
	DeviceStatus string `json:deviceStatus"`
	ProductKey   string `json:"productKey"`
	TimeCreated  string `json:"timeCreated"`
	TimeModified string `json:"timeModified"`
}

// updateDevice update the identity of a device in the identity registry of
// an IoT Hub
func updateDevice(c echo.Context) error {
	logInfo(c, "updateDevice(%s) called", c.Param("id"))
	// Get product
	req := new(updateDeviceRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	// Connect with registry
	r, err := db.NewRegistry(*c.(*base.ApiContext))
	if err != nil {
		logFatal(c, "Registry connection failed")
		return err
	}
	defer r.Release()

	// Insert device into registry, the created product
	// will be modified to retrieve specific information sucha as
	// product.id and creation time
	dp := db.Device{
		Id:           c.Param("id"),
		Name:         req.DeviceName,
		ProductKey:   req.ProductKey,
		DeviceStatus: req.DeviceStatus,
		DeviceSecret: req.DeviceSecret,
		TimeModified: time.Now().String(),
	}
	rcp := ResponseCommonParameter{
		RequestId:    uuid.NewV4().String(),
		Success:      true,
		ErrorMessage: "",
	}
	err = r.UpdateDevice(&dp)
	if err != nil {
		rcp.Success = false
		rcp.ErrorMessage = err.Error()
		return c.JSON(http.StatusOK, rcp)
	}
	rsp := &updateDeviceResponse{
		ResponseCommonParameter: rcp,
		DeviceId:                dp.Id,
		DeviceName:              dp.Name,
		ProductKey:              dp.ProductKey,
		DeviceSecret:            dp.DeviceSecret,
		TimeCreated:             dp.TimeCreated,
		TimeModified:            dp.TimeModified,
	}
	return c.JSON(http.StatusOK, rsp)
}

// Delete all the pending commands for this devices from the IoT hub
func purgeCommandQueue(c echo.Context) error {
	return nil
}

// Query an IoT hub to retrieve information regarding device twis
// using a SQL-like language
func queryDevices(c echo.Context) error {
	return nil
}

// Get the identifies of multiple devices from The IoT hub
func getMultipleDevices(c echo.Context) error {
	return nil
}
