// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package encoding implements encoding for Lua.
package encoding

import (
	"github.com/jefurry/gola/lua/libs/encoding/base32"
	"github.com/jefurry/gola/lua/libs/encoding/base64"
	"github.com/jefurry/gola/lua/libs/encoding/binary"
	"github.com/jefurry/gola/lua/libs/encoding/hex"
	"github.com/yuin/gopher-lua"
)

const (
	EncodingLibName = "encoding"
)

func Open(L *lua.LState) {
	L.PreloadModule(EncodingLibName, Loader)

	base64.Open(L)
	base32.Open(L)
	hex.Open(L)
	binary.Open(L)
}

func Loader(L *lua.LState) int {
	encodingmod := L.SetFuncs(L.NewTable(), encodingFuncs)
	L.Push(encodingmod)

	return 1
}

var encodingFuncs = map[string]lua.LGFunction{}

var encodingFields = map[string]lua.LValue{}
