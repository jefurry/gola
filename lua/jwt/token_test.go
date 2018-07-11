// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package jwt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenNone(t *testing.T) {
	token, err := New(SIGNING_METHOD_NONE)
	if !assert.NoError(t, err, "NewToken should succeed") {
		return
	}

	s, err := token.Signed("")
	if !assert.NoError(t, err, "Signed should succeed") {
		return
	}

	if !assert.NotEqual(t, "", s, "s should be not equals to empty string") {
		return
	}
}

func TestTokenHS(t *testing.T) {
	token, err := New(SIGNING_METHOD_HS256)
	if !assert.NoError(t, err, "NewToken should succeed") {
		return
	}

	s, err := token.Signed("any key xxxx")
	if !assert.NoError(t, err, "Signed should succeed") {
		return
	}

	if !assert.NotEqual(t, "", s, "s should be not equals to empty string") {
		return
	}
}

func TestTokenES(t *testing.T) {
	token, err := New(SIGNING_METHOD_ES384)
	if !assert.NoError(t, err, "NewToken should succeed") {
		return
	}

	s, err := token.Signed("any key xxxx")
	if !assert.Error(t, err, "Signed should succeed") {
		return
	}

	if !assert.Equal(t, "", s, "s should be equals to empty string") {
		return
	}
}

func TestTokenRS(t *testing.T) {
	token, err := New(SIGNING_METHOD_RS512)
	if !assert.NoError(t, err, "NewToken should succeed") {
		return
	}

	s, err := token.Signed("any key xxxx")
	if !assert.Error(t, err, "Signed should succeed") {
		return
	}

	if !assert.Equal(t, "", s, "s should be equals to empty string") {
		return
	}
}
