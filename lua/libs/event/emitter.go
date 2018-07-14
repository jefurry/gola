// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package event

import (
	"github.com/yuin/gopher-lua"
)

const (
	eventEmitterTypeName = EventLibName + ".EMITTER*"
	eventEventTypeName   = EventLibName + ".EVENT*"
)

const (
	defaultPriority = 1000

	// By default Emitter will print a warning if more than 10 listeners are
	// added to it. This is a useful default which helps finding memory leaks.
	defaultMaxListeners = 10
)

func newEvent(L *lua.LState, data *lua.LTable, ctxs ...lua.LValue) *lua.LTable {
	evt := L.CreateTable(0, 2)
	if data == nil {
		evt.RawSetH(lua.LString("data"), lua.LNil)
	} else {
		evt.RawSetH(lua.LString("data"), data)
	}

	if len(ctxs) > 0 {
		evt.RawSetH(lua.LString("context"), ctxs[0])
	} else {
		evt.RawSetH(lua.LString("context"), lua.LNil)
	}

	return evt
}

func eventEventNew(L *lua.LState) int {
	data := L.OptTable(1, nil)
	L.CheckTypes(2, lua.LTTable, lua.LTUserData, lua.LTNil)
	ctx := L.Get(2)

	evt := newEvent(L, data, ctx)

	L.SetMetatable(evt, L.GetTypeMetatable(eventEventTypeName))
	L.Push(evt)

	return 1
}

func eventEmitterNew(L *lua.LState) int {
	maxListeners := L.OptInt(1, -1)
	if maxListeners < 0 {
		maxListeners = defaultMaxListeners
	}

	emit := newEmitter(maxListeners)
	ud := L.NewUserData()
	ud.Value = emit

	L.SetMetatable(ud, L.GetTypeMetatable(eventEmitterTypeName))
	L.Push(ud)

	return 1
}

func eventEmitterOn(L *lua.LState) int {
	emit := checkEmitter(L, 1)
	etype := L.CheckString(2)
	handler := L.CheckFunction(3)
	priority := L.OptInt(4, -1)
	if priority < 0 {
		priority = defaultPriority
	}

	err := emit.emitterOn(lua.LString(etype), handler, priority)
	if err != nil {
		L.Push(lua.LString(err.Error()))

		return 1
	}

	return 0
}

func eventEmitterOnce(L *lua.LState) int {
	emit := checkEmitter(L, 1)
	etype := L.CheckString(2)
	handler := L.CheckFunction(3)
	priority := L.OptInt(4, -1)
	if priority < 0 {
		priority = defaultPriority
	}

	err := emit.emitterOnce(L, lua.LString(etype), handler, priority)
	if err != nil {
		L.Push(lua.LString(err.Error()))

		return 1
	}

	return 0
}

func eventEmitterOff(L *lua.LState) int {
	emit := checkEmitter(L, 1)
	etype := L.CheckString(2)
	handler := L.CheckFunction(3)
	emit.emitterOff(lua.LString(etype), handler)

	return 0
}

func eventEmitterFire(L *lua.LState) int {
	emit := checkEmitter(L, 1)
	etype := L.CheckString(2)
	data := L.OptTable(3, nil)
	L.CheckTypes(4, lua.LTTable, lua.LTUserData, lua.LTNil)
	ctx := L.Get(4)

	emit.emitterFire(L, lua.LString(etype), data, ctx)

	return 0
}

func eventEmitterGetListeners(L *lua.LState) int {
	emit := checkEmitter(L, 1)
	etype := L.CheckString(2)
	L.Push(emit.getListeners(L, lua.LString(etype)))

	return 1
}

func eventEmitterSetMaxListeners(L *lua.LState) int {
	emit := checkEmitter(L, 1)
	n := L.OptInt(2, -1)

	err := emit.setMaxListeners(n)
	if err != nil {
		L.Push(lua.LString(err.Error()))

		return 1
	}

	return 0
}

func eventRegisterEmitterMetatype(L *lua.LState) {
	// meta table
	mt := L.NewTypeMetatable(eventEmitterTypeName)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), eventEmitterFuncs))
}

func eventRegisterEventMetatype(L *lua.LState) {
	// meta table
	mt := L.NewTypeMetatable(eventEventTypeName)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), eventEventFuncs))
}

var eventEmitterFuncs = map[string]lua.LGFunction{
	"on":              eventEmitterOn,
	"once":            eventEmitterOnce,
	"off":             eventEmitterOff,
	"fire":            eventEmitterFire,
	"getListeners":    eventEmitterGetListeners,
	"setMaxListeners": eventEmitterSetMaxListeners,
}

var eventEventFuncs = map[string]lua.LGFunction{}
