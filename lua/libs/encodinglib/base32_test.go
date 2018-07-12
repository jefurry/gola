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

func TestBase32EncodingWithPadding(t *testing.T) {
	L := lua.NewState()
	OpenBase32(L)
	defer L.Close()

	code := `
	local base32 = require('encoding.base32')

	local str = "StdEncoding is the standard base32 encoding, as defined in RFC 4648."
	local stdBS = "KN2GIRLOMNXWI2LOM4QGS4ZAORUGKIDTORQW4ZDBOJSCAYTBONSTGMRAMVXGG33ENFXGOLBAMFZSAZDFMZUW4ZLEEBUW4ICSIZBSANBWGQ4C4==="
	local nopadBS = "KN2GIRLOMNXWI2LOM4QGS4ZAORUGKIDTORQW4ZDBOJSCAYTBONSTGMRAMVXGG33ENFXGOLBAMFZSAZDFMZUW4ZLEEBUW4ICSIZBSANBWGQ4C4"
	local encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	local encodeHex = "0123456789ABCDEFGHIJKLMNOPQRSTUV"

	local enc = base32.newEncoding(encodeStd)
	if enc:encode(str) ~= stdBS then
		return false
	end

	local encPad = enc:withPadding(base32.NO_PADDING)
	if encPad:encode(str) ~= nopadBS then
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

func TestBase32EncodingStd(t *testing.T) {
	L := lua.NewState()
	OpenBase32(L)
	defer L.Close()

	code := `
	local base32 = require('encoding.base32')

	local str = "StdEncoding is the standard base32 encoding, as defined in RFC 4648."
	local stdBS = "KN2GIRLOMNXWI2LOM4QGS4ZAORUGKIDTORQW4ZDBOJSCAYTBONSTGMRAMVXGG33ENFXGOLBAMFZSAZDFMZUW4ZLEEBUW4ICSIZBSANBWGQ4C4==="
	local nopadBS = "KN2GIRLOMNXWI2LOM4QGS4ZAORUGKIDTORQW4ZDBOJSCAYTBONSTGMRAMVXGG33ENFXGOLBAMFZSAZDFMZUW4ZLEEBUW4ICSIZBSANBWGQ4C4"
	local encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	local encodeHex = "0123456789ABCDEFGHIJKLMNOPQRSTUV"

	if base32.StdEncoding:encode(str) ~= stdBS then
		return false
	end

	if base32.StdEncoding:decode(stdBS) ~= str then
		return false
	end

	if base32.encode(str) ~= stdBS then
		return false
	end

	if base32.decode(stdBS) ~= str then
		return false
	end

	local enc = base32.newEncoding(encodeStd)
	if enc:encode(str) ~= stdBS then
		return false
	end

	local enc = base32.newEncoding(encodeStd, base32.NO_PADDING)
	if enc:encode(str) ~= nopadBS then
		return false
	end

	local enc = base32.newEncoding(encodeStd, base32.STD_PADDING)
	if enc:encode(str) ~= stdBS then
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

func TestBase32EncodingHex(t *testing.T) {
	L := lua.NewState()
	OpenBase32(L)
	defer L.Close()

	code := `
	local base32 = require('encoding.base32')

	local str = "StdEncoding is the standard base32 encoding, as defined in RFC 4648."
	local hexBS = "ADQ68HBECDNM8QBECSG6ISP0EHK6A83JEHGMSP31E9I20OJ1EDIJ6CH0CLN66RR4D5N6EB10C5PI0P35CPKMSPB441KMS82I8P1I0D1M6GS2S==="
	local nopadBS = "ADQ68HBECDNM8QBECSG6ISP0EHK6A83JEHGMSP31E9I20OJ1EDIJ6CH0CLN66RR4D5N6EB10C5PI0P35CPKMSPB441KMS82I8P1I0D1M6GS2S"
	local encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	local encodeHex = "0123456789ABCDEFGHIJKLMNOPQRSTUV"

	if base32.HexEncoding:encode(str) ~= hexBS then
		return false
	end

	if base32.HexEncoding:decode(hexBS) ~= str then
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
