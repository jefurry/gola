// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pm

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"github.com/yuin/gopher-lua"
)

func TestNewLState_1(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	ls, err := newLState(ctx, 1, "1h", 3600, lua.Options{}, nil)

	if !assert.NoError(t, err, "newLState should succeed") {
		return
	}

	defer ls.close()
	defer cancel()

	code := `
	local clock = os.clock
	function sleep(n)  -- seconds
		local t0 = clock()
		while clock() - t0 <= n do end
	end
	
	sleep(3)
	`

	err = ls.L.DoString(code)
	if !assert.Error(t, err, "DoString should be not succeed") {
		return
	}
}

func TestNewLState_2(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	ls, err := newLState(ctx, 1, "1h", 3600, lua.Options{}, nil)

	if !assert.NoError(t, err, "newLState should succeed") {
		return
	}

	go func() {
		code := `
		local clock = os.clock
		function sleep(n)  -- seconds
			local t0 = clock()
			while clock() - t0 <= n do end
		end

		sleep(3)
		`

		err = ls.L.DoString(code)
		if !assert.Error(t, err, "DoString should be not succeed") {
			return
		}
	}()

	time.Sleep(2)
	cancel()
}

func TestNewLState(t *testing.T) {
	var ls *lState
	var err error

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	for _, v := range []struct {
		maxRequest  int
		idleTimeout string
		seconds     int
		whenNew     NewFunc
	}{
		{20, "1d", 24 * 60 * 60, nil},
		{100, "3h", 3600, nil},
		{45, "1m", 60 * 60, nil},
		{77, "10s", 10, nil},
	} {
		ls, err = newLState(ctx, v.maxRequest, v.idleTimeout, v.seconds, lua.Options{}, v.whenNew)
		if !assert.NoError(t, err, "newLState should succeed") {
			return
		}

		for i := 0; i < v.maxRequest; i++ {
			if !assert.Equal(t, false, ls.mustTerminate(), "value mismatching") {
				return
			}

			ls.incRequestNum()
		}

		if !assert.Equal(t, true, ls.mustTerminate(), "value mismatching") {
			return
		}

		ls.close()
	}
}
