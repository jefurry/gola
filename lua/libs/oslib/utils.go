// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package oslib implements os for Lua.
package oslib

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	"os"
	"syscall"
)

func newFile(L *lua.LState, file *os.File) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = file

	L.SetMetatable(ud, L.GetTypeMetatable(osFileTypeName))

	return ud
}

func checkFile(L *lua.LState, n int) *os.File {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*os.File); ok {
		return v
	}

	L.ArgError(n, fmt.Sprintf("%s expected, got %s", osFileTypeName, ud.Type()))

	return nil
}

func newFileInfo(L *lua.LState, fileInfo os.FileInfo) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = fileInfo

	L.SetMetatable(ud, L.GetTypeMetatable(osFileInfoTypeName))

	return ud
}

func checkFileInfo(L *lua.LState, n int) os.FileInfo {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(os.FileInfo); ok {
		return v
	}

	L.ArgError(n, fmt.Sprintf("%s expected, got %s", osFileInfoTypeName, ud.Type()))

	return nil
}

func newFileMode(L *lua.LState, fileMode os.FileMode) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = fileMode

	L.SetMetatable(ud, L.GetTypeMetatable(osFileModeTypeName))

	return ud
}

func checkFileMode(L *lua.LState, n int) os.FileMode {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(os.FileMode); ok {
		return v
	}

	L.ArgError(n, fmt.Sprintf("%s expected, got %s", osFileModeTypeName, ud.Type()))

	return 0
}

func newProcess(L *lua.LState, process *os.Process) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = process

	L.SetMetatable(ud, L.GetTypeMetatable(osProcessTypeName))

	return ud
}

func checkProcess(L *lua.LState, n int) *os.Process {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*os.Process); ok {
		return v
	}

	L.ArgError(n, fmt.Sprintf("%s expected, got %s", osProcessTypeName, ud.Type()))

	return nil
}

// TODO: process_attr.go
func tableToProcessAttr(L *lua.LState, tb *lua.LTable) *os.ProcAttr {
	pa := &os.ProcAttr{
		Dir:   "",
		Env:   []string{},
		Files: nil,
		Sys: &syscall.SysProcAttr{
			Chroot: "",
			Credential: &syscall.Credential{
				Uid:         0,
				Gid:         0,
				Groups:      []uint32{},
				NoSetGroups: false,
			},
			Ptrace:     false,
			Setsid:     false,
			Setpgid:    false,
			Setctty:    false,
			Noctty:     false,
			Ctty:       0,
			Foreground: false,
			Pgid:       0,
		},
	}

	if tb == nil {
		return pa
	}

	return pa
}
