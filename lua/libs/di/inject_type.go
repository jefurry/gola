// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package di

import (
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/jefurry/gola/lua/cb"
	"github.com/pkg/errors"
	"github.com/yuin/gopher-lua"
)

const (
	diInjectorTypeName  = DiLibName + ".INJECTOR*"
	defaultInjectorSize = 10
)

const (
	diInjectorTypeType    lua.LString = "type"
	diInjectorTypeFactory lua.LString = "factory"
	diInjectorTypeValue   lua.LString = "value"
)

const (
	diInjectorSkipClassMethodParamName lua.LString = "self"
	diInjectorInjectClassMethodName    string      = "injector"
)

type (
	diInjector struct {
		currentlyResolving *arraylist.List
		providers          map[lua.LValue][]lua.LValue
		instances          map[lua.LValue]lua.LValue
	}
)

func newInjector(size int) *diInjector {
	return &diInjector{
		currentlyResolving: arraylist.New(),
		providers:          make(map[lua.LValue][]lua.LValue, size),
		instances:          make(map[lua.LValue]lua.LValue, size),
	}
}

func (dii *diInjector) add(L *lua.LState, providers *lua.LTable) error {
	if providers == nil {
		return errors.Errorf("%s expected, got %s", lua.LTTable, lua.LTNil)
	}

	key := lua.LNil
	for {
		lk, lv := providers.Next(key)
		key = lk
		if lk == lua.LNil {
			break
		}

		v, ok := lv.(*lua.LTable)
		if !ok {
			return errors.Errorf("%s expected, got %s", lua.LTTable, lv.Type())
		}

		if v == nil {
			return errors.Errorf("%s expected, got %s", lua.LTTable, lua.LTNil)
		}

		if v.Len() < 2 {
			return errors.New("the table requires at least 2 lengths")
		}

		typ := v.RawGetInt(1)
		val := v.RawGetInt(2)
		if typ != diInjectorTypeType && typ != diInjectorTypeFactory && typ != diInjectorTypeValue {
			return errors.Errorf("the type value must be '%s', '%s' or '%s'",
				diInjectorTypeType, diInjectorTypeFactory, diInjectorTypeValue)
		}

		dii.providers[lk] = []lua.LValue{typ, val}
	}

	return nil
}

func (dii *diInjector) get(L *lua.LState, name lua.LValue, locals map[lua.LValue]lua.LValue) (lua.LValue, error) {
	if ins, ok := dii.instances[name]; ok {
		return ins, nil
	}

	provider, ok := dii.providers[name]
	if !ok {
		return lua.LNil, errors.Errorf("no provider for '%s'", name)
	}

	if dii.currentlyResolving.IndexOf(name) > -1 {
		dii.currentlyResolving.Add(name)
		return lua.LNil, errors.Errorf("cannot resolve circular dependency for '%s'", name)
	}

	dii.currentlyResolving.Add(name)

	var err error
	var ins lua.LValue = lua.LNil
	typ, val := provider[0], provider[1]
	switch typ {
	case diInjectorTypeType:
		ins, err = dii.instantiate(L, val, locals)
		if err != nil {
			return lua.LNil, err
		}
	case diInjectorTypeFactory:
		ins, err = dii.invoke(L, val, locals)
		if err != nil {
			return lua.LNil, err
		}
	case diInjectorTypeValue:
		ins, _ = dii.value(L, val, nil)
	}

	dii.instances[name] = ins
	dii.currentlyResolving.Remove(dii.currentlyResolving.Size() - 1)

	return ins, nil
}

func (dii *diInjector) fnDef(L *lua.LState, self lua.LValue, fn *lua.LFunction,
	inject *lua.LTable, locals map[lua.LValue]lua.LValue) ([]lua.LValue, error) {
	val := L.CreateTable(2, 0)
	val.RawSetInt(1, lua.LNil)
	val.RawSetInt(2, fn)

	if inject != nil {
		err := assoc(L, val, inject)
		if err != nil {
			return nil, err
		}
	} else {
		inj, err := claim(L, val)
		if err != nil && err == cb.ErrInvalidCallable {
			return nil, err
		}

		inject = inj
		if inject == nil {
			inj, err := parse(L, val)
			if err != nil && err == cb.ErrInvalidCallable {
				return nil, err
			}

			if inj != nil {
				err := assoc(L, val, inj)
				if err != nil {
					return nil, err
				}
			}

			inject = inj
		}
	}

	if inject == nil {
		return nil, nil
	}

	key := lua.LNil
	deps := make([]lua.LValue, 0, inject.Len())
	for {
		lk, lv := inject.Next(key)
		key = lk
		if lk == lua.LNil {
			break
		}

		i, ok := lk.(lua.LNumber)
		if !ok {
			continue
		}

		if i == lua.LNumber(1) && lv == diInjectorSkipClassMethodParamName && self != lua.LNil {
			continue
		}

		if dep, ok := locals[lv]; ok {
			deps = append(deps, dep)
		} else {
			dep, err := dii.get(L, lv, nil)
			if err != nil {
				return nil, err
			}

			deps = append(deps, dep)
		}
	}

	return deps, nil
}

func (dii *diInjector) instantiate(L *lua.LState, val lua.LValue,
	locals map[lua.LValue]lua.LValue) (lua.LValue, error) {
	if val.Type() != lua.LTTable {
		return lua.LNil, errors.Errorf("%s expected, got %s", lua.LTTable, val.Type())
	}

	tb := val.(*lua.LTable)
	if tb == nil {
		return lua.LNil, errors.New("attempt to index a non-table")
	}

	var cls *lua.LTable
	var inject *lua.LTable = nil
	if isClass(L, tb) {
		cls = tb
	} else {
		size := tb.Len()
		if size == 0 {
			return lua.LNil, errors.New("attempt to index a non-table")
		}

		c := tb.RawGetInt(size)
		cc, ok := c.(*lua.LTable)
		if !ok {
			return lua.LNil, errors.Errorf("%s expected, got %s", lua.LTTable, c.Type())
		}

		if cc == nil {
			return lua.LNil, errors.New("attempt to index a non-table")
		}

		if !isClass(L, cc) {
			return lua.LNil, errors.Errorf("class expected, got %s", cls.Type())
		}

		tb.Remove(size)
		inject = tb
		cls = cc
	}

	o, err := newClass(L, cls)
	if err != nil {
		return lua.LNil, err
	}

	injectFunc := getMethod(L, o, diInjectorInjectClassMethodName)
	if injectFunc == nil {
		return o, nil
	}

	deps, err := dii.fnDef(L, o, injectFunc, inject, locals)
	if err != nil {
		return lua.LNil, err
	}

	L.Push(injectFunc)
	L.Push(o)

	if deps != nil {
		for _, dep := range deps {
			L.Push(dep)
		}
	}

	L.Call(len(deps)+1, 0)

	return o, nil
}

func (dii *diInjector) invoke(L *lua.LState, val lua.LValue,
	locals map[lua.LValue]lua.LValue) (lua.LValue, error) {
	typ := val.Type()
	if typ != lua.LTFunction && typ != lua.LTTable {
		return lua.LNil, errors.Errorf("%s or %s expected, got %s", lua.LTFunction, lua.LTTable, typ)
	}

	var callable *cb.Callable
	var inject *lua.LTable = nil
	if typ == lua.LTTable {
		tb, ok := val.(*lua.LTable)
		if !ok || tb == nil || tb.Len() == 0 {
			return lua.LNil, errors.New("attempt to index a non-table")
		}

		size := tb.Len()
		cc := tb.RawGetInt(size)
		tb.Remove(size)
		inject = tb

		ca, err := cb.New(L, cc)
		if err != nil {
			return lua.LNil, err
		}

		callable = ca
	} else {
		fn, ok := val.(*lua.LFunction)
		if !ok || fn == nil {
			return lua.LNil, errors.Errorf("%s expected, got %s", lua.LTFunction, typ)
		}

		ca, err := cb.New(L, fn)
		if err != nil {
			return lua.LNil, err
		}

		callable = ca
	}

	o := callable.Ref()
	fn, err := callable.ObjFn(L)
	if err != nil {
		return lua.LNil, err
	}

	deps, err := dii.fnDef(L, o, fn, inject, locals)
	if err != nil {
		return lua.LNil, err
	}

	n := len(deps)
	L.Push(fn)
	if o != lua.LNil {
		L.Push(o)
		n += 1
	}

	if deps != nil {
		for _, dep := range deps {
			L.Push(dep)
		}
	}

	L.Call(n, 1)

	return L.Get(-1), nil
}

func (dii *diInjector) value(L *lua.LState, val lua.LValue,
	locals map[lua.LValue]lua.LValue) (lua.LValue, error) {
	return val, nil
}
