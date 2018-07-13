// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package bit32 implements bit32 for Lua.
package bit32

import (
	"github.com/BixData/gluabit32"
	"github.com/yuin/gopher-lua"
)

const (
	Bit32LibName = "bit32"
)

func Open(L *lua.LState) {
	L.PreloadModule(Bit32LibName, gluabit32.Loader)
}
