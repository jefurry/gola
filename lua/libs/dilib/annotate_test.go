// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package dilib

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
	"testing"
)

func TestAnnotate_1(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	fn := L.NewFunction(func(l *lua.LState) int {
		arg1 := l.CheckString(1)
		arg2 := l.CheckString(2)

		l.Push(lua.LString(arg1 + " " + arg2))

		return 1
	})

	tb := L.CreateTable(3, 0)
	tb.RawSetInt(1, lua.LString("Good"))
	tb.RawSetInt(2, lua.LString("Luck"))
	tb.RawSetInt(3, fn)

	// annotate
	callable, args, err := annotate(L, tb)
	if !assert.NoError(t, err, "annotate should succeed") {
		return
	}

	f, err := callable.ObjFn(L)
	if !assert.NoError(t, err, "getObjFn should succeed") {
		return
	}

	if !assert.Equal(t, f, fn, "should be equals") {
		return
	}

	if !assert.Equal(t, 2, args.Len(), "should be equals") {
		return
	}
}

func TestAnnotate_2(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	fn := L.NewFunction(func(l *lua.LState) int {
		arg1 := l.CheckString(1)
		arg2 := l.CheckString(2)

		l.Push(lua.LString(arg1 + " " + arg2))

		return 1
	})

	table := L.CreateTable(3, 0)
	table.RawSetInt(1, lua.LString("Good"))
	table.RawSetInt(2, lua.LString("Luck"))
	table.RawSetInt(3, fn)

	// annotate
	callable, args, err := annotate(L, table)
	if !assert.NoError(t, err, "annotate should succeed") {
		return
	}

	fun, err := callable.ObjFn(L)
	if !assert.NoError(t, err, "getObjFn should succeed") {
		return
	}

	if !assert.Equal(t, fun, fn, "should be equals") {
		return
	}

	if !assert.Equal(t, 2, args.Len(), "should be equals") {
		return
	}

	// assoc
	if !assert.NoError(t, assoc(L, fn, args), "assoc should succeed") {
		return
	}

	// claim
	tb, err := claim(L, fn)
	if !assert.NoError(t, err, "claim should succeed") {
		return
	}

	if !assert.NotEqual(t, nil, tb, "should be not equals") {
		return
	}

	if !assert.Equal(t, 2, tb.Len(), "should be equals") {
		return
	}

	arg1 := tb.RawGetInt(1)
	if !assert.Equal(t, lua.LTString, arg1.Type(), fmt.Sprintf("%s expected, got %s", lua.LTString, arg1.Type())) {
		return
	}
	a1, _ := arg1.(lua.LString)
	if !assert.Equal(t, "Good", string(a1), "should be equals") {
		return
	}

	arg2 := tb.RawGetInt(2)
	if !assert.Equal(t, lua.LTString, arg2.Type(), fmt.Sprintf("%s expected, got %s", lua.LTString, arg2.Type())) {
		return
	}
	a2, _ := arg2.(lua.LString)
	if !assert.Equal(t, "Luck", string(a2), "should be equals") {
		return
	}

	// claim
	tba, err := claim(L, fun)
	if !assert.NoError(t, err, "claim should succeed") {
		return
	}

	if !assert.NotEqual(t, nil, tba, "should be not equals") {
		return
	}

	if !assert.Equal(t, 2, tba.Len(), "should be equals") {
		return
	}

	arg11 := tba.RawGetInt(1)
	if !assert.Equal(t, lua.LTString, arg11.Type(), fmt.Sprintf("%s expected, got %s", lua.LTString, arg11.Type())) {
		return
	}
	a11, _ := arg11.(lua.LString)
	if !assert.Equal(t, "Good", string(a11), "should be equals") {
		return
	}

	arg21 := tba.RawGetInt(2)
	if !assert.Equal(t, lua.LTString, arg21.Type(), fmt.Sprintf("%s expected, got %s", lua.LTString, arg21.Type())) {
		return
	}
	a21, _ := arg21.(lua.LString)
	if !assert.Equal(t, "Luck", string(a21), "should be equals") {
		return
	}

	// claim
	val := L.CreateTable(2, 0)
	val.RawSetInt(1, lua.LNil)
	val.RawSetInt(2, fn)

	tbb, err := claim(L, val)
	if !assert.NoError(t, err, "claim should succeed") {
		return
	}

	if !assert.NotEqual(t, nil, tbb, "should be not equals") {
		return
	}

	if !assert.Equal(t, 2, tbb.Len(), "should be equals") {
		return
	}

	arg12 := tbb.RawGetInt(1)
	if !assert.Equal(t, lua.LTString, arg12.Type(), fmt.Sprintf("%s expected, got %s", lua.LTString, arg12.Type())) {
		return
	}
	a12, _ := arg12.(lua.LString)
	if !assert.Equal(t, "Good", string(a12), "should be equals") {
		return
	}

	arg22 := tbb.RawGetInt(2)
	if !assert.Equal(t, lua.LTString, arg22.Type(), fmt.Sprintf("%s expected, got %s", lua.LTString, arg22.Type())) {
		return
	}
	a22, _ := arg22.(lua.LString)
	if !assert.Equal(t, "Luck", string(a22), "should be equals") {
		return
	}
}

func TestAnnotate_3(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	fn := L.NewFunction(func(l *lua.LState) int {
		arg1 := l.CheckString(1)
		arg2 := l.CheckString(2)

		l.Push(lua.LString(arg1 + " " + arg2))

		return 1
	})

	val := L.CreateTable(2, 0)
	val.RawSetInt(1, lua.LNil)
	val.RawSetInt(2, fn)

	table := L.CreateTable(3, 0)
	table.RawSetInt(1, lua.LString("Good"))
	table.RawSetInt(2, lua.LString("Luck"))
	table.RawSetInt(3, val)

	// annotate
	callable, args, err := annotate(L, table)
	if !assert.NoError(t, err, "annotate should succeed") {
		return
	}

	fun, err := callable.ObjFn(L)
	if !assert.NoError(t, err, "getObjFn should succeed") {
		return
	}

	if !assert.Equal(t, fun, fn, "should be equals") {
		return
	}

	if !assert.Equal(t, 2, args.Len(), "should be equals") {
		return
	}

	// assoc
	if !assert.NoError(t, assoc(L, fn, args), "assoc should succeed") {
		return
	}

	// claim
	tb, err := claim(L, fun)
	if !assert.NoError(t, err, "claim should succeed") {
		return
	}

	if !assert.NotEqual(t, nil, tb, "should be not equals") {
		return
	}

	if !assert.Equal(t, 2, tb.Len(), "should be equals") {
		return
	}

	arg1 := tb.RawGetInt(1)
	if !assert.Equal(t, lua.LTString, arg1.Type(), fmt.Sprintf("%s expected, got %s", lua.LTString, arg1.Type())) {
		return
	}
	a1, _ := arg1.(lua.LString)
	if !assert.Equal(t, "Good", string(a1), "should be equals") {
		return
	}

	arg2 := tb.RawGetInt(2)
	if !assert.Equal(t, lua.LTString, arg2.Type(), fmt.Sprintf("%s expected, got %s", lua.LTString, arg2.Type())) {
		return
	}
	a2, _ := arg2.(lua.LString)
	if !assert.Equal(t, "Luck", string(a2), "should be equals") {
		return
	}

	// claim
	tba, err := claim(L, fn)
	if !assert.NoError(t, err, "claim should succeed") {
		return
	}

	if !assert.NotEqual(t, nil, tba, "should be not equals") {
		return
	}

	if !assert.Equal(t, 2, tba.Len(), "should be equals") {
		return
	}

	arg11 := tba.RawGetInt(1)
	if !assert.Equal(t, lua.LTString, arg11.Type(), fmt.Sprintf("%s expected, got %s", lua.LTString, arg11.Type())) {
		return
	}
	a11, _ := arg11.(lua.LString)
	if !assert.Equal(t, "Good", string(a11), "should be equals") {
		return
	}

	arg21 := tba.RawGetInt(2)
	if !assert.Equal(t, lua.LTString, arg21.Type(), fmt.Sprintf("%s expected, got %s", lua.LTString, arg21.Type())) {
		return
	}
	a21, _ := arg21.(lua.LString)
	if !assert.Equal(t, "Luck", string(a21), "should be equals") {
		return
	}

	// claim
	tbb, err := claim(L, val)
	if !assert.NoError(t, err, "claim should succeed") {
		return
	}

	if !assert.NotEqual(t, nil, tbb, "should be not equals") {
		return
	}

	if !assert.Equal(t, 2, tbb.Len(), "should be equals") {
		return
	}

	arg12 := tba.RawGetInt(1)
	if !assert.Equal(t, lua.LTString, arg12.Type(), fmt.Sprintf("%s expected, got %s", lua.LTString, arg12.Type())) {
		return
	}
	a12, _ := arg12.(lua.LString)
	if !assert.Equal(t, "Good", string(a12), "should be equals") {
		return
	}

	arg22 := tba.RawGetInt(2)
	if !assert.Equal(t, lua.LTString, arg22.Type(), fmt.Sprintf("%s expected, got %s", lua.LTString, arg22.Type())) {
		return
	}
	a22, _ := arg21.(lua.LString)
	if !assert.Equal(t, "Luck", string(a22), "should be equals") {
		return
	}
}

func TestAnnotate_4(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	fn := L.NewFunction(func(l *lua.LState) int {
		arg1 := l.CheckString(1)
		arg2 := l.CheckString(2)

		l.Push(lua.LString(arg1 + " " + arg2))

		return 1
	})

	val := L.CreateTable(2, 0)
	val.RawSetInt(1, lua.LNil)
	val.RawSetInt(2, fn)

	table := L.CreateTable(3, 0)
	table.RawSetInt(1, lua.LString("Good"))
	table.RawSetInt(2, lua.LString("Luck"))
	table.RawSetInt(3, val)

	// annotate
	callable, args, err := annotate(L, table)
	if !assert.NoError(t, err, "annotate should succeed") {
		return
	}

	fun, err := callable.ObjFn(L)
	if !assert.NoError(t, err, "getObjFn should succeed") {
		return
	}

	if !assert.Equal(t, fun, fn, "should be equals") {
		return
	}

	if !assert.Equal(t, 2, args.Len(), "should be equals") {
		return
	}

	// assoc
	if !assert.NoError(t, assoc(L, fn, args), "assoc should succeed") {
		return
	}

	// claim
	tb, err := claim(L, fn)
	if !assert.NoError(t, err, "claim should succeed") {
		return
	}

	if !assert.NotEqual(t, nil, tb, "should be not equals") {
		return
	}

	if !assert.Equal(t, 2, tb.Len(), "should be equals") {
		return
	}

	arg1 := tb.RawGetInt(1)
	if !assert.Equal(t, lua.LTString, arg1.Type(), fmt.Sprintf("%s expected, got %s", lua.LTString, arg1.Type())) {
		return
	}
	a1, _ := arg1.(lua.LString)
	if !assert.Equal(t, "Good", string(a1), "should be equals") {
		return
	}

	arg2 := tb.RawGetInt(2)
	if !assert.Equal(t, lua.LTString, arg2.Type(), fmt.Sprintf("%s expected, got %s", lua.LTString, arg2.Type())) {
		return
	}
	a2, _ := arg2.(lua.LString)
	if !assert.Equal(t, "Luck", string(a2), "should be equals") {
		return
	}

	// claim
	tba, err := claim(L, fun)
	if !assert.NoError(t, err, "claim should succeed") {
		return
	}

	if !assert.NotEqual(t, nil, tba, "should be not equals") {
		return
	}

	if !assert.Equal(t, 2, tba.Len(), "should be equals") {
		return
	}

	arg11 := tba.RawGetInt(1)
	if !assert.Equal(t, lua.LTString, arg11.Type(), fmt.Sprintf("%s expected, got %s", lua.LTString, arg11.Type())) {
		return
	}
	a11, _ := arg11.(lua.LString)
	if !assert.Equal(t, "Good", string(a11), "should be equals") {
		return
	}

	arg21 := tba.RawGetInt(2)
	if !assert.Equal(t, lua.LTString, arg21.Type(), fmt.Sprintf("%s expected, got %s", lua.LTString, arg21.Type())) {
		return
	}
	a21, _ := arg21.(lua.LString)
	if !assert.Equal(t, "Luck", string(a21), "should be equals") {
		return
	}

	// claim
	tbb, err := claim(L, val)
	if !assert.NoError(t, err, "claim should succeed") {
		return
	}

	if !assert.NotEqual(t, nil, tbb, "should be not equals") {
		return
	}

	if !assert.Equal(t, 2, tbb.Len(), "should be equals") {
		return
	}

	arg12 := tba.RawGetInt(1)
	if !assert.Equal(t, lua.LTString, arg12.Type(), fmt.Sprintf("%s expected, got %s", lua.LTString, arg12.Type())) {
		return
	}
	a12, _ := arg12.(lua.LString)
	if !assert.Equal(t, "Good", string(a12), "should be equals") {
		return
	}

	arg22 := tba.RawGetInt(2)
	if !assert.Equal(t, lua.LTString, arg22.Type(), fmt.Sprintf("%s expected, got %s", lua.LTString, arg22.Type())) {
		return
	}
	a22, _ := arg21.(lua.LString)
	if !assert.Equal(t, "Luck", string(a22), "should be equals") {
		return
	}

	// dissoc
	if !assert.NoError(t, dissoc(L, val), "dissoc should succeed") {
		return
	}

	// claim
	_, err = claim(L, val)
	if !assert.Error(t, err, "claim should succeed") {
		return
	}

	// claim
	_, err = claim(L, fn)
	if !assert.Error(t, err, "claim should succeed") {
		return
	}

	// claim
	_, err = claim(L, fun)
	if !assert.Error(t, err, "claim should succeed") {
		return
	}
}

func TestAnnotate_5(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local di = require('di')

	local version = 1.0

	function test(a, b, c)
		local v = version
	end

	if di.annotate("good", "luck", "!", test) ~= true then
		return false
	end

	local args = di.claim(test)
	if table.maxn(args) ~= 3 or #args ~= 3 then
		return false
	end
	if args[1] ~= "good" then
		return false
	end
	if args[2] ~= "luck" then
		return false
	end
	if args[3] ~= "!" then
		return false
	end

	return true
	`

	err := L.DoString(code)
	if !assert.NoError(t, err, `L.DoString should succeed`) {
		return
	}

	if !assert.Equal(t, 1, L.GetTop(), "L.GetTop mismatching") {
		return
	}

	ret := L.Get(-1)
	if !assert.Equal(t, lua.LTBool, ret.Type(), "type mismatching") {
		return
	}

	if !assert.Equal(t, lua.LTrue, ret.(lua.LBool), "value mismatching") {
		return
	}
}

func TestAnnotate_6(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local di = require('di')

	local version = 1.0

	function test(a, b, c)
		local v = version
	end

	if di.annotate{"good", "luck", "!", test} ~= true then
		return false
	end

	local args = di.claim{test}
	if table.maxn(args) ~= 3 or #args ~= 3 then
		return false
	end
	if args[1] ~= "good" then
		return false
	end
	if args[2] ~= "luck" then
		return false
	end
	if args[3] ~= "!" then
		return false
	end

	return true
	`

	err := L.DoString(code)
	if !assert.NoError(t, err, `L.DoString should succeed`) {
		return
	}

	if !assert.Equal(t, 1, L.GetTop(), "L.GetTop mismatching") {
		return
	}

	ret := L.Get(-1)
	if !assert.Equal(t, lua.LTBool, ret.Type(), "type mismatching") {
		return
	}

	if !assert.Equal(t, lua.LTrue, ret.(lua.LBool), "value mismatching") {
		return
	}
}

func TestAnnotate_7(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local di = require('di')

	local version = 1.0

	function test(a, b, c)
		local v = version
	end

	local tb = {
		t=test
	}

	if di.annotate("good", "luck", "!", {tb, "t"}) ~= true then
		return false
	end

	local args = di.claim({tb, "t"})
	if table.maxn(args) ~= 3 or #args ~= 3 then
		return false
	end
	if args[1] ~= "good" then
		return false
	end
	if args[2] ~= "luck" then
		return false
	end
	if args[3] ~= "!" then
		return false
	end

	return true
	`

	err := L.DoString(code)
	if !assert.NoError(t, err, `L.DoString should succeed`) {
		return
	}

	if !assert.Equal(t, 1, L.GetTop(), "L.GetTop mismatching") {
		return
	}

	ret := L.Get(-1)
	if !assert.Equal(t, lua.LTBool, ret.Type(), "type mismatching") {
		return
	}

	if !assert.Equal(t, lua.LTrue, ret.(lua.LBool), "value mismatching") {
		return
	}
}

func TestAnnotate_8(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local di = require('di')

	local version = 1.0

	function test(a, b, c)
		local v = version
	end

	local tb = {
		t=test
	}

	if di.annotate{"good", "luck", "!", {tb, "t"}} ~= true then
		return false
	end

	local args = di.claim{tb, "t"}
	if table.maxn(args) ~= 3 or #args ~= 3 then
		return false
	end
	if args[1] ~= "good" then
		return false
	end
	if args[2] ~= "luck" then
		return false
	end
	if args[3] ~= "!" then
		return false
	end

	return true
	`

	err := L.DoString(code)
	if !assert.NoError(t, err, `L.DoString should succeed`) {
		return
	}

	if !assert.Equal(t, 1, L.GetTop(), "L.GetTop mismatching") {
		return
	}

	ret := L.Get(-1)
	if !assert.Equal(t, lua.LTBool, ret.Type(), "type mismatching") {
		return
	}

	if !assert.Equal(t, lua.LTrue, ret.(lua.LBool), "value mismatching") {
		return
	}
}

func TestAssoc(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local di = require('di')

	local version = 1.0

	function test(a, b, c)
		local v = version
	end

	local tb = {
		t=test
	}

	if di.assoc({tb, "t"}, {"good", "luck", "!"}) ~= true then
		return false
	end

	local args = di.claim{tb, "t"}
	if table.maxn(args) ~= 3 or #args ~= 3 then
		return false
	end
	if args[1] ~= "good" then
		return false
	end
	if args[2] ~= "luck" then
		return false
	end
	if args[3] ~= "!" then
		return false
	end

	return true
	`

	err := L.DoString(code)
	if !assert.NoError(t, err, `L.DoString should succeed`) {
		return
	}

	if !assert.Equal(t, 1, L.GetTop(), "L.GetTop mismatching") {
		return
	}

	ret := L.Get(-1)
	if !assert.Equal(t, lua.LTBool, ret.Type(), "type mismatching") {
		return
	}

	if !assert.Equal(t, lua.LTrue, ret.(lua.LBool), "value mismatching") {
		return
	}
}

func TestDissoc_1(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local di = require('di')

	local version = 1.0

	function test(a, b, c)
		local v = version
	end

	if di.annotate("good", "luck", "!", test) ~= true then
		return false
	end

	local args = di.claim(test)
	if table.maxn(args) ~= 3 or #args ~= 3 then
		return false
	end
	if args[1] ~= "good" then
		return false
	end
	if args[2] ~= "luck" then
		return false
	end
	if args[3] ~= "!" then
		return false
	end

	if di.dissoc(test) ~= true then
		return false
	end
	local args = di.claim(test)
	if args ~= nil then
		return false
	end

	return true
	`

	err := L.DoString(code)
	if !assert.NoError(t, err, `L.DoString should succeed`) {
		return
	}

	if !assert.Equal(t, 1, L.GetTop(), "L.GetTop mismatching") {
		return
	}

	ret := L.Get(-1)
	if !assert.Equal(t, lua.LTBool, ret.Type(), "type mismatching") {
		return
	}

	if !assert.Equal(t, lua.LTrue, ret.(lua.LBool), "value mismatching") {
		return
	}
}

func TestDissoc_2(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local di = require('di')

	local version = 1.0

	function test(a, b, c)
		local v = version
	end

	local tb = {
		t=test
	}

	if di.assoc({tb, "t"}, {"good", "luck", "!"}) ~= true then
		return false
	end

	local args = di.claim{tb, "t"}
	if table.maxn(args) ~= 3 or #args ~= 3 then
		return false
	end
	if args[1] ~= "good" then
		return false
	end
	if args[2] ~= "luck" then
		return false
	end
	if args[3] ~= "!" then
		return false
	end

	if di.dissoc{tb, "t"} ~= true then
		return false
	end
	local args = di.claim{tb, "t"}
	if args ~= nil then
		return false
	end

	return true
	`

	err := L.DoString(code)
	if !assert.NoError(t, err, `L.DoString should succeed`) {
		return
	}

	if !assert.Equal(t, 1, L.GetTop(), "L.GetTop mismatching") {
		return
	}

	ret := L.Get(-1)
	if !assert.Equal(t, lua.LTBool, ret.Type(), "type mismatching") {
		return
	}

	if !assert.Equal(t, lua.LTrue, ret.(lua.LBool), "value mismatching") {
		return
	}
}
