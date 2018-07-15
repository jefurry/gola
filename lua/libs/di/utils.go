// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package di

import (
	"fmt"
	"github.com/yuin/gopher-lua"
)

func checkInjector(L *lua.LState, n int) *diInjector {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*diInjector); ok {
		return v
	}

	L.ArgError(n, fmt.Sprintf("%s expected, got %s", diInjectorTypeName, ud.Type()))

	return nil
}
