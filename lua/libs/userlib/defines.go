// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package userlib

import (
	"github.com/yuin/gopher-lua"
	"os/user"
)

func userCurrent(L *lua.LState) int {
	u, err := user.Current()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(newUser(L, u))

	return 1
}

func userLookup(L *lua.LState) int {
	username := L.CheckString(1)

	u, err := user.Lookup(username)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(newUser(L, u))

	return 1
}

func userLookupId(L *lua.LState) int {
	uid := L.CheckString(1)

	u, err := user.LookupId(uid)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(newUser(L, u))

	return 1
}

func userLookupGroup(L *lua.LState) int {
	name := L.CheckString(1)

	g, err := user.LookupGroup(name)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(newGroup(L, g))

	return 1
}

func userLookupGroupId(L *lua.LState) int {
	gid := L.CheckString(1)

	g, err := user.LookupGroupId(gid)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(newGroup(L, g))

	return 1
}
