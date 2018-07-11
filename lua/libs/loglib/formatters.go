// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package loglib

import (
	"github.com/jefurry/gola/lua/cb"
	"github.com/jefurry/logrus"
	"github.com/pkg/errors"
	"github.com/yuin/gopher-lua"
)

type (
	formatter struct {
		l  *lua.LState
		uv lua.LValue
	}
)

// Format renders a single log entry
func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var err error
	var callable *cb.Callable

	L := f.l
	if f.uv.Type() == lua.LTFunction {
		callable, err = cb.New(L, f.uv)
	} else {
		callable, err = cb.With(L, f.uv, lua.LString("format"))
	}

	if err != nil {
		return nil, err
	}

	objFn, err := callable.ObjFn(L)
	if err != nil {
		return nil, err
	}

	L.Push(objFn)

	n := 1
	ref := callable.Ref()
	if ref != lua.LNil {
		L.Push(ref)
		n += 1
	}

	// TODO: implements Buffer
	// push entry
	ud := L.NewUserData()
	ud.Value = entry
	L.SetMetatable(ud, L.GetTypeMetatable(logEntryTypeName))
	L.Push(ud)

	L.Call(n, 2)

	str, errstr := L.Get(-2), L.Get(-1)
	if errstr != lua.LNil {
		es, ok := errstr.(lua.LString)
		if !ok {
			return nil, errors.New("Failed to obtain reader")
		}

		return nil, errors.New(string(es))
	}

	s, ok := str.(lua.LString)
	if !ok {
		return nil, errors.New("Failed to obtain reader")
	}

	return []byte(s), nil
}
