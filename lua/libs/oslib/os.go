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
	"github.com/jefurry/gola/lua/cb"
	"github.com/yuin/gopher-lua"
	"os"
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
	"O_RDONLY": lua.LNumber(os.O_RDONLY),
	"O_WRONLY": lua.LNumber(os.O_WRONLY),
	"O_RDWR":   lua.LNumber(os.O_RDWR),
	// The remaining values may be or'ed in to control behavior.
	"O_APPEND": lua.LNumber(os.O_APPEND),
	"O_CREATE": lua.LNumber(os.O_CREATE),
	"O_EXCL":   lua.LNumber(os.O_EXCL),
	"O_SYNC":   lua.LNumber(os.O_SYNC),
	"O_TRUNC":  lua.LNumber(os.O_TRUNC),

	"PathSeparator":     lua.LString(os.PathSeparator),
	"PathListSeparator": lua.LString(os.PathListSeparator),
	"DevNull":           lua.LString(os.DevNull),
}

var osFileModeFields = map[string]os.FileMode{
	// The single letters are the abbreviations
	// used by the String method's formatting.
	"ModeDir":        os.ModeDir,        // d: is a directory
	"ModeAppend":     os.ModeAppend,     // a: append-only
	"ModeExclusive":  os.ModeExclusive,  // l: exclusive use
	"ModeTemporary":  os.ModeTemporary,  // T: temporary file; Plan 9 only
	"ModeSymlink":    os.ModeSymlink,    // L: symbolic link
	"ModeDevice":     os.ModeDevice,     // D: device file
	"ModeNamedPipe":  os.ModeNamedPipe,  // p: named pipe (FIFO)
	"ModeSocket":     os.ModeSocket,     // S: Unix domain socket
	"ModeSetuid":     os.ModeSetuid,     // u: setuid
	"ModeSetgid":     os.ModeSetgid,     // g: setgid
	"ModeCharDevice": os.ModeCharDevice, // c: Unix character device, when ModeDevice is set
	"ModeSticky":     os.ModeSticky,     // t: sticky
	// Mask for the type bits. For regular files, none will be set.
	"ModeType": os.ModeType,
	"ModePerm": os.ModePerm, // Unix permission bits
}

func osChdir(L *lua.LState) int {
	dir := L.CheckString(1)

	if err := os.Chdir(dir); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LFalse)

	return 1
}

func osOpen(L *lua.LState) int {
	name := L.CheckString(1)

	file, err := os.Open(name)
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

	file, err := os.Create(name)
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

	file, err := os.OpenFile(name, flag, os.FileMode(perm))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(newFile(L, file))

	return 1
}

func osTempDir(L *lua.LState) int {
	L.Push(lua.LString(os.TempDir()))

	return 1
}

func osChmod(L *lua.LState) int {
	name := L.CheckString(1)
	mode := L.CheckInt(2)

	if err := os.Chmod(name, os.FileMode(mode)); err != nil {
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

	if err := os.Chown(name, uid, gid); err != nil {
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

	ss := os.Expand(str, func(s string) string {
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

	L.Push(lua.LString(os.ExpandEnv(s)))

	return 1
}

func osLookupEnv(L *lua.LState) int {
	key := L.CheckString(1)

	v, ok := os.LookupEnv(key)

	L.Push(lua.LString(v))
	L.Push(lua.LBool(ok))

	return 2
}

func osSetenv(L *lua.LState) int {
	key := L.CheckString(1)
	value := L.CheckString(2)

	if err := os.Setenv(key, value); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func osUnsetenv(L *lua.LState) int {
	key := L.CheckString(1)

	if err := os.Unsetenv(key); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func osClearenv(L *lua.LState) int {
	os.Clearenv()

	return 0
}

func osChtimes(L *lua.LState) int {
	name := L.CheckString(1)
	atime := L.CheckInt(2)
	mtime := L.CheckInt(3)

	if err := os.Chtimes(name, time.Unix(int64(atime), 0), time.Unix(int64(mtime), 0)); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func osEnviron(L *lua.LState) int {
	env := os.Environ()

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

	if err := os.Truncate(name, int64(size)); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func osGetwd(L *lua.LState) int {
	dir, err := os.Getwd()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LString(dir))

	return 1
}

func osExecutable(L *lua.LState) int {
	exe, err := os.Executable()
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

	fileInfo, err := os.Lstat(name)
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

	fileInfo, err := os.Stat(name)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(newFileInfo(L, fileInfo))

	return 1
}

func osGetpagesize(L *lua.LState) int {
	n := os.Getpagesize()

	L.Push(lua.LNumber(n))

	return 1
}

func osSameFile(L *lua.LState) int {
	fileInfo1 := checkFileInfo(L, 1)
	fileInfo2 := checkFileInfo(L, 2)

	L.Push(lua.LBool(os.SameFile(fileInfo1, fileInfo2)))

	return 1
}

func osLchown(L *lua.LState) int {
	name := L.CheckString(1)
	uid := L.CheckInt(2)
	gid := L.CheckInt(3)

	if err := os.Lchown(name, uid, gid); err != nil {
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

	if err := os.Link(oldname, newname); err != nil {
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

	if err := os.Symlink(oldname, newname); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func osReadlink(L *lua.LState) int {
	name := L.CheckString(1)
	s, err := os.Readlink(name)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LString(s))

	return 1
}
