// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package di implements DI(Dependency Injection) for Lua.
package di

import (
	"github.com/yuin/gopher-lua"
)

const (
	DiLibName = "di"
)

func Open(L *lua.LState) {
	L.PreloadModule(DiLibName, Loader)
}

func Loader(L *lua.LState) int {
	diRegisterClassMetatype(L)
	diRegisterInjectorMetatype(L)

	dimod := L.SetFuncs(L.NewTable(), diFuncs)
	L.Push(dimod)

	return 1
}

var diFuncs = map[string]lua.LGFunction{
	"createClass": diCreateClass,
	"isClass":     diClassIsClass,
	"instanceof":  diClassInstanceof,
	"getMethod":   diClassGetMethod,
	"parse":       diParse,
	"annotate":    diAnnotate,
	"assoc":       diAssoc,
	"claim":       diClaim,
	"dissoc":      diDissoc,
	"newInjector": diInjectorNew,
	"call":        diCall,
	"apply":       diApply,
	"bind":        diBind,
}
