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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenNone(t *testing.T) {
	token, err := New(SIGNING_METHOD_NONE)
	if !assert.NoError(t, err, "NewToken should succeed") {
		return
	}

	es, err := token.Signed("")
	if !assert.NoError(t, err, "Signed should succeed") {
		return
	}

	if !assert.NotEqual(t, "", es, "es should be not equals to empty string") {
		return
	}

	newToken, err := Parse(es, func(t *djwt.Token) (interface{}, error) {
		return djwt.UnsafeAllowNoneSignatureType, nil
	})
	if !assert.NoError(t, err, "Parse should succeed") {
		return
	}

	if !assert.NotEqual(t, newToken, nil, t, "newToken should be not equals to nil") {
		return
	}
}

func TestTokenHS(t *testing.T) {
	token, err := New(SIGNING_METHOD_HS256)
	if !assert.NoError(t, err, "NewToken should succeed") {
		return
	}

	key := "any key xxxx"

	es, err := token.Signed(key)
	if !assert.NoError(t, err, "Signed should succeed") {
		return
	}

	if !assert.NotEqual(t, "", es, "es should be not equals to empty string") {
		return
	}

	newToken, err := Parse(es, func(t *djwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if !assert.NoError(t, err, "Parse should succeed") {
		return
	}

	if !assert.NotEqual(t, newToken, nil, t, "newToken should be not equals to nil") {
		return
	}
}

func TestTokenES(t *testing.T) {
	token, err := New(SIGNING_METHOD_ES384)
	if !assert.NoError(t, err, "NewToken should succeed") {
		return
	}

	key := "any key xxxx"

	es, err := token.Signed(key)
	if !assert.Error(t, err, "Signed should succeed") {
		return
	}

	if !assert.Equal(t, "", es, "es should be equals to empty string") {
		return
	}

	_, err = Parse(es, func(t *djwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if !assert.Error(t, err, "Parse should succeed") {
		return
	}
}

func TestTokenRS(t *testing.T) {
	token, err := New(SIGNING_METHOD_RS512)
	if !assert.NoError(t, err, "NewToken should succeed") {
		return
	}

	key := "any key xxxx"

	es, err := token.Signed(key)
	if !assert.Error(t, err, "Signed should succeed") {
		return
	}

	if !assert.Equal(t, "", es, "es should be equals to empty string") {
		return
	}

	_, err = Parse(es, func(t *djwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if !assert.Error(t, err, "Parse should succeed") {
		return
	}
}
