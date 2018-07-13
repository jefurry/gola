// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package base implements base for Lua.
package base

import (
	bbase "github.com/jefurry/gola/lua/base"
	"github.com/yuin/gopher-lua"
	"os"
)

func Open(L *lua.LState) {
	packagemod, ok := L.GetGlobal(lua.LoadLibName).(*lua.LTable)
	if !ok {
		L.RaiseError("module(%v) not exists", lua.LoadLibName)
	}

	L.SetFuncs(packagemod, baseFuncs)
}

var baseFuncs = map[string]lua.LGFunction{
	"setDefaultPath":    baseSetDefaultPath,
	"addDefaultPath":    baseAddDefaultPath,
	"removeDefaultPath": baseRemoveDefaultPath,
	"resetDefaultPath":  baseResetDefaultPath,
}

func baseSetDefaultPath(L *lua.LState) int {
	p := L.CheckString(1)
	err := bbase.SetDefaultPath(p)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	restore(L)
	L.Push(lua.LTrue)

	return 1
}

func baseAddDefaultPath(L *lua.LState) int {
	p := L.CheckString(1)
	err := bbase.AddDefaultPath(p)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	restore(L)
	L.Push(lua.LTrue)

	return 1
}

func baseRemoveDefaultPath(L *lua.LState) int {
	p := L.CheckString(1)
	err := bbase.RemoveDefaultPath(p)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	restore(L)
	L.Push(lua.LTrue)

	return 1
}

func baseResetDefaultPath(L *lua.LState) int {
	err := bbase.ResetDefaultPath()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	restore(L)
	L.Push(lua.LTrue)

	return 1
}

func restore(L *lua.LState) {
	path := os.Getenv(lua.LuaPath)
	if len(path) == 0 {
		path = lua.LuaPathDefault
	}

	packagemod := L.GetGlobal(lua.LoadLibName).(*lua.LTable)
	L.SetField(packagemod, "path", lua.LString(path))
}
