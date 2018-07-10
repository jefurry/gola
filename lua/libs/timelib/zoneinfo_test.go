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

func TestLocation(t *testing.T) {
	L := lua.NewState()
	L.PreloadModule(TimeLibName, Loader)
	defer L.Close()

	code := `
	local time = require('time')

	if tostring(time.UTC) ~= "UTC" then
		return false
	end

	if tostring(time.Local) ~= "Local" then
		return false
	end

	if time.UTC ~= time.UTC then
		return false
	end

	if time.Local ~= time.Local then
		return false
	end

	if tostring(time.fixedZone("BeiJing", 8)) ~= "BeiJing" then
		return false
	end

	if tostring(time.loadLocation("America/New_York")) ~= "America/New_York" then
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
