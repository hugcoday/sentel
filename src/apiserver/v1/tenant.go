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
	"apiserver/api"
	"apiserver/db"
	"net/http"

	"github.com/labstack/echo"
)

// addTenant add a new tenant
func addTenant(c echo.Context) error {
	// Make security check, for add content, no security policy

	// Get registry store instance by context
	ctx := *c.(*api.ApiContext)
	r, _ := db.NewRegistry(ctx)
	defer r.Release()

	//id := c.Param("id")
	//	r.DeleteDevice(id)
	return c.NoContent(http.StatusNoContent)
}

func deleteTenant(c echo.Context) error {
	return nil
}

func getTenant(c echo.Context) error {
	return nil
}
