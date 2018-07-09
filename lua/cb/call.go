// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cb

import (
	"github.com/yuin/gopher-lua"
)

func Call(L *lua.LState, val lua.LValue, args ...lua.LValue) (lua.LValue, error) {
	callable, err := New(L, val)
	if err != nil {
		return lua.LNil, err
	}

	fn, err := callable.ObjFn(L)
	if err != nil {
		return lua.LNil, err
	}

	n := len(args)
	ref := callable.Ref()

	L.Push(fn)
	if ref != lua.LNil {
		L.Push(ref)
		n += 1
	}

	for _, v := range args {
		L.Push(v)
	}

	L.Call(n, 1)

	return L.Get(-1), nil
}

func Apply(L *lua.LState, val lua.LValue, args []lua.LValue) (lua.LValue, error) {
	return Call(L, val, args...)
}
