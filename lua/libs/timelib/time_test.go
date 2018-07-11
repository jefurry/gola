// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package timelib

import (
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
	"testing"
)

func TestTimeDefaultDate(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local time = require('time')

	local d = time.date()
	if d:year() ~= 1970 then
		return false
	end

	if d:month() ~= 1 then
		return false
	end

	if d:day() ~= 1 then
		return false
	end

	if d:hour() ~= 0 then
		return false
	end

	if d:minute() ~= 0 then
		return false
	end

	if d:second() ~= 0 then
		return false
	end

	if d:nanosecond() ~= 0 then
		return false
	end

	if d:weekday() ~= 4 then
		return false
	end

	if d:yearDay() ~= 1 then
		return false
	end

	if tostring(d:UTC()) ~= "1970-01-01 00:00:00 +0000 UTC" then
		return false
	end

	if tostring(d:Local()) ~= "1970-01-01 08:00:00 +0800 CST" then
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

func TestTimeDate(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local time = require('time')

	local d = time.date{year=2018, month=7, day=10, hour=11, min=30, sec=0, nsec=0, loc=time.UTC}
	if d:year() ~= 2018 then
		return false
	end

	if d:month() ~= 7 then
		return false
	end

	if d:day() ~= 10 then
		return false
	end

	if d:hour() ~= 11 then
		return false
	end

	if d:minute() ~= 30 then
		return false
	end

	if d:second() ~= 0 then
		return false
	end

	if d:nanosecond() ~= 0 then
		return false
	end

	if d:weekday() ~= 2 then
		return false
	end

	if d:yearDay() ~= 191 then
		return false
	end

	if tostring(d:UTC()) ~= "2018-07-10 11:30:00 +0000 UTC" then
		return false
	end

	if tostring(d:Local()) ~= "2018-07-10 19:30:00 +0800 CST" then
		return false
	end

	if d:zone() ~= "UTC" then
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

func TestTimeLocalDate(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local time = require('time')

	local loc, msg = time.loadLocation("Asia/ShangHai")
	if loc == nil then
		return false
	end

	local d = time.date{year=2018, month=7, day=10, hour=11, min=30, sec=0, nsec=0, loc=loc}
	if d:year() ~= 2018 then
		return false
	end

	if d:month() ~= 7 then
		return false
	end

	if d:day() ~= 10 then
		return false
	end

	if d:hour() ~= 11 then
		return false
	end

	if d:minute() ~= 30 then
		return false
	end

	if d:second() ~= 0 then
		return false
	end

	if d:nanosecond() ~= 0 then
		return false
	end

	if d:weekday() ~= 2 then
		return false
	end

	if d:yearDay() ~= 191 then
		return false
	end

	if tostring(d:UTC()) ~= "2018-07-10 03:30:00 +0000 UTC" then
		return false
	end

	if tostring(d:Local()) ~= "2018-07-10 11:30:00 +0800 CST" then
		return false
	end

	if d:zone() ~= "CST" then
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
