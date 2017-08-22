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
	"apiserver/types"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
)

type productAddRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type productAddResponse struct {
	types.ResponseCommonParameter
	types.Product
}

func addProduct(c echo.Context) error {
	logInfo(c, "addProduct called")
	// Get product
	p := new(productAddRequest)
	if err := c.Bind(p); err != nil {
		return err
	}
	// Connect with registry
	ctx := *c.(*base.ApiContext)
	r, err := db.NewRegistry(ctx)
	if err != nil {
		logFatal(c, "Registry connection failed")
		return err
	}
	defer r.Release()

	// Insert product into registry, the created product
	// will be modified to retrieve specific information sucha as
	// product.id and creation time
	product := types.Product{
		Name:        p.Name,
		Description: p.Description,
		TimeCreated: time.Now().String(),
	}
	rcp := types.ResponseCommonParameter{
		RequestId:    uuid.NewV4().String(),
		Success:      true,
		ErrorMessage: "",
	}
	err = r.AddProduct(&product)
	if err != nil {
		rcp.Success = false
		rcp.ErrorMessage = err.Error()
		return c.JSON(http.StatusOK, rcp)
	}
	rsp := &productAddResponse{
		ResponseCommonParameter: rcp,
		Product:                 product,
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
	p := new(productUpdateRequest)
	if err := c.Bind(p); err != nil {
		return err
	}
	// Connect with registry
	ctx := *c.(*base.ApiContext)
	registry, err := db.NewRegistry(ctx)
	if err != nil {
		return err
	}
	defer registry.Release()

	// Update product into registry
	product := types.Product{
		Id:           p.Id,
		Name:         p.Name,
		Description:  p.Description,
		CategoryId:   p.CategoryId,
		TimeModified: time.Now().String(),
	}
	rcp := types.ResponseCommonParameter{
		RequestId:    uuid.NewV4().String(),
		Success:      true,
		ErrorMessage: "",
	}
	err = registry.UpdateProduct(&product)
	if err != nil {
		rcp.Success = false
		rcp.ErrorMessage = err.Error()
	}
	return c.JSON(http.StatusOK, rcp)
}

// deleteProduct delete product from registry store
func deleteProduct(c echo.Context) error {
	logInfo(c, "deleteProduct:%s", c.Param("id"))

	// Connect with registry
	registry, err := db.NewRegistry(*c.(*base.ApiContext))
	if err != nil {
		logFatal(c, "Registry connection failed")
		return err
	}
	defer registry.Release()

	// Update product into registry
	rcp := types.ResponseCommonParameter{
		RequestId:    uuid.NewV4().String(),
		Success:      true,
		ErrorMessage: "",
	}
	err = registry.DeleteProduct(c.Param("id"))
	if err != nil {
		rcp.Success = false
		rcp.ErrorMessage = err.Error()
	}
	return c.JSON(http.StatusOK, rcp)
}

type getProductResponse struct {
	types.ResponseCommonParameter
	types.Product
}

// getProduct retrieve production information from registry store
func getProduct(c echo.Context) error {
	logInfo(c, "getProduct:%s", c.Param("id"))
	// Connect with registry
	registry, err := db.NewRegistry(*c.(*base.ApiContext))
	if err != nil {
		logFatal(c, "Registry connection failed")
		return err
	}
	defer registry.Release()

	// Update product into registry
	rcp := types.ResponseCommonParameter{
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
	rsp := getProductResponse{ResponseCommonParameter: rcp, Product: *p}
	return c.JSON(http.StatusOK, rsp)
}

type getProductDevicesResponse struct {
	types.ResponseCommonParameter
	devices []types.Device
}

// getProductDevices retrieve product devices list from registry store
func getProductDevices(c echo.Context) error {
	logInfo(c, "getProductDevices:%s", c.Param("id"))
	// Connect with registry
	registry, err := db.NewRegistry(*c.(*base.ApiContext))
	if err != nil {
		logFatal(c, "Registry connection failed")
		return err
	}
	defer registry.Release()

	// Update product into registry
	rcp := types.ResponseCommonParameter{
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
		ResponseCommonParameter: rcp, devices: []types.Device{}}
	for _, device := range devices {
		rsp.devices = append(rsp.devices, device)
	}
	return c.JSON(http.StatusOK, rsp)

}
