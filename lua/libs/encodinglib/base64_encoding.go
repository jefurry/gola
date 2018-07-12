// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package encodinglib

import (
	"github.com/yuin/gopher-lua"
)

const (
	encodingBase64EncodingTypeName = EncodingBase64LibName + ".ENCODING*"
)

func encodingBase64EncodeToString(L *lua.LState) int {
	enc := checkBase64Encoding(L, 1)
	src := L.CheckString(2)

	L.Push(lua.LString(enc.EncodeToString([]byte(src))))

	return 1
}

func encodingBase64DecodeString(L *lua.LState) int {
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

func encodingBase64EncodingStrict(L *lua.LState) int {
	enc := checkBase64Encoding(L, 1)

	strictEnc := enc.Strict()
	L.Push(newBase64Encoding(L, strictEnc))

	return 1
}

func encodingRegisterBase64EncodingMetatype(L *lua.LState) {
	// meta table
	mt := L.NewTypeMetatable(encodingBase64EncodingTypeName)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), encodingBase64EncodingFuncs))
}

var encodingBase64EncodingFuncs = map[string]lua.LGFunction{
	"encode": encodingBase64EncodeToString,
	"decode": encodingBase64DecodeString,
	"strict": encodingBase64EncodingStrict,
}
