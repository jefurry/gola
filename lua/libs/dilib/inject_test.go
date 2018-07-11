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

func TestInjectorInstantiate_1(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local di = require('di')

	-- define class
	local Animal = di.createClass{className="Animal", name="Cat"}
	function Animal:init()
		--print("call init method of Animal class by di.")
	end

	local Person = di.createClass{className="Person", name="Jeff"}
	function Person:inject(animal)
		self.animal = animal
	end
	function Person:say(message)
		return string.format("%s: my name is %s, i love %s.", message, self.name, self.animal.name)
	end

	-- create injector
	local injector = di.newInjector{animal={"type", Animal}}
	injector:add{person={"type", Person}}
	injector:add{message={"value", "he said"}}

	-- call
	local person = injector:get("person")
	if person == nil then
		return false
	end

	local msg = injector:invoke{"message", {person, "say"}}
	if msg ~= "he said: my name is Jeff, i love Cat." then
		return false
	end

	local person = injector:instantiate{Person}
	if person == nil then
		return false
	end

	local msg = injector:invoke(function(animal, message)
		return string.format("%s: %s's className is %s.", message, animal.className, animal.name)
	end)

	if msg ~= "he said: Animal's className is Cat." then
		return false
	end

	local msg = injector:invoke{"animal", "message", function(animal, message)
		return string.format("%s: %s's className is %s.", message, animal.className, animal.name)
	end}

	if msg ~= "he said: Animal's className is Cat." then
		return false
	end

	local msg = injector:invoke({"animal", "message", function(animal, message)
		return string.format("%s: %s's className is %s.", message, animal.className, animal.name)
	end}, {message="you said"})

	if msg ~= "you said: Animal's className is Cat." then
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

func TestInjectorInstantiate_2(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local di = require('di')

	-- define class
	local Cat = di.createClass{className="Cat", name="Cat"}
	function Cat:init()
		--print("call init method of Cat class by di.")
	end
	function Cat:inject(catName)
		self.name = catName
	end

	local Dog = di.createClass{className="Dog", name="Dog"}
	function Dog:init()
		--print("call init method of Dog class by di.")
	end
	function Dog:inject(dogName)
		self.name = dogName
	end

	local Person = di.createClass{className="Person", name="Jeff"}
	function Person:inject(cat, dog)
		self.cat = cat
		self.dog = dog
	end
	function Person:say(message)
		return string.format("%s: my name is %s, i love %s and %s.", message, self.name, self.cat.name, self.dog.name)
	end

	-- create injector
	local injector = di.newInjector{cat={"type", Cat}, dog={"type", Dog}}
	injector:add{person={"type", Person}}
	injector:add{message={"value", "he said"}}
	injector:add{catName={"value", "ZeroCat"}}
	injector:add{dogName={"value", "HotDog"}}

	-- call
	local person = injector:instantiate(Person)
	if person == nil then
		return false
	end

	local msg = injector:invoke{"message", {person, "say"}}
	if msg ~= "he said: my name is Jeff, i love ZeroCat and HotDog." then
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

func TestInjectorInstantiate_3(t *testing.T) {
	L := lua.NewState()
	Open(L)
	defer L.Close()

	code := `
	local di = require('di')

	-- define class
	local Cat = di.createClass{className="Cat", name="Cat"}
	function Cat:init()
		--print("call init method of Cat class by di.")
	end
	function Cat:inject(person)
		self.person = person
	end

	local Person = di.createClass{className="Person", name="Jeff"}
	function Person:inject(cat)
		self.cat = cat
	end

	-- create injector
	local injector = di.newInjector{person={"type", Person}, cat={"type", Cat}}

	-- call
	local persion, msg = injector:get("person")
	if person ~= nil or msg ~= "cannot resolve circular dependency for 'person'" then
		return false
	end

	local cat, msg = injector:get("cat")
	if cat ~= nil or msg ~= "cannot resolve circular dependency for 'cat'" then
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
