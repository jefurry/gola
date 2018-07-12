// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package jwtlib

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/yuin/gopher-lua"
)

const (
	jwtClaimsTypeName = JwtLibName + ".CLAIMS*"
)

type (
	jwtClaims struct {
		mc jwt.MapClaims
	}
)

func jwtClaimsNew(L *lua.LState) int {
	val := L.CheckTable(1)

	claims := make(jwt.MapClaims, 7)
	val.ForEach(func(key, value lua.LValue) {
		if key.Type() == lua.LTString {
			k, _ := key.(lua.LString)
			switch value.(type) {
			case lua.LNumber:
				v, _ := value.(lua.LNumber)
				claims[string(k)] = float64(v)
			default:
				claims[string(k)] = value
			}
		}
	})

	if err := claims.Valid(); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	ud := newClaims(L, &jwtClaims{mc: claims})

	L.Push(ud)

	return 1
}

func jwtClaimsVerifyAudience(L *lua.LState) int {
	claims := checkClaims(L, 1)
	cmp := L.CheckString(2)
	req := L.CheckBool(3)

	L.Push(lua.LBool(claims.mc.VerifyAudience(cmp, req)))

	return 1
}

func jwtClaimsVerifyExpiresAt(L *lua.LState) int {
	claims := checkClaims(L, 1)
	cmp := L.CheckInt(2)
	req := L.CheckBool(3)

	L.Push(lua.LBool(claims.mc.VerifyExpiresAt(int64(cmp), req)))

	return 1
}

func jwtClaimsVerifyIssuedAt(L *lua.LState) int {
	claims := checkClaims(L, 1)
	cmp := L.CheckInt(2)
	req := L.CheckBool(3)

	L.Push(lua.LBool(claims.mc.VerifyIssuedAt(int64(cmp), req)))

	return 1
}

func jwtClaimsVerifyIssuer(L *lua.LState) int {
	claims := checkClaims(L, 1)
	cmp := L.CheckString(1)
	req := L.CheckBool(2)

	L.Push(lua.LBool(claims.mc.VerifyIssuer(cmp, req)))

	return 1
}

func jwtClaimsVerifyNotBefore(L *lua.LState) int {
	claims := checkClaims(L, 1)
	cmp := L.CheckInt(2)
	req := L.CheckBool(3)

	L.Push(lua.LBool(claims.mc.VerifyNotBefore(int64(cmp), req)))

	return 1
}

func jwtClaimsValid(L *lua.LState) int {
	claims := checkClaims(L, 1)

	if err := claims.mc.Valid(); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func jwtRegisterClaimsMetatype(L *lua.LState) {
	// meta table
	mt := L.NewTypeMetatable(jwtClaimsTypeName)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), jwtClaimsFuncs))
}

var jwtClaimsFuncs = map[string]lua.LGFunction{
	"verifyAudience":  jwtClaimsVerifyAudience,
	"verifyExpiresAt": jwtClaimsVerifyExpiresAt,
	"verifyIssuedAt":  jwtClaimsVerifyIssuedAt,
	"verifyIssuer":    jwtClaimsVerifyIssuer,
	"verifyNotBefore": jwtClaimsVerifyNotBefore,
	"valid":           jwtClaimsValid,
}
