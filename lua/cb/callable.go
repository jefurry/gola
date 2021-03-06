// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package cb implements callable for Lua.
package cb

import (
	"fmt"
	"github.com/yuin/gopher-lua"
)

type (
	Callable                 []lua.LValue
	errInvalidCallableString struct {
		s string
	}
)

var (
	ErrInvalidCallable = newErrInvalidCallable()
)

func newErrInvalidCallable() *errInvalidCallableString {
	return &errInvalidCallableString{
		s: "invalid callable: ",
	}
}

func (e *errInvalidCallableString) errorf(format string, a ...interface{}) *errInvalidCallableString {
	e.s = e.s + fmt.Sprintf(format, a...)

	return e
}

func (e *errInvalidCallableString) Error() string {
	return e.s
}

func With(L *lua.LState, obj, fn lua.LValue) (*Callable, error) {
	val := L.CreateTable(2, 0)
	val.RawSetInt(1, obj)
	val.RawSetInt(2, fn)

	return New(L, val)
}

func New(L *lua.LState, val lua.LValue) (*Callable, error) {
	dic := make(Callable, 2, 2)
	switch val.Type() {
	case lua.LTFunction:
		f := val.(*lua.LFunction)
		if f == nil {
			return nil, ErrInvalidCallable.errorf("attempt to call a non-function")
		}

		dic[0] = lua.LNil
		dic[1] = val

		return &dic, nil
	case lua.LTTable:
		tb, _ := val.(*lua.LTable)
		if tb == nil {
			return nil, ErrInvalidCallable.errorf("attempt to index a non-table")
		}

		size := tb.Len()
		if size == 0 {
			return nil, ErrInvalidCallable.errorf("attempt to index a non-table")
		} else if size > 1 {
			ref := tb.RawGetInt(1)
			fun := tb.RawGetInt(2)

			typ := ref.Type()
			if typ != lua.LTNil && typ != lua.LTTable && typ != lua.LTUserData {
				return nil, ErrInvalidCallable.errorf("%s or %s or %s expected, got %s", lua.LTTable, lua.LTUserData, lua.LTNil, typ)
			}

			if typ == lua.LTNil {
				f, ok := fun.(*lua.LFunction)
				if !ok {
					return nil, ErrInvalidCallable.errorf("%s expected, got %s", lua.LTFunction, fun.Type())
				}

				if f == nil {
					return nil, ErrInvalidCallable.errorf("attempt to call a non-function")
				}

				dic[0] = ref
				dic[1] = fun

				return &dic, nil
			} else {
				if typ == lua.LTTable {
					obj, _ := ref.(*lua.LTable)
					if obj == nil {
						return nil, ErrInvalidCallable.errorf("attempt to index a non-table")
					}
				} else if typ == lua.LTUserData {
					obj, _ := ref.(*lua.LUserData)
					if obj == nil {
						return nil, ErrInvalidCallable.errorf("attempt to index a non-userdata")
					}
				}

				f, ok := fun.(lua.LString)
				if !ok {
					return nil, ErrInvalidCallable.errorf("%s expected, got %s", lua.LTString, fun.Type())
				}

				if f == lua.LString("") {
					return nil, ErrInvalidCallable.errorf("attempt to call a non-function")
				}

				dic[0] = ref
				dic[1] = fun

				return &dic, nil
			}
		} else if size > 0 {
			ref := lua.LNil
			fun := tb.RawGetInt(1)

			f, ok := fun.(*lua.LFunction)
			if !ok {
				return nil, ErrInvalidCallable.errorf("%s expected, got %s", lua.LTFunction, fun.Type())
			}

			if f == nil {
				return nil, ErrInvalidCallable.errorf("attempt to call a non-function")
			}

			dic[0] = ref
			dic[1] = fun

			return &dic, nil
		}
	}

	return nil, ErrInvalidCallable.errorf("%s or %s expected, got %s", lua.LTTable, lua.LTFunction, val.Type())
}

func (dic *Callable) size() int {
	return len(*dic)
}

func (dic *Callable) Ref() lua.LValue {
	ic := *dic

	return ic[0]
}

func (dic *Callable) Fn() lua.LValue {
	ic := *dic

	return ic[1]
}

func (dic *Callable) ObjFn(L *lua.LState) (*lua.LFunction, error) {
	ref := dic.Ref()
	fun := dic.Fn()

	if ref == lua.LNil {
		fn, ok := fun.(*lua.LFunction)
		if !ok {
			return nil, ErrInvalidCallable.errorf("%s expected, got %s", lua.LTFunction, fun.Type())
		}

		if fn == nil {
			return nil, ErrInvalidCallable.errorf("attempt to call a non-function")
		}

		return fn, nil
	}

	switch ref.Type() {
	case lua.LTTable:
		obj := ref.(*lua.LTable)
		f := L.GetField(obj, string(fun.(lua.LString)))

		fn, ok := f.(*lua.LFunction)
		if !ok {
			return nil, ErrInvalidCallable.errorf("%s expected, got %s", lua.LTFunction, f.Type())
		}

		if fn == nil {
			return nil, ErrInvalidCallable.errorf("attempt to call a non-function")
		}

		return fn, nil
	case lua.LTUserData:
		obj := ref.(*lua.LUserData)
		f := L.GetField(obj, string(fun.(lua.LString)))

		fn, ok := f.(*lua.LFunction)
		if !ok {
			return nil, ErrInvalidCallable.errorf("%s expected, got %s", lua.LTFunction, f.Type())
		}

		if fn == nil {
			return nil, ErrInvalidCallable.errorf("attempt to call a non-function")
		}

		return fn, nil
	}

	return nil, ErrInvalidCallable.errorf("attempt to index a non-table")
}
