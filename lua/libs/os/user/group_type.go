// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
	"github.com/yuin/gopher-lua"
)

const (
	userGroupTypeName = UserLibName + ".GROUP*"
)

func userRegisterGroupMetatype(L *lua.LState) {
	// meta table
	mt := L.NewTypeMetatable(userGroupTypeName)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), userGroupFuncs))
}

var userGroupFuncs = map[string]lua.LGFunction{
	"gid":  userGroupGid,
	"name": userGroupName,
}

func userGroupGid(L *lua.LState) int {
	g := checkGroup(L, 1)

	L.Push(lua.LString(g.Gid))

	return 1
}

func userGroupName(L *lua.LState) int {
	g := checkGroup(L, 1)

	L.Push(lua.LString(g.Name))

	return 1
}
