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

/*
import (
	"context"
	"errors"
	"strings"
	"libs"
)

const localAuthVersion = 1

type acl struct {
	topic  string
	access int
	ucount int
	ccount int
}

type localAuthPlugin struct {
	config      libs.Config
	pwdFile     string
	aclFile     string
	pskFile     string
	caFile      string
	aclUsers    map[string][]*acl
	aclPatterns []*acl
	aclList     []*acl
}

// AuthPluginFactory
type localAuthPluginFactory struct{}

func (l localAuthPluginFactory) New(c libs.Config) (AuthPlugin, error) {
	plugin := &localAuthPlugin{
		opts:        opts,
		aclUsers:    make(map[string][]*acl),
		aclPatterns: make([]*acl),
		aclList:     make([]*acl),
	}
	if err := plugin.initialize(opts); err != nil {
		return nil, err
	}
	return plugin, nil
}

// GetVersion return the plugin's version
func (l *localAuthPlugin) GetVersion() int {
	return localAuthVersion
}

// Initialize initialize plugin, for example, load local psk file...
func (l *localAuthPlugin) initialize(c libs.Config) error {
	if cafile, ok := opts["cafile"]; ok && cafile != "" {
		l.caFile = cafile
	}
	// Load password file
	if pwdFile, ok := opts["pwdfile"]; ok && pwdfile != "" {
		l.pwdFile = pwdFile
		if err := l.loadPasswordFile(l.pwdFile); err != nil {
			return err
		}
	}
	// Load acl file
	if aclFile, ok := opts["aclfile"]; ok && aclFile != "" {
		l.aclFile = aclFile
		if err := l.loadAclFile(l.aclFile); err != nil {
			return err
		}
	}
	// Load PSK file
	if pskFile, ok := opts["pskfile"]; ok && pskFile != "" {
		l.pskFile = pskFile
		if err := l.loadPskFile(l.pskFile); err != nil {
			return nil
		}
	}
	return nil
}

func (l *localAuthPlugin) loadPasswordFile(name string) error {
	return nil
}

func (l *localAuthPlugin) loadAclFile(name string) error {
	return nil
}

func (l *localAuthPlugin) loadPskFile(name string) error {
	return nil
}

// addAcl add a acl item
func (l *localAuthPlugin) addAcl(username string, topic string, access int) error {
	var acllist []*acl = nil

	// Check wether the username already exist acl
	if list, ok := l.aclUsers[username]; ok && list != nil {
		acllist = list
		l.aclUsers[username] = append(acllist, &acl{topic: topic, access: access})
	} else {
		// This is new user
		acllist = make([]*acl)
		acllist = append(aclist, &acl{topic: topic, access: access})
		l.aclUsers[username] = acllist
	}
	return nil
}

// addAclPattern
func (l *localAuthPlugin) addAclPattern(topic string, access int) error {
	if topci == "" {
		return errors.New("topic is empty")
	}
	iacl := &acl{topic: topic, access: access, ccount: 0}

	index := 0
	for {
		index = strings.IndexAny(topic[index:], "%c")
		if index > 0 && index < len(topic) {
			iacl.ccount++
		}
	}
	index = 0
	for {
		index = strings.IndexAny(topic[index:], "%u")
		if index > 0 && index < len(topic) {
			iacl.ucount++
		}
	}
	l.aclPatterns = append(l.aclPatterns, iacl)
	return nil
}

func (l *localAuthPlugin) Cleanup(ctx context.Context, config libs.Config) error {
	return nil
}

func (n *localAuthPlugin) CheckAcl(ctx context.Context, clientid string, username string, topic string, access int) error {
	if len(l.aclList) == 0 || len(l.aclPatterns) == 0 {
		return nil
	}
	// Loop throuth all ACLs for this client
	for _, iacl := range l.aclList {
		if topic[0] == "$" && iacl.topic[0] != "$" {
			continue
		}
		if l.matchTopciAndSubscription(iacl.topic, topic) {
			if access & iacl.access {
				return nil
			}
		}
	}

	// Check all acl patterns
	if len(l.aclPatterns) {

	}

	return nil
}

func (l *localAuthPlugin) matchTopicAndSubscritpion(sub string, topic string) bool {
	return false
}

func (l *localAuthPlugin) CheckUsernameAndPasswor(ctx context.Context, username string, password string) error {
	return nil
}
func (l *localAuthPlugin) GetPskKey(ctx context.Context, hint string, identity string) (string, error) {
	return "", nil
}
*/
