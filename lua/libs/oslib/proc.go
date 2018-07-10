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
)

func osGetuid(L *lua.LState) int {
	L.Push(lua.LNumber(os.Getuid()))

	return 1
}

func osGeteuid(L *lua.LState) int {
	L.Push(lua.LNumber(os.Geteuid()))

	return 1
}

func osGetgid(L *lua.LState) int {
	L.Push(lua.LNumber(os.Getgid()))

	return 1
}

func osGetgroups(L *lua.LState) int {
	groups, err := os.Getgroups()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	tb := L.CreateTable(len(groups), 0)
	for i, v := range groups {
		tb.RawSetInt(i, lua.LNumber(v))
	}

	L.Push(tb)

	return 1
}
