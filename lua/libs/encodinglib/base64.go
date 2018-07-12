// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package encodinglib

import (
	"encoding/base64"
	"github.com/yuin/gopher-lua"
)

const (
	EncodingBase64LibName = EncodingLibName + ".base64"
)

func OpenBase64(L *lua.LState) {
	L.PreloadModule(EncodingBase64LibName, Base64Loader)
}

func Base64Loader(L *lua.LState) int {
	encbase64mod := L.SetFuncs(L.NewTable(), encodingBase64EncodingFuncs)
	L.Push(encbase64mod)

	encodingRegisterBase64EncodingMetatype(L)

	L.SetFuncs(encbase64mod, encodingBase64Funcs)
	L.SetField(encbase64mod, "StdEncoding", newBase64Encoding(L, base64.StdEncoding))
	L.SetField(encbase64mod, "URLEncoding", newBase64Encoding(L, base64.URLEncoding))
	L.SetField(encbase64mod, "RawStdEncoding", newBase64Encoding(L, base64.RawStdEncoding))
	L.SetField(encbase64mod, "RawURLEncoding", newBase64Encoding(L, base64.RawURLEncoding))

	for k, v := range encodingBase64Fields {
		encbase64mod.RawSetString(k, v)
	}

	return 1
}

func encodingBase64EncodingNew(L *lua.LState) int {
	encoder := L.CheckString(1)
	padding := L.OptInt(2, int(base64.StdPadding))

	enc := base64.NewEncoding(encoder)
	pad := rune(padding)
	if pad != base64.StdPadding {
		enc = enc.WithPadding(pad)
	}

	ud := newBase64Encoding(L, enc)

	L.Push(ud)

	return 1
}

func encodingBase64Encode(L *lua.LState) int {
	src := L.CheckString(1)

	L.Push(lua.LString(base64.StdEncoding.EncodeToString([]byte(src))))

	return 1
}

func encodingBase64Decode(L *lua.LState) int {
	s := L.CheckString(1)

	bs, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LString(string(bs)))

	return 1
}

var encodingBase64Funcs = map[string]lua.LGFunction{
	"newEncoding": encodingBase64EncodingNew,
	"encode":      encodingBase64Encode,
	"decode":      encodingBase64Decode,
}

var encodingBase64Fields = map[string]lua.LValue{
	"STD_PADDING": lua.LNumber(base64.StdPadding), // Standard padding character
	"NO_PADDING":  lua.LNumber(base64.NoPadding),  // No padding
}
