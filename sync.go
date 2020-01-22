// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package longroutine

import "sync"

// Starter should start no more than one (potentially long-running) concurrent go routine for a given key
type Starter interface {
	// Start no more than one instance of (potentially long-running) routine `f` gets started for a given `key`
	Start(key string, f func())
}

// NewStarter creates a Starter
func NewStarter() Starter {
	return &syncStarter{
		m: map[string]struct{}{},
	}
}

type syncStarter struct {
	sync.Mutex

	m map[string]struct{}
}

func (s *syncStarter) Start(key string, f func()) {
	s.Lock()
	defer s.Unlock()

	if _, exists := s.m[key]; !exists {
		go s.run(key, f)
		s.m[key] = struct{}{}
	}
}

func (s *syncStarter) run(key string, f func()) {
	defer func() {
		s.Lock()
		defer s.Unlock()

		delete(s.m, key)
	}()
	f()
}
