// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package libs

import (
	"github.com/jefurry/gola/lua/libs/base"
	"github.com/jefurry/gola/lua/libs/bit32"
	"github.com/jefurry/gola/lua/libs/di"
	"github.com/jefurry/gola/lua/libs/event"
	"github.com/jefurry/gola/lua/libs/http"
	"github.com/jefurry/gola/lua/libs/json"
	"github.com/jefurry/gola/lua/libs/lfs"
	"github.com/jefurry/gola/lua/libs/moon"
	"github.com/jefurry/gola/lua/libs/os"
	"github.com/jefurry/gola/lua/libs/re"
	//"github.com/jefurry/gola/lua/libs/scrape"
	"github.com/jefurry/gola/lua/libs/encoding"
	"github.com/jefurry/gola/lua/libs/jwt"
	"github.com/jefurry/gola/lua/libs/log"
	"github.com/jefurry/gola/lua/libs/path"
	"github.com/jefurry/gola/lua/libs/socket"
	"github.com/jefurry/gola/lua/libs/sys"
	"github.com/jefurry/gola/lua/libs/time"
	"github.com/jefurry/gola/lua/libs/url"
	"github.com/jefurry/gola/lua/libs/xmlpath"
	"github.com/jefurry/gola/lua/libs/yaml"
	"github.com/yuin/gopher-lua"
)

func OpenLibs(L *lua.LState) {
	base.Open(L)
	bit32.Open(L)
	os.Open(L)
	sys.Open(L)
	path.Open(L)
	time.Open(L)
	encoding.Open(L)
	di.Open(L)
	log.Open(L)
	event.Open(L)
	re.Open(L)
	http.Open(L)
	json.Open(L)
	yaml.Open(L)
	url.Open(L)
	jwt.Open(L)
	//scrapelib.Open(L)
	xmlpath.Open(L)
	socket.Open(L)
	lfs.Open(L)

	moon.Open(L)
}
