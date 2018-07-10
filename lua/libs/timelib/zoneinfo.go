// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package timelib

import (
	"github.com/yuin/gopher-lua"
	"time"
)

const (
	timeLocationTypeName = TimeLibName + ".LOCATION*"
)

func timeFixedZone(L *lua.LState) int {
	name := L.CheckString(1)
	offset := L.CheckInt(2)

	l := time.FixedZone(name, offset)
	ud := newLocation(L, l)

	L.Push(ud)

	return 1
}

func timeLoadLocation(L *lua.LState) int {
	name := L.CheckString(1)

	l, err := time.LoadLocation(name)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	ud := newLocation(L, l)

	L.Push(ud)

	return 1
}

func timeLoadLocationFromTZData(L *lua.LState) int {
	name := L.CheckString(1)
	data := L.CheckString(2)

	loc, err := time.LoadLocationFromTZData(name, []byte(data))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(newLocation(L, loc))

	return 1
}

func timeLocationString(L *lua.LState) int {
	l := checkLocation(L, 1)

	L.Push(lua.LString(l.String()))

	return 1
}

func timeRegisterLocationMetatype(L *lua.LState) {
	// meta table
	mt := L.NewTypeMetatable(timeLocationTypeName)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), timeLocationFuncs))
	L.SetField(mt, "__tostring", L.NewFunction(timeLocationString))
}

var timeLocationFuncs = map[string]lua.LGFunction{}
