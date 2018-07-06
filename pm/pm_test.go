// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pm

import (
	"github.com/jefurry/gola/libs"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
	"testing"
)

func TestDefault(t *testing.T) {
	pm := Default()
	defer pm.Shutdown()

	err := pm.Start(func(L *lua.LState) error {
		libs.OpenLibs(L)

		return nil
	})

	if !assert.NoError(t, err, "Start should succeed") {
		return
	}

	if !assert.Equal(t, DefaultStartNum, pm.Size(), "size mismatching") {
		return
	}

	if !assert.Equal(t, DefaultMaxNum, pm.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, DefaultStartNum, pm.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, DefaultMaxRequest, pm.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, DefaultIdleTimeout, pm.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 3600, pm.seconds, "seconds mismatching") {
		return
	}

	ls, err := pm.Get()
	if !assert.NoError(t, err, "Get should succeed") {
		return
	}

	if !assert.Equal(t, DefaultStartNum, pm.Size(), "size mismatching") {
		return
	}

	if !assert.Equal(t, 1, pm.ServingNum(), "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 1, ls.requestedNum, "requestedNum mismatching") {
		return
	}

	ls, err = pm.Get()
	if !assert.Equal(t, "lua state pool fulled", err.Error(), "Get should failed") {
		return
	}

	if !assert.Equal(t, ((*lState)(nil)), ls, "ls mismatching") {
		return
	}
}

func TestNew(t *testing.T) {
	var pm *PoolManager
	var err error

	pm, err = New(45, -1, 230, "1d")
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}
	if !assert.Equal(t, 45, pm.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 23, pm.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 230, pm.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "1d", pm.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 24*60*60, pm.seconds, "seconds mismatching") {
		return
	}
	pm.Shutdown()

	pm, err = New(45, -1, 230, "-2d")
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}
	if !assert.Equal(t, 45, pm.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 23, pm.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 230, pm.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "1d", pm.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 24*60*60, pm.seconds, "seconds mismatching") {
		return
	}
	pm.Shutdown()

	pm, err = New(50, 30, 500, "2h")
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}
	if !assert.Equal(t, 50, pm.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 30, pm.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 500, pm.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "2h", pm.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 2*60*60, pm.seconds, "seconds mismatching") {
		return
	}
	pm.Shutdown()

	pm, err = New(50, 30, 500, "-3h")
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}
	if !assert.Equal(t, 50, pm.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 30, pm.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 500, pm.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "1h", pm.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 60*60, pm.seconds, "seconds mismatching") {
		return
	}
	pm.Shutdown()

	pm, err = New(83, 0, 450, "78m")
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}
	if !assert.Equal(t, 83, pm.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 42, pm.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 450, pm.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "78m", pm.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 78*60, pm.seconds, "seconds mismatching") {
		return
	}
	pm.Shutdown()

	pm, err = New(83, 0, 450, "-98m")
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}
	if !assert.Equal(t, 83, pm.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 42, pm.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 450, pm.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "1m", pm.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 60, pm.seconds, "seconds mismatching") {
		return
	}
	pm.Shutdown()

	pm, err = New(59, 100, 763, "1583s")
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}
	if !assert.Equal(t, 59, pm.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 30, pm.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 763, pm.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "1583s", pm.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 1583, pm.seconds, "seconds mismatching") {
		return
	}
	pm.Shutdown()

	pm, err = New(59, 100, 763, "-583s")
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}
	if !assert.Equal(t, 59, pm.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 30, pm.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 763, pm.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "1s", pm.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 1, pm.seconds, "seconds mismatching") {
		return
	}
	pm.Shutdown()
}

func TestNewLuaCode(t *testing.T) {
	pm, err := New(45, -1, 230, "1d")
	if !assert.NoError(t, err, "New should succeed") {
		return
	}

	err = pm.Start(func(L *lua.LState) error {
		libs.OpenLibs(L)

		return nil
	})
	if !assert.NoError(t, err, "Start should succeed") {
		return
	}

	defer pm.Shutdown()

	if !assert.Equal(t, 23, pm.Size(), "Size mismatching") {
		return
	}

	if !assert.Equal(t, 0, pm.ServingNum(), "ServingNum mismatching") {
		return
	}

	if !assert.Equal(t, 0, pm.TotalRequestedNum(), "TotalRequestedNum mismatching") {
		return
	}

	ch := make(chan bool, 0)
	go func() {
		ls, err := pm.Get()
		if !assert.NoError(t, err, "Get should succeed") {
			close(ch)
			return
		}

		if !assert.NoError(t, err, "Get should succeed") {
			close(ch)
			return
		}

		code := `
		return true
		`

		err = ls.L.DoString(code)
		if !assert.NoError(t, err, `L.DoString should succeed`) {
			//close(ch)
			ch <- false
			return
		}

		if !assert.Equal(t, 1, ls.L.GetTop(), "L.GetTop mismatching") {
			close(ch)
			return
		}

		ret := ls.L.Get(-1)
		if !assert.Equal(t, lua.LTBool, ret.Type(), "type mismatching") {
			close(ch)
			return
		}

		if !assert.Equal(t, lua.LTrue, ret.(lua.LBool), "value mismatching") {
			close(ch)
			return
		}

		if !assert.Equal(t, 1, pm.ServingNum(), "ServingNum mismatching") {
			close(ch)
			return
		}

		pm.Put(ls)

		ch <- true
	}()

	for b := range ch {
		if !assert.Equal(t, true, b, "should equals true") {
			return
		}

		if !assert.Equal(t, 0, pm.ServingNum(), "ServingNum mismatching") {
			return
		}

		break
	}
}
