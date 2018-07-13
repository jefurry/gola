// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package log

import (
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestLog(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	tempDir := func(L *lua.LState) int {
		dir := L.CheckString(1)
		prefix := L.CheckString(2)

		name, err := ioutil.TempDir(dir, prefix)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))

			return 2
		}

		L.Push(lua.LString(name))

		return 1
	}

	join := func(L *lua.LState) int {
		top := L.GetTop()
		if top < 1 {
			L.Push(lua.LString(""))

			return 1
		}

		params := make([]string, 0, top)
		for i := 1; i <= top; i++ {
			params = append(params, L.CheckString(i))
		}

		L.Push(lua.LString(filepath.Join(params...)))

		return 1
	}

	removeAll := func(L *lua.LState) int {
		if err := os.RemoveAll(L.CheckString(1)); err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(err.Error()))

			return 2
		}

		L.Push(lua.LTrue)

		return 1
	}

	L.SetGlobal("tempDir", L.NewFunction(tempDir))
	L.SetGlobal("join", L.NewFunction(join))
	L.SetGlobal("removeAll", L.NewFunction(removeAll))

	code := `
	local log = require('log')

	local dir = tempDir("", "file-rotatelog-test")
	if dir == nil then
		return false
	end

	-- text formatter
	local formatter = log.newJSONFormatter("3:04PM")

	-- rotatelog hook
	local rotatelogHook = log.newRotatelogHook(join(dir, "gola.log.%Y%m%d%H%M%S"), {withClock="UTC", withLocation="Asia/ShangHai"})
	rotatelogHook:setFormatter(log.newTextFormatter("3:04PM"))

	-- logger
	local logger = log.newLogger()
	logger:setOut(log.DISCARD)
	logger:setFormatter(formatter)
	logger:addHook(rotatelogHook)
	logger:info("Good Luck")

	removeAll(dir)
	
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
