// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package oslib

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	"os"
	"time"
)

const (
	osFileTypeName = lua.OsLibName + ".FILE*"
)

func osRegisterFileMetatype(L *lua.LState) {
	// meta table
	mt := L.NewTypeMetatable(osFileTypeName)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), osFileFuncs))
}

var osFileFuncs = map[string]lua.LGFunction{
	"fileRead":         osFileRead,
	"readAt":           osFileReadAt,
	"write":            osFileWrite,
	"writeAt":          osFileWriteAt,
	"seek":             osFileSeek,
	"close":            osFileClose,
	"chmod":            osFileChmod,
	"name":             osFileName,
	"setDeadline":      osFileSetDeadline,
	"setReadDeadline":  osFileSetReadDeadline,
	"setWriteDeadline": osFileSetWriteDeadline,
	"readdir":          osFileReaddir,
	"Readdirnames":     osFileReaddirnames,
	"stat":             osFileStat,
	"sync":             osFileSync,
	"truncate":         osFileTruncate,
}

func osFileRead(L *lua.LState) int {
	file := checkFile(L, 1)
	n := L.CheckInt(2)

	if n < 0 {
		L.ArgError(1, fmt.Sprintf("n(%v) must be a positive number", n))
	}

	b := make([]byte, 0, n)
	nn, err := file.Read(b)
	if err != nil {
		L.Push(lua.LNumber(0))
		L.Push(lua.LString(err.Error()))

		return 2
	}
	L.Push(lua.LNumber(nn))

	return 1
}

func osFileReadAt(L *lua.LState) int {
	file := checkFile(L, 1)
	n := L.CheckInt(2)
	off := L.CheckInt(3)

	if n < 0 {
		L.ArgError(1, fmt.Sprintf("n(%v) must be a positive number", n))
	}

	if off < 0 {
		L.ArgError(1, fmt.Sprintf("off(%v) must be a positive number", off))
	}

	b := make([]byte, 0, n)
	nn, err := file.ReadAt(b, int64(off))
	if err != nil {
		L.Push(lua.LNumber(0))
		L.Push(lua.LString(err.Error()))

		return 2
	}
	L.Push(lua.LNumber(nn))

	return 1
}

func osFileWrite(L *lua.LState) int {
	file := checkFile(L, 1)
	s := L.CheckString(2)

	n, err := file.WriteString(s)
	if err != nil {
		L.Push(lua.LNumber(0))
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LNumber(n))

	return 1
}

func osFileWriteAt(L *lua.LState) int {
	file := checkFile(L, 1)
	s := L.CheckString(2)
	off := L.CheckInt(3)

	if off < 0 {
		L.ArgError(1, fmt.Sprintf("off(%v) must be a positive number", off))
	}

	nn, err := file.WriteAt([]byte(s), int64(off))
	if err != nil {
		L.Push(lua.LNumber(0))
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LNumber(nn))

	return 1
}

func osFileSeek(L *lua.LState) int {
	file := checkFile(L, 1)
	offset := L.CheckInt(2)
	whence := L.CheckInt(3)

	if offset < 0 {
		L.ArgError(1, fmt.Sprintf("offset(%v) must be a positive number", offset))
	}

	ret, err := file.Seek(int64(offset), whence)
	if err != nil {
		L.Push(lua.LNumber(0))
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LNumber(ret))

	return 1
}

func osFileChmod(L *lua.LState) int {
	file := checkFile(L, 1)
	mode := L.CheckInt(1)

	if err := file.Chmod(os.FileMode(mode)); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func osFileName(L *lua.LState) int {
	file := checkFile(L, 1)

	L.Push(lua.LString(file.Name()))

	return 1
}

func osFileClose(L *lua.LState) int {
	file := checkFile(L, 1)
	if err := file.Close(); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func osFileSetDeadline(L *lua.LState) int {
	file := checkFile(L, 1)
	timestamp := L.CheckInt(2)

	if timestamp < 0 {
		L.ArgError(1, fmt.Sprintf("timestamp(%v) must be a positive number", timestamp))
	}

	t := time.Unix(int64(timestamp), 0)
	if err := file.SetDeadline(t); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func osFileSetReadDeadline(L *lua.LState) int {
	file := checkFile(L, 1)
	timestamp := L.CheckInt(2)

	if timestamp < 0 {
		L.ArgError(1, fmt.Sprintf("timestamp(%v) must be a positive number", timestamp))
	}

	t := time.Unix(int64(timestamp), 0)
	if err := file.SetReadDeadline(t); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func osFileSetWriteDeadline(L *lua.LState) int {
	file := checkFile(L, 1)
	timestamp := L.CheckInt(2)

	if timestamp < 0 {
		L.ArgError(1, fmt.Sprintf("timestamp(%v) must be a positive number", timestamp))
	}

	t := time.Unix(int64(timestamp), 0)
	if err := file.SetWriteDeadline(t); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func osFileReaddir(L *lua.LState) int {
	file := checkFile(L, 1)
	n := L.CheckInt(2)

	fileInfos, err := file.Readdir(n)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	tb := L.CreateTable(len(fileInfos), 0)
	for i, fileInfo := range fileInfos {
		tb.RawSetInt(i, newFileInfo(L, fileInfo))
	}

	L.Push(tb)

	return 1
}

func osFileReaddirnames(L *lua.LState) int {
	file := checkFile(L, 1)
	n := L.CheckInt(2)

	names, err := file.Readdirnames(n)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	tb := L.CreateTable(len(names), 0)
	for i, name := range names {
		tb.RawSetInt(i, lua.LString(name))
	}

	L.Push(tb)

	return 1
}

func osFileStat(L *lua.LState) int {
	file := checkFile(L, 1)

	fileInfo, err := file.Stat()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(newFileInfo(L, fileInfo))

	return 1
}

func osFileSync(L *lua.LState) int {
	file := checkFile(L, 1)
	if err := file.Sync(); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func osFileTruncate(L *lua.LState) int {
	file := checkFile(L, 1)
	size := L.CheckInt(2)

	if err := file.Truncate(int64(size)); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}
