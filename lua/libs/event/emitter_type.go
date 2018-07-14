// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package event

import (
	"fmt"
	"github.com/jefurry/gola/config"
	"github.com/pkg/errors"
	"github.com/yuin/gopher-lua"
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

func newListener(priority lua.LNumber, handler *lua.LFunction) *listener {
	return &listener{
		priority: priority,
		handler:  handler,
		pointer:  handlerPointer(handler),
	}
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
