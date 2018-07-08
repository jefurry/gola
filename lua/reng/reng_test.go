// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package reng

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
	"testing"
	"time"
)

func xTestDefault(t *testing.T) {
	ctx, _ := context.WithTimeout(context.TODO(), 1*time.Second)
	lpm, err := Default(ctx)
	if !assert.NoError(t, err, "Default should succeed") {
		return
	}

	if !assert.Equal(t, 1, lpm.Len(), "length mismatching") {
		return
	}

	if !assert.Equal(t, 0, lpm.Cap(), "cap mismatching") {
		return
	}

	if !assert.Equal(t, "Running", lpm.StatusString(), "status string mismatching") {
		return
	}

	options := lpm.Config().Options()
	if !assert.Equal(t, true, options.SkipOpenLibs, "") {
		return
	}

	lpm.Shutdown()
}

func TestWithLuaCode(t *testing.T) {
	ctx, _ := context.WithTimeout(context.TODO(), 1*time.Second)
	lpm, err := Default(ctx)
	if !assert.NoError(t, err, "Default should succeed") {
		return
	}

	if !assert.Equal(t, 1, lpm.Len(), "length mismatching") {
		return
	}

	if !assert.Equal(t, 0, lpm.Cap(), "cap mismatching") {
		return
	}

	if !assert.Equal(t, "Running", lpm.StatusString(), "status string mismatching") {
		return
	}

	options := lpm.Config().Options()

	if !assert.Equal(t, true, options.SkipOpenLibs, "") {
		return
	}

	code := `
	-- preload module assert
	assert(io == nil, "io should be nil")
	assert(os == nil, "os should be nil")
	assert(debug == nil, "debug should be nil")
	assert(channel == nil, "channel should be nil")
	assert(coroutine == nil, "coroutine should be nil")

	-- functions assert
	assert(collectgarbage == nil, "collectgarbage should be nil")
	assert(dofile == nil, "dofile should be nil")
	assert(error == nil, "error should be nil")
	assert(getfenv == nil, "getfenv should be nil")
	assert(getmetatable == nil, "getmetatable should be nil")
	assert(load == nil, "load should be nil")
	assert(loadfile == nil, "loadfile should be nil")
	assert(loadstring == nil, "loadstring should be nil")
	assert(pcall == nil, "pcall should be nil")
	assert(print == nil, "print should be nil")
	assert(rawequal == nil, "rawequal should be nil")
	assert(rawget == nil, "rawget should be nil")
	assert(rawset == nil, "rawset should be nil")
	assert(_printregs == nil, "_printregs should be nil")
	assert(setfenv == nil, "setfenv should be nil")
	assert(setmetatable == nil, "setmetatable should be nil")
	assert(xpcall == nil, "xpcall should be nil")
	assert(module == nil, "module should be nil")
	assert(require == nil, "require should be nil")
	assert(newproxy == nil, "newproxy should be nil")


	assert(assert ~= nil, "assert should be not nil")
	assert(next ~= nil, "next should be not nil")
	assert(select ~= nil, "select should be not nil")
	assert(tonumber ~= nil, "tonumber should be not nil")
	assert(tostring ~= nil, "tostring should be not nil")
	assert(type ~= nil, "type should be not nil")
	assert(unpack ~= nil, "unpack should be not nil")


	return true
	`

	lv, err := lpm.DoString(ctx, code, func(L *lua.LState) (lua.LValue, error) {
		if !assert.Equal(t, 1, L.GetTop(), "L.GetTop mismatching") {
			return lua.LFalse, nil
		}

		yesno := lua.LVAsBool(L.Get(-1))

		return lua.LBool(yesno), nil
	})

	if !assert.NoError(t, err, `L.DoString should succeed`) {
		return
	}

	if !assert.Equal(t, lua.LTrue, lv, "value mismatching") {
		return
	}

	if !assert.Equal(t, 0, lpm.ServingNum(), "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 1, lpm.Len(), "length mismatching") {
		return
	}

	lpm.Shutdown()

	if !assert.Equal(t, 0, lpm.ServingNum(), "servingNum mismatching") {
		return
	}

	if !assert.Equal(t, 0, lpm.Len(), "length mismatching") {
		return
	}
}
