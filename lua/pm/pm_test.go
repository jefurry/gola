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
	"fmt"
	"github.com/jefurry/gola/lua/libs"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	ctx, _ := context.WithCancel(context.TODO())
	pm, err := Default(ctx)
	if !assert.NoError(t, err, "Start should succeed") {
		return
	}

	if !assert.Equal(t, DefaultStartNum, pm.Len(), "length mismatching") {
		return
	}

	if !assert.Equal(t, DefaultMaxNum, pm.config.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, DefaultStartNum, pm.config.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, DefaultMaxRequest, pm.config.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, DefaultIdleTimeout, pm.config.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 3600, pm.config.seconds, "seconds mismatching") {
		return
	}

	ls1, err := pm.get(ctx)
	if !assert.NoError(t, err, "get should succeed") {
		return
	}

	if !assert.Equal(t, DefaultStartNum, pm.Len(), "lenth mismatching") {
		return
	}

	if !assert.Equal(t, 1, pm.ServingNum(), "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 1, ls1.requestedNum, "requestedNum mismatching") {
		return
	}

	ls2, err := pm.get(ctx)
	if !assert.NoError(t, err, "get should succeed") {
		return
	}

	if !assert.Equal(t, 2, pm.TotalRequestedNum(), "length mismatching") {
		return
	}

	if !assert.Equal(t, 2, pm.ServingNum(), "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 2, pm.Len(), "requestedNum mismatching") {
		return
	}

	pm.put(ls2)
	go pm.Shutdown()
	go pm.Close(ls1)

	<-time.After(3 * time.Second)

	if !assert.Equal(t, 0, pm.ServingNum(), "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 0, pm.Len(), "requestedNum mismatching") {
		return
	}
}

func TestServingNumWithDefault(t *testing.T) {
	ctx, _ := context.WithCancel(context.TODO())
	pm, err := Default(ctx)
	if !assert.NoError(t, err, "Default should succeed") {
		return
	}

	if !assert.Equal(t, 0, pm.servingNum, "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 1, pm.length, "length mismatching") {
		return
	}

	if !assert.Equal(t, 0, pm.Cap(), "Cap mismatching") {
		return
	}

	ls1, err := pm.get(ctx)
	if !assert.NoError(t, err, "get should succeed") {
		return
	}

	if !assert.Equal(t, 1, pm.servingNum, "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 1, pm.length, "length mismatching") {
		return
	}

	ls2, err := pm.get(ctx)
	if !assert.NoError(t, err, "get should succeed") {
		return
	}

	if !assert.Equal(t, 2, pm.servingNum, "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 2, pm.length, "length mismatching") {
		return
	}

	err = pm.put(ls1)
	if !assert.NoError(t, err, "put should succeed") {
		return
	}

	if !assert.Equal(t, 1, pm.servingNum, "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 2, pm.length, "length mismatching") {
		return
	}

	ls3, err := pm.get(ctx)
	if !assert.NoError(t, err, "get should succeed") {
		return
	}

	if !assert.Equal(t, 2, pm.servingNum, "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 2, pm.length, "length mismatching") {
		return
	}

	pm.Close(ls2)
	err = pm.put(ls2)
	if !assert.Error(t, err, "put should be not succeed") {
		return
	}

	if !assert.Equal(t, 1, pm.servingNum, "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 1, pm.length, "length mismatching") {
		return
	}

	ls4, err := pm.get(ctx)
	if !assert.NoError(t, err, "get should succeed") {
		return
	}

	if !assert.Equal(t, 2, pm.servingNum, "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 2, pm.length, "length mismatching") {
		return
	}

	err = pm.put(ls3)
	if !assert.NoError(t, err, "put should be not succeed") {
		return
	}

	if !assert.Equal(t, 1, pm.servingNum, "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 2, pm.length, "length mismatching") {
		return
	}

	go pm.Shutdown()
	go pm.Close(ls4)

	<-time.After(3 * time.Second)

	if !assert.Equal(t, true, ls1.closed, "closed mismatching") {
		return
	}

	if !assert.Equal(t, true, ls2.closed, "closed mismatching") {
		return
	}

	if !assert.Equal(t, true, ls3.closed, "closed mismatching") {
		return
	}

	if !assert.Equal(t, true, ls4.closed, "closed mismatching") {
		return
	}

	if !assert.Equal(t, 0, pm.servingNum, "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 0, pm.length, "length mismatching") {
		return
	}
}

func TestNew_1(t *testing.T) {
	config, err := NewConfig(1, -1, 230, 120, "1d")
	if !assert.NoError(t, err, "NewConfig should succeed") {
		return
	}

	ctx, _ := context.WithCancel(context.TODO())
	pm, err := New(ctx, config)
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}

	ls1, err := pm.get(ctx)
	if !assert.NoError(t, err, "get should succeed") {
		return
	}

	ls2, err := pm.get(ctx)
	if !assert.Equal(t, "lua state pool fulled", err.Error(), "get should failed") {
		return
	}

	if !assert.Equal(t, ((*lState)(nil)), ls2, "ls mismatching") {
		return
	}

	if !assert.Equal(t, 1, pm.servingNum, "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 1, pm.Len(), "length mismatching") {
		return
	}

	fmt.Println(pm.servingNum, pm.Len())

	go pm.Shutdown()
	go pm.Close(ls1)

	<-time.After(3 * time.Second)

	if !assert.Equal(t, 0, pm.servingNum, "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 0, pm.Len(), "length mismatching") {
		return
	}
}

func TestNew_2(t *testing.T) {
	var pm *PM
	var err error

	config1, err := NewConfig(45, -1, 230, 120, "1d")
	if !assert.NoError(t, err, "NewConfig should succeed") {
		return
	}

	ctx1, _ := context.WithCancel(context.TODO())
	pm, err = New(ctx1, config1)
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}

	if !assert.Equal(t, 45, pm.config.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 23, pm.config.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 230, pm.config.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "1d", pm.config.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 24*60*60, pm.config.seconds, "seconds mismatching") {
		return
	}
	pm.Shutdown()

	config2, err := NewConfig(45, -1, 230, 120, "-2d")
	if !assert.NoError(t, err, "NewConfig should succeed") {
		return
	}

	ctx2, _ := context.WithCancel(context.TODO())
	pm, err = New(ctx2, config2)
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}

	if !assert.Equal(t, 45, pm.config.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 23, pm.config.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 230, pm.config.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "1d", pm.config.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 24*60*60, pm.config.seconds, "seconds mismatching") {
		return
	}
	pm.Shutdown()

	config3, err := NewConfig(50, 30, 500, 120, "2h")
	if !assert.NoError(t, err, "NewConfig should succeed") {
		return
	}

	ctx3, _ := context.WithCancel(context.TODO())
	pm, err = New(ctx3, config3)
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}

	if !assert.Equal(t, 50, pm.config.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 30, pm.config.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 500, pm.config.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "2h", pm.config.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 2*60*60, pm.config.seconds, "seconds mismatching") {
		return
	}
	pm.Shutdown()

	config4, err := NewConfig(50, 30, 500, 120, "-3h")
	if !assert.NoError(t, err, "NewConfig should succeed") {
		return
	}

	ctx4, _ := context.WithCancel(context.TODO())
	pm, err = New(ctx4, config4)
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}

	if !assert.Equal(t, 50, pm.config.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 30, pm.config.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 500, pm.config.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "1h", pm.config.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 60*60, pm.config.seconds, "seconds mismatching") {
		return
	}
	pm.Shutdown()

	config5, err := NewConfig(83, 0, 450, 120, "78m")
	if !assert.NoError(t, err, "NewConfig should succeed") {
		return
	}

	ctx5, _ := context.WithCancel(context.TODO())
	pm, err = New(ctx5, config5)
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}

	if !assert.Equal(t, 83, pm.config.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 42, pm.config.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 450, pm.config.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "78m", pm.config.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 78*60, pm.config.seconds, "seconds mismatching") {
		return
	}
	pm.Shutdown()

	config6, err := NewConfig(83, 0, 450, 120, "-98m")
	if !assert.NoError(t, err, "NewConfig should succeed") {
		return
	}

	ctx6, _ := context.WithCancel(context.TODO())
	pm, err = New(ctx6, config6)
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}

	if !assert.Equal(t, 83, pm.config.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 42, pm.config.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 450, pm.config.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "1m", pm.config.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 60, pm.config.seconds, "seconds mismatching") {
		return
	}
	pm.Shutdown()

	config7, err := NewConfig(59, 100, 763, 120, "1583s")
	if !assert.NoError(t, err, "NewConfig should succeed") {
		return
	}

	ctx7, _ := context.WithCancel(context.TODO())
	pm, err = New(ctx7, config7)
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}

	if !assert.Equal(t, 59, pm.config.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 30, pm.config.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 763, pm.config.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "1583s", pm.config.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 1583, pm.config.seconds, "seconds mismatching") {
		return
	}
	pm.Shutdown()

	config8, err := NewConfig(59, 100, 763, 120, "-583s")
	if !assert.NoError(t, err, "NewConfig should succeed") {
		return
	}

	ctx8, _ := context.WithCancel(context.TODO())
	pm, err = New(ctx8, config8)
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}

	if !assert.Equal(t, 59, pm.config.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 30, pm.config.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 763, pm.config.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "1s", pm.config.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 1, pm.config.seconds, "seconds mismatching") {
		return
	}
	pm.Shutdown()
}

func TestNewLuaCode(t *testing.T) {
	config, err := NewConfig(45, -1, 230, 120, "1d")
	if !assert.NoError(t, err, "NewConfig should succeed") {
		return
	}

	ctx, _ := context.WithCancel(context.TODO())
	pm, err := New(ctx, config, func(L *lua.LState) error {
		libs.OpenLibs(L)

		return nil
	})

	if !assert.NoError(t, err, "New should succeed") {
		return
	}

	defer pm.Shutdown()

	if !assert.Equal(t, 23, pm.Len(), "length mismatching") {
		return
	}

	if !assert.Equal(t, 0, pm.ServingNum(), "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 0, pm.TotalRequestedNum(), "totalRequestedNum mismatching") {
		return
	}

	ch := make(chan bool, 0)
	go func() {
		ls, err := pm.get(ctx)
		if !assert.NoError(t, err, "get should succeed") {
			close(ch)
			return
		}

		if !assert.NoError(t, err, "get should succeed") {
			close(ch)
			return
		}

		code := `
		return true
		`

		err = ls.L.DoString(code)
		if !assert.NoError(t, err, `L.DoString should succeed`) {
			close(ch)
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

		if !assert.Equal(t, 1, pm.ServingNum(), "servingNum mismatching") {
			close(ch)
			return
		}

		pm.put(ls)

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
