// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package loglib

import (
	"github.com/jefurry/logrus"
	"github.com/yuin/gopher-lua"
)

const (
	logEntryTypeName = LogLibName + ".ENTRY*"
)

func logEntryNew(L *lua.LState) int {
	logger := checkLogger(L, 1)

	ud := L.NewUserData()
	ud.Value = logrus.NewEntry(logger)

	L.SetMetatable(ud, L.GetTypeMetatable(logEntryTypeName))
	L.Push(ud)

	return 1
}

func logEntryWithError(L *lua.LState) int {
	entry := checkEntry(L, 1)
	value := L.CheckString(2)

	ud := L.NewUserData()
	ud.Value = entry.WithField(logrus.ErrorKey, value)
	L.SetMetatable(ud, L.GetTypeMetatable(logEntryTypeName))
	L.Push(ud)

	return 1
}

func logEntryWithField(L *lua.LState) int {
	entry := checkEntry(L, 1)
	key := L.CheckString(2)
	value := L.CheckAny(3)

	ud := L.NewUserData()
	ud.Value = entry.WithField(key, value)
	L.SetMetatable(ud, L.GetTypeMetatable(logEntryTypeName))
	L.Push(ud)

	return 1
}

func logEntryWithFields(L *lua.LState) int {
	entry := checkEntry(L, 1)
	tb := L.CheckTable(2)

	fields := make(logrus.Fields)
	tb.ForEach(func(key, value lua.LValue) {
		fields[key.String()] = value
	})

	ud := L.NewUserData()
	ud.Value = entry.WithFields(fields)
	L.SetMetatable(ud, L.GetTypeMetatable(logEntryTypeName))
	L.Push(ud)

	return 1
}

func logEntryDebug(L *lua.LState) int {
	entry, vals := entryWrap(L)

	entry.Debug(vals...)

	return 0
}

func logEntryInfo(L *lua.LState) int {
	entry, vals := entryWrap(L)

	entry.Info(vals...)

	return 0
}

func logEntryPrint(L *lua.LState) int {
	entry, vals := entryWrap(L)

	entry.Print(vals...)

	return 0
}

func logEntryWarn(L *lua.LState) int {
	entry, vals := entryWrap(L)

	entry.Warn(vals...)

	return 0
}

func logEntryWarning(L *lua.LState) int {
	entry, vals := entryWrap(L)

	entry.Warning(vals...)

	return 0
}

func logEntryError(L *lua.LState) int {
	entry, vals := entryWrap(L)

	entry.Error(vals...)

	return 0
}

func logEntryFatal(L *lua.LState) int {
	entry, vals := entryWrap(L)

	entry.Fatal(vals...)

	return 0
}

func logEntryPanic(L *lua.LState) int {
	entry, vals := entryWrap(L)

	entry.Panic(vals...)

	return 0
}

func logEntryDebugf(L *lua.LState) int {
	entry, format, vals := entryWrapformat(L)

	entry.Debugf(format, vals...)

	return 0
}

func logEntryInfof(L *lua.LState) int {
	entry, format, vals := entryWrapformat(L)

	entry.Infof(format, vals...)

	return 0
}

func logEntryPrintf(L *lua.LState) int {
	entry, format, vals := entryWrapformat(L)

	entry.Printf(format, vals...)

	return 0
}

func logEntryWarnf(L *lua.LState) int {
	entry, format, vals := entryWrapformat(L)

	entry.Warnf(format, vals...)

	return 0
}

func logEntryWarningf(L *lua.LState) int {
	entry, format, vals := entryWrapformat(L)

	entry.Warningf(format, vals...)

	return 0
}

func logEntryErrorf(L *lua.LState) int {
	entry, format, vals := entryWrapformat(L)

	entry.Errorf(format, vals...)

	return 0
}

func logEntryFatalf(L *lua.LState) int {
	entry, format, vals := entryWrapformat(L)

	entry.Fatalf(format, vals...)

	return 0
}

func logEntryPanicf(L *lua.LState) int {
	entry, format, vals := entryWrapformat(L)

	entry.Panicf(format, vals...)

	return 0
}

func logEntryDebugln(L *lua.LState) int {
	entry, vals := entryWrap(L)

	entry.Debugln(vals...)

	return 0
}

func logEntryInfoln(L *lua.LState) int {
	entry, vals := entryWrap(L)

	entry.Infoln(vals...)

	return 0
}

func logEntryPrintln(L *lua.LState) int {
	entry, vals := entryWrap(L)

	entry.Println(vals...)

	return 0
}

func logEntryWarnln(L *lua.LState) int {
	entry, vals := entryWrap(L)

	entry.Warnln(vals...)

	return 0
}

func logEntryWarningln(L *lua.LState) int {
	entry, vals := entryWrap(L)

	entry.Warningln(vals...)

	return 0
}

func logEntryErrorln(L *lua.LState) int {
	entry, vals := entryWrap(L)

	entry.Errorln(vals...)

	return 0
}

func logEntryFatalln(L *lua.LState) int {
	entry, vals := entryWrap(L)

	entry.Fatalln(vals...)

	return 0
}

func logEntryPanicln(L *lua.LState) int {
	entry, vals := entryWrap(L)

	entry.Panicln(vals...)

	return 0
}

func logEntryString(L *lua.LState) int {
	entry := checkEntry(L, 1)
	str, err := entry.String()
	if err != nil {
	}

	L.Push(lua.LString(str))

	return 1
}

func logRegisterEntryMetatype(L *lua.LState) {
	// meta table
	mt := L.NewTypeMetatable(logEntryTypeName)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), logEntryFuncs))
}

var logEntryFuncs = map[string]lua.LGFunction{
	"withError":  logEntryWithError,
	"withField":  logEntryWithField,
	"withFields": logEntryWithFields,

	"debug":  logEntryDebug,
	"info":   logEntryInfo,
	"print":  logEntryPrint,
	"warn":   logEntryWarn,
	"waring": logEntryWarning,
	"error":  logEntryError,
	"fatal":  logEntryFatal,
	"panic":  logEntryPanic,

	"debugf":  logEntryDebugf,
	"infof":   logEntryInfof,
	"printf":  logEntryPrintf,
	"warnf":   logEntryWarnf,
	"waringf": logEntryWarningf,
	"errorf":  logEntryErrorf,
	"fatalf":  logEntryFatalf,
	"panicf":  logEntryPanicf,

	"debugln":  logEntryDebugln,
	"infoln":   logEntryInfoln,
	"println":  logEntryPrintln,
	"warnln":   logEntryWarnln,
	"waringln": logEntryWarningln,
	"errorln":  logEntryErrorln,
	"fatalln":  logEntryFatalln,
	"panicln":  logEntryPanicln,
}
