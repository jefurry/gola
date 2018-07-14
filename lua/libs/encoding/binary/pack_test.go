// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package binary

import (
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
	"testing"
)

func TestPack(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local binary = require('encoding.binary')

	local binstr, msg = binary.pack("2?iL5sf", true, false, 32, 89, "hello", 4.5682)
	if binstr == nil or msg ~= nil then
		return false
	end

	local tb, msg = binary.unpack("2?iL5sf", binstr)
	if tb == nil or msg ~= nil then
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

func TestUnpack(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local binary = require('encoding.binary')
	
	local binstr, msg = binary.pack("?0L3s10s2ifhHiIlLfdqQ", true, 35, "one", "two fields", 0xff, 56, 9.8, -28, 45, 123, 45, -99, 20, 3.5, 5.4, 321, 567)
	if binstr == nil or msg ~= nil then
		return false
	end

	if #binstr ~= 74 then
		return false
	end

	local tb, msg = binary.unpack("?3s10s2ifhHiIlLfdqQ", binstr)
	if tb == nil or msg ~= nil then
		return false
	end

	if #tb ~= 16 then
		return false
	end

	if tb[1] ~= true then
		return false
	end

	if tb[2] ~= "one" then
		return false
	end

	if tb[3] ~= "two fields" then
		return false
	end

	if tb[4] ~= 255 then
		return false
	end

	if tb[5] ~= 56 then
		return false
	end

	--if tb[6] ~= 9.8 then
		--return false
	--end

	if tb[7] ~= -28 then
		return false
	end

	if tb[8] ~= 45 then
		return false
	end

	if tb[9] ~= 123 then
		return false
	end

	if tb[10] ~= 45 then
		return false
	end

	if tb[11] ~= -99 then
		return false
	end

	if tb[12] ~= 20 then
		return false
	end

	if tb[13] ~= 3.5 then
		return false
	end

	if tb[14] ~= 5.4 then
		return false
	end

	if tb[15] ~= 321 then
		return false
	end

	if tb[16] ~= 567 then
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
