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

func TestCreateClass_1(t *testing.T) {
	L := lua.NewState()
	L.PreloadModule(DiLibName, Loader)
	defer L.Close()

	code := `
	local di = require('di')

	local Animal = di.createClass({name="Animal"})

	function Animal:init()
		self.name = "@" .. self.name
	end

	function Animal:say(msg)
		return string.format("%s say: %s", self.name, msg)
	end

	local dog = Animal:new({name="Dog", age=5})

	local cat = Animal:new{name="Cat"}

	if dog.name ~= "@Dog" then
		return false
	end

	if cat.name ~= "@Cat" then
		return false
	end

	if dog:say("good night") ~= "@Dog say: good night" then
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

func TestCreateClass_2(t *testing.T) {
	L := lua.NewState()
	L.PreloadModule(DiLibName, Loader)
	defer L.Close()

	code := `
	local di = require('di')

	local init_called = false

	local Animal = di.createClass({name="Animal"})

	function Animal:init()
		init_called = true
	end

	function Animal:say(msg)
		return string.format("%s say: %s", self.name, msg)
	end

	function Animal:change(name)
		self.name = name
	end

	local dog = Animal:new({name="Dog", age=5})
	local cat = Animal:new{name="Cat"}

	if not init_called then
		return false
	end

	if dog.age ~= 5 and cat.age ~= nil then
		return false
	end

	if dog:say("good night") ~= "Dog say: good night" then
		return false
	end

	if cat:say("good night") ~= "Cat say: good night" then
		return false
	end

	dog:change("@" .. dog.name)
	if dog.name ~= "@Dog" and cat.name ~= "Cat" then
		return false
	end

	cat:change("@" .. cat.name)
	if dog.name ~= "@Dog" and cat.name ~= "@Cat" then
		return false
	end

	if dog:say("good night") ~= "@Dog say: good night" then
		return false
	end

	if cat:say("good night") ~= "@Cat say: good night" then
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

func TestCreateClass_3(t *testing.T) {
	L := lua.NewState()
	L.PreloadModule(DiLibName, Loader)
	defer L.Close()

	code := `
	local di = require('di')

	local Animal = di.createClass({name="Animal"})

	function Animal:init()
		self.name = "@" .. self.name
	end

	function Animal:say(msg)
		return string.format("%s say: %s", self.name, msg)
	end

	local dog = Animal:new({name="Dog", age=5})

	if dog.name ~= "@Dog" then
		return false
	end

	if dog:say("good night") ~= "@Dog say: good night" then
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

func TestCreateClass_4(t *testing.T) {
	L := lua.NewState()
	L.PreloadModule(DiLibName, Loader)
	defer L.Close()

	code := `
	local di = require('di')

	local Animal = di.createClass({name="Animal"})

	function Animal:init()
		self.name = "@" .. self.name
	end

	function Animal:say(msg)
		return string.format("%s say: %s", self.name, msg)
	end

	local dog = Animal:new({name="Dog", age=5})

	if di.isClass(Animal) ~= true then
		return false
	end
	if di.isClass(dog) ~= false then
		return false
	end
	if di.instanceof(dog, Animal) ~= true then
		return false
	end
	if di.instanceof(Animal, dog) ~= false then
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
