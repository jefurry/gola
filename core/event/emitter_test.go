// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package event

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmitter_On(t *testing.T) {
	const EventClientType = "CLICK"

	var lis []*listener
	var ok bool
	var emit *Emitter

	var f1 Handler = func(*Event) bool {
		return true
	}

	emit = New()
	emit.On(EventClientType, f1)

	lis, ok = emit.listeners[EventClientType]
	if !assert.Equal(t, true, ok, "ok mismatching") {
		return
	}

	if !assert.Equal(t, 1, len(lis), "lis mismatching") {
		return
	}

	emit.Off(EventClientType, f1)

	lis, ok = emit.listeners[EventClientType]
	if !assert.Equal(t, true, ok, "ok mismatching") {
		return
	}

	if !assert.Equal(t, 0, len(lis), "lis mismatching") {
		return
	}
}

func TestEmitter_Once(t *testing.T) {
	const EventClientType = "CLICK"

	var lis []*listener
	var ok bool
	var emit *Emitter

	var f1 Handler = func(*Event) bool {
		return true
	}

	emit = New()
	emit.Once(EventClientType, f1)

	lis, ok = emit.listeners[EventClientType]
	if !assert.Equal(t, true, ok, "ok mismatching") {
		return
	}

	if !assert.Equal(t, 1, len(lis), "lis mismatching") {
		return
	}

	emit.Fire(EventClientType, nil)

	lis, ok = emit.listeners[EventClientType]
	if !assert.Equal(t, true, ok, "ok mismatching") {
		return
	}

	if !assert.Equal(t, 0, len(lis), "lis mismatching") {
		return
	}
}

func TestEmitter_Off_1(t *testing.T) {
	const EventClientType = "CLICK"
	const AddEventCount = 9

	var lis []*listener
	var ok bool
	var emit *Emitter

	var f1 Handler = func(*Event) bool {
		return true
	}

	emit = New()
	for i := 0; i < AddEventCount; i++ {
		emit.On(EventClientType, f1)
	}

	lis, ok = emit.listeners[EventClientType]
	if !assert.Equal(t, true, ok, "ok mismatching") {
		return
	}

	if !assert.Equal(t, 9, len(lis), "lis mismatching") {
		return
	}

	emit.Off(EventClientType, f1)

	lis, ok = emit.listeners[EventClientType]
	if !assert.Equal(t, true, ok, "ok mismatching") {
		return
	}

	if !assert.Equal(t, 0, len(lis), "lis mismatching") {
		return
	}
}

func TestEmitter_Off_2(t *testing.T) {
	const EventClientType = "CLICK"

	var lis []*listener
	var ok bool
	var emit *Emitter

	var f1 Handler = func(*Event) bool {
		return true
	}

	var f2 Handler = func(*Event) bool {
		return false
	}

	emit = New()
	emit.On(EventClientType, f1)
	emit.On(EventClientType, f2)
	emit.On(EventClientType, f2)
	emit.On(EventClientType, f1)
	emit.On(EventClientType, f2)
	emit.On(EventClientType, f1)
	emit.On(EventClientType, f2)
	emit.On(EventClientType, f1)
	emit.On(EventClientType, f2)

	lis, ok = emit.listeners[EventClientType]
	if !assert.Equal(t, true, ok, "ok mismatching") {
		return
	}

	if !assert.Equal(t, 9, len(lis), "lis mismatching") {
		return
	}

	emit.Off(EventClientType, f1)

	lis, ok = emit.listeners[EventClientType]
	if !assert.Equal(t, true, ok, "ok mismatching") {
		return
	}

	if !assert.Equal(t, 5, len(lis), "lis mismatching") {
		return
	}

	emit.Off(EventClientType, f2)

	lis, ok = emit.listeners[EventClientType]
	if !assert.Equal(t, true, ok, "ok mismatching") {
		return
	}

	if !assert.Equal(t, 0, len(lis), "lis mismatching") {
		return
	}
}

func TestEmitter_Off_3(t *testing.T) {
	const EventClientType = "CLICK"

	var lis []*listener
	var ok bool
	var emit *Emitter

	var f1 Handler = func(*Event) bool {
		return true
	}

	var f2 Handler = func(*Event) bool {
		return false
	}

	emit = New()
	emit.On(EventClientType, f1)
	emit.On(EventClientType, f2)
	emit.On(EventClientType, f2)
	emit.On(EventClientType, f1)
	emit.On(EventClientType, f2)
	emit.On(EventClientType, f1)
	emit.On(EventClientType, f2)
	emit.On(EventClientType, f1)
	emit.On(EventClientType, f2)

	lis, ok = emit.listeners[EventClientType]
	if !assert.Equal(t, true, ok, "ok mismatching") {
		return
	}

	if !assert.Equal(t, 9, len(lis), "lis mismatching") {
		return
	}

	emit.Off(EventClientType, f2)

	lis, ok = emit.listeners[EventClientType]
	if !assert.Equal(t, true, ok, "ok mismatching") {
		return
	}

	if !assert.Equal(t, 4, len(lis), "lis mismatching") {
		return
	}

	emit.Off(EventClientType, f1)

	lis, ok = emit.listeners[EventClientType]
	if !assert.Equal(t, true, ok, "ok mismatching") {
		return
	}

	if !assert.Equal(t, 0, len(lis), "lis mismatching") {
		return
	}
}
