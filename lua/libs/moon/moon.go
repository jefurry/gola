// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package moon implements moon script for Lua.
package moon

import (
	"github.com/rucuriousyet/gmoonscript"
	"github.com/yuin/gopher-lua"
)

const (
	MoonLibName = "moon"
)

func Open(L *lua.LState) {
	L.PreloadModule(MoonLibName, gmoonscript.Loader)
}
