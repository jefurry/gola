// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package urllib implements url for Lua.
package urllib

import (
	"github.com/cjoudrey/gluaurl"
	"github.com/yuin/gopher-lua"
)

const (
	UrlLibName = "url"
)

func Open(L *lua.LState) {
	L.PreloadModule(UrlLibName, gluaurl.Loader)
}
