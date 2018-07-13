// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package di

import (
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
	"testing"
)

func TestParse_1(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local di = require('di')

	local version = 1.0

	function test(a, b, c)
		local v = version
	end

	local args, msg = di.parse(test)
	if args == nil then
		return false
	end

	if type(args) ~= "table" then
		return false
	end

	if table.maxn(args) ~= 3 or #args ~= 3 then
		return false
	end

	if args[1] ~= "a" then
		return false
	end

	if args[2] ~= "b" then
		return false
	end

	if args[3] ~= "c" then
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

func TestParse_2(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local di = require('di')

	local version = 1.0

	function test(a, b, c)
		local v = version
	end

	local args, msg = di.parse{test}
	if args == nil then
		return false
	end

	if type(args) ~= "table" then
		return false
	end

	if table.maxn(args) ~= 3 or #args ~= 3 then
		return false
	end

	if args[1] ~= "a" then
		return false
	end

	if args[2] ~= "b" then
		return false
	end

	if args[3] ~= "c" then
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

func TestParse_3(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local di = require('di')

	local version = 1.0

	function test(a, b, c)
		local v = version
	end

	local args, msg = di.parse{nil, test}
	if args == nil then
		return false
	end

	if type(args) ~= "table" then
		return false
	end

	if table.maxn(args) ~= 3 or #args ~= 3 then
		return false
	end

	if args[1] ~= "a" then
		return false
	end

	if args[2] ~= "b" then
		return false
	end

	if args[3] ~= "c" then
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

func TestParse_4(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local di = require('di')

	local version = 1.0

	function test(a, b, c)
		local v = version
	end

	local tb = {
		t=test
	}

	local args, msg = di.parse{tb, "t"}
	if args == nil then
		return false
	end

	if type(args) ~= "table" then
		return false
	end

	if table.maxn(args) ~= 3 or #args ~= 3 then
		return false
	end

	if args[1] ~= "a" then
		return false
	end

	if args[2] ~= "b" then
		return false
	end

	if args[3] ~= "c" then
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

func TestParseDotParams_1(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local di = require('di')

	local version = 1.0

	function test(...)
		local v = version
	end

	local args, msg = di.parse(test)
	if args ~= nil then
		return false
	end

	if msg ~= "maybe includes invalid parameters for '...'" then
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

func TestParseDotParams_2(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local di = require('di')

	local version = 1.0

	function test(a, b, c, ...)
		local v = version
	end

	local args, msg = di.parse(test)
	if args == nil then
		return false
	end

	if type(args) ~= "table" then
		return false
	end
	
	if table.maxn(args) ~= 3 or #args ~= 3 then
		return false
	end

	if args[1] ~= "a" then
		return false
	end

	if args[2] ~= "b" then
		return false
	end

	if args[3] ~= "c" then
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
