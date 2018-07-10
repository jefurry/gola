// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package timelib

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	"time"
)

func isLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func checkTime(L *lua.LState, n int) time.Time {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(time.Time); ok {
		return v
	}

	L.ArgError(n, fmt.Sprintf("%s expected, got %s", timeLocationTypeName, ud.Type()))

	return time.Time{}
}

func checkLocation(L *lua.LState, n int) *time.Location {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*time.Location); ok {
		return v
	}

	L.ArgError(n, fmt.Sprintf("%s expected, got %s", timeLocationTypeName, ud.Type()))

	return nil
}

func newTime(L *lua.LState, t time.Time) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = t

	L.SetMetatable(ud, L.GetTypeMetatable(timeTimeTypeName))

	return ud
}

func newLocation(L *lua.LState, l *time.Location) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = l

	L.SetMetatable(ud, L.GetTypeMetatable(timeLocationTypeName))

	return ud
}
