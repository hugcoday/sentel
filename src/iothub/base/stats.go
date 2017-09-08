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

package base

import "sync"

// Stats declration
type Stats struct {
	stats map[string]uint64
	mutex *sync.Mutex
}

func NewStats(withlock bool) *Stats {
	if withlock {
		return &Stats{
			stats: make(map[string]uint64),
			mutex: &sync.Mutex{},
		}
	} else {
		return &Stats{
			stats: make(map[string]uint64),
			mutex: nil,
		}
	}
}

func (s *Stats) Get() map[string]uint64 {
	if s.mutex != nil {
		s.mutex.Lock()
		defer s.mutex.Unlock()
	}
	return s.stats
}

func (s *Stats) AddStat(name string, value uint64) {
	if s.mutex != nil {
		s.mutex.Lock()
		defer s.mutex.Unlock()
	}
	s.stats[name] += value
}

func (s *Stats) AddStats(stats *Stats) {
	if s.mutex != nil {
		s.mutex.Lock()
		defer s.mutex.Unlock()
	}
	for k, v := range stats.Get() {
		if _, ok := s.stats[k]; !ok {
			s.stats[k] = v
		} else {
			s.stats[k] += v
		}
	}
}
