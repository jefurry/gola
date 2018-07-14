// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package event

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	"reflect"
)

func handlerPointer(handler *lua.LFunction) uintptr {
	return reflect.ValueOf(handler).Pointer()
}

func checkEmitter(L *lua.LState, n int) *emitter {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*emitter); ok {
		return v
	}

	L.ArgError(n, fmt.Sprintf("%s expected, got %s", eventEmitterTypeName, ud.Type()))

	return nil
}
