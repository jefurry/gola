// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package event

import (
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
	"testing"
)

func TestEvent(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	data := L.CreateTable(0, 2)
	data.RawSetH(lua.LString("name"), lua.LString("Jeff"))
	ctx := lua.LNil

	evt := newEvent(L, data, ctx)
	if !assert.NotEqual(t, nil, evt.RawGetH(lua.LString("data")), "data should matching") {
		return
	}

	if !assert.Equal(t, lua.LNil, evt.RawGetH(lua.LString("context")), "context should matching") {
		return
	}

	dat := evt.RawGetH(lua.LString("data"))
	if !assert.Equal(t, lua.LTTable, dat.Type(), "type should matching") {
		return
	}

	d, ok := dat.(*lua.LTable)
	if !assert.Equal(t, true, ok, "ok should matching") {
		return
	}

	if !assert.NotEqual(t, nil, d, "d should mismatching") {
		return
	}

	if !assert.Equal(t, lua.LString("Jeff"), d.RawGetH(lua.LString("name")), "name should matching") {
		return
	}
}

func TestEventLuaCode(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local event = require('event')

	local evt = event.newEvent({name="Jeff"}, {})
	if evt.data == nil then
		return false
	end

	if evt.data.name ~= "Jeff" then
		return false
	end

	if evt.context == nil then
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
