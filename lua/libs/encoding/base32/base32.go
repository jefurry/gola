// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package base32 implements base32 encoding for Lua.
package base32

import (
	b32 "encoding/base32"
	"github.com/yuin/gopher-lua"
)

const (
	Base32LibName = "encoding.base32"
)

func Open(L *lua.LState) {
	L.PreloadModule(Base32LibName, Loader)
}

func Loader(L *lua.LState) int {
	base32mod := L.SetFuncs(L.NewTable(), base32Funcs)
	L.Push(base32mod)

	base32RegisterEncodingMetatype(L)

	L.SetField(base32mod, "StdEncoding", newBase32Encoding(L, b32.StdEncoding))
	L.SetField(base32mod, "HexEncoding", newBase32Encoding(L, b32.HexEncoding))

	for k, v := range base32Fields {
		base32mod.RawSetString(k, v)
	}

	return 1
}

func base32EncodingNew(L *lua.LState) int {
	encoder := L.CheckString(1)
	padding := L.OptInt(2, int(b32.StdPadding))

	enc := b32.NewEncoding(encoder)
	pad := rune(padding)
	if pad != b32.StdPadding {
		enc = enc.WithPadding(pad)
	}

	ud := newBase32Encoding(L, enc)

	L.Push(ud)

	return 1
}

func base32Encode(L *lua.LState) int {
	src := L.CheckString(1)

	L.Push(lua.LString(b32.StdEncoding.EncodeToString([]byte(src))))

	return 1
}

func base32Decode(L *lua.LState) int {
	s := L.CheckString(1)

	bs, err := b32.StdEncoding.DecodeString(s)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LString(string(bs)))

	return 1
}

var base32Funcs = map[string]lua.LGFunction{
	"newEncoding": base32EncodingNew,
	"encode":      base32Encode,
	"decode":      base32Decode,
}

var base32Fields = map[string]lua.LValue{
	"STD_PADDING": lua.LNumber(b32.StdPadding), // Standard padding character
	"NO_PADDING":  lua.LNumber(b32.NoPadding),  // No padding
}
