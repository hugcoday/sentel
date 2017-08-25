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

package config

import "errors"

// Config interface
type Config interface {
	Bool(section string, key string) (bool, error)
	Int(section string, key string) (int, error)
	String(section string, key string) (string, error)
	MustBool(section string, key string) bool
	MustInt(section string, key string) int
	MustString(section string, key string) string
	SetValue(section string, key string, val string)
}

var (
	ErrorInvalidConfiguration = errors.New("Invalid configuration")
)
