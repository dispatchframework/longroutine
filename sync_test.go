// Copyright (c) 2020 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package longroutine_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dispatchframework/longroutine"
)

func TestSyncStarter_Start(t *testing.T) {
	starter := longroutine.NewStarter()

	started := &sync.WaitGroup{}
	started.Add(1)

	run := &sync.WaitGroup{}
	run.Add(1)

	c := make(chan struct{}, 10)

	for range [10]struct{}{} {
		starter.Start("one-long", func() {
			c <- struct{}{}
			started.Done() // will rightfully panic if more than one of these are called concurrently

			run.Wait() // simulating a long running routine
		})
	}

	started.Wait()
	require.Equal(t, 1, len(c))
	<-c
	require.Equal(t, 0, len(c))

	for range [10]struct{}{} {
		started.Wait()
		started.Add(1)
		starter.Start("many-short", func() {
			c <- struct{}{}
			started.Done()
		})
	}

	started.Wait()
	require.Equal(t, 10, len(c))
}
