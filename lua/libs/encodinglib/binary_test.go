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

func TestBinary(t *testing.T) {
	L := lua.NewState()
	OpenBinary(L)
	defer L.Close()

	code := `
	local binary = require('encoding.binary')

	local testData = {
		{
			value=45,
			dataType=binary.INT8,
			byteOrder=binary.LITTLE_ENDIAN,
		};
		{
			value=197,
			dataType=binary.INT16,
			byteOrder=binary.BIG_ENDIAN,
		};
		{
			value=3.14157,
			dataType=binary.FLOAT32,
			byteOrder=binary.BIG_ENDIAN,
		};
		{
			value="Good Luck",
			dataType=binary.STRING,
			byteOrder=9,
		};
		{
			value=298,
			dataType=binary.UINT16,
			byteOrder=binary.LITTLE_ENDIAN,
		};
	}

	for _, item in ipairs(testData) do
		local wr = binary.newWriter()

		local ret, msg = wr:write(item.value, item.dataType, item.byteOrder)
		if ret == nil or msg ~= nil then
			return false
		end

		local rdr = binary.newReader(tostring(wr))
		local ret, msg = rdr:read(item.dataType, item.byteOrder)
		if ret == nil or msg ~= nil then
			return false
		end
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
