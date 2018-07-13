// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package hex implements hex encoding for Lua.
package hex

import (
	hhex "encoding/hex"
	"github.com/yuin/gopher-lua"
)

const (
	HexLibName = "encoding.hex"
)

func Open(L *lua.LState) {
	L.PreloadModule(HexLibName, Loader)
}

func Loader(L *lua.LState) int {
	hexmod := L.SetFuncs(L.NewTable(), hexFuncs)
	L.Push(hexmod)

	return 1
}

func hexEncodeToString(L *lua.LState) int {
	src := L.CheckString(1)

	s := hhex.EncodeToString([]byte(src))

	L.Push(lua.LString(s))

	return 1
}

func hexDecodeString(L *lua.LState) int {
	s := L.CheckString(1)

	bs, err := hhex.DecodeString(s)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LString(string(bs)))

	return 1
}

func hexDump(L *lua.LState) int {
	src := L.CheckString(1)

	s := hhex.Dump([]byte(src))

	L.Push(lua.LString(s))

	return 1
}

func hexEncodedLen(L *lua.LState) int {
	n := L.CheckInt(1)

	L.Push(lua.LNumber(hhex.EncodedLen(n)))

	return 1
}

func hexDecodedLen(L *lua.LState) int {
	x := L.CheckInt(1)

	L.Push(lua.LNumber(hhex.DecodedLen(x)))

	return 1
}

var hexFuncs = map[string]lua.LGFunction{
	"encode":     hexEncodeToString,
	"decode":     hexDecodeString,
	"dump":       hexDump,
	"encodedLen": hexEncodedLen,
	"decodedLen": hexDecodedLen,
}

var hexFields = map[string]lua.LValue{}
