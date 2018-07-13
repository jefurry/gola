// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package base64 implements base64 encoding for Lua.
package base64

import (
	b64 "encoding/base64"
	"github.com/yuin/gopher-lua"
)

const (
	Base64LibName = "encoding.base64"
)

func Open(L *lua.LState) {
	L.PreloadModule(Base64LibName, Loader)
}

func Loader(L *lua.LState) int {
	base64mod := L.SetFuncs(L.NewTable(), base64Funcs)
	L.Push(base64mod)

	base64RegisterEncodingMetatype(L)

	L.SetField(base64mod, "StdEncoding", newBase64Encoding(L, b64.StdEncoding))
	L.SetField(base64mod, "URLEncoding", newBase64Encoding(L, b64.URLEncoding))
	L.SetField(base64mod, "RawStdEncoding", newBase64Encoding(L, b64.RawStdEncoding))
	L.SetField(base64mod, "RawURLEncoding", newBase64Encoding(L, b64.RawURLEncoding))

	for k, v := range base64Fields {
		base64mod.RawSetString(k, v)
	}

	return 1
}

func base64EncodingNew(L *lua.LState) int {
	encoder := L.CheckString(1)
	padding := L.OptInt(2, int(b64.StdPadding))

	enc := b64.NewEncoding(encoder)
	pad := rune(padding)
	if pad != b64.StdPadding {
		enc = enc.WithPadding(pad)
	}

	ud := newBase64Encoding(L, enc)

	L.Push(ud)

	return 1
}

func base64Encode(L *lua.LState) int {
	src := L.CheckString(1)

	L.Push(lua.LString(b64.StdEncoding.EncodeToString([]byte(src))))

	return 1
}

func base64Decode(L *lua.LState) int {
	s := L.CheckString(1)

	bs, err := b64.StdEncoding.DecodeString(s)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LString(string(bs)))

	return 1
}

var base64Funcs = map[string]lua.LGFunction{
	"newEncoding": base64EncodingNew,
	"encode":      base64Encode,
	"decode":      base64Decode,
}

var base64Fields = map[string]lua.LValue{
	"STD_PADDING": lua.LNumber(b64.StdPadding), // Standard padding character
	"NO_PADDING":  lua.LNumber(b64.NoPadding),  // No padding
}
