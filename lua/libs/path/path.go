// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package path implements path for Lua.
package path

import (
	"github.com/jefurry/gola/lua/libs/path/filepath"
	"github.com/yuin/gopher-lua"
	ppath "path"
)

const (
	PathLibName = "path"
)

func Open(L *lua.LState) {
	L.PreloadModule(PathLibName, Loader)

	filepath.Open(L)
}

func Loader(L *lua.LState) int {
	pathmod := L.SetFuncs(L.NewTable(), pathFuncs)
	L.Push(pathmod)

	return 1
}

var pathFuncs = map[string]lua.LGFunction{
	"base":  pathBase,
	"clean": pathClean,
	"dir":   pathDir,
	"ext":   pathExt,
	"isAbs": pathIsAbs,
	"join":  pathJoin,
	"match": pathMatch,
	"split": pathSplit,
}

func pathBase(L *lua.LState) int {
	p := L.CheckString(1)

	L.Push(lua.LString(ppath.Base(p)))

	return 1
}

func pathClean(L *lua.LState) int {
	p := L.CheckString(1)

	L.Push(lua.LString(ppath.Clean(p)))

	return 1
}

func pathDir(L *lua.LState) int {
	p := L.CheckString(1)

	L.Push(lua.LString(ppath.Dir(p)))

	return 1
}

func pathExt(L *lua.LState) int {
	p := L.CheckString(1)

	L.Push(lua.LString(ppath.Ext(p)))

	return 1
}

func pathIsAbs(L *lua.LState) int {
	p := L.CheckString(1)

	L.Push(lua.LBool(ppath.IsAbs(p)))

	return 1
}

func pathJoin(L *lua.LState) int {
	top := L.GetTop()
	if top < 1 {
		L.Push(lua.LString(""))

		return 1
	}

	ps := make([]string, 0, top)
	for i := 1; i <= top; i++ {
		ps = append(ps, L.CheckString(i))
	}

	L.Push(lua.LString(ppath.Join(ps...)))

	return 1
}

func pathMatch(L *lua.LState) int {
	pattern := L.CheckString(1)
	name := L.CheckString(2)

	matched, err := ppath.Match(pattern, name)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LBool(matched))

	return 1
}

func pathSplit(L *lua.LState) int {
	p := L.CheckString(1)

	dir, file := ppath.Split(p)

	L.Push(lua.LString(dir))
	L.Push(lua.LString(file))

	return 2
}
