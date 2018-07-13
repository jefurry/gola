// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package sys implements syscall for Lua.
package sys

import (
	"github.com/yuin/gopher-lua"
	"syscall"
)

const (
	SysLibName = "sys"
)

func Open(L *lua.LState) {
	L.PreloadModule(SysLibName, Loader)
}

func Loader(L *lua.LState) int {
	sysmod := L.SetFuncs(L.NewTable(), sysFuncs)
	L.Push(sysmod)

	for k, v := range sysFields {
		sysmod.RawSetString(k, lua.LNumber(v))
	}

	return 1
}

var sysFuncs = map[string]lua.LGFunction{
	"kill":   sysKill,
	"getsid": sysGetsid,
}

var sysFields = map[string]syscall.Signal{
	"SIGABRT":   syscall.SIGABRT,
	"SIGALRM":   syscall.SIGALRM,
	"SIGBUS":    syscall.SIGBUS,
	"SIGCHLD":   syscall.SIGCHLD,
	"SIGCONT":   syscall.SIGCONT,
	"SIGEMT":    syscall.SIGEMT,
	"SIGFPE":    syscall.SIGFPE,
	"SIGHUP":    syscall.SIGHUP,
	"SIGILL":    syscall.SIGILL,
	"SIGINFO":   syscall.SIGINFO,
	"SIGINT":    syscall.SIGINT,
	"SIGIO":     syscall.SIGIO,
	"SIGIOT":    syscall.SIGIOT,
	"SIGKILL":   syscall.SIGKILL,
	"SIGPIPE":   syscall.SIGPIPE,
	"SIGPROF":   syscall.SIGPROF,
	"SIGQUIT":   syscall.SIGQUIT,
	"SIGSEGV":   syscall.SIGSEGV,
	"SIGSTOP":   syscall.SIGSTOP,
	"SIGSYS":    syscall.SIGSYS,
	"SIGTERM":   syscall.SIGTERM,
	"SIGTRAP":   syscall.SIGTRAP,
	"SIGTSTP":   syscall.SIGTSTP,
	"SIGTTIN":   syscall.SIGTTIN,
	"SIGTTOU":   syscall.SIGTTOU,
	"SIGURG":    syscall.SIGURG,
	"SIGUSR1":   syscall.SIGUSR1,
	"SIGUSR2":   syscall.SIGUSR2,
	"SIGVTALRM": syscall.SIGVTALRM,
	"SIGWINCH":  syscall.SIGWINCH,
	"SIGXCPU":   syscall.SIGXCPU,
	"SIGXFSZ":   syscall.SIGXFSZ,
}

func sysKill(L *lua.LState) int {
	pid := L.CheckInt(1)
	sig := L.CheckInt(2)

	if err := syscall.Kill(pid, syscall.Signal(sig)); err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func sysGetsid(L *lua.LState) int {
	pid := L.CheckInt(1)

	sid, err := syscall.Getsid(pid)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LNumber(sid))

	return 1
}
