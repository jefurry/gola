// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package reng implements Rule Engine.
package reng

import (
	"context"
	golua "github.com/jefurry/gola/lua"
	"github.com/jefurry/gola/lua/pm"
	"github.com/yuin/gopher-lua"
)

var luaLibs = []golua.LuaLib{
	golua.LuaLib{lua.BaseLibName, lua.OpenBase},
	golua.LuaLib{lua.TabLibName, lua.OpenTable},
	golua.LuaLib{lua.StringLibName, lua.OpenString},
	golua.LuaLib{lua.MathLibName, lua.OpenMath},
}

func whenNew(L *lua.LState) error {
	for _, lib := range luaLibs {
		L.Push(L.NewFunction(lib.LibFunc))
		L.Push(lua.LString(lib.LibName))
		L.Call(1, 0)
	}

	excludeBaseLibFuncs(L)

	return nil
}

func Default(ctx context.Context) (*pm.LPM, error) {
	options := lua.Options{}
	options.SkipOpenLibs = true

	config, err := pm.NewConfig(pm.DefaultMaxNum, pm.DefaultStartNum,
		pm.DefaultMaxRequest, pm.DefaultRequestTerminateTimeout, pm.DefaultIdleTimeout, options)
	if err != nil {
		return nil, err
	}

	return pm.New(ctx, config, whenNew)
}

func New(ctx context.Context, config *pm.Config) (*pm.LPM, error) {
	if config == nil {
		return Default(ctx)
	}

	options := config.Options()
	options.SkipOpenLibs = true
	config.SetOptions(options)

	return pm.New(ctx, config, whenNew)
}

func excludeBaseLibFuncs(L *lua.LState) {
	g := L.GetGlobal("_G").(*lua.LTable)

	for _, v := range baseLibFuncs {
		g.RawSetString(v, lua.LNil)
	}
}

// keeps "next", "select", "tonumber", "tostring"
// type", "unpack", "ipairs", "pairs", "assert"
var baseLibFuncs = []string{
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
