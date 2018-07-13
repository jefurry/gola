// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package user

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	uuser "os/user"
)

func newUser(L *lua.LState, user *uuser.User) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = user

	L.SetMetatable(ud, L.GetTypeMetatable(osuserUserTypeName))

	return ud
}

func checkUser(L *lua.LState, n int) *uuser.User {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*uuser.User); ok {
		return v
	}

	L.ArgError(n, fmt.Sprintf("%s expected, got %s", osuserUserTypeName, ud.Type()))

	return nil
}

func newGroup(L *lua.LState, user *uuser.Group) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = user

	L.SetMetatable(ud, L.GetTypeMetatable(osuserGroupTypeName))

	return ud
}

func checkGroup(L *lua.LState, n int) *uuser.Group {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*uuser.Group); ok {
		return v
	}

	L.ArgError(n, fmt.Sprintf("%s expected, got %s", osuserGroupTypeName, ud.Type()))

	return nil
}
