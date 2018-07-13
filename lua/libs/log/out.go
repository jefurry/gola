// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package log

import (
	"github.com/jefurry/gola/lua/cb"
	"github.com/pkg/errors"
	"github.com/yuin/gopher-lua"
)

type (
	ioOut struct {
		l  *lua.LState
		uv lua.LValue
	}
)

// Write write logs.
func (o ioOut) Write(p []byte) (int, error) {
	var err error
	var callable *cb.Callable

	L := o.l
	if o.uv.Type() == lua.LTFunction {
		callable, err = cb.New(L, o.uv)
	} else {
		callable, err = cb.With(L, o.uv, lua.LString("write"))
	}

	if err != nil {
		return 0, err
	}

	objFn, err := callable.ObjFn(L)
	if err != nil {
		return 0, err
	}

	L.Push(objFn)

	n := 1
	ref := callable.Ref()
	if ref != lua.LNil {
		L.Push(ref)
		n += 1
	}

	L.Push(lua.LString(string(p)))

	L.Call(n, 2)

	ret, errstr := L.Get(-2), L.Get(-1)
	if ret != lua.LTrue {
		es, ok := errstr.(lua.LString)
		if !ok {
			return 0, errors.New("Failed to obtain writer")
		}

		return 0, errors.New(string(es))
	}

	return len(p), nil
}
