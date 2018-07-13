// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package jwt implements json web token for Lua.
package jwt

import (
	djwt "github.com/dgrijalva/jwt-go"
	gjwt "github.com/jefurry/gola/core/jwt"
	"github.com/yuin/gopher-lua"
)

const (
	JwtLibName = "jwt"
)

func Open(L *lua.LState) {
	L.PreloadModule(JwtLibName, Loader)
}

func Loader(L *lua.LState) int {
	jwtmod := L.SetFuncs(L.NewTable(), jwtFuncs)
	L.Push(jwtmod)

	jwtRegisterClaimsMetatype(L)
	jwtRegisterTokenMetatype(L)

	for k, v := range jwtSigningMethodTypeFields {
		jwtmod.RawSetString(k, lua.LNumber(v))
	}

	for k, v := range jwtSigningMethodFields {
		jwtmod.RawSetString(k, lua.LNumber(v))
	}

	return 1
}

func jwtParse(L *lua.LState) int {
	tokenString := L.CheckString(1)
	key := L.OptString(2, "")

	token, err := gjwt.Parse(tokenString, func(t *djwt.Token) (interface{}, error) {
		alg := t.Method.Alg()

		switch alg {
		case "none":
			return djwt.UnsafeAllowNoneSignatureType, nil
		case "HS256", "HS384", "HS512":
			return []byte(key), nil
		case "ES256", "ES384", "ES512":
			return djwt.ParseECPublicKeyFromPEM([]byte(key))
		case "RS256", "RS384", "RS512":
			return djwt.ParseRSAPublicKeyFromPEM([]byte(key))
		case "PS256", "PS384", "PS512":
			return djwt.ParseRSAPublicKeyFromPEM([]byte(key))
		}

		return []byte(key), nil
	})

	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	ud := newToken(L, token)

	L.Push(ud)

	return 1
}

func jwtParseWithClaims(L *lua.LState) int {
	tokenString := L.CheckString(1)
	claims := checkClaims(L, 2)
	key := L.OptString(3, "")

	token, err := gjwt.ParseWithClaims(tokenString, claims.mc, func(t *djwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	ud := newToken(L, token)

	L.Push(ud)

	return 1
}

var jwtFuncs = map[string]lua.LGFunction{
	"newClaims":       jwtClaimsNew,
	"newToken":        jwtTokenNew,
	"parse":           jwtParse,
	"parseWithClaims": jwtParseWithClaims,
}

var jwtSigningMethodTypeFields = map[string]gjwt.SigningMethodType{
	"SIGNING_METHOD_INVALID_TYPE": gjwt.SIGNING_METHOD_INVALID_TYPE,
	"SIGNING_METHOD_NONE_TYPE":    gjwt.SIGNING_METHOD_NONE_TYPE,
	"SIGNING_METHOD_HS_TYPE":      gjwt.SIGNING_METHOD_HS_TYPE,
	"SIGNING_METHOD_ES_TYPE":      gjwt.SIGNING_METHOD_ES_TYPE,
	"SIGNING_METHOD_RS_TYPE":      gjwt.SIGNING_METHOD_RS_TYPE,
	"SIGNING_METHOD_PS_TYPE":      gjwt.SIGNING_METHOD_PS_TYPE,
}

var jwtSigningMethodFields = map[string]gjwt.SigningMethod{
	"SIGNING_METHOD_INVALID": gjwt.SIGNING_METHOD_INVALID,

	"SIGNING_METHOD_NONE": gjwt.SIGNING_METHOD_NONE,

	"SIGNING_METHOD_HS256": gjwt.SIGNING_METHOD_HS256,
	"SIGNING_METHOD_HS384": gjwt.SIGNING_METHOD_HS384,
	"SIGNING_METHOD_HS512": gjwt.SIGNING_METHOD_HS512,

	"SIGNING_METHOD_ES256": gjwt.SIGNING_METHOD_ES256,
	"SIGNING_METHOD_ES384": gjwt.SIGNING_METHOD_ES384,
	"SIGNING_METHOD_ES512": gjwt.SIGNING_METHOD_ES512,

	"SIGNING_METHOD_RS256": gjwt.SIGNING_METHOD_RS256,
	"SIGNING_METHOD_RS384": gjwt.SIGNING_METHOD_RS384,
	"SIGNING_METHOD_RS512": gjwt.SIGNING_METHOD_RS512,

	"SIGNING_METHOD_PS256": gjwt.SIGNING_METHOD_PS256,
	"SIGNING_METHOD_PS384": gjwt.SIGNING_METHOD_PS384,
	"SIGNING_METHOD_PS512": gjwt.SIGNING_METHOD_PS512,
}
