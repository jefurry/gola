// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package encodinglib

import (
	"encoding/hex"
	"github.com/yuin/gopher-lua"
)

const (
	EncodingHexLibName = EncodingLibName + ".hex"
)

func OpenHex(L *lua.LState) {
	L.PreloadModule(EncodingHexLibName, HexLoader)
}

func HexLoader(L *lua.LState) int {
	hexmod := L.SetFuncs(L.NewTable(), encodingHexFuncs)
	L.Push(hexmod)

	encodingRegisterBase32EncodingMetatype(L)

	L.SetFuncs(hexmod, encodingHexFuncs)
	//L.SetField(hexmod, "StdEncoding", newBase32Encoding(L, base32.StdEncoding))
	//L.SetField(hexmod, "HexEncoding", newBase32Encoding(L, base32.HexEncoding))

	for k, v := range encodingBase32Fields {
		hexmod.RawSetString(k, v)
	}

	return 1
}

func encodingHexEncodeToString(L *lua.LState) int {
	src := L.CheckString(1)

	s := hex.EncodeToString([]byte(src))

	L.Push(lua.LString(s))

	return 1
}

func encodingHexDecodeString(L *lua.LState) int {
	s := L.CheckString(1)

	bs, err := hex.DecodeString(s)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LString(string(bs)))

	return 1
}

func encodingHexDump(L *lua.LState) int {
	src := L.CheckString(1)

	s := hex.Dump([]byte(src))

	L.Push(lua.LString(s))

	return 1
}

func encodingHexEncodedLen(L *lua.LState) int {
	n := L.CheckInt(1)

	L.Push(lua.LNumber(hex.EncodedLen(n)))

	return 1
}

func encodingHexDecodedLen(L *lua.LState) int {
	x := L.CheckInt(1)

	L.Push(lua.LNumber(hex.DecodedLen(x)))

	return 1
}

var encodingHexFuncs = map[string]lua.LGFunction{
	"encode":     encodingHexEncodeToString,
	"decode":     encodingHexDecodeString,
	"dump":       encodingHexDump,
	"encodedLen": encodingHexEncodedLen,
	"decodedLen": encodingHexDecodedLen,
}

var encodingHexFields = map[string]lua.LValue{}
