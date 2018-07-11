// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package eventlib

import (
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
	"testing"
)

func TestEmitter_OnOff_1(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local event = require('event')

	local EventClientType = "CLICK"

	local f1 = function(evt)
		return true
	end

	local emitter = event.newEmitter()
	emitter:on(EventClientType, f1)
	if table.maxn(emitter:getListeners(EventClientType)) ~= 1 then
		return false
	end

	emitter:off(EventClientType, f1)
	if table.maxn(emitter:getListeners(EventClientType)) ~= 0 then
		return false
	end

	return true
	`

	err := L.DoString(code)
	if !assert.NoError(t, err, `L.DoString should succeed`) {
		return
	}

	if !assert.Equal(t, 1, L.GetTop(), "L.GetTop mismatching") {
		return
	}

	ret := L.Get(-1)
	if !assert.Equal(t, lua.LTBool, ret.Type(), "type mismatching") {
		return
	}

	if !assert.Equal(t, lua.LTrue, ret.(lua.LBool), "value mismatching") {
		return
	}
}

func TestEmitter_OnOff_2(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local event = require('event')

	local EventClientType = "CLICK"

	local f1 = function(evt)
		return true
	end

	local emitter = event.newEmitter()
	emitter:on(EventClientType, f1)
	emitter:on(EventClientType, f1)
	emitter:on(EventClientType, f1)
	emitter:on(EventClientType, f1)
	emitter:on(EventClientType, f1)
	emitter:on(EventClientType, f1)
	emitter:on(EventClientType, f1)
	emitter:on(EventClientType, f1)
	emitter:on(EventClientType, f1)
	if table.maxn(emitter:getListeners(EventClientType)) ~= 9 then
		return false
	end

	emitter:off(EventClientType, f1)
	if table.maxn(emitter:getListeners(EventClientType)) ~= 0 then
		return false
	end

	return true
	`

	err := L.DoString(code)
	if !assert.NoError(t, err, `L.DoString should succeed`) {
		return
	}

	if !assert.Equal(t, 1, L.GetTop(), "L.GetTop mismatching") {
		return
	}

	ret := L.Get(-1)
	if !assert.Equal(t, lua.LTBool, ret.Type(), "type mismatching") {
		return
	}

	if !assert.Equal(t, lua.LTrue, ret.(lua.LBool), "value mismatching") {
		return
	}
}

func TestEmitter_OnOff_3(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local event = require('event')

	local EventClientType = "CLICK"

	local f1 = function(evt)
		return true
	end

	local f2 = function(evt)
		return false
	end

	local emitter = event.newEmitter()
	emitter:on(EventClientType, f1)
	emitter:on(EventClientType, f2)
	emitter:on(EventClientType, f2)
	emitter:on(EventClientType, f1)
	emitter:on(EventClientType, f2)
	emitter:on(EventClientType, f1)
	emitter:on(EventClientType, f2)
	emitter:on(EventClientType, f1)
	emitter:on(EventClientType, f1)
	if table.maxn(emitter:getListeners(EventClientType)) ~= 9 then
		return false
	end

	emitter:off(EventClientType, f1)
	if table.maxn(emitter:getListeners(EventClientType)) ~= 4 then
		return false
	end

	emitter:off(EventClientType, f2)
	if table.maxn(emitter:getListeners(EventClientType)) ~= 0 then
		return false
	end

	return true
	`

	err := L.DoString(code)
	if !assert.NoError(t, err, `L.DoString should succeed`) {
		return
	}

	if !assert.Equal(t, 1, L.GetTop(), "L.GetTop mismatching") {
		return
	}

	ret := L.Get(-1)
	if !assert.Equal(t, lua.LTBool, ret.Type(), "type mismatching") {
		return
	}

	if !assert.Equal(t, lua.LTrue, ret.(lua.LBool), "value mismatching") {
		return
	}
}

func TestEmitter_Fire_1(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local event = require('event')

	local EventClientType = "CLICK"

	local f1 = function(evt)
		assert(evt.data ~= nil, "data should be not nil")
		assert(type(evt.data) == "table", "data's type should be table")
		assert(evt.data.name == "Jeff", "data's value should be Jeff")
		assert(evt.context ~= nil, "context's value should be nil")

		return true
	end

	local emitter = event.newEmitter()
	emitter:on(EventClientType, f1)

	if table.maxn(emitter:getListeners(EventClientType)) ~= 1 then
		return false
	end

	emitter:fire(EventClientType, {name="Jeff"}, {})

	return true
	`

	err := L.DoString(code)
	if !assert.NoError(t, err, `L.DoString should succeed`) {
		return
	}

	if !assert.Equal(t, 1, L.GetTop(), "L.GetTop mismatching") {
		return
	}

	ret := L.Get(-1)
	if !assert.Equal(t, lua.LTBool, ret.Type(), "type mismatching") {
		return
	}

	if !assert.Equal(t, lua.LTrue, ret.(lua.LBool), "value mismatching") {
		return
	}
}

func TestEmitter_Once(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local event = require('event')

	local EventClientType = "CLICK"

	local f1 = function(evt)
		assert(evt.data ~= nil, "data should be not nil")
		assert(type(evt.data) == "table", "data's type should be table")
		assert(evt.data.name == "Jeff", "data's value should be Jeff")
		assert(evt.context ~= nil, "context's value should be nil")

		return true
	end

	local emitter = event.newEmitter()
	emitter:once(EventClientType, f1)

	if table.maxn(emitter:getListeners(EventClientType)) ~= 1 then
		return false
	end

	emitter:fire(EventClientType, {name="Jeff"}, {})

	if table.maxn(emitter:getListeners(EventClientType)) ~= 0 then
		return false
	end

	return true
	`

	err := L.DoString(code)
	if !assert.NoError(t, err, `L.DoString should succeed`) {
		return
	}

	if !assert.Equal(t, 1, L.GetTop(), "L.GetTop mismatching") {
		return
	}

	ret := L.Get(-1)
	if !assert.Equal(t, lua.LTBool, ret.Type(), "type mismatching") {
		return
	}

	if !assert.Equal(t, lua.LTrue, ret.(lua.LBool), "value mismatching") {
		return
	}
}
