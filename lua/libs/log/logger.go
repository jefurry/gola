// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package log

import (
	"github.com/jefurry/logrus"
	"github.com/yuin/gopher-lua"
	"io/ioutil"
)

const (
	logLoggerTypeName = LogLibName + ".LOGGER*"
)

func logLoggerNew(L *lua.LState) int {
	logger := logrus.New()
	ud := L.NewUserData()
	ud.Value = logger

	L.SetMetatable(ud, L.GetTypeMetatable(logLoggerTypeName))
	L.Push(ud)

	return 1
}

func logLoggerWithError(L *lua.LState) int {
	logger := checkLogger(L, 1)
	value := L.CheckString(2)

	ud := L.NewUserData()
	ud.Value = logger.WithField(logrus.ErrorKey, value)
	L.Push(ud)

	return 1
}

func logLoggerWithField(L *lua.LState) int {
	logger := checkLogger(L, 1)
	key := L.CheckString(2)
	value := L.CheckAny(3)

	ud := L.NewUserData()
	ud.Value = logger.WithField(key, value)
	L.Push(ud)

	return 1
}

func logLoggerWithFields(L *lua.LState) int {
	logger := checkLogger(L, 1)
	tb := L.CheckTable(2)
	fields := make(logrus.Fields)
	tb.ForEach(func(key, value lua.LValue) {
		fields[key.String()] = value
	})

	ud := L.NewUserData()
	ud.Value = logger.WithFields(fields)
	L.Push(ud)

	return 1
}

func logLoggerDebug(L *lua.LState) int {
	logger, vals := loggerWrap(L)

	logger.Debug(vals...)

	return 0
}

func logLoggerInfo(L *lua.LState) int {
	logger, vals := loggerWrap(L)

	logger.Info(vals...)

	return 0
}

func logLoggerPrint(L *lua.LState) int {
	logger, vals := loggerWrap(L)
	logger.Print(vals...)

	return 0
}

func logLoggerWarn(L *lua.LState) int {
	logger, vals := loggerWrap(L)
	logger.Warn(vals...)

	return 0
}

func logLoggerWarning(L *lua.LState) int {
	logger, vals := loggerWrap(L)
	logger.Warning(vals...)

	return 0
}

func logLoggerError(L *lua.LState) int {
	logger, vals := loggerWrap(L)
	logger.Error(vals...)

	return 0
}

func logLoggerFatal(L *lua.LState) int {
	logger, vals := loggerWrap(L)
	logger.Fatal(vals...)

	return 0
}

func logLoggerPanic(L *lua.LState) int {
	logger, vals := loggerWrap(L)
	logger.Panic(vals...)

	return 0
}

func logLoggerDebugf(L *lua.LState) int {
	logger, format, vals := loggerWrapformat(L)

	logger.Debugf(format, vals...)

	return 0
}

func logLoggerInfof(L *lua.LState) int {
	logger, format, vals := loggerWrapformat(L)

	logger.Infof(format, vals...)

	return 0
}

func logLoggerPrintf(L *lua.LState) int {
	logger, format, vals := loggerWrapformat(L)

	logger.Printf(format, vals...)

	return 0
}

func logLoggerWarnf(L *lua.LState) int {
	logger, format, vals := loggerWrapformat(L)

	logger.Warnf(format, vals...)

	return 0
}

func logLoggerWarningf(L *lua.LState) int {
	logger, format, vals := loggerWrapformat(L)

	logger.Warningf(format, vals...)

	return 0
}

func logLoggerErrorf(L *lua.LState) int {
	logger, format, vals := loggerWrapformat(L)

	logger.Errorf(format, vals...)

	return 0
}

func logLoggerFatalf(L *lua.LState) int {
	logger, format, vals := loggerWrapformat(L)

	logger.Fatalf(format, vals...)

	return 0
}

func logLoggerPanicf(L *lua.LState) int {
	logger, format, vals := loggerWrapformat(L)

	logger.Panicf(format, vals...)

	return 0
}

func logLoggerDebugln(L *lua.LState) int {
	logger, vals := loggerWrap(L)

	logger.Debugln(vals...)

	return 0
}

func logLoggerInfoln(L *lua.LState) int {
	logger, vals := loggerWrap(L)

	logger.Infoln(vals...)

	return 0
}

func logLoggerPrintln(L *lua.LState) int {
	logger, vals := loggerWrap(L)

	logger.Println(vals...)

	return 0
}

func logLoggerWarnln(L *lua.LState) int {
	logger, vals := loggerWrap(L)

	logger.Warnln(vals...)

	return 0
}

func logLoggerWarningln(L *lua.LState) int {
	logger, vals := loggerWrap(L)

	logger.Warningln(vals...)

	return 0
}

func logLoggerErrorln(L *lua.LState) int {
	logger, vals := loggerWrap(L)

	logger.Errorln(vals...)

	return 0
}

func logLoggerFatalln(L *lua.LState) int {
	logger, vals := loggerWrap(L)

	logger.Fatalln(vals...)

	return 0
}

func logLoggerPanicln(L *lua.LState) int {
	logger, vals := loggerWrap(L)

	logger.Panicln(vals...)

	return 0
}

func logLoggerGetLevel(L *lua.LState) int {
	logger := checkLogger(L, 1)

	L.Push(lua.LNumber(logger.Level))
	L.Push(lua.LString(logger.Level.String()))

	return 2
}

func logLoggerGetLevels(L *lua.LState) int {
	tb := L.CreateTable(len(logrus.AllLevels), 0)
	for i, v := range logrus.AllLevels {
		tb.RawSetInt(i, lua.LNumber(v))
	}

	L.Push(tb)

	return 1
}

func logLoggerSetLevel(L *lua.LState) int {
	logger := checkLogger(L, 1)
	level := L.CheckNumber(2)

	logger.SetLevel(logrus.Level(level))

	return 0
}

func logLoggerBlockingLevel(L *lua.LState) int {
	logger := checkLogger(L, 1)
	lv := lua.LNumber(logger.GetBlockingLevel())
	L.Push(lv)

	return 1
}

func logLoggerSetNoLock(L *lua.LState) int {
	logger := checkLogger(L, 1)
	logger.SetNoLock()

	return 0
}

func logLoggerSetOut(L *lua.LState) int {
	logger := checkLogger(L, 1)
	uv := L.CheckAny(2)
	if uv == logDiscard {
		logger.Out = ioutil.Discard

		return 0
	}

	L.CheckTypes(2, lua.LTUserData, lua.LTTable, lua.LTFunction)
	logger.Out = ioOut{l: L, uv: uv}

	return 0
}

func logLoggerSetFormatter(L *lua.LState) int {
	logger := checkLogger(L, 1)
	L.CheckTypes(2, lua.LTUserData, lua.LTTable, lua.LTFunction)
	uv := L.CheckAny(2)

	if uv.Type() == lua.LTUserData {
		ud, ok := uv.(*lua.LUserData)
		if ok {
			tft, ok := ud.Value.(*logrus.TextFormatter)
			if ok {
				logger.SetFormatter(tft)

				return 0
			}

			jft, ok := ud.Value.(*logrus.JSONFormatter)
			if ok {
				logger.SetFormatter(jft)

				return 0
			}
		}
	}

	logger.SetFormatter(&formatter{l: L, uv: uv})

	return 0
}

func logLoggerExit(L *lua.LState) int {
	logger := checkLogger(L, 1)
	code := L.OptInt(2, 0)
	logger.Exit(code)

	return 0
}

func logRegisterLoggerMetatype(L *lua.LState) {
	// meta table
	mt := L.NewTypeMetatable(logLoggerTypeName)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), logLoggerFuncs))
}

var logLoggerFuncs = map[string]lua.LGFunction{
	"withError":  logLoggerWithError,
	"withField":  logLoggerWithField,
	"withFields": logLoggerWithFields,
	//When file is opened with appending mode, it's safe to
	//write concurrently to a file (within 4k message on Linux).
	//In these cases user can choose to disable the lock.
	"getLevel":      logLoggerGetLevel,
	"getLevels":     logLoggerGetLevels,
	"blockingLevel": logLoggerBlockingLevel,
	"setNoLock":     logLoggerSetNoLock,
	"setLevel":      logLoggerSetLevel,
	"setOut":        logLoggerSetOut,
	"setFormatter":  logLoggerSetFormatter,
	"exit":          logLoggerExit,
	"addHook":       logLoggerAddHook,

	"debug":   logLoggerDebug,
	"info":    logLoggerInfo,
	"print":   logLoggerPrint,
	"warn":    logLoggerWarn,
	"warning": logLoggerWarning,
	"error":   logLoggerError,
	"fatal":   logLoggerFatal,
	"panic":   logLoggerPanic,

	"debugf":   logLoggerDebugf,
	"infof":    logLoggerInfof,
	"printf":   logLoggerPrintf,
	"warnf":    logLoggerWarnf,
	"warningf": logLoggerWarningf,
	"errorf":   logLoggerErrorf,
	"fatalf":   logLoggerFatalf,
	"panicf":   logLoggerPanicf,

	"debufln":   logLoggerDebugln,
	"infoln":    logLoggerInfoln,
	"println":   logLoggerPrintln,
	"warnln":    logLoggerWarnln,
	"warningln": logLoggerWarningln,
	"errorln":   logLoggerErrorln,
	"fatalln":   logLoggerFatalln,
	"panicln":   logLoggerPanicln,
}
