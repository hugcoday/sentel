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

package authlet

import (
	"context"
	"errors"
	"fmt"
	"libs"

	"github.com/golang/glog"
)

var (
	ErrorAclDenied = errors.New("Acl denied")
)

var (
	_authServices = make(map[string]AuthServiceFactory)
)

// AuthService interface for security
type AuthService interface {
	GetVersion() int
	CheckAcl(ctx context.Context, clientid string, username string, topic string, access int) error
	CheckUsernameAndPasswor(ctx context.Context, username string, password string) error
	GetPskKey(ctx context.Context, hint string, identity string) (string, error)
	Cleanup(ctx context.Context) error
}

// AuthServiceFactory
type AuthServiceFactory interface {
	New(c libs.Config) (AuthService, error)
}

// RegisterAuthService register a auth plugin
func RegisterAuthService(name string, factory AuthServiceFactory) {
	if _authServices[name] != nil {
		glog.Errorf("AuthService '%s' is already registered")
		return
	}
	_authServices[name] = factory
}

// LoadAuthService load a authService
func LoadAuthService(name string, c libs.Config) (AuthService, error) {
	// Default authentication is 'none'
	if name == "" {
		glog.Warning("No authentication method is specified, using none authentication")
		name = "none"
	}
	if _authServices[name] == nil {
		return nil, fmt.Errorf("AuthService '%s' is not registered", name)
	}
	return _authServices[name].New(c)
}
