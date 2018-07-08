// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package libs

import (
	"github.com/jefurry/gola/libs/baselib"
	"github.com/jefurry/gola/libs/dilib"
	"github.com/jefurry/gola/libs/eventlib"
	"github.com/yuin/gopher-lua"
)

func OpenLibs(L *lua.LState) {
	baselib.Open(L)
	dilib.Open(L)
	eventlib.Open(L)
}
