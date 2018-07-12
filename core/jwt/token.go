// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package jwt

import (
	"github.com/dgrijalva/jwt-go"
)

type (
	Token struct {
		tk *jwt.Token
		mt SigningMethodType
	}
)

func makeToken(method SigningMethod, claims ...jwt.Claims) (*jwt.Token, SigningMethodType, error) {
	m, mt, err := signingMethod(method)
	if err != nil {
		return nil, mt, err
	}

	if len(claims) > 0 {
		return jwt.NewWithClaims(m, claims[0]), mt, nil
	}

	return jwt.New(m), mt, nil
}

func New(method SigningMethod, claims ...jwt.Claims) (*Token, error) {
	tk, mt, err := makeToken(method, claims...)
	if err != nil {
		return nil, err
	}

	return &Token{tk: tk, mt: mt}, nil
}

// Get the complete, signed token
func (t *Token) Signed(key string) (string, error) {
	switch t.mt {
	case SIGNING_METHOD_NONE_TYPE:
		return t.tk.SignedString(jwt.UnsafeAllowNoneSignatureType)
	case SIGNING_METHOD_HS_TYPE:
		return t.tk.SignedString([]byte(key))
	case SIGNING_METHOD_ES_TYPE:
		k, err := jwt.ParseECPrivateKeyFromPEM([]byte(key))
		if err != nil {
			return "", err
		}

		return t.tk.SignedString(k)
	case SIGNING_METHOD_RS_TYPE:
		k, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(key))
		if err != nil {
			return "", err
		}

		return t.tk.SignedString(k)
	}

	return "", ErrInvalidSigningMethod
}

func (t *Token) GetToken() *jwt.Token {
	return t.tk
}

func (t *Token) GetClaims() jwt.Claims {
	return t.tk.Claims
}

func Parse(tokenString string, keyFunc jwt.Keyfunc) (*Token, error) {
	tk, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return nil, err
	}

	return &Token{tk: tk, mt: methodType(tk.Method)}, nil
}

func ParseWithClaims(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*Token, error) {
	tk, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return nil, err
	}

	return &Token{tk: tk, mt: methodType(tk.Method)}, nil
}

func methodType(method jwt.SigningMethod) SigningMethodType {
	alg := method.Alg()
	switch alg {
	case "none":
		return SIGNING_METHOD_NONE_TYPE
	case "HS256", "HS384", "HS512":
		return SIGNING_METHOD_HS_TYPE
	case "ES256", "ES384", "ES512":
		return SIGNING_METHOD_ES_TYPE
	case "RS256", "RS384", "RS512":
		return SIGNING_METHOD_RS_TYPE
	}

	return SIGNING_METHOD_INVALID_TYPE
}
