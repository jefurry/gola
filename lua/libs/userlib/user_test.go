// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package userlib

import (
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
	"testing"
)

func TestExec(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local user = require('os/user')

	assert(type(user) == "table", "user should be a table")
	assert(type(user.current) == "function", "user.current should be a function")
	assert(type(user.lookup) == "function", "user.lookup should be a function")
	assert(type(user.lookupId) == "function", "user.lookupId should be a function")
	assert(type(user.lookupGroup) == "function", "user.lookupGroup should be a function")
	assert(type(user.lookupGroupId) == "function", "user.lookupGroupId should be a function")

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
