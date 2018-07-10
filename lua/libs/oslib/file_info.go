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
	osFileInfoTypeName = lua.OsLibName + ".FILE_INFO*"
)

func osRegisterFileInfoMetatype(L *lua.LState) {
	// meta table
	mt := L.NewTypeMetatable(osFileInfoTypeName)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), osFileInfoFuncs))
}

var osFileInfoFuncs = map[string]lua.LGFunction{
	"name":    osFileInfoName,
	"size":    osFileInfoSize,
	"mode":    osFileInfoMode,
	"modTime": osFileInfoModTime,
	"isDir":   osFileInfoIsDir,
	"sys":     osFileInfoSys,
}

// base name of the file
func osFileInfoName(L *lua.LState) int {
	fileInfo := checkFileInfo(L, 1)

	name := fileInfo.Name()

	L.Push(lua.LString(name))

	return 1
}

// length in bytes for regular files; system-dependent for others
func osFileInfoSize(L *lua.LState) int {
	fileInfo := checkFileInfo(L, 1)

	size := fileInfo.Size()

	L.Push(lua.LNumber(size))

	return 1
}

// file mode bits
func osFileInfoMode(L *lua.LState) int {
	fileInfo := checkFileInfo(L, 1)

	mode := fileInfo.Mode()

	L.Push(newFileMode(L, mode))

	return 1
}

// modification time
func osFileInfoModTime(L *lua.LState) int {
	fileInfo := checkFileInfo(L, 1)

	t := fileInfo.ModTime()

	L.Push(lua.LNumber(t.Unix()))

	return 1
}

// abbreviation for Mode().IsDir()
func osFileInfoIsDir(L *lua.LState) int {
	fileInfo := checkFileInfo(L, 1)

	isDir := fileInfo.IsDir()

	L.Push(lua.LBool(isDir))

	return 1
}

// underlying data source (can return nil)
func osFileInfoSys(L *lua.LState) int {
	L.Push(lua.LNil)

	return 1
}
