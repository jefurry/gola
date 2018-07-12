// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package libs

import (
	"github.com/jefurry/gola/lua/libs/baselib"
	"github.com/jefurry/gola/lua/libs/bit32lib"
	"github.com/jefurry/gola/lua/libs/dilib"
	"github.com/jefurry/gola/lua/libs/eventlib"
	"github.com/jefurry/gola/lua/libs/execlib"
	"github.com/jefurry/gola/lua/libs/httplib"
	"github.com/jefurry/gola/lua/libs/jsonlib"
	"github.com/jefurry/gola/lua/libs/lfslib"
	"github.com/jefurry/gola/lua/libs/moonlib"
	"github.com/jefurry/gola/lua/libs/oslib"
	"github.com/jefurry/gola/lua/libs/relib"
	//"github.com/jefurry/gola/lua/libs/scrapelib"
	"github.com/jefurry/gola/lua/libs/encodinglib"
	"github.com/jefurry/gola/lua/libs/filepathlib"
	"github.com/jefurry/gola/lua/libs/jwtlib"
	"github.com/jefurry/gola/lua/libs/loglib"
	"github.com/jefurry/gola/lua/libs/pathlib"
	"github.com/jefurry/gola/lua/libs/socketlib"
	"github.com/jefurry/gola/lua/libs/syslib"
	"github.com/jefurry/gola/lua/libs/timelib"
	"github.com/jefurry/gola/lua/libs/urllib"
	"github.com/jefurry/gola/lua/libs/userlib"
	"github.com/jefurry/gola/lua/libs/xmlpathlib"
	"github.com/jefurry/gola/lua/libs/yamllib"
	"github.com/yuin/gopher-lua"
)

func OpenLibs(L *lua.LState) {
	baselib.Open(L)
	oslib.Open(L)
	execlib.Open(L)
	userlib.Open(L)
	syslib.Open(L)
	pathlib.Open(L)
	filepathlib.Open(L)
	timelib.Open(L)
	encodinglib.Open(L)
	dilib.Open(L)
	loglib.Open(L)
	eventlib.Open(L)
	relib.Open(L)
	httplib.Open(L)
	jsonlib.Open(L)
	yamllib.Open(L)
	lfslib.Open(L)
	urllib.Open(L)
	//scrapelib.Open(L)
	xmlpathlib.Open(L)
	socketlib.Open(L)
	bit32lib.Open(L)
	moonlib.Open(L)
	jwtlib.Open(L)
}
