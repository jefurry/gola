// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package socketlib implements socket for Lua.
package socketlib

import (
	"github.com/BixData/gluasocket"
	"github.com/yuin/gopher-lua"
)

const (
	SocketLibName = "socket"
)

func Open(L *lua.LState) {
	gluasocket.Preload(L)
}
