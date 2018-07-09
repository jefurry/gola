// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
	"testing"
)

func testCallTableSay(l *lua.LState) int {
	top := l.GetTop()
	args := l.CreateTable(top, 0)
	for i := 1; i <= top; i++ {
		args.RawSetInt(i, l.Get(i))
	}

	l.Push(args)

	return 1
}

func TestCall_1(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.CreateTable(0, 1)
	table.RawSetString("say", L.NewFunction(testCallTableSay))

	for _, v := range []struct {
		ref     lua.LValue
		fn      lua.LValue
		args    []lua.LValue
		objCall bool
	}{
		{
			lua.LNil,
			L.NewFunction(testCallTableSay),
			[]lua.LValue{lua.LString("Hello"),
				lua.LString("World")},
			false,
		},
		{
			table,
			lua.LString("say"),
			[]lua.LValue{lua.LString("Good"),
				lua.LString("Luck")},
			true,
		},
	} {
		val := L.CreateTable(2, 0)
		val.RawSetInt(1, v.ref)
		val.RawSetInt(2, v.fn)

		ret, err := Call(L, val, v.args...)
		if !assert.NoError(t, err, "Call should succeed") {
			return
		}

		if !assert.Equal(t, lua.LTTable, ret.Type(), fmt.Sprintf("%s expected, got %s", lua.LTTable, ret.Type())) {
			return
		}

		tb := ret.(*lua.LTable)

		if !assert.NotEqual(t, nil, tb, fmt.Sprintf("%s expected, got %s", lua.LTTable, lua.LNil)) {
			return
		}

		args := v.args
		if v.objCall {
			args = append([]lua.LValue{table}, args...)
		}

		if !assert.Equal(t, len(args), tb.Len(), "length no match") {
			return
		}

		for i, v := range args {
			if !assert.Equal(t, v, tb.RawGetInt(i+1), "value mismatching") {
				return
			}
		}
	}
}

func TestCall_2(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := L.CreateTable(0, 1)
	table.RawSetString("say", L.NewFunction(testCallTableSay))

	uv1 := L.CreateTable(2, 0)
	uv1.RawSetInt(1, table)
	uv1.RawSetInt(2, lua.LString("say"))

	uv2 := L.CreateTable(2, 0)
	uv2.RawSetInt(1, lua.LNil)
	uv2.RawSetInt(2, L.NewFunction(testCallTableSay))

	uv3 := L.NewFunction(testCallTableSay)

	for _, v := range []struct {
		uv      lua.LValue
		args    []lua.LValue
		objCall bool
	}{
		{
			uv1,
			[]lua.LValue{lua.LString("Hello"),
				lua.LString("World")},
			true,
		},
		{uv2,
			[]lua.LValue{lua.LString("Good"),
				lua.LString("Luck")},
			false,
		},
		{uv3,
			[]lua.LValue{lua.LString("Test")},
			false,
		},
	} {
		ret, err := Call(L, v.uv, v.args...)
		if !assert.NoError(t, err, "Call should succeed") {
			return
		}

		if !assert.Equal(t, lua.LTTable, ret.Type(), fmt.Sprintf("%s expected, got %s", lua.LTTable, ret.Type())) {
			return
		}

		tb := ret.(*lua.LTable)

		if !assert.NotEqual(t, nil, tb, fmt.Sprintf("%s expected, got %s", lua.LTTable, lua.LNil)) {
			return
		}

		args := v.args
		if v.objCall {
			args = append([]lua.LValue{table}, args...)
		}

		if !assert.Equal(t, len(args), tb.Len(), "length no match") {
			return
		}

		for i, v := range args {
			if !assert.Equal(t, v, tb.RawGetInt(i+1), "value mismatching") {
				return
			}
		}
	}
}
