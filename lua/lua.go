// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package lua implements Lua features.
package lua

import (
	"github.com/yuin/gopher-lua"
)

type (
	LuaLib struct {
		LibName string
		LibFunc lua.LGFunction
	}
)

func ExcludeBaseLib(L *lua.LState) {
	g := L.GetGlobal("_G").(*lua.LTable)

	for _, v := range ExcludeBaseLibFuncs {
		g.RawSetString(v, lua.LNil)
	}
}

// keeps "next", "select", "tonumber", "tostring"
// type", "unpack", "ipairs", "pairs", "assert"
var ExcludeBaseLibFuncs = []string{
	//"assert",
	"collectgarbage",
	"dofile",
	"error",
	"getfenv",
	"getmetatable",
	"load",
	"loadfile",
	"loadstring",
	//"next",
	"pcall",
	"print",
	"rawequal",
	"rawget",
	"rawset",
	//"select",
	"_printregs",
	"setfenv",
	"setmetatable",
	//"tonumber",
	//"tostring",
	//"type",
	//"unpack",
	"xpcall",
	// loadlib
	"module",
	"require",
	// hidden features
	"newproxy",
}
