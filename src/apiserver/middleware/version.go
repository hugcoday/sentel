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

package middleware

import (
	"apiserver/util"

	"github.com/labstack/echo"
)

type (
	// BodyLimitConfig defines the config for BodyLimit middleware.
	ApiVersionConfig struct {
		Version string
	}
)

var (
	// DefaultBodyLimitConfig is the default Gzip middleware config.
	DefaultApiVersionConfig = ApiVersionConfig{
		Version: "v1",
	}
)

func ApiVersion(version string) echo.MiddlewareFunc {
	c := DefaultApiVersionConfig
	c.Version = version
	return ApiVersionWithConfig(c)
}

func ApiVersionWithConfig(config ApiVersionConfig) echo.MiddlewareFunc {
	// Defaults
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			version := c.Param("api-version")
			if version != config.Version {
				return c.JSON(util.ErrorInvalidApiVersion.Code, util.ErrorInvalidApiVersion)
			}
			return next(c)
		}
	}
}
