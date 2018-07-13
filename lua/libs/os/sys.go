// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package os

import (
	"github.com/yuin/gopher-lua"
	oos "os"
)

func osHostname(L *lua.LState) int {
	name, err := oos.Hostname()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	L.Push(lua.LString(name))

	return 1
}
