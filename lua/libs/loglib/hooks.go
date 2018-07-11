// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package loglib

import (
	"github.com/jefurry/logrus"
	"github.com/yuin/gopher-lua"
)

func logLoggerAddHook(L *lua.LState) int {
	logger := checkLogger(L, 1)
	ud := L.CheckUserData(2)

	hook, ok := ud.Value.(logrus.Hook)
	if !ok {
		L.ArgError(2, "invalid logger hook")
	}

	logger.AddHook(hook)

	return 0
}
