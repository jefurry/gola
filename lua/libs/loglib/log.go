// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package loglib implements logger for Lua.
package loglib

import (
	"github.com/jefurry/logrus"
	"github.com/yuin/gopher-lua"
)

const (
	LogLibName = "log"
)

func Open(L *lua.LState) {
	L.PreloadModule(LogLibName, Loader)
}

func Loader(L *lua.LState) int {
	logmod := L.SetFuncs(L.NewTable(), logFuncs)
	logRegister(L, logmod)
	L.Push(logmod)

	logRegisterLoggerMetatype(L)
	logRegisterEntryMetatype(L)

	// formatters
	logRegisterTextFormatterMetatype(L)
	logRegisterJSONFormatterMetatype(L)

	// hooks
	logRegisterRotatelogMetatype(L)

	return 1
}

func logParseLevel(L *lua.LState) int {
	lv := L.CheckString(1)

	l, err := logrus.ParseLevel(lv)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LNumber(l))

	return 1
}

func logGetAllLevels(L *lua.LState) int {
	tb := L.CreateTable(len(logrus.AllLevels), 0)
	for i, v := range logrus.AllLevels {
		tb.RawSetInt(i, lua.LNumber(v))
	}

	L.Push(tb)

	return 1
}

func logRegister(L *lua.LState, module *lua.LTable) {
	// static attributes
	for k, v := range logFields {
		L.SetField(module, k, v)
	}
}

var logFuncs = map[string]lua.LGFunction{
	"parseLevel":       logParseLevel,
	"getAllLevels":     logGetAllLevels,
	"newLogger":        logLoggerNew,
	"newEntry":         logEntryNew,
	"newTextFormatter": logTextFormatterNew,
	"newJSONFormatter": logJSONFormatterNew,
	"newRotatelogHook": logRotatelogHookNew,
}
