// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package oslib

import (
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
	"testing"
)

func TestFile(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local os = require('os')

	assert(type(os.chdir) == "function", "os.chdir should be a function")
	assert(type(os.open) == "function", "os.open should be a function")
	assert(type(os.openFile) == "function", "os.openFile should be a function")
	assert(type(os.tempDir) == "function", "os.tempDir should be a function")
	assert(type(os.chmod) == "function", "os.chmod should be a function")
	assert(type(os.chown) == "function", "os.chown should be a function")
	assert(type(os.expand) == "function", "os.expand should be a function")
	assert(type(os.expandEnv) == "function", "os.expandEnv should be a function")
	assert(type(os.lookupEnv) == "function", "os.lookupEnv should be a function")
	assert(type(os.setenv) == "function", "os.setenv should be a function")
	assert(type(os.unsetenv) == "function", "os.unsetenv should be a function")
	assert(type(os.clearenv) == "function", "os.clearenv should be a function")
	assert(type(os.chtimes) == "function", "os.chtimes should be a function")
	assert(type(os.environ) == "function", "os.environ should be a function")
	assert(type(os.truncate) == "function", "os.truncate should be a function")
	assert(type(os.getwd) == "function", "os.getwd should be a function")
	assert(type(os.executable) == "function", "os.executable should be a function")
	assert(type(os.lstat) == "function", "os.lstat should be a function")
	assert(type(os.stat) == "function", "os.stat should be a function")
	assert(type(os.getpagesize) == "function", "os.getpagesize should be a function")
	assert(type(os.sameFile) == "function", "os.sameFile should be a function")
	assert(type(os.lchown) == "function", "os.lchown should be a function")
	assert(type(os.link) == "function", "os.link should be a function")
	assert(type(os.symlink) == "function", "os.symlink should be a function")
	assert(type(os.readlink) == "function", "os.readlink should be a function")
	assert(type(os.mkdir) == "function", "os.mkdir should be a function")
	assert(type(os.remove) == "function", "os.remove should be a function")
	assert(type(os.mkdirAll) == "function", "os.mkdirAll should be a function")
	assert(type(os.removeAll) == "function", "os.removeAll should be a function")
	assert(type(os.isNotExist) == "function", "os.isNotExist should be a function")
	assert(type(os.isExist) == "function", "os.isExist should be a function")
	assert(type(os.getuid) == "function", "os.getuid should be a function")
	assert(type(os.geteuid) == "function", "os.geteuid should be a function")
	assert(type(os.getgid) == "function", "os.getgid should be a function")
	assert(type(os.getgroups) == "function", "os.getgroups should be a function")
	assert(type(os.hostname) == "function", "os.hostname should be a function")
	assert(type(os.getpid) == "function", "os.getpid should be a function")
	assert(type(os.getppid) == "function", "os.getppid should be a function")
	assert(type(os.findProcess) == "function", "os.findProcess should be a function")

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
