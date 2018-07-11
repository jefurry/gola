// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package userlib implements os/user for Lua.
package userlib

import (
	"github.com/yuin/gopher-lua"
)

const (
	UserLibName = "os.user"
)

func Open(L *lua.LState) {
	L.PreloadModule(UserLibName, Loader)
}

func Loader(L *lua.LState) int {
	usermod := L.SetFuncs(L.NewTable(), userFuncs)
	L.Push(usermod)

	userRegisterUserMetatype(L)

	return 1
}

var userFuncs = map[string]lua.LGFunction{
	"current":       userCurrent,
	"lookup":        userLookup,
	"lookupId":      userLookupId,
	"lookupGroup":   userLookupGroup,
	"lookupGroupId": userLookupGroupId,
}

var userFields = map[string]lua.LValue{}
