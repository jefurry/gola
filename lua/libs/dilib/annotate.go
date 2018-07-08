// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package dilib

import (
	"github.com/jefurry/gola/lua/cb"
	"github.com/pkg/errors"
	"github.com/yuin/gopher-lua"
)

type (
	diInjectWrapType struct {
		inject *lua.LTable
	}
)

func diAnnotate(L *lua.LState) int {
	val := L.CheckAny(1)

	var tb *lua.LTable
	if val != lua.LNil && val.Type() == lua.LTTable {
		tb, _ = val.(*lua.LTable)
	} else {
		top := L.GetTop()
		tb = L.CreateTable(top, 0)
		for i := 1; i <= top; i++ {
			tb.RawSetInt(i, L.CheckAny(i))
		}
	}

	callable, inject, err := annotate(L, tb)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	fn, err := callable.ObjFn(L)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	assoc(L, fn, inject)

	L.Push(lua.LTrue)

	return 1
}

func diAssoc(L *lua.LState) int {
	L.CheckTypes(1, lua.LTFunction, lua.LTTable)
	val := L.CheckAny(1)
	inject := L.CheckTable(2)

	err := assoc(L, val, inject)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func diClaim(L *lua.LState) int {
	L.CheckTypes(1, lua.LTFunction, lua.LTTable)
	val := L.CheckAny(1)

	inject, err := claim(L, val)
	if err != nil && err == cb.ErrInvalidCallable {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	if inject == nil {
		L.Push(lua.LNil)
	} else {
		L.Push(inject)
	}

	return 1
}

func diDissoc(L *lua.LState) int {
	L.CheckTypes(1, lua.LTFunction, lua.LTTable)
	val := L.CheckAny(1)

	err := dissoc(L, val)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LTrue)

	return 1
}

func annotate(L *lua.LState, inject *lua.LTable) (*cb.Callable, *lua.LTable, error) {
	if inject == nil {
		return nil, nil, errors.New("attempt to index a non-table")
	}

	size := inject.Len()
	if size <= 0 {
		return nil, nil, errors.New("attempt to index a non-table")
	}

	callable, err := cb.New(L, inject.RawGetInt(size))
	if err != nil {
		return nil, nil, err
	}

	inject.Remove(size)

	return callable, inject, nil
}

func assoc(L *lua.LState, val lua.LValue, inject *lua.LTable) error {
	callable, err := cb.New(L, val)
	if err != nil {
		return err
	}

	fn, err := callable.ObjFn(L)
	if err != nil {
		return err
	}

	ud := L.NewUserData()
	ud.Value = &diInjectWrapType{inject: inject}

	uvLen := len(fn.Upvalues)
	fn.Upvalues = append(fn.Upvalues, &lua.Upvalue{})
	fn.Upvalues[uvLen].Close()
	fn.Upvalues[uvLen].SetValue(ud)

	return nil
}

func claim(L *lua.LState, val lua.LValue) (*lua.LTable, error) {
	callable, err := cb.New(L, val)
	if err != nil {
		return nil, err
	}

	fn, err := callable.ObjFn(L)
	if err != nil {
		return nil, err
	}

	if fn.Upvalues == nil {
		fn.Upvalues = make([]*lua.Upvalue, 0, 0)
	}

	uvLen := len(fn.Upvalues)
	if uvLen == 0 {
		return nil, errors.New("attempt to index a non-upvalues object")
	}

	uv := fn.Upvalues[uvLen-1].Value()
	ud, ok := uv.(*lua.LUserData)
	if !ok {
		return nil, errors.Errorf("%s expected, got %s", lua.LTUserData, uv.Type())
	}

	iw, ok := ud.Value.(*diInjectWrapType)
	if !ok {
		return nil, errors.New("invalid injector data")
	}

	return iw.inject, nil
}

func dissoc(L *lua.LState, val lua.LValue) error {
	callable, err := cb.New(L, val)
	if err != nil {
		return err
	}

	fn, err := callable.ObjFn(L)
	if err != nil {
		return err
	}

	if fn.Upvalues == nil {
		return nil
	}

	uvLen := len(fn.Upvalues)
	if uvLen == 0 {
		return nil
	}

	uv := fn.Upvalues[uvLen-1].Value()
	ud, ok := uv.(*lua.LUserData)
	if !ok {
		return nil
	}

	_, ok = ud.Value.(*diInjectWrapType)
	if !ok {
		return nil
	}

	fn.Upvalues = fn.Upvalues[:uvLen-1]

	return nil
}
