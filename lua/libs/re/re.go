// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package re implements re for Lua.
package re

import (
	"github.com/yuin/gluare"
	"github.com/yuin/gopher-lua"
)

const (
	ReLibName = "re"
)

func Open(L *lua.LState) {
	L.PreloadModule(ReLibName, gluare.Loader)
}
