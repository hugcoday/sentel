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
	"net/http"

	"github.com/cloustone/sentel/apiserver/base"
	"github.com/cloustone/sentel/apiserver/db"

	"github.com/labstack/echo"
)

// addTenant add a new tenant
func addTenant(ctx echo.Context) error {
	// Make security check, for add content, no security policy

	// Get registry store instance by context
	config := ctx.(*base.ApiContext).Config
	r, _ := db.NewRegistry(config)
	defer r.Release()

	//id := c.Param("id")
	//	r.DeleteDevice(id)
	return ctx.NoContent(http.StatusNoContent)
}

func deleteTenant(ctx echo.Context) error {
	return nil
}

func getTenant(ctx echo.Context) error {
	return nil
}
