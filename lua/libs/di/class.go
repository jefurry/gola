// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package di

import (
	"crypto/md5"
	"fmt"
	"github.com/pkg/errors"
	"github.com/yuin/gopher-lua"
	"strings"
	"sync/atomic"
	"time"
)

const (
	diClassTypeName      = DiLibName + ".CLASS*"
	diClassClassFlagName = "@class"
)

var (
	diCreateClassKey      string
	diCreateClassKeyIndex uint32
)

func init() {
	s1 := fmt.Sprintf("%p", &(struct{}{}))
	s2 := fmt.Sprintf("%d", time.Now().UnixNano())
	s := fmt.Sprintf("%s#%s", s1, s2)
	diCreateClassKey = fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func diCreateClass(L *lua.LState) int {
	o := L.OptTable(1, nil)
	tb := L.NewTable()

	idx := atomic.AddUint32(&diCreateClassKeyIndex, 1)
	ckey := keyIndex(idx, diCreateClassKey, true)
	key := keyIndex(idx, diCreateClassKey, false)

	nmt := L.NewTable()
	L.SetField(nmt, "__metatable", lua.LString(ckey))
	L.SetMetatable(tb, nmt)

	mt := L.GetTypeMetatable(diClassTypeName).(*lua.LTable)
	index := mt.RawGetString("__index").(*lua.LTable)
	index.ForEach(func(k, v lua.LValue) {
		tb.RawSet(k, v)
	})

	upvalues := []lua.LValue{o, lua.LString(key)}
	tb.RawSetString("new", L.NewClosure(diClassNew, upvalues...))
	L.Push(tb)

	return 1
}

func diClassNew(L *lua.LState) int {
	self := L.CheckTable(1)
	key := L.Get(lua.UpvalueIndex(2)).(lua.LString)
	if key == lua.LNil {
		L.ArgError(1, "attempt to call a non-class object")
	}

	uv := L.Get(lua.UpvalueIndex(1)).(*lua.LTable)
	o := L.ToTable(2)

	tb := L.NewTable()
	if uv != nil {
		uv.ForEach(func(k, v lua.LValue) {
			tb.RawSet(k, v)
		})
	}
	if o != nil {
		o.ForEach(func(k, v lua.LValue) {
			tb.RawSet(k, v)
		})
	}

	L.SetField(self, "__index", self)
	L.SetField(self, "__tostring", L.NewFunction(diClassToString))
	L.SetField(self, "__metatable", key)
	L.SetMetatable(tb, self)

	initFunc := self.RawGetString("init")
	if f, ok := initFunc.(*lua.LFunction); ok && f != nil {
		L.Push(f)
		L.Push(tb)
		L.Call(1, 0)
	}

	L.Push(tb)

	return 1
}

func diClassToString(L *lua.LState) int {
	tb := L.Get(-1).(*lua.LTable)

	L.Push(lua.LString(fmt.Sprintf("<class table: %p>", &tb)))

	return 1
}

func diClassIsClass(L *lua.LState) int {
	tb := L.CheckTable(1)

	if isClass(L, tb) {
		L.Push(lua.LTrue)
	} else {
		L.Push(lua.LFalse)
	}

	return 1
}

func diClassInstanceof(L *lua.LState) int {
	tbch := L.CheckTable(1)
	tbsu := L.CheckTable(2)

	if instanceof(L, tbch, tbsu) {
		L.Push(lua.LTrue)
	} else {
		L.Push(lua.LFalse)
	}

	return 1
}

func diClassGetMethod(L *lua.LState) int {
	L.CheckTypes(1, lua.LTTable, lua.LTUserData)
	obj := L.CheckAny(1)
	m := L.CheckString(2)

	fn := getMethod(L, obj, m)
	L.Push(fn)

	return 1
}

func diRegisterClassMetatype(L *lua.LState) {
	// meta table
	mt := L.NewTypeMetatable(diClassTypeName)

	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), diClassFuncs))
}

func keyIndex(idx uint32, key string, issuper bool) string {
	su := ""
	if issuper {
		su = diClassClassFlagName
	}
	return fmt.Sprintf("%s#%d%s", key, idx, su)
}

func isClass(L *lua.LState, tb *lua.LTable) bool {
	mt := L.GetMetatable(tb)
	if mt.Type() != lua.LTString {
		return false
	}

	mts := string(mt.(lua.LString))

	if strings.HasPrefix(mts, diCreateClassKey) && strings.HasSuffix(mts, diClassClassFlagName) {
		return true
	}

	return false
}

func newClass(L *lua.LState, cls *lua.LTable) (*lua.LTable, error) {
	fn := getMethod(L, cls, "new")
	if fn == nil {
		return nil, errors.New("attempt to call a non-function object")
	}

	L.Push(fn)
	L.Push(cls)
	L.Call(1, 1)

	obj := L.Get(-1)
	o, ok := obj.(*lua.LTable)
	if !ok {
		return nil, errors.Errorf("%s expected, got %s", lua.LTTable, obj.Type())
	}

	if o == nil || !instanceof(L, o, cls) {
		return nil, errors.New("make class instance failed")
	}

	return o, nil
}

func getMethod(L *lua.LState, cls lua.LValue, method string) *lua.LFunction {
	fun := L.GetField(cls, method)
	fn, ok := fun.(*lua.LFunction)
	if !ok {
		return nil
	}
	if fn == nil {
		return nil
	}

	return fn
}

func instanceof(L *lua.LState, tbch, tbsu *lua.LTable) bool {
	mtch := L.GetMetatable(tbch)
	if mtch.Type() != lua.LTString {
		return false
	}

	mtsu := L.GetMetatable(tbsu)
	if mtsu.Type() != lua.LTString {
		return false
	}

	ch := string(mtch.(lua.LString))
	su := string(mtsu.(lua.LString))
	if strings.HasPrefix(su, ch) && (ch+diClassClassFlagName) == su {
		return true
	}

	return false
}

var diClassFuncs = map[string]lua.LGFunction{}
