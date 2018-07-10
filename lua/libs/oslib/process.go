// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package oslib

import (
	"github.com/yuin/gopher-lua"
	"os"
	"syscall"
)

const (
	osProcessTypeName = lua.OsLibName + ".PROCESS*"
)

func osRegisterProcessMetatype(L *lua.LState) {
	// meta table
	mt := L.NewTypeMetatable(osProcessTypeName)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), osProcessFuncs))
}

var osProcessFuncs = map[string]lua.LGFunction{
	"kill":   osProcessKill,
	"signal": osProcessSignal,
}

func osProcessKill(L *lua.LState) int {
	process := checkProcess(L, 1)
	if err := process.Kill(); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func osProcessSignal(L *lua.LState) int {
	process := checkProcess(L, 1)
	sig := L.CheckInt(2)
	if err := process.Signal(syscall.Signal(sig)); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func osGetpid(L *lua.LState) int {
	L.Push(lua.LNumber(os.Getpid()))

	return 1
}

func osGetppid(L *lua.LState) int {
	L.Push(lua.LNumber(os.Getppid()))

	return 1
}

func osFindProcess(L *lua.LState) int {
	pid := L.CheckInt(1)
	process, err := os.FindProcess(pid)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(newProcess(L, process))

	return 1
}
