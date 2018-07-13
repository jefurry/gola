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
	"fmt"
	"github.com/yuin/gopher-lua"
)

func newBase64Encoding(L *lua.LState, enc *b64.Encoding) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = enc

	L.SetMetatable(ud, L.GetTypeMetatable(base64EncodingTypeName))

	return ud
}

func checkBase64Encoding(L *lua.LState, n int) *b64.Encoding {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*b64.Encoding); ok {
		return v
	}

	L.ArgError(n, fmt.Sprintf("%s expected, got %s", base64EncodingTypeName, ud.Type()))

	return nil
}
