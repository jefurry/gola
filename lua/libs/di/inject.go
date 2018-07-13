// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package di

import (
	"github.com/yuin/gopher-lua"
)

func diInjectorNew(L *lua.LState) int {
	providers := L.OptTable(1, nil)
	size := L.OptInt(2, defaultInjectorSize)

	dii := newInjector(size)

	if providers != nil {
		err := dii.add(L, providers)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))

			return 2
		}
	}

	ud := L.NewUserData()
	ud.Value = dii

	L.SetMetatable(ud, L.GetTypeMetatable(diInjectorTypeName))
	L.Push(ud)

	return 1
}

func diInjectorAdd(L *lua.LState) int {
	dii := checkInjector(L, 1)
	tb := L.CheckTable(2)

	err := dii.add(L, tb)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

// Return a named service.
func diInjectorGet(L *lua.LState) int {
	dii := checkInjector(L, 1)
	name := L.CheckAny(2)
	deps := L.OptTable(3, nil)

	var locals map[lua.LValue]lua.LValue
	if deps != nil {
		locals = make(map[lua.LValue]lua.LValue, 3)

		deps.ForEach(func(k, v lua.LValue) {
			locals[k] = v
		})
	}

	ins, err := dii.get(L, name, locals)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(ins)

	return 1
}

func diInjectorInvoke(L *lua.LState) int {
	dii := checkInjector(L, 1)
	L.CheckTypes(2, lua.LTTable, lua.LTFunction)
	val := L.CheckAny(2)
	deps := L.OptTable(3, nil)

	var locals map[lua.LValue]lua.LValue
	if deps != nil {
		locals = make(map[lua.LValue]lua.LValue, 3)

		deps.ForEach(func(k, v lua.LValue) {
			locals[k] = v
		})
	}

	ins, err := dii.invoke(L, val, locals)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 1
	}

	L.Push(ins)

	return 1
}

func diInjectorInstantiate(L *lua.LState) int {
	dii := checkInjector(L, 1)
	val := L.CheckTable(2)
	deps := L.OptTable(3, nil)

	var locals map[lua.LValue]lua.LValue
	if deps != nil {
		locals = make(map[lua.LValue]lua.LValue, 3)

		deps.ForEach(func(k, v lua.LValue) {
			locals[k] = v
		})
	}

	ins, err := dii.instantiate(L, val, locals)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 1
	}

	L.Push(ins)

	return 1
}

func diRegisterInjectorMetatype(L *lua.LState) {
	// meta table
	mt := L.NewTypeMetatable(diInjectorTypeName)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), diInjectorFuncs))
}

var diInjectorFuncs = map[string]lua.LGFunction{
	"add":         diInjectorAdd,
	"get":         diInjectorGet,
	"invoke":      diInjectorInvoke,
	"instantiate": diInjectorInstantiate,
}
