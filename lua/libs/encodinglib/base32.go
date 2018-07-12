// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package encodinglib

import (
	"encoding/base32"
	"github.com/yuin/gopher-lua"
)

const (
	EncodingBase32LibName = EncodingLibName + ".base32"
)

func OpenBase32(L *lua.LState) {
	L.PreloadModule(EncodingBase32LibName, Base32Loader)
}

func Base32Loader(L *lua.LState) int {
	encbase32mod := L.SetFuncs(L.NewTable(), encodingBase32Funcs)
	L.Push(encbase32mod)

	encodingRegisterBase32EncodingMetatype(L)

	L.SetField(encbase32mod, "StdEncoding", newBase32Encoding(L, base32.StdEncoding))
	L.SetField(encbase32mod, "HexEncoding", newBase32Encoding(L, base32.HexEncoding))

	for k, v := range encodingBase32Fields {
		encbase32mod.RawSetString(k, v)
	}

	return 1
}

func encodingBase32EncodingNew(L *lua.LState) int {
	encoder := L.CheckString(1)
	padding := L.OptInt(2, int(base32.StdPadding))

	enc := base32.NewEncoding(encoder)
	pad := rune(padding)
	if pad != base32.StdPadding {
		enc = enc.WithPadding(pad)
	}

	ud := newBase32Encoding(L, enc)

	L.Push(ud)

	return 1
}

func encodingBase32Encode(L *lua.LState) int {
	src := L.CheckString(1)

	L.Push(lua.LString(base32.StdEncoding.EncodeToString([]byte(src))))

	return 1
}

func encodingBase32Decode(L *lua.LState) int {
	s := L.CheckString(1)

	bs, err := base32.StdEncoding.DecodeString(s)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LString(string(bs)))

	return 1
}

var encodingBase32Funcs = map[string]lua.LGFunction{
	"newEncoding": encodingBase32EncodingNew,
	"encode":      encodingBase32Encode,
	"decode":      encodingBase32Decode,
}

var encodingBase32Fields = map[string]lua.LValue{
	"STD_PADDING": lua.LNumber(base32.StdPadding), // Standard padding character
	"NO_PADDING":  lua.LNumber(base32.NoPadding),  // No padding
}
