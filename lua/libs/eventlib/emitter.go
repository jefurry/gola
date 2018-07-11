// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package eventlib implements event emitter for Lua.
package eventlib

import (
	"fmt"
	"github.com/jefurry/gola/config"
	"github.com/pkg/errors"
	"github.com/yuin/gopher-lua"
	"reflect"
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

type (
	listener struct {
		priority lua.LNumber
		handler  *lua.LFunction
		pointer  uintptr
	}

	emitter struct {
		maxListeners lua.LNumber
		listeners    map[lua.LString][]*listener
	}
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

func newListener(priority lua.LNumber, handler *lua.LFunction) *listener {
	return &listener{
		priority: priority,
		handler:  handler,
		pointer:  handlerPointer(handler),
	}
}

func handlerPointer(handler *lua.LFunction) uintptr {
	return reflect.ValueOf(handler).Pointer()
}

func newEmitter(maxListeners int) *emitter {
	if maxListeners < 0 {
		maxListeners = defaultMaxListeners
	}

	return &emitter{
		maxListeners: lua.LNumber(maxListeners),
		listeners:    make(map[lua.LString][]*listener, defaultMaxListeners),
	}
}

// setMaxListeners obviously not all Emitters should be limited to 10. This function allows
// that to be increased. Set to zero for unlimited.
func (emit *emitter) setMaxListeners(n int) error {
	if n < 0 {
		return errors.Errorf("n(%v) must be a positive number", n)
	}

	emit.maxListeners = lua.LNumber(n)

	return nil
}

func (emit *emitter) emitterOn(etype lua.LString, handler *lua.LFunction, priority ...int) error {
	pri := defaultPriority
	if len(priority) > 0 {
		pri = priority[0]
		if pri < 0 {
			return errors.Errorf("priority(%v) must be a positive number", pri)
		}
	}

	li := newListener(lua.LNumber(pri), handler)
	lis, ok := emit.listeners[etype]
	if !ok {
		lis = make([]*listener, 0, defaultMaxListeners)
		lis = append(lis, li)
		emit.listeners[etype] = lis

		return nil
	}

	size := len(lis)
	if size >= int(emit.maxListeners) {
		fmt.Printf(`
			(%s) warning: possible Emitter memory leak detected. %d listeners added. 
			Use emitter.setMaxListeners() to increase limit.\n`,
			config.PROJECT_NAME, size)
	}

	for i := 0; i < size; i++ {
		oldLi := lis[i]
		if oldLi.priority < li.priority {
			lis = append(lis[:i], li)
			lis = append(lis, lis[i:]...)
			emit.listeners[etype] = lis

			return nil
		}
	}

	lis = append(lis, li)
	emit.listeners[etype] = lis

	return nil
}

func (emit *emitter) emitterOnce(L *lua.LState, etype lua.LString, handler *lua.LFunction, priority ...int) error {
	var wrappedFunc *lua.LFunction
	gfn := func(L *lua.LState) int {
		emit.emitterOff(etype, wrappedFunc)

		L.Push(handler)
		L.Push(L.CheckTable(1))
		L.Call(1, 1)

		L.Push(L.Get(-1))

		return 1
	}

	wrappedFunc = L.NewFunction(gfn)

	return emit.emitterOn(etype, wrappedFunc, priority...)
}

func (emit *emitter) emitterOff(etype lua.LString, handler *lua.LFunction) {
	lis, ok := emit.listeners[etype]
	if !ok || handler == nil {
		return
	}

	size := len(lis)
	// move through listeners from back to front
	// and remove matching listeners
	for i := size - 1; i >= 0; i-- {
		li := lis[i]
		if li.pointer == handlerPointer(handler) {
			if i >= size-1 {
				lis = lis[:i]
			} else {
				lis = append(lis[:i], lis[i+1:]...)
			}
		}
	}

	emit.listeners[etype] = lis
}

func (emit *emitter) emitterFire(L *lua.LState, etype lua.LString, data *lua.LTable, cxts ...lua.LValue) {
	lis, ok := emit.listeners[etype]
	if !ok {
		return
	}

	evt := newEvent(L, data, cxts...)

	for i := 0; i < len(lis); i++ {
		li := lis[i]
		L.Push(li.handler)
		L.Push(evt)
		L.Call(1, 1)
		// stopPropagation
		if !lua.LVAsBool(L.Get(-1)) {
			break
		}
	}
}

func (emit *emitter) getListeners(L *lua.LState, etype lua.LString) *lua.LTable {
	listeners, ok := emit.listeners[etype]
	if !ok {
		return nil
	}
	lis := L.CreateTable(len(listeners), 0)
	for i, listener := range listeners {
		li := L.CreateTable(0, 3)
		li.RawSetH(lua.LString("priority"), listener.priority)
		li.RawSetH(lua.LString("handler"), listener.handler)
		li.RawSetH(lua.LString("pointer"), lua.LNumber(listener.pointer))
		lis.RawSetInt(i+1, li)
	}

	L.Push(lis)

	return lis
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

func checkEmitter(L *lua.LState, n int) *emitter {
	ud := L.CheckUserData(n)
	if emit, ok := ud.Value.(*emitter); ok {
		return emit
	}

	L.ArgError(n, fmt.Sprintf("%s expected, got %s", eventEmitterTypeName, ud.Type()))

	return nil
}
