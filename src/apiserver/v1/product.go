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
	// Get product
	p := new(productAddRequest)
	if err := c.Bind(p); err != nil {
		return err
	}
	// Connect with registry
	ctx := *c.(*base.ApiContext)
	r, err := db.NewRegistry(ctx)
	if err != nil {
		return err
	}
	defer r.Release()

	// Insert product into registry, the created product
	// will be modified to retrieve specific information sucha as
	// product.id and creation time
	product := types.Product{
		Name: p.Name, Description: p.Description}
	rcp := types.ResponseCommonParameter{
		RequestId:    uuid.NewV4().String(),
		Success:      true,
		ErrorMessage: "",
	}

	err = r.AddProduct(&product)
	if err != nil {
		rsp := &productAddResponse{
			ResponseCommonParameter: rcp,
			Product:                 product,
		}
		return c.JSON(http.StatusOK, rsp)
	}
	rcp.Success = false
	rcp.ErrorMessage = err.Error()
	return c.JSON(http.StatusOK, rcp)
}

func deleteProduct(c echo.Context) error {
	return nil
}

func getProduct(c echo.Context) error {
	return nil
}

func getProductDevices(c echo.Context) error {
	return nil
}
