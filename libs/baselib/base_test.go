// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package baselib

import (
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
	"testing"
)

func TestBaseLib(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local os = require('os')

	os.setenv("HOME", "/home/gola")

	package.setDefaultPath("~/modules")
	if package.path ~= "./?.lua;./?/init.lua;/home/gola/modules/?.lua;/home/gola/modules/?/init.lua" then
		return false
	end

	package.setDefaultPath("")
	if package.path ~= "./?.lua;./?/init.lua" then
		return false
	end

	package.resetDefaultPath()
	if package.path ~= "./?.lua;/usr/local/share/lua/5.1/?.lua;/usr/local/share/lua/5.1/?/init.lua" then
		return false
	end

	package.addDefaultPath("~/modules")
	if package.path ~= "./?.lua;./?/init.lua;/home/gola/modules/?.lua;/home/gola/modules/?/init.lua;/usr/local/share/lua/5.1/?.lua;/usr/local/share/lua/5.1/?/init.lua" then
		return false
	end

	package.removeDefaultPath("~/modules")
	if package.path ~= "./?.lua;./?/init.lua;/usr/local/share/lua/5.1/?.lua;/usr/local/share/lua/5.1/?/init.lua" then
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
	if !assert.Equal(t, lua.LTrue, ret, "value mismatching") {
		return
	}
}
