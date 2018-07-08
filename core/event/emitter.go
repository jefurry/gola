// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package event implements event manager.
package event

import (
	"fmt"
	"github.com/jefurry/gola/config"
	"reflect"
)

const (
	DefaultPriority = 1000

	// By default Emitter will print a warning if more than 10 listeners are
	// added to it. This is a useful default which helps finding memory leaks.
	DefaultMaxListeners = 10
)

type (
	Handler func(*Event) bool

	listener struct {
		priority int
		handler  Handler
		pointer  uintptr
	}

	Emitter struct {
		maxListeners int
		listeners    map[string][]*listener
	}
)

func newListener(priority int, handler Handler) *listener {
	return &listener{
		priority: priority,
		handler:  handler,
		pointer:  handlerPointer(handler),
	}
}

func handlerPointer(handler Handler) uintptr {
	return reflect.ValueOf(handler).Pointer()
}

func NewEmitter() *Emitter {
	return &Emitter{
		maxListeners: DefaultMaxListeners,
		listeners:    make(map[string][]*listener, 10),
	}
}

// SetMaxListeners obviously not all Emitters should be limited to 10. This function allows
// that to be increased. Set to zero for unlimited.
func (emit *Emitter) SetMaxListeners(n int) {
	if n < 0 {
		panic("n must be a positive number")
	}

	emit.maxListeners = n
}

// On register an event listener that is executed only once.
func (emit *Emitter) On(etype string, handler Handler, priority ...int) {
	pri := DefaultPriority
	if len(priority) > 0 {
		if priority[0] < 0 {
			panic("priority must be a number")
		}

		pri = priority[0]
	}

	li := newListener(pri, handler)
	lis, ok := emit.listeners[etype]
	if !ok {
		lis := make([]*listener, 0, DefaultMaxListeners)
		lis = append(lis, li)
		emit.listeners[etype] = lis

		return
	}

	size := len(emit.listeners)
	if size >= emit.maxListeners {
		fmt.Printf(`
			(%s) warning: possible Emitter memory leak detected. %d listeners added. 
			Use emitter.setMaxListeners() to increase limit.\n`,
			config.PROJECT_NAME, size)
	}

	for i := 0; i < len(lis); i++ {
		oldLi := lis[i]
		if oldLi.priority < li.priority {
			lis = append(lis[:i], li)
			lis = append(lis, lis[i:]...)
			emit.listeners[etype] = lis

			return
		}
	}

	lis = append(lis, li)
	emit.listeners[etype] = lis
}

// Once register an event listener that is executed only once.
func (emit *Emitter) Once(etype string, handler Handler, priority ...int) {
	var wrapppedFunc Handler
	wrapppedFunc = func(*Event) bool {
		emit.Off(etype, wrapppedFunc)

		return true
	}

	emit.On(etype, wrapppedFunc, priority...)
}

// Off removes event listeners by event type and handler.
// If no handler is given, all listeners for a given event type are being removed.
func (emit *Emitter) Off(etype string, handler Handler) {
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

func (emit *Emitter) Fire(etype string, data interface{}, cxts ...interface{}) {
	lis, ok := emit.listeners[etype]
	if !ok {
		return
	}

	evt := NewEvent(data, cxts...)

	for i := 0; i < len(lis); i++ {
		li := lis[i]
		// stopPropagation
		if !li.handler(evt) {
			break
		}
	}
}
