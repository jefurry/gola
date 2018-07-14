// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package binary

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	"strings"
)

const (
	binaryUnpackRetbTypeName = "encoding." + BinaryLibName + ".UNPACK_RETB*"
)

func newUnpackRetb(L *lua.LState, tb *lua.LTable) *lua.LTable {
	if tb == nil {
		tb = L.CreateTable(0, 0)
	}

	L.SetMetatable(tb, L.GetTypeMetatable(binaryUnpackRetbTypeName))

	return tb
}

func binaryUnpackRetbToString(L *lua.LState) int {
	tb := L.Get(-1).(*lua.LTable)

	out := make([]string, 0, 5)
	tb.ForEach(func(_, v lua.LValue) {
		out = append(out, fmt.Sprintf("%#v", v))
	})

	s := "<binary.unpack: { "
	s += strings.Join(out, ", ")
	s += " }>"

	L.Push(lua.LString(s))

	return 1
}

func binaryRegisterUnpackRetbMetatype(L *lua.LState) {
	// meta table
	mt := L.NewTypeMetatable(binaryUnpackRetbTypeName)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), binaryUnpackRetbFuncs))
	L.SetField(mt, "__tostring", L.NewFunction(binaryUnpackRetbToString))
}

var binaryUnpackRetbFuncs = map[string]lua.LGFunction{}
