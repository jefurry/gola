// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package execlib implements os/exec for Lua.
package execlib

import (
	"github.com/yuin/gopher-lua"
	"os/exec"
)

const (
	ExecLibName = "exec"
)

func Open(L *lua.LState) {
	L.PreloadModule(ExecLibName, Loader)
}

func Loader(L *lua.LState) int {
	execmod := L.SetFuncs(L.NewTable(), execFuncs)
	L.Push(execmod)

	return 1
}

var execFuncs = map[string]lua.LGFunction{
	"lookPath": execLookPath,
}

var execFields = map[string]lua.LValue{}

func execLookPath(L *lua.LState) int {
	file := L.CheckString(1)

	path, err := exec.LookPath(file)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LString(path))

	return 1
}
