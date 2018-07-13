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
	"github.com/jefurry/logrus/hooks/rotatelog"
	"github.com/yuin/gopher-lua"
)

const (
	logRotatelogTypeName = LogLibName + ".ROTATELOG_HOOK*"
)

func logRotatelogHookNew(L *lua.LState) int {
	p := L.CheckString(1)
	val := L.OptTable(2, nil)

	opts := checkRotatelogHookOptions(L, val)

	hk, err := rotatelog.NewHook(p, opts...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	ud := L.NewUserData()
	ud.Value = hk

	L.SetMetatable(ud, L.GetTypeMetatable(logRotatelogTypeName))
	L.Push(ud)

	return 1
}

func logRotatelogSetFormatter(L *lua.LState) int {
	Rotateloghook := checkRotatelog(L, 1)
	L.CheckTypes(2, lua.LTUserData, lua.LTTable, lua.LTFunction)
	uv := L.CheckAny(2)

	if uv.Type() == lua.LTUserData {
		ud, ok := uv.(*lua.LUserData)
		if ok {
			tft, ok := ud.Value.(*logrus.TextFormatter)
			if ok {
				Rotateloghook.SetFormatter(tft)

				return 0
			}

			jft, ok := ud.Value.(*logrus.JSONFormatter)
			if ok {
				Rotateloghook.SetFormatter(jft)

				return 0
			}
		}
	}

	Rotateloghook.SetFormatter(&formatter{l: L, uv: uv})

	return 0
}

func logRegisterRotatelogMetatype(L *lua.LState) {
	// meta table
	mt := L.NewTypeMetatable(logRotatelogTypeName)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), logRotatelogFuncs))
}

var logRotatelogFuncs = map[string]lua.LGFunction{
	"setFormatter": logRotatelogSetFormatter,
}
