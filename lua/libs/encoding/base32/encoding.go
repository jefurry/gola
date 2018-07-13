// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package base32

import (
	b32 "encoding/base32"
	"github.com/yuin/gopher-lua"
)

const (
	base32EncodingTypeName = Base32LibName + ".ENCODING*"
)

func base32EncodeWithPadding(L *lua.LState) int {
	encoder := checkBase32Encoding(L, 1)
	padding := L.OptInt(2, int(b32.StdPadding))

	pad := rune(padding)
	enc := encoder.WithPadding(pad)

	ud := newBase32Encoding(L, enc)

	L.Push(ud)

	return 1
}

func base32EncodeToString(L *lua.LState) int {
	enc := checkBase32Encoding(L, 1)
	src := L.CheckString(2)

	L.Push(lua.LString(enc.EncodeToString([]byte(src))))

	return 1
}

func base32DecodeString(L *lua.LState) int {
	enc := checkBase32Encoding(L, 1)
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

func base32RegisterEncodingMetatype(L *lua.LState) {
	// meta table
	mt := L.NewTypeMetatable(base32EncodingTypeName)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), base32EncodingFuncs))
}

var base32EncodingFuncs = map[string]lua.LGFunction{
	"encode":      base32EncodeToString,
	"decode":      base32DecodeString,
	"withPadding": base32EncodeWithPadding,
}
