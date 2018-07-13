// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package di

import (
	"github.com/jefurry/gola/lua/cb"
	"github.com/yuin/gopher-lua"
)

func diCall(L *lua.LState) int {
	L.CheckTypes(1, lua.LTFunction, lua.LTTable)
	val := L.CheckAny(1) // callable

	top := L.GetTop()
	args := make([]lua.LValue, 0, top-1)
	for i := 2; i <= top; i++ {
		args = append(args, L.CheckAny(i))
	}

	v, err := cb.Call(L, val, args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(v)

	return 1
}

func diApply(L *lua.LState) int {
	L.CheckTypes(1, lua.LTFunction, lua.LTTable)
	val := L.CheckAny(1) // callable
	tb := L.OptTable(2, nil)

	var args []lua.LValue
	if tb == nil {
		args = make([]lua.LValue, 0, 0)
	} else {
		args = make([]lua.LValue, 0, tb.Len())
		tb.ForEach(func(k, v lua.LValue) {
			args = append(args, v)
		})
	}

	v, err := cb.Call(L, val, args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(v)

	return 1
}

func diBind(L *lua.LState) int {
	L.CheckTypes(1, lua.LTNil, lua.LTTable, lua.LTUserData)
	obj := L.CheckAny(1)
	val := L.CheckAny(2) // callable
	tb := L.OptTable(3, nil)

	callable, err := cb.New(L, val)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	fn, err := callable.ObjFn(L)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	var args []lua.LValue
	if tb == nil {
		args = make([]lua.LValue, 0, 0)
	} else {
		args = make([]lua.LValue, 0, tb.Len())
		tb.ForEach(func(k, v lua.LValue) {
			args = append(args, v)
		})
	}

	if obj == lua.LNil {
		obj = callable.Ref()
	}

	L.Push(fn)
	n := len(args)
	if obj != lua.LNil {
		L.Push(obj)
		n += 1
	}

	for _, v := range args {
		L.Push(v)
	}

	L.Call(n, 1)

	L.Push(L.Get(-1))

	return 1
}
