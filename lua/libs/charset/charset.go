// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package charset implements charset for Lua.
package charset

import (
	"github.com/yuin/charsetutil"
	"github.com/yuin/gopher-lua"
)

const (
	CharsetLibName = "charset"
)

func Open(L *lua.LState) {
	L.PreloadModule(CharsetLibName, Loader)
}

func Loader(L *lua.LState) int {
	charsetmod := L.SetFuncs(L.NewTable(), charsetFuncs)
	L.Push(charsetmod)

	return 1
}

func charsetEncode(L *lua.LState) int {
	s, err := charsetutil.EncodeString(L.CheckString(1), L.CheckString(2))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LString(s))

	return 1
}

func charsetDecode(L *lua.LState) int {
	bytes, err := charsetutil.DecodeString(L.CheckString(1), L.CheckString(2))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LString(string(bytes)))

	return 1
}

var charsetFuncs = map[string]lua.LGFunction{
	"encode": charsetEncode,
	"eecode": charsetDecode,
}
