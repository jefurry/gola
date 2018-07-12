// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package encodinglib implements encoding for Lua.
package encodinglib

import (
	"github.com/yuin/gopher-lua"
)

const (
	EncodingLibName = "encoding"
)

func Open(L *lua.LState) {
	L.PreloadModule(EncodingLibName, Loader)

	OpenBase64(L)
	OpenBase32(L)
	OpenHex(L)
}

func Loader(L *lua.LState) int {
	encodingmod := L.SetFuncs(L.NewTable(), encodingFuncs)
	L.Push(encodingmod)

	return 1
}

var encodingFuncs = map[string]lua.LGFunction{}

var encodingFields = map[string]lua.LValue{}
