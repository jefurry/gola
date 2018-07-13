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
	"github.com/pkg/errors"
	"github.com/yuin/gopher-lua"
)

func diParse(L *lua.LState) int {
	L.CheckTypes(1, lua.LTFunction, lua.LTTable)
	val := L.CheckAny(1)

	tb, err := parse(L, val)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(tb)

	return 1
}

func parse(L *lua.LState, val lua.LValue) (*lua.LTable, error) {
	callable, err := cb.New(L, val)
	if err != nil {
		return nil, err
	}

	fn, err := callable.ObjFn(L)
	if err != nil {
		return nil, err
	}

	var dbg *lua.Debug = &lua.Debug{}
	fun, err := L.GetInfo(">f", dbg, fn)
	if err != nil {
		return nil, err
	}

	fn, ok := fun.(*lua.LFunction)
	if !ok {
		return nil, errors.New("function reflection failed")
	}

	var nArgs int = int(fn.Proto.NumParameters)
	var tb *lua.LTable
	if nArgs > 0 {
		var i int = 1
		tb = L.CreateTable(nArgs, 0)
		for _, linfo := range fn.Proto.DbgLocals {
			if i > nArgs {
				break
			}

			tb.Append(lua.LString(linfo.Name))
			i += 1
		}
	}

	if tb == nil {
		return nil, errors.New("maybe includes invalid parameters for '...'")
	}

	return tb, nil
}
