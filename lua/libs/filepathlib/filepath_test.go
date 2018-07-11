// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package filepathlib

import (
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
	"testing"
)

func TestFilepath(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local filepath = require('path.filepath')

	assert(type(filepath) == "table", "filepath should be a table")
	assert(type(filepath.abs) == "function", "filepath.abs should be a function")
	assert(type(filepath.base) == "function", "filepath.base should be a function")
	assert(type(filepath.clean) == "function", "filepath.clean should be a function")
	assert(type(filepath.dir) == "function", "filepath.dir should be a function")
	assert(type(filepath.evalSymlinks) == "function", "filepath.evalSymlinks should be a function")
	assert(type(filepath.ext) == "function", "filepath.ext should be a function")
	assert(type(filepath.fromSlash) == "function", "filepath.fromSlash should be a function")
	assert(type(filepath.glob) == "function", "filepath.glob should be a function")
	assert(type(filepath.hasPrefix) == "function", "filepath.hasPrefix should be a function")
	assert(type(filepath.isAbs) == "function", "filepath.isAbs should be a function")
	assert(type(filepath.join) == "function", "filepath.join should be a function")
	assert(type(filepath.match) == "function", "filepath.match should be a function")
	assert(type(filepath.rel) == "function", "filepath.rel should be a function")
	assert(type(filepath.split) == "function", "filepath.split should be a function")
	assert(type(filepath.splitList) == "function", "filepath.splitList should be a function")
	assert(type(filepath.toSlash) == "function", "filepath.toSlash should be a function")
	assert(type(filepath.volumeName) == "function", "filepath.volumeName should be a function")

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
