// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package binary

import (
	"fmt"
	"github.com/yuin/gopher-lua"
)

func newBinaryReader(L *lua.LState, br *binaryReader) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = br

	L.SetMetatable(ud, L.GetTypeMetatable(binaryReaderTypeName))

	return ud
}

func checkBinaryReader(L *lua.LState, n int) *binaryReader {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*binaryReader); ok {
		return v
	}

	L.ArgError(n, fmt.Sprintf("%s expected, got %s", binaryReaderTypeName, ud.Type()))

	return nil
}

func newBinaryWriter(L *lua.LState, wr *binaryWriter) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = wr

	L.SetMetatable(ud, L.GetTypeMetatable(binaryWriterTypeName))

	return ud
}

func checkBinaryWriter(L *lua.LState, n int) *binaryWriter {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*binaryWriter); ok {
		return v
	}

	L.ArgError(n, fmt.Sprintf("%s expected, got %s", binaryWriterTypeName, ud.Type()))

	return nil
}

func checkPackVal(L *lua.LState, n int) lua.LValue {
	L.CheckTypes(n, lua.LTNumber, lua.LTBool, lua.LTString)

	return L.CheckAny(n)
}
