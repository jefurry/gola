// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package dilib

import (
	"github.com/stretchr/testify/assert"
	"github.com/yuin/gopher-lua"
	"testing"
)

func TestCall_3(t *testing.T) {
	L := lua.NewState()
	L.PreloadModule(DiLibName, Loader)
	defer L.Close()

	code := `
	local di = require('di')

	function test(a, b, c)
		return string.format("%s-%s-%s", a, b, c)
	end

	local Animal = di.createClass{name="Animal"}
	function Animal:test(a, b, c)
		return string.format("%s: %s-%s-%s", self.name, a, b, c)
	end

	local cat = Animal:new{name="Cat"}
	local args = {"arg a", "arg b", "arg c"}

	local items = {
		test, "arg a-arg b-arg c";
		{cat, "test"}, "Cat: arg a-arg b-arg c";
		{nil, test}, "arg a-arg b-arg c";
	}

	for i = 1, #items, 2 do
		if di.call(items[i], unpack(args)) ~= items[i+1] then
			return false
		end

		if di.apply(items[i], args) ~= items[i+1] then
			return false
		end
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

func TestCall_4(t *testing.T) {
	L := lua.NewState()
	L.PreloadModule(DiLibName, Loader)
	defer L.Close()

	code := `
	local di = require('di')

	function getName(self, who)
		return string.format("%s say: my name is %s", who, self.name)
	end

	local Animal = di.createClass{name="Animal"}
	function Animal:change(name)
		self.name = name
	end
	function Animal:test(a, b, c)
		return string.format("%s: %s-%s-%s", self.name, a, b, c)
	end

	local cat = Animal:new{name="Cat"}
	local args = {"arg a", "arg b", "arg c"}

	local items = {
		{nil, getName, {"Jeff"}}, "nil say: my name is nil";
		{cat, getName, {"Goal"}}, "Goal say: my name is Cat";
		{nil, {cat, "change"}, {"Dog"}}, nil;
		{cat, getName, {"Goal"}}, "Goal say: my name is Dog";
		{cat, {cat, "test"}, {"arg a", "arg b", "arb c"}}, "Dog: arg a-arg b-arb c";
		{cat, {cat, "change"}, {"Bee"}}, nil;
		{cat, getName, {"Goal"}}, "Goal say: my name is Bee";
		{nil, {cat, "test"}, {"arg a", "arg b", "arb c"}}, "Bee: arg a-arg b-arb c";
	}

	for i = 1, #items, 2 do
		item = items[i]
		if di.bind(item[1], item[2], item[3]) ~= items[i+1] then
			return false
		end
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
