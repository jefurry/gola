// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package jwt

import (
	djwt "github.com/dgrijalva/jwt-go"
	gjwt "github.com/jefurry/gola/core/jwt"
	"github.com/yuin/gopher-lua"
)

const (
	jwtTokenTypeName = JwtLibName + ".TOKEN*"
)

func jwtTokenNew(L *lua.LState) int {
	method := L.CheckInt(1)
	val := L.OptUserData(2, nil)

	var claims djwt.MapClaims
	if val == nil {
		claims = make(djwt.MapClaims, 0)
	} else {
		c, ok := val.Value.(*jwtClaims)
		if !ok {
			L.Push(lua.LNil)
			L.Push(lua.LString("invalid claims"))

			return 2
		}

		claims = c.mc
	}

	token, err := gjwt.New(gjwt.SigningMethod(method), claims)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	ud := newToken(L, token)

	L.Push(ud)

	return 1
}

func jwtTokenGetClaims(L *lua.LState) int {
	token := checkToken(L, 1)

	tb := L.CreateTable(0, 7)
	claims := token.GetClaims()
	if mc, ok := claims.(djwt.MapClaims); ok {
		for key, value := range mc {
			var v lua.LValue
			switch value.(type) {
			case string:
				val, _ := value.(string)
				v = lua.LString(val)
			case int8, int16, int32, int64, int, uint8, uint16,
				uint32, uint64, uint, float32, float64:
				val, _ := value.(float64)
				v = lua.LNumber(val)
			case bool:
				val, _ := value.(bool)
				v = lua.LBool(val)
			default:
				continue
			}

			tb.RawSetString(string(key), v)
		}
	}

	L.Push(tb)

	return 1
}

func jwtTokenSigned(L *lua.LState) int {
	token := checkToken(L, 1)
	key := L.OptString(2, "")
	password := L.OptString(3, "")

	s, err := token.Signed(key, password)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LString(s))

	return 1
}

func jwtTokenValid(L *lua.LState) int {
	token := checkToken(L, 1)

	L.Push(lua.LBool(token.Valid()))

	return 1
}

func jwtRegisterTokenMetatype(L *lua.LState) {
	// meta table
	mt := L.NewTypeMetatable(jwtTokenTypeName)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), jwtTokenFuncs))
}

var jwtTokenFuncs = map[string]lua.LGFunction{
	"signed":    jwtTokenSigned,
	"getClaims": jwtTokenGetClaims,
	"valid":     jwtTokenValid,
}
