// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package base

import (
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
	"os"
	"testing"
)

func TestLuaPath(t *testing.T) {
	var err error

	err = os.Setenv("HOME", "/home/gola")
	if !assert.NoError(t, err, `os.Setenv should succeed`) {
		return
	}

	err = SetDefaultPath("")
	if !assert.NoError(t, err, "SetDefaultPath should succeed") {
		return
	}

	if !assert.Equal(t, "./?.lua;./?/init.lua", lua.LuaPathDefault, "LuaPathDefault mismatching") {
		return
	}

	err = SetDefaultPath("~/modules")
	if !assert.NoError(t, err, "SetDefaultPath should succeed") {
		return
	}
	if !assert.Equal(t, "./?.lua;./?/init.lua;/home/gola/modules/?.lua;/home/gola/modules/?/init.lua", lua.LuaPathDefault, "LuaPathDefault mismatching") {
		return
	}

	err = AddDefaultPath("~/lualibs")
	if !assert.NoError(t, err, "AddDefaultPath should succeed") {
		return
	}
	if !assert.Equal(t, "./?.lua;./?/init.lua;/home/gola/lualibs/?.lua;/home/gola/lualibs/?/init.lua;/home/gola/modules/?.lua;/home/gola/modules/?/init.lua", lua.LuaPathDefault, "LuaPathDefault mismatching") {
		return
	}

	lua.LuaPathDefault = OldLuaPathDefault
	if !assert.Equal(t, "./?.lua;/usr/local/share/lua/5.1/?.lua;/usr/local/share/lua/5.1/?/init.lua", lua.LuaPathDefault, "LuaPathDefault mismatching") {
		return
	}

	err = AddDefaultPath("~/modules")
	if !assert.NoError(t, err, "AddDefaultPath should succeed") {
		return
	}
	if !assert.Equal(t, "./?.lua;./?/init.lua;/home/gola/modules/?.lua;/home/gola/modules/?/init.lua;/usr/local/share/lua/5.1/?.lua;/usr/local/share/lua/5.1/?/init.lua", lua.LuaPathDefault, "LuaPathDefault mismatching") {
		return
	}

	err = RemoveDefaultPath("~/modules")
	if !assert.NoError(t, err, "RemoveDefaultPath should succeed") {
		return
	}
	if !assert.Equal(t, "./?.lua;./?/init.lua;/usr/local/share/lua/5.1/?.lua;/usr/local/share/lua/5.1/?/init.lua", lua.LuaPathDefault, "LuaPathDefault mismatching") {
		return
	}

	ResetDefaultPath()
	if !assert.Equal(t, "./?.lua;/usr/local/share/lua/5.1/?.lua;/usr/local/share/lua/5.1/?/init.lua", lua.LuaPathDefault, "LuaPathDefault mismatching") {
		return
	}
}
