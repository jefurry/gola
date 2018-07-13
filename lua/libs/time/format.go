// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package time

import (
	"github.com/yuin/gopher-lua"
	ttime "time"
)

func timeParse(L *lua.LState) int {
	layout := L.CheckString(1)
	value := L.CheckString(2)

	t, err := ttime.Parse(layout, value)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(newTime(L, t))

	return 1
}

func timeParseInLocation(L *lua.LState) int {
	layout := L.CheckString(1)
	value := L.CheckString(2)
	loc := checkLocation(L, 3)

	t, err := ttime.ParseInLocation(layout, value, loc)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(newTime(L, t))

	return 1
}

func timeParseDuration(L *lua.LState) int {
	s := L.CheckString(1)

	d, err := ttime.ParseDuration(s)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LNumber(d))

	return 1
}
