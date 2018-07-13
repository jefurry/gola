// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package os

import (
	"github.com/yuin/gopher-lua"
	oos "os"
)

func osIsExist(L *lua.LState) int {
	name := L.CheckString(1)
	if _, err := oos.Stat(name); oos.IsNotExist(err) {
		L.Push(lua.LFalse)
	} else {
		L.Push(lua.LTrue)
	}

	return 1
}

func osIsNotExist(L *lua.LState) int {
	name := L.CheckString(1)
	if _, err := oos.Stat(name); oos.IsNotExist(err) {
		L.Push(lua.LTrue)
	} else {
		L.Push(lua.LFalse)
	}

	return 1
}

func osMkdir(L *lua.LState) int {
	name := L.CheckString(1)
	perm := L.CheckInt(2)

	if err := oos.Mkdir(name, oos.FileMode(perm)); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func osMkdirAll(L *lua.LState) int {
	path := L.CheckString(1)
	perm := L.CheckInt(2)

	if err := oos.MkdirAll(path, oos.FileMode(perm)); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func osRemove(L *lua.LState) int {
	name := L.CheckString(1)
	if err := oos.Remove(name); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func osRemoveAll(L *lua.LState) int {
	path := L.CheckString(1)

	if err := oos.RemoveAll(path); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}
