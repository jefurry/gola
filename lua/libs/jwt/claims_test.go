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
	"github.com/yuin/gopher-lua"
	"testing"
)

// Normal
func TestJwtClaims_1(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local jwt = require('jwt')
	local os = require('os')

	local t = os.time()
	local exp = 5

	local claims, msg = jwt.newClaims{nbf=t-exp, exp=t+exp, iat=t, id=1, iss="Jeff", sub="Gola framework", aud="Gola"}
	if claims == nil or msg ~= nil then
		return false
	end

	if claims:valid() ~= true then
		return false
	end

	return true
	`

	err := L.DoString(code)
	if !assert.NoError(t, err, `L.DoString should succeed`) {
		return
	}

	if !assert.Equal(t, 1, L.GetTop(), "L.GetTop mismatching") {
		return
	}

	ret := L.Get(-1)
	if !assert.Equal(t, lua.LTBool, ret.Type(), "type mismatching") {
		return
	}

	if !assert.Equal(t, lua.LTrue, ret.(lua.LBool), "value mismatching") {
		return
	}
}

func TestJwtClaims_2(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local jwt = require('jwt')
	local os = require('os')

	local t = os.time()
	local exp = 5

	local claims, msg = jwt.newClaims{nbf=t+exp, exp=t+exp, iat=t, id=1, iss="Jeff", sub="Gola framework", aud="Gola"}
	if claims ~= nil or msg == nil then
		return false
	end

	if msg ~= "Token is not valid yet" then
		return false
	end

	return true
	`

	err := L.DoString(code)
	if !assert.NoError(t, err, `L.DoString should succeed`) {
		return
	}

	if !assert.Equal(t, 1, L.GetTop(), "L.GetTop mismatching") {
		return
	}

	ret := L.Get(-1)
	if !assert.Equal(t, lua.LTBool, ret.Type(), "type mismatching") {
		return
	}

	if !assert.Equal(t, lua.LTrue, ret.(lua.LBool), "value mismatching") {
		return
	}
}

func TestJwtClaims_3(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local jwt = require('jwt')
	local os = require('os')

	local t = os.time()
	local exp = 5

	local claims, msg = jwt.newClaims{nbf=t-exp, exp=t-exp, iat=t, id=1, iss="Jeff", sub="Gola framework", aud="Gola"}
	if claims ~= nil or msg == nil then
		return false
	end

	if msg ~= "Token is expired" then
		return false
	end

	return true
	`

	err := L.DoString(code)
	if !assert.NoError(t, err, `L.DoString should succeed`) {
		return
	}

	if !assert.Equal(t, 1, L.GetTop(), "L.GetTop mismatching") {
		return
	}

	ret := L.Get(-1)
	if !assert.Equal(t, lua.LTBool, ret.Type(), "type mismatching") {
		return
	}

	if !assert.Equal(t, lua.LTrue, ret.(lua.LBool), "value mismatching") {
		return
	}
}

func TestJwtClaims_4(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local jwt = require('jwt')
	local os = require('os')

	local t = os.time()
	local exp = 5

	local claims, msg = jwt.newClaims{nbf=t-exp, exp=t+exp, iat=t+exp, id=1, iss="Jeff", sub="Gola framework", aud="Gola"}
	if claims ~= nil or msg == nil then
		return false
	end

	if msg ~= "Token used before issued" then
		return false
	end

	return true
	`

	err := L.DoString(code)
	if !assert.NoError(t, err, `L.DoString should succeed`) {
		return
	}

	if !assert.Equal(t, 1, L.GetTop(), "L.GetTop mismatching") {
		return
	}

	ret := L.Get(-1)
	if !assert.Equal(t, lua.LTBool, ret.Type(), "type mismatching") {
		return
	}

	if !assert.Equal(t, lua.LTrue, ret.(lua.LBool), "value mismatching") {
		return
	}
}
