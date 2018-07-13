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
	"fmt"
	"github.com/yuin/gopher-lua"
)

func newBase32Encoding(L *lua.LState, enc *b32.Encoding) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = enc

	L.SetMetatable(ud, L.GetTypeMetatable(base32EncodingTypeName))

	return ud
}

func checkBase32Encoding(L *lua.LState, n int) *b32.Encoding {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*b32.Encoding); ok {
		return v
	}

	L.ArgError(n, fmt.Sprintf("%s expected, got %s", base32EncodingTypeName, ud.Type()))

	return nil
}
