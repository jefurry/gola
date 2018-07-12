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

func TestEncodingStd(t *testing.T) {
	L := lua.NewState()
	OpenBase64(L)
	defer L.Close()

	code := `
	local base64 = require('encoding.base64')

	local encodeStd = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	local str = "StdEncoding is the standard base64 encoding, as defined in RFC 4648."
	local stdBS = "U3RkRW5jb2RpbmcgaXMgdGhlIHN0YW5kYXJkIGJhc2U2NCBlbmNvZGluZywgYXMgZGVmaW5lZCBpbiBSRkMgNDY0OC4="
	local nopadBS = "U3RkRW5jb2RpbmcgaXMgdGhlIHN0YW5kYXJkIGJhc2U2NCBlbmNvZGluZywgYXMgZGVmaW5lZCBpbiBSRkMgNDY0OC4"

	if base64.StdEncoding:encode(str) ~= stdBS then
		return false
	end

	if base64.StdEncoding:decode(stdBS) ~= str then
		return false
	end

	if base64.encode(str) ~= stdBS then
		return false
	end

	if base64.decode(stdBS) ~= str then
		return false
	end

	local enc = base64.newEncoding(encodeStd)
	if enc:encode(str) ~= stdBS then
		return false
	end

	local enc = base64.newEncoding(encodeStd, base64.NO_PADDING)
	if enc:encode(str) ~= nopadBS then
		return false
	end

	local enc = base64.newEncoding(encodeStd, base64.STD_PADDING)
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

func TestEncodingURL(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local base64 = require('encoding.base64')

	local encodeURL = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	local str = "StdEncoding is - the standard base64_encoding, as defined in RFC 4648."
	local urlBS = "U3RkRW5jb2RpbmcgaXMgLSB0aGUgc3RhbmRhcmQgYmFzZTY0X2VuY29kaW5nLCBhcyBkZWZpbmVkIGluIFJGQyA0NjQ4Lg=="
	local nopadBS = "U3RkRW5jb2RpbmcgaXMgLSB0aGUgc3RhbmRhcmQgYmFzZTY0X2VuY29kaW5nLCBhcyBkZWZpbmVkIGluIFJGQyA0NjQ4Lg"

	if base64.URLEncoding:encode(str) ~= urlBS then
		return false
	end

	if base64.URLEncoding:decode(urlBS) ~= str then
		return false
	end

	if base64.encode(str) ~= urlBS then
		return false
	end

	if base64.decode(urlBS) ~= str then
		return false
	end

	local enc = base64.newEncoding(encodeURL, base64.NO_PADDING)
	if enc:encode(str) ~= nopadBS then
		return false
	end

	local enc = base64.newEncoding(encodeURL, base64.STD_PADDING)
	if enc:encode(str) ~= urlBS then
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

func TestEncodingRawStd(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local base64 = require('encoding.base64')

	local encodeURL = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	local str = "StdEncoding is - the standard base64_encoding, as defined in RFC 4648."
	local nopadBS = "U3RkRW5jb2RpbmcgaXMgLSB0aGUgc3RhbmRhcmQgYmFzZTY0X2VuY29kaW5nLCBhcyBkZWZpbmVkIGluIFJGQyA0NjQ4Lg"

	if base64.RawStdEncoding:encode(str) ~= nopadBS then
		return false
	end

	if base64.RawStdEncoding:decode(nopadBS) ~= str then
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

func TestEncodingRawURL(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local base64 = require('encoding.base64')

	local encodeURL = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	local str = "StdEncoding is - the standard base64_encoding, as defined in RFC 4648."
	local nopadBS = "U3RkRW5jb2RpbmcgaXMgLSB0aGUgc3RhbmRhcmQgYmFzZTY0X2VuY29kaW5nLCBhcyBkZWZpbmVkIGluIFJGQyA0NjQ4Lg"

	if base64.RawURLEncoding:encode(str) ~= nopadBS then
		return false
	end

	if base64.RawURLEncoding:decode(nopadBS) ~= str then
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
