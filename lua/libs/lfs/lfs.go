// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package lfs implements lfs for Lua.
package lfs

import (
	glualfs "github.com/layeh/gopher-lfs"
	"github.com/yuin/gopher-lua"
)

const (
	LfsLibName = "lfs"
)

func Open(L *lua.LState) {
	glualfs.Preload(L)
}
