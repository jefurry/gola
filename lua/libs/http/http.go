// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package http implements http request for Lua.
package http

import (
	"github.com/cjoudrey/gluahttp"
	"github.com/yuin/gopher-lua"
	"net/http"
)

const (
	HttpLibName = "http"
)

func Open(L *lua.LState) {
	L.PreloadModule(HttpLibName, gluahttp.NewHttpModule(&http.Client{}).Loader)
}
