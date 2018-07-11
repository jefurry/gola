// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package xmlpathlib implements xmlpath for Lua.
package xmlpathlib

import (
	"github.com/ailncode/gluaxmlpath"
	"github.com/yuin/gopher-lua"
)

const (
	XmlpathLibName = "xmlpath"
)

func Open(L *lua.LState) {
	L.PreloadModule(XmlpathLibName, gluaxmlpath.Loader)
}
