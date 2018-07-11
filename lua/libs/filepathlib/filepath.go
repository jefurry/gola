// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package filepathlib implements path/filepath for Lua.
package filepathlib

import (
	"github.com/jefurry/gola/lua/cb"
	"github.com/yuin/gopher-lua"
	"os"
	"path/filepath"
)

const (
	FilepathLibName = "path.filepath"
)

func Open(L *lua.LState) {
	L.PreloadModule(FilepathLibName, Loader)
}

func Loader(L *lua.LState) int {
	filepathmod := L.SetFuncs(L.NewTable(), filepathFuncs)
	L.Push(filepathmod)

	for k, v := range filepathFields {
		filepathmod.RawSetString(k, v)
	}

	return 1
}

var filepathFuncs = map[string]lua.LGFunction{
	"abs":          filepathAbs,
	"base":         filepathBase,
	"clean":        filepathClean,
	"dir":          filepathDir,
	"evalSymlinks": filepathEvalSymlinks,
	"ext":          filepathExt,
	"fromSlash":    filepathFromSlash,
	"glob":         filepathGlob,
	"hasPrefix":    filepathHasPrefix,
	"isAbs":        filepathIsAbs,
	"join":         filepathJoin,
	"match":        filepathMatch,
	"rel":          filepathRel,
	"split":        filepathSplit,
	"splitList":    filepathSplitList,
	"toSlash":      filepathToSlash,
	"volumeName":   filepathVolumeName,
}

var filepathFields = map[string]lua.LValue{
	"Separator":     lua.LString(filepath.Separator),
	"ListSeparator": lua.LString(filepath.ListSeparator),
}

func filepathAbs(L *lua.LState) int {
	p := L.CheckString(1)

	s, err := filepath.Abs(p)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LString(s))

	return 1
}

func filepathBase(L *lua.LState) int {
	p := L.CheckString(1)

	L.Push(lua.LString(filepath.Base(p)))

	return 1
}

func filepathClean(L *lua.LState) int {
	p := L.CheckString(1)

	L.Push(lua.LString(filepath.Clean(p)))

	return 1
}

func filepathDir(L *lua.LState) int {
	p := L.CheckString(1)

	L.Push(lua.LString(filepath.Dir(p)))

	return 1
}

func filepathEvalSymlinks(L *lua.LState) int {
	p := L.CheckString(1)

	s, err := filepath.EvalSymlinks(p)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LString(s))

	return 1
}

func filepathExt(L *lua.LState) int {
	p := L.CheckString(1)

	L.Push(lua.LString(filepath.Ext(p)))

	return 1
}

func filepathFromSlash(L *lua.LState) int {
	p := L.CheckString(1)

	L.Push(lua.LString(filepath.FromSlash(p)))

	return 1
}

func filepathGlob(L *lua.LState) int {
	pattern := L.CheckString(1)

	matches, err := filepath.Glob(pattern)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	ms := L.CreateTable(len(matches), 0)
	for i, m := range matches {
		ms.RawSetInt(i, lua.LString(m))
	}

	L.Push(ms)

	return 1
}

func filepathHasPrefix(L *lua.LState) int {
	p := L.CheckString(1)
	prefix := L.CheckString(2)

	L.Push(lua.LBool(filepath.HasPrefix(p, prefix)))

	return 1
}

func filepathIsAbs(L *lua.LState) int {
	p := L.CheckString(1)

	L.Push(lua.LBool(filepath.IsAbs(p)))

	return 1
}

func filepathJoin(L *lua.LState) int {
	top := L.GetTop()
	if top < 1 {
		L.Push(lua.LString(""))

		return 1
	}

	ps := make([]string, 0, top)
	for i := 1; i <= top; i++ {
		ps = append(ps, L.CheckString(i))
	}

	L.Push(lua.LString(filepath.Join(ps...)))

	return 1
}

func filepathMatch(L *lua.LState) int {
	pattern := L.CheckString(1)
	name := L.CheckString(2)

	matched, err := filepath.Match(pattern, name)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LBool(matched))

	return 1
}

func filepathRel(L *lua.LState) int {
	basepath := L.CheckString(1)
	targpath := L.CheckString(2)

	s, err := filepath.Rel(basepath, targpath)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LString(s))

	return 1
}

func filepathSplit(L *lua.LState) int {
	p := L.CheckString(1)

	dir, file := filepath.Split(p)

	L.Push(lua.LString(dir))
	L.Push(lua.LString(file))

	return 2
}

func filepathSplitList(L *lua.LState) int {
	p := L.CheckString(1)

	list := filepath.SplitList(p)
	tb := L.CreateTable(len(list), 0)

	for i, v := range list {
		tb.RawSetInt(i, lua.LString(v))
	}

	L.Push(tb)

	return 1
}

func filepathToSlash(L *lua.LState) int {
	p := L.CheckString(1)

	L.Push(lua.LString(filepath.ToSlash(p)))

	return 1
}

func filepathVolumeName(L *lua.LState) int {
	p := L.CheckString(1)

	L.Push(lua.LString(filepath.VolumeName(p)))

	return 1
}

func filepathWalk(L *lua.LState) int {
	root := L.CheckString(1)
	val := L.CheckAny(2)

	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		callable, er := cb.New(L, val)
		if er != nil {
			return er
		}

		objFn, e := callable.ObjFn(L)
		if e != nil {
			return e
		}

		infoTb := L.CreateTable(0, 6)
		infoTb.RawSetString("name", lua.LString(info.Name()))
		infoTb.RawSetString("size", lua.LNumber(info.Size()))
		infoTb.RawSetString("mode", lua.LNumber(info.Mode()))
		infoTb.RawSetString("modTime", lua.LNumber(info.ModTime().Unix()))
		infoTb.RawSetString("isDir", lua.LBool(info.IsDir()))
		infoTb.RawSetString("sys", lua.LNil)

		L.Push(objFn)

		n := 2
		if callable.Ref() != lua.LNil {
			n += 1
			L.Push(callable.Ref())
		}

		L.Push(lua.LString(p))
		L.Push(infoTb)

		L.Call(n, 0)

		return nil
	})

	if err != nil {
		L.Push(lua.LString(err.Error()))

		return 1
	}

	return 0
}
