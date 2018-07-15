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
	userUserTypeName = UserLibName + ".USER*"
)

func userRegisterUserMetatype(L *lua.LState) {
	// meta table
	mt := L.NewTypeMetatable(userUserTypeName)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), userUserFuncs))
}

var userUserFuncs = map[string]lua.LGFunction{
	"uid":      userUserUid,
	"gid":      userUserGid,
	"username": userUserUsername,
	"name":     userUserName,
	"homeDir":  userUserHomeDir,
	"groupIds": userUserGroupIds,
}

func userUserUid(L *lua.LState) int {
	u := checkUser(L, 1)

	L.Push(lua.LString(u.Uid))

	return 1
}

func userUserGid(L *lua.LState) int {
	u := checkUser(L, 1)

	L.Push(lua.LString(u.Gid))

	return 1
}

func userUserUsername(L *lua.LState) int {
	u := checkUser(L, 1)

	L.Push(lua.LString(u.Username))

	return 1
}

func userUserName(L *lua.LState) int {
	u := checkUser(L, 1)

	L.Push(lua.LString(u.Name))

	return 1
}

func userUserHomeDir(L *lua.LState) int {
	u := checkUser(L, 1)

	L.Push(lua.LString(u.HomeDir))

	return 1
}

func userUserGroupIds(L *lua.LState) int {
	u := checkUser(L, 1)

	ids, err := u.GroupIds()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))

		return 2
	}

	tb := L.CreateTable(len(ids), 0)
	for i, v := range ids {
		tb.RawSetInt(i, lua.LString(v))
	}

	L.Push(tb)

	return 1
}
