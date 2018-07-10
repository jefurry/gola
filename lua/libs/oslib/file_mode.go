// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package oslib

import (
	"github.com/yuin/gopher-lua"
)

const (
	osFileModeTypeName = lua.OsLibName + ".FILE_MODE*"
)

func osRegisterFileModeMetatype(L *lua.LState) {
	// meta table
	mt := L.NewTypeMetatable(osFileModeTypeName)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), osFileModeFuncs))
	L.SetField(mt, "__tostring", L.NewFunction(osFileModeString))
}

var osFileModeFuncs = map[string]lua.LGFunction{
	"isDir":     osFileModeIsDir,
	"isRegular": osFileModeIsRegular,
	"perm":      osFileModePerm,
}

func osFileModeString(L *lua.LState) int {
	fileMode := checkFileMode(L, 1)

	L.Push(lua.LString(fileMode.String()))

	return 1
}

func osFileModeIsDir(L *lua.LState) int {
	fileMode := checkFileMode(L, 1)

	isDir := fileMode.IsDir()

	L.Push(lua.LBool(isDir))

	return 1
}

func osFileModeIsRegular(L *lua.LState) int {
	fileMode := checkFileMode(L, 1)

	isRegular := fileMode.IsRegular()

	L.Push(lua.LBool(isRegular))

	return 1
}

func osFileModePerm(L *lua.LState) int {
	fileMode := checkFileMode(L, 1)

	perm := fileMode.Perm()

	L.Push(lua.LNumber(perm))

	return 1
}
