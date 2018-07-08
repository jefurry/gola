// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package dilib

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
	"testing"
)

func testCallableTableSay(l *lua.LState) int {
	l.Push(lua.LString("Good Luck"))

	return 1
}

func TestCallable_1(t *testing.T) {
	L := lua.NewState()
	L.PreloadModule(DiLibName, Loader)
	defer L.Close()

	table := L.CreateTable(0, 1)
	table.RawSetString("say", L.NewFunction(testCallableTableSay))

	for _, v := range []struct {
		ref lua.LValue
		fn  lua.LValue
		ret lua.LString
	}{
		{lua.LNil, L.NewFunction(func(l *lua.LState) int {
			l.Push(lua.LString("Hello World"))

			return 1
		}), lua.LString("Hello World")},

		{lua.LNil, L.NewFunction(testCallableTableSay), lua.LString("Good Luck")},

		{table, lua.LString("say"), lua.LString("Good Luck")},
	} {
		val := L.CreateTable(2, 0)
		val.RawSetInt(1, v.ref)
		val.RawSetInt(2, v.fn)

		callable, err := newDiCallable(L, val)
		if !assert.NoError(t, err, "newDiCallable should succeed") {
			return
		}

		if !assert.Equal(t, v.ref, callable.getRef(), "v.ref must be equals to getRef()") {
			return
		}

		objFn, err := callable.getObjFn(L)
		if !assert.NoError(t, err, "getObjFn should succeed") {
			return
		}

		if !assert.NotEqual(t, objFn, nil, "getObjFn should succeed") {
			return
		}

		L.Push(objFn)
		L.Call(0, 1)
		ret := L.Get(-1)

		r, ok := ret.(lua.LString)
		if !assert.Equal(t, true, ok, fmt.Sprintf("objFn return type must be %s", lua.LTString)) {
			return
		}

		if !assert.Equal(t, string(v.ret), string(r), "objFn should succeed") {
			return
		}
	}
}

func TestCallable_2(t *testing.T) {
	L := lua.NewState()
	L.PreloadModule(DiLibName, Loader)
	defer L.Close()

	table := L.CreateTable(0, 1)
	table.RawSetString("say", L.NewFunction(testCallableTableSay))

	uv1 := L.CreateTable(2, 0)
	uv1.RawSetInt(1, table)
	uv1.RawSetInt(2, lua.LString("say"))

	uv2 := L.CreateTable(2, 0)
	uv2.RawSetInt(1, lua.LNil)
	uv2.RawSetInt(2, L.NewFunction(testCallableTableSay))

	uv3 := L.NewFunction(testCallableTableSay)

	for _, v := range []struct {
		uv  lua.LValue
		ref lua.LValue
		ret lua.LString
	}{
		{uv1, table, lua.LString("Good Luck")},
		{uv2, lua.LNil, lua.LString("Good Luck")},
		{uv3, lua.LNil, lua.LString("Good Luck")},
	} {
		callable, err := newDiCallable(L, v.uv)
		if !assert.NoError(t, err, "newDiCallable should succeed") {
			return
		}

		if !assert.Equal(t, v.ref, callable.getRef(), "v.ref must be equals to getRef()") {
			return
		}

		objFn, err := callable.getObjFn(L)
		if !assert.NoError(t, err, "getObjFn should succeed") {
			return
		}

		if !assert.NotEqual(t, objFn, nil, "getObjFn should succeed") {
			return
		}

		L.Push(objFn)
		L.Call(0, 1)
		ret := L.Get(-1)

		r, ok := ret.(lua.LString)
		if !assert.Equal(t, true, ok, fmt.Sprintf("objFn return type must be %s", lua.LTString)) {
			return
		}

		if !assert.Equal(t, string(v.ret), string(r), "objFn should succeed") {
			return
		}
	}
}
