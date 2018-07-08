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
	lpm, err := Default(ctx)
	if !assert.NoError(t, err, "Start should succeed") {
		return
	}

	if !assert.Equal(t, DefaultStartNum, lpm.Len(), "length mismatching") {
		return
	}

	if !assert.Equal(t, DefaultMaxNum, lpm.config.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, DefaultStartNum, lpm.config.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, DefaultMaxRequest, lpm.config.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, DefaultIdleTimeout, lpm.config.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 3600, lpm.config.seconds, "seconds mismatching") {
		return
	}

	ls1, err := lpm.get(ctx)
	if !assert.NoError(t, err, "get should succeed") {
		return
	}

	if !assert.Equal(t, DefaultStartNum, lpm.Len(), "lenth mismatching") {
		return
	}

	if !assert.Equal(t, 1, lpm.ServingNum(), "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 1, ls1.requestedNum, "requestedNum mismatching") {
		return
	}

	ls2, err := lpm.get(ctx)
	if !assert.NoError(t, err, "get should succeed") {
		return
	}

	if !assert.Equal(t, 2, lpm.TotalRequestedNum(), "length mismatching") {
		return
	}

	if !assert.Equal(t, 2, lpm.ServingNum(), "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 2, lpm.Len(), "requestedNum mismatching") {
		return
	}

	lpm.put(ls2)
	go lpm.Shutdown()
	go lpm.Close(ls1)

	<-time.After(3 * time.Second)

	if !assert.Equal(t, 0, lpm.ServingNum(), "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 0, lpm.Len(), "requestedNum mismatching") {
		return
	}
}

func TestServingNumWithDefault(t *testing.T) {
	ctx, _ := context.WithCancel(context.TODO())
	lpm, err := Default(ctx)
	if !assert.NoError(t, err, "Default should succeed") {
		return
	}

	if !assert.Equal(t, 0, lpm.servingNum, "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 1, lpm.length, "length mismatching") {
		return
	}

	if !assert.Equal(t, 0, lpm.Cap(), "Cap mismatching") {
		return
	}

	ls1, err := lpm.get(ctx)
	if !assert.NoError(t, err, "get should succeed") {
		return
	}

	if !assert.Equal(t, 1, lpm.servingNum, "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 1, lpm.length, "length mismatching") {
		return
	}

	ls2, err := lpm.get(ctx)
	if !assert.NoError(t, err, "get should succeed") {
		return
	}

	if !assert.Equal(t, 2, lpm.servingNum, "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 2, lpm.length, "length mismatching") {
		return
	}

	err = lpm.put(ls1)
	if !assert.NoError(t, err, "put should succeed") {
		return
	}

	if !assert.Equal(t, 1, lpm.servingNum, "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 2, lpm.length, "length mismatching") {
		return
	}

	ls3, err := lpm.get(ctx)
	if !assert.NoError(t, err, "get should succeed") {
		return
	}

	if !assert.Equal(t, 2, lpm.servingNum, "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 2, lpm.length, "length mismatching") {
		return
	}

	lpm.Close(ls2)
	err = lpm.put(ls2)
	if !assert.Error(t, err, "put should be not succeed") {
		return
	}

	if !assert.Equal(t, 1, lpm.servingNum, "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 1, lpm.length, "length mismatching") {
		return
	}

	ls4, err := lpm.get(ctx)
	if !assert.NoError(t, err, "get should succeed") {
		return
	}

	if !assert.Equal(t, 2, lpm.servingNum, "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 2, lpm.length, "length mismatching") {
		return
	}

	err = lpm.put(ls3)
	if !assert.NoError(t, err, "put should be not succeed") {
		return
	}

	if !assert.Equal(t, 1, lpm.servingNum, "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 2, lpm.length, "length mismatching") {
		return
	}

	go lpm.Shutdown()
	go lpm.Close(ls4)

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

	if !assert.Equal(t, 0, lpm.servingNum, "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 0, lpm.length, "length mismatching") {
		return
	}
}

func TestNew_1(t *testing.T) {
	config, err := NewConfig(1, -1, 230, 120, "1d")
	if !assert.NoError(t, err, "NewConfig should succeed") {
		return
	}

	ctx, _ := context.WithCancel(context.TODO())
	lpm, err := New(ctx, config)
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}

	ls1, err := lpm.get(ctx)
	if !assert.NoError(t, err, "get should succeed") {
		return
	}

	ls2, err := lpm.get(ctx)
	if !assert.Equal(t, "lua state pool fulled", err.Error(), "get should failed") {
		return
	}

	if !assert.Equal(t, ((*lState)(nil)), ls2, "ls mismatching") {
		return
	}

	if !assert.Equal(t, 1, lpm.servingNum, "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 1, lpm.Len(), "length mismatching") {
		return
	}

	fmt.Println(lpm.servingNum, lpm.Len())

	go lpm.Shutdown()
	go lpm.Close(ls1)

	<-time.After(3 * time.Second)

	if !assert.Equal(t, 0, lpm.servingNum, "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 0, lpm.Len(), "length mismatching") {
		return
	}
}

func TestNew_2(t *testing.T) {
	var lpm *LPM
	var err error

	config1, err := NewConfig(45, -1, 230, 120, "1d")
	if !assert.NoError(t, err, "NewConfig should succeed") {
		return
	}

	ctx1, _ := context.WithCancel(context.TODO())
	lpm, err = New(ctx1, config1)
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}

	if !assert.Equal(t, 45, lpm.config.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 23, lpm.config.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 230, lpm.config.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "1d", lpm.config.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 24*60*60, lpm.config.seconds, "seconds mismatching") {
		return
	}
	lpm.Shutdown()

	config2, err := NewConfig(45, -1, 230, 120, "-2d")
	if !assert.NoError(t, err, "NewConfig should succeed") {
		return
	}

	ctx2, _ := context.WithCancel(context.TODO())
	lpm, err = New(ctx2, config2)
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}

	if !assert.Equal(t, 45, lpm.config.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 23, lpm.config.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 230, lpm.config.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "1d", lpm.config.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 24*60*60, lpm.config.seconds, "seconds mismatching") {
		return
	}
	lpm.Shutdown()

	config3, err := NewConfig(50, 30, 500, 120, "2h")
	if !assert.NoError(t, err, "NewConfig should succeed") {
		return
	}

	ctx3, _ := context.WithCancel(context.TODO())
	lpm, err = New(ctx3, config3)
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}

	if !assert.Equal(t, 50, lpm.config.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 30, lpm.config.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 500, lpm.config.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "2h", lpm.config.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 2*60*60, lpm.config.seconds, "seconds mismatching") {
		return
	}
	lpm.Shutdown()

	config4, err := NewConfig(50, 30, 500, 120, "-3h")
	if !assert.NoError(t, err, "NewConfig should succeed") {
		return
	}

	ctx4, _ := context.WithCancel(context.TODO())
	lpm, err = New(ctx4, config4)
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}

	if !assert.Equal(t, 50, lpm.config.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 30, lpm.config.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 500, lpm.config.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "1h", lpm.config.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 60*60, lpm.config.seconds, "seconds mismatching") {
		return
	}
	lpm.Shutdown()

	config5, err := NewConfig(83, 0, 450, 120, "78m")
	if !assert.NoError(t, err, "NewConfig should succeed") {
		return
	}

	ctx5, _ := context.WithCancel(context.TODO())
	lpm, err = New(ctx5, config5)
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}

	if !assert.Equal(t, 83, lpm.config.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 42, lpm.config.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 450, lpm.config.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "78m", lpm.config.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 78*60, lpm.config.seconds, "seconds mismatching") {
		return
	}
	lpm.Shutdown()

	config6, err := NewConfig(83, 0, 450, 120, "-98m")
	if !assert.NoError(t, err, "NewConfig should succeed") {
		return
	}

	ctx6, _ := context.WithCancel(context.TODO())
	lpm, err = New(ctx6, config6)
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}

	if !assert.Equal(t, 83, lpm.config.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 42, lpm.config.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 450, lpm.config.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "1m", lpm.config.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 60, lpm.config.seconds, "seconds mismatching") {
		return
	}
	lpm.Shutdown()

	config7, err := NewConfig(59, 100, 763, 120, "1583s")
	if !assert.NoError(t, err, "NewConfig should succeed") {
		return
	}

	ctx7, _ := context.WithCancel(context.TODO())
	lpm, err = New(ctx7, config7)
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}

	if !assert.Equal(t, 59, lpm.config.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 30, lpm.config.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 763, lpm.config.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "1583s", lpm.config.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 1583, lpm.config.seconds, "seconds mismatching") {
		return
	}
	lpm.Shutdown()

	config8, err := NewConfig(59, 100, 763, 120, "-583s")
	if !assert.NoError(t, err, "NewConfig should succeed") {
		return
	}

	ctx8, _ := context.WithCancel(context.TODO())
	lpm, err = New(ctx8, config8)
	if !assert.NoError(t, err, `New should succeed`) {
		return
	}

	if !assert.Equal(t, 59, lpm.config.maxNum, "maxNum mismatching") {
		return
	}

	if !assert.Equal(t, 30, lpm.config.startNum, "startNum mismatching") {
		return
	}

	if !assert.Equal(t, 763, lpm.config.maxRequest, "maxRequest mismatching") {
		return
	}

	if !assert.Equal(t, "1s", lpm.config.idleTimeout, "idleTimeout mismatching") {
		return
	}

	if !assert.Equal(t, 1, lpm.config.seconds, "seconds mismatching") {
		return
	}
	lpm.Shutdown()
}

func TestNewLuaCode(t *testing.T) {
	config, err := NewConfig(45, -1, 230, 120, "1d")
	if !assert.NoError(t, err, "NewConfig should succeed") {
		return
	}

	ctx, _ := context.WithCancel(context.TODO())
	lpm, err := New(ctx, config, func(L *lua.LState) error {
		libs.OpenLibs(L)

		return nil
	})

	if !assert.NoError(t, err, "New should succeed") {
		return
	}

	defer lpm.Shutdown()

	if !assert.Equal(t, 23, lpm.Len(), "length mismatching") {
		return
	}

	if !assert.Equal(t, 0, lpm.ServingNum(), "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 0, lpm.TotalRequestedNum(), "totalRequestedNum mismatching") {
		return
	}

	ch := make(chan bool, 0)
	go func() {
		ls, err := lpm.get(ctx)
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

		if !assert.Equal(t, 1, lpm.ServingNum(), "servingNum mismatching") {
			close(ch)
			return
		}

		lpm.put(ls)

		ch <- true
	}()

	for b := range ch {
		if !assert.Equal(t, true, b, "should equals true") {
			return
		}

		if !assert.Equal(t, 0, lpm.ServingNum(), "ServingNum mismatching") {
			return
		}

		break
	}
}
