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

package auth

import (
	"golang.org/x/net/context"
)

type dummyAuthService struct {}

func NewDummyAuthService() IAuthAPI {
	return &dummyAuthService{}
}

// Get version of Authlet service
func (d *dummyAuthService) GetVersion(ctx context.Context) int {
	return 1
}
// Check acl based on client id, user name and topic
func (d *dummyAuthService) CheckAcl(ctx context.Context, clientid string, username string, topic string, access string) error {
	return nil
}
// Check user name and password
func (d *dummyAuthService) CheckUserNameAndPassword(ctx context.Context, username string, password string) error {
	return nil
}
// Get PSK key
func (d *dummyAuthService) GetPskKey(ctx context.Context, hint string, identity string) (string, error) {
	return "", nil
}