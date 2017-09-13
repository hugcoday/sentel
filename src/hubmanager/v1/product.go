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
	"hubmanager/base"
	"hubmanager/db"
	"hubmanager/util"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
)

// Product internal definition
type product struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	TimeCreated  string `json:"timeCreated"`
	TimeModified string `json:"timeModified"`
	CategoryId   string `json:"categoryId"`
}
type productAddRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type productAddResponse struct {
	ResponseCommonParameter
	product
}

func registerProduct(c echo.Context) error {
	logInfo(c, "addProduct called")
	// Get product
	req := new(productAddRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	// Connect with registry
	config := c.(*base.ApiContext).Config
	r, err := db.NewRegistry(config)
	if err != nil {
		logFatal(c, "Registry connection failed")
		return err
	}
	defer r.Release()

	// Insert product into registry, the created product
	// will be modified to retrieve specific information sucha as
	// product.id and creation time
	dp := db.Product{
		Name:        req.Name,
		Description: req.Description,
		TimeCreated: time.Now().String(),
	}
	rcp := ResponseCommonParameter{
		RequestId:    uuid.NewV4().String(),
		Success:      true,
		ErrorMessage: "",
	}
	err = r.RegisterProduct(&dp)
	if err != nil {
		rcp.Success = false
		rcp.ErrorMessage = err.Error()
		return c.JSON(http.StatusOK, rcp)
	}

	// Notify kafka
	base.AsyncProduceMessage(c, util.TopicNameProduct,
		&util.ProductTopic{
			ProductId:   dp.Id,
			ProductName: dp.Name,
			Action:      util.ObjectActionRegister,
		})
	// Send Reply to client
	p := product{
		Id:          dp.Id,
		Name:        dp.Name,
		Description: dp.Description,
		TimeCreated: dp.TimeCreated,
	}
	rsp := &productAddResponse{
		ResponseCommonParameter: rcp,
		product:                 p,
	}
	return c.JSON(http.StatusOK, rsp)
}

type productUpdateRequest struct {
	Id          string `json:productId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CategoryId  string `json:categoryId"`
}

// updateProduct update product information in registry
func updateProduct(c echo.Context) error {
	logInfo(c, "updateProduct called")

	// Get product
	req := new(productUpdateRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	// Connect with registry
	config := c.(*base.ApiContext).Config
	registry, err := db.NewRegistry(config)
	if err != nil {
		return err
	}
	defer registry.Release()

	// Update product into registry
	dp := db.Product{
		Id:           req.Id,
		Name:         req.Name,
		Description:  req.Description,
		CategoryId:   req.CategoryId,
		TimeModified: time.Now().String(),
	}
	rcp := ResponseCommonParameter{
		RequestId:    uuid.NewV4().String(),
		Success:      true,
		ErrorMessage: "",
	}
	err = registry.UpdateProduct(&dp)
	if err != nil {
		logError(c, "Registry.UpdateProduct(%s) failed", req.Id)
		rcp.Success = false
		rcp.ErrorMessage = err.Error()
	}
	// Notify kafka
	base.AsyncProduceMessage(c, util.TopicNameProduct,
		&util.ProductTopic{
			ProductId:   req.Id,
			ProductName: req.Name,
			Action:      util.ObjectActionUpdate,
		})

	return c.JSON(http.StatusOK, rcp)
}

// deleteProduct delete product from registry store
func deleteProduct(c echo.Context) error {
	logInfo(c, "deleteProduct:%s", c.Param("id"))

	// Connect with registry
	config := c.(*base.ApiContext).Config
	registry, err := db.NewRegistry(config)
	if err != nil {
		logFatal(c, "Registry connection failed")
		return err
	}
	defer registry.Release()

	// Update product into registry
	rcp := ResponseCommonParameter{
		RequestId:    uuid.NewV4().String(),
		Success:      true,
		ErrorMessage: "",
	}
	err = registry.DeleteProduct(c.Param("id"))
	if err != nil {
		rcp.Success = false
		rcp.ErrorMessage = err.Error()
	}
	// Notify kafka
	base.AsyncProduceMessage(c, util.TopicNameProduct,
		&util.ProductTopic{
			ProductId: c.Param("id"),
			Action:    util.ObjectActionDelete,
		})

	return c.JSON(http.StatusOK, rcp)
}

type getProductResponse struct {
	ResponseCommonParameter
	product
}

// getProduct retrieve production information from registry store
func getProduct(c echo.Context) error {
	logInfo(c, "getProduct:%s", c.Param("id"))
	// Connect with registry
	config := c.(*base.ApiContext).Config
	registry, err := db.NewRegistry(config)
	if err != nil {
		logFatal(c, "Registry connection failed")
		return err
	}
	defer registry.Release()

	// Update product into registry
	rcp := ResponseCommonParameter{
		RequestId:    uuid.NewV4().String(),
		Success:      true,
		ErrorMessage: "",
	}
	p, err := registry.GetProduct(c.Param("id"))
	if err != nil {
		rcp.Success = false
		rcp.ErrorMessage = err.Error()
		return c.JSON(http.StatusOK, rcp)
	}
	rsp := getProductResponse{
		ResponseCommonParameter: rcp,
		product: product{
			Id:           p.Id,
			Name:         p.Name,
			TimeCreated:  p.TimeCreated,
			TimeModified: p.TimeModified,
			Description:  p.Description,
		}}
	return c.JSON(http.StatusOK, rsp)
}

type device struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

type getProductDevicesResponse struct {
	ResponseCommonParameter
	devices []device `json:"devices"`
}

// getProductDevices retrieve product devices list from registry store
func getProductDevices(c echo.Context) error {
	logInfo(c, "getProductDevices:%s", c.Param("id"))

	// Connect with registry
	config := c.(*base.ApiContext).Config
	registry, err := db.NewRegistry(config)
	if err != nil {
		logFatal(c, "Registry connection failed")
		return err
	}
	defer registry.Release()

	// Update product into registry
	rcp := ResponseCommonParameter{
		RequestId:    uuid.NewV4().String(),
		Success:      true,
		ErrorMessage: "",
	}
	devices, err := registry.GetProductDevices(c.Param("id"))
	if err != nil {
		logDebug(c, "Registry.getProductDevices(%s) failed:%v", c.Param("id"), err)
		rcp.Success = false
		rcp.ErrorMessage = err.Error()
		return c.JSON(http.StatusOK, rcp)
	}
	rsp := getProductDevicesResponse{
		ResponseCommonParameter: rcp, devices: []device{}}
	for _, dev := range devices {
		rsp.devices = append(rsp.devices,
			device{Id: dev.Id,
				Status: dev.DeviceStatus,
			})
	}
	return c.JSON(http.StatusOK, rsp)

}
