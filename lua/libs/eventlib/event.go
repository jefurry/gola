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
	"github.com/yuin/gopher-lua"
)

const (
	EventLibName = "event"
)

func Open(L *lua.LState) {
	L.PreloadModule(EventLibName, Loader)
}

func Loader(L *lua.LState) int {
	eventRegisterEmitterMetatype(L)
	eventRegisterEventMetatype(L)

	eventmod := L.SetFuncs(L.NewTable(), eventFuncs)
	L.Push(eventmod)

	return 1
}

var eventFuncs = map[string]lua.LGFunction{
	"newEvent":   eventEventNew,
	"newEmitter": eventEmitterNew,
}
