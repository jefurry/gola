// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package jwtlib

import (
	"fmt"
	gjwt "github.com/jefurry/gola/core/jwt"
	"github.com/yuin/gopher-lua"
)

func newToken(L *lua.LState, token *gjwt.Token) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = token

	L.SetMetatable(ud, L.GetTypeMetatable(jwtTokenTypeName))

	return ud
}

func checkToken(L *lua.LState, n int) *gjwt.Token {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*gjwt.Token); ok {
		return v
	}

	L.ArgError(n, fmt.Sprintf("%s expected, got %s", jwtTokenTypeName, ud.Type()))

	return nil
}

func newClaims(L *lua.LState, claims *jwtClaims) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = claims

	L.SetMetatable(ud, L.GetTypeMetatable(jwtClaimsTypeName))

	return ud
}

func checkClaims(L *lua.LState, n int) *jwtClaims {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*jwtClaims); ok {
		return v
	}

	L.ArgError(n, fmt.Sprintf("%s expected, got %s", jwtClaimsTypeName, ud.Type()))

	return nil
}
