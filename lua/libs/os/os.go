// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package os implements os for Lua.
package os

import (
	"github.com/jefurry/gola/lua/cb"
	"github.com/jefurry/gola/lua/libs/os/exec"
	"github.com/jefurry/gola/lua/libs/os/user"
	"github.com/yuin/gopher-lua"
	oos "os"
	"time"
)

func Open(L *lua.LState) {
	osmod, ok := L.GetGlobal(lua.OsLibName).(*lua.LTable)
	if !ok {
		L.RaiseError("module(%v) not exists", lua.OsLibName)
	}

	osRegisterFileMetatype(L)
	osRegisterFileInfoMetatype(L)
	osRegisterFileModeMetatype(L)
	osRegisterProcessMetatype(L)

	L.SetFuncs(osmod, osFuncs)

	for k, v := range osFields {
		osmod.RawSetString(k, v)
	}

	// FileMode fields
	for k, v := range osFileModeFields {
		osmod.RawSetString(k, newFileMode(L, v))
	}

	user.Open(L)
	exec.Open(L)
}

var osFuncs = map[string]lua.LGFunction{
	"chdir":       osChdir,
	"open":        osOpen,
	"create":      osCreate,
	"openFile":    osOpenFile,
	"tempDir":     osTempDir,
	"chmod":       osChmod,
	"chown":       osChown,
	"expand":      osExpand,
	"expandEnv":   osExpandEnv,
	"lookupEnv":   osLookupEnv,
	"setenv":      osSetenv,
	"unsetenv":    osUnsetenv,
	"clearenv":    osClearenv,
	"chtimes":     osChtimes,
	"environ":     osEnviron,
	"truncate":    osTruncate,
	"getwd":       osGetwd,
	"executable":  osExecutable,
	"lstat":       osLstat,
	"stat":        osStat,
	"getpagesize": osGetpagesize,
	"sameFile":    osSameFile,
	"lchown":      osLchown,
	"link":        osLink,
	"symlink":     osSymlink,
	"readlink":    osReadlink,

	// path.go
	"mkdir":      osMkdir,
	"remove":     osRemove,
	"mkdirAll":   osMkdirAll,
	"removeAll":  osRemoveAll,
	"isNotExist": osIsNotExist,
	"isExist":    osIsExist,

	// proc.go
	"getuid":    osGetuid,
	"geteuid":   osGeteuid,
	"getgid":    osGetgid,
	"getgroups": osGetgroups,

	// sys.go
	"hostname": osHostname,

	// process.go
	"getpid":      osGetpid,
	"getppid":     osGetppid,
	"findProcess": osFindProcess,
}

var osFields = map[string]lua.LValue{
	// Exactly one of O_RDONLY, O_WRONLY, or O_RDWR must be specified.
	"O_RDONLY": lua.LNumber(oos.O_RDONLY),
	"O_WRONLY": lua.LNumber(oos.O_WRONLY),
	"O_RDWR":   lua.LNumber(oos.O_RDWR),
	// The remaining values may be or'ed in to control behavior.
	"O_APPEND": lua.LNumber(oos.O_APPEND),
	"O_CREATE": lua.LNumber(oos.O_CREATE),
	"O_EXCL":   lua.LNumber(oos.O_EXCL),
	"O_SYNC":   lua.LNumber(oos.O_SYNC),
	"O_TRUNC":  lua.LNumber(oos.O_TRUNC),

	"PathSeparator":     lua.LString(oos.PathSeparator),
	"PathListSeparator": lua.LString(oos.PathListSeparator),
	"DevNull":           lua.LString(oos.DevNull),
}

var osFileModeFields = map[string]oos.FileMode{
	// The single letters are the abbreviations
	// used by the String method's formatting.
	"MODE_DIR":         oos.ModeDir,        // d: is a directory
	"MODE_APPEND":      oos.ModeAppend,     // a: append-only
	"MODE_EXCLUSIVE":   oos.ModeExclusive,  // l: exclusive use
	"MODE_TEMPORARY":   oos.ModeTemporary,  // T: temporary file; Plan 9 only
	"MODE_SYMLINK":     oos.ModeSymlink,    // L: symbolic link
	"MODE_DEVICE":      oos.ModeDevice,     // D: device file
	"MODE_NAME_PIPE":   oos.ModeNamedPipe,  // p: named pipe (FIFO)
	"MODE_SOCKET":      oos.ModeSocket,     // S: Unix domain socket
	"MODE_SET_UID":     oos.ModeSetuid,     // u: setuid
	"MODE_SET_GID":     oos.ModeSetgid,     // g: setgid
	"MODE_CHAR_DEVICE": oos.ModeCharDevice, // c: Unix character device, when ModeDevice is set
	"MODE_STICKY":      oos.ModeSticky,     // t: sticky
	// Mask for the type bits. For regular files, none will be set.
	"MODE_TYPE": oos.ModeType,
	"MODE_PERM": oos.ModePerm, // Unix permission bits
}

func osChdir(L *lua.LState) int {
	dir := L.CheckString(1)

	if err := oos.Chdir(dir); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LFalse)

	return 1
}

func osOpen(L *lua.LState) int {
	name := L.CheckString(1)

	file, err := oos.Open(name)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(newFile(L, file))

	return 1
}

func osCreate(L *lua.LState) int {
	name := L.CheckString(1)

	file, err := oos.Create(name)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(newFile(L, file))

	return 1
}

func osOpenFile(L *lua.LState) int {
	name := L.CheckString(1)
	flag := L.CheckInt(2)
	perm := L.CheckInt(3)

	file, err := oos.OpenFile(name, flag, oos.FileMode(perm))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(newFile(L, file))

	return 1
}

func osTempDir(L *lua.LState) int {
	L.Push(lua.LString(oos.TempDir()))

	return 1
}

func osChmod(L *lua.LState) int {
	name := L.CheckString(1)
	mode := L.CheckInt(2)

	if err := oos.Chmod(name, oos.FileMode(mode)); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func osChown(L *lua.LState) int {
	name := L.CheckString(1)
	uid := L.CheckInt(2)
	gid := L.CheckInt(3)

	if err := oos.Chown(name, uid, gid); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func osExpand(L *lua.LState) int {
	str := L.CheckString(1)
	mapping := L.CheckFunction(2)

	callable, err := cb.New(L, mapping)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	objFn, err := callable.ObjFn(L)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	ss := oos.Expand(str, func(s string) string {
		n := 1
		L.Push(objFn)
		if callable.Ref() != lua.LNil {
			n += 1
			L.Push(callable.Ref())
		}

		L.Push(lua.LString(s))
		L.Call(n, 1)

		ret := L.Get(-1)
		v, ok := ret.(lua.LString)
		if !ok {
			return ""
		}

		return string(v)
	})

	L.Push(lua.LString(ss))

	return 1
}

func osExpandEnv(L *lua.LState) int {
	s := L.CheckString(1)

	L.Push(lua.LString(oos.ExpandEnv(s)))

	return 1
}

func osLookupEnv(L *lua.LState) int {
	key := L.CheckString(1)

	v, ok := oos.LookupEnv(key)

	L.Push(lua.LString(v))
	L.Push(lua.LBool(ok))

	return 2
}

func osSetenv(L *lua.LState) int {
	key := L.CheckString(1)
	value := L.CheckString(2)

	if err := oos.Setenv(key, value); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func osUnsetenv(L *lua.LState) int {
	key := L.CheckString(1)

	if err := oos.Unsetenv(key); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func osClearenv(L *lua.LState) int {
	oos.Clearenv()

	return 0
}

func osChtimes(L *lua.LState) int {
	name := L.CheckString(1)
	atime := L.CheckInt(2)
	mtime := L.CheckInt(3)

	if err := oos.Chtimes(name, time.Unix(int64(atime), 0), time.Unix(int64(mtime), 0)); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func osEnviron(L *lua.LState) int {
	env := oos.Environ()

	tb := L.CreateTable(len(env), 0)
	if env != nil {
		for i, v := range env {
			tb.RawSetInt(i, lua.LString(v))
		}
	}

	L.Push(tb)

	return 1
}

func osTruncate(L *lua.LState) int {
	name := L.CheckString(1)
	size := L.CheckInt(2)

	if err := oos.Truncate(name, int64(size)); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func osGetwd(L *lua.LState) int {
	dir, err := oos.Getwd()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LString(dir))

	return 1
}

func osExecutable(L *lua.LState) int {
	exe, err := oos.Executable()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LString(exe))

	return 1
}

func osLstat(L *lua.LState) int {
	name := L.CheckString(1)

	fileInfo, err := oos.Lstat(name)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(newFileInfo(L, fileInfo))

	return 1
}

func osStat(L *lua.LState) int {
	name := L.CheckString(1)

	fileInfo, err := oos.Stat(name)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(newFileInfo(L, fileInfo))

	return 1
}

func osGetpagesize(L *lua.LState) int {
	n := oos.Getpagesize()

	L.Push(lua.LNumber(n))

	return 1
}

func osSameFile(L *lua.LState) int {
	fileInfo1 := checkFileInfo(L, 1)
	fileInfo2 := checkFileInfo(L, 2)

	L.Push(lua.LBool(oos.SameFile(fileInfo1, fileInfo2)))

	return 1
}

func osLchown(L *lua.LState) int {
	name := L.CheckString(1)
	uid := L.CheckInt(2)
	gid := L.CheckInt(3)

	if err := oos.Lchown(name, uid, gid); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func osLink(L *lua.LState) int {
	oldname := L.CheckString(1)
	newname := L.CheckString(2)

	if err := oos.Link(oldname, newname); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func osSymlink(L *lua.LState) int {
	oldname := L.CheckString(1)
	newname := L.CheckString(2)

	if err := oos.Symlink(oldname, newname); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func osReadlink(L *lua.LState) int {
	name := L.CheckString(1)
	s, err := oos.Readlink(name)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LString(s))

	return 1
}
