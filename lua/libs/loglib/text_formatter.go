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
	logTextFormatterTypeName = LogLibName + ".TEXT_FORMATTER*"
)

func logTextFormatterNew(L *lua.LState) int {
	timestampFormat := L.OptString(1, "")

	ft := new(logrus.TextFormatter)
	if timestampFormat != "" {
		ft.TimestampFormat = timestampFormat
	}

	ud := L.NewUserData()
	ud.Value = ft

	L.SetMetatable(ud, L.GetTypeMetatable(logTextFormatterTypeName))
	L.Push(ud)

	return 1
}

func logTextFormatterSetTimestampFormat(L *lua.LState) int {
	ft := checkTextFormatter(L, 1)
	timestampFormat := L.CheckString(2)
	ft.TimestampFormat = timestampFormat

	return 0
}

func logTextFormatterGetTimestampFormat(L *lua.LState) int {
	ft := checkTextFormatter(L, 1)
	L.Push(lua.LString(ft.TimestampFormat))

	return 1
}

func logRegisterTextFormatterMetatype(L *lua.LState) {
	// meta table
	mt := L.NewTypeMetatable(logTextFormatterTypeName)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), logTextFormatterFuncs))
}

var logTextFormatterFuncs = map[string]lua.LGFunction{
	"setTimestampFormat": logTextFormatterSetTimestampFormat,
	"getTimestampFormat": logTextFormatterGetTimestampFormat,
}
