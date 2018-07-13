// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package base64

import (
	b64 "encoding/base64"
	"github.com/yuin/gopher-lua"
)

const (
	base64EncodingTypeName = Base64LibName + ".ENCODING*"
)

func base64EncodeWithPadding(L *lua.LState) int {
	encoder := checkBase64Encoding(L, 1)
	padding := L.OptInt(2, int(b64.StdPadding))

	pad := rune(padding)
	enc := encoder.WithPadding(pad)

	ud := newBase64Encoding(L, enc)

	L.Push(ud)

	return 1
}

func base64EncodeToString(L *lua.LState) int {
	enc := checkBase64Encoding(L, 1)
	src := L.CheckString(2)

	L.Push(lua.LString(enc.EncodeToString([]byte(src))))

	return 1
}

func base64DecodeString(L *lua.LState) int {
	enc := checkBase64Encoding(L, 1)
	s := L.CheckString(2)

	bs, err := enc.DecodeString(s)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LString(string(bs)))

	return 1
}

func base64EncodingStrict(L *lua.LState) int {
	enc := checkBase64Encoding(L, 1)

	strictEnc := enc.Strict()
	L.Push(newBase64Encoding(L, strictEnc))

	return 1
}

func base64RegisterEncodingMetatype(L *lua.LState) {
	// meta table
	mt := L.NewTypeMetatable(base64EncodingTypeName)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), base64EncodingFuncs))
}

var base64EncodingFuncs = map[string]lua.LGFunction{
	"encode":      base64EncodeToString,
	"decode":      base64DecodeString,
	"strict":      base64EncodingStrict,
	"withPadding": base64EncodeWithPadding,
}
