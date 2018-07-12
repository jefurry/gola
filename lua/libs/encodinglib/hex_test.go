// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package encodinglib

import (
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
	"testing"
)

func TestHex(t *testing.T) {
	L := lua.NewState()
	OpenHex(L)
	defer L.Close()

	code := `
	local hex = require('encoding.hex')

	local str = "Package hex implements hexadecimal encoding and decoding."
	local enstr = "5061636b6167652068657820696d706c656d656e74732068657861646563696d616c20656e636f64696e6720616e64206465636f64696e672e"
	local enstrLen = 14

	if hex.encode(str) ~= enstr then
		return false
	end
	
	if hex.decode(enstr) ~= str then
		return false
	end

	if hex.encodedLen(#str) ~= #enstr then
		return false
	end

	if hex.decodedLen(#enstr) ~= #str then
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
