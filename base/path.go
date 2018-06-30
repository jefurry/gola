// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package base

import (
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/yuin/gopher-lua"
	"os"
	"strings"
)

const (
	pathPrefix1 = "." + noDotPathPrefix1
	pathPrefix2 = "." + noDotPathPrefix2
)

var (
	errMethodNotAllow = errors.New("the method is not allowed after setting the environment variable LUA_PATH")
)

var OldLuaPathDefault string

func init() {
	OldLuaPathDefault = lua.LuaPathDefault

	trimPathPrefix()
	addPathPrefix()
}

func trimPathPrefix() {
	lua.LuaPathDefault = strings.TrimPrefix(lua.LuaPathDefault, pathPrefix1)
	lua.LuaPathDefault = strings.TrimPrefix(lua.LuaPathDefault, ";")
	lua.LuaPathDefault = strings.TrimPrefix(lua.LuaPathDefault, pathPrefix2)
	lua.LuaPathDefault = strings.TrimPrefix(lua.LuaPathDefault, ";")
}

func addPathPrefix() {
	lpd := lua.LuaPathDefault
	lua.LuaPathDefault = pathPrefix1 + ";" + pathPrefix2
	if lpd != "" {
		lua.LuaPathDefault = lua.LuaPathDefault + ";" + lpd
	}
}

func checkError() error {
	path := os.Getenv(lua.LuaPath)
	if len(path) != 0 {
		return errMethodNotAllow
	}

	return nil
}

// SetDefaultPath set path to lua.LuaPathDefault.
func SetDefaultPath(p string) error {
	if err := checkError(); err != nil {
		return err
	}

	lua.LuaPathDefault = ""

	if p == "" {
		addPathPrefix()

		return nil
	}

	return AddDefaultPath(p)
}

// AddDefaultPath add path to lua.LuaPathDefault.
func AddDefaultPath(p string) error {
	if err := checkError(); err != nil {
		return err
	}

	if p == "" {
		return nil
	}

	p, err := homedir.Expand(p)
	if err != nil {
		return err
	}

	p = strings.TrimRight(p, "/\\")
	if strings.Contains(lua.LuaPathDefault, p+noDotPathPrefix1+";") {
		return nil
	}

	trimPathPrefix()
	p = p + noDotPathPrefix1 + ";" + p + noDotPathPrefix2 + ";" + lua.LuaPathDefault
	lua.LuaPathDefault = strings.TrimRight(p, ";")
	addPathPrefix()

	return nil
}

// RemoveDefaultPath remote path from lua.LuaPathDefault.
func RemoveDefaultPath(p string) error {
	if err := checkError(); err != nil {
		return err
	}

	p, err := homedir.Expand(p)
	if err != nil {
		return err
	}

	trimPathPrefix()
	path := lua.LuaPathDefault
	p = strings.TrimRight(p, "/\\")
	s := p + noDotPathPrefix1 + ";"
	if i := strings.Index(path, s); i > -1 {
		path = strings.Replace(path, s+p+noDotPathPrefix2, "", -1)
		path = strings.TrimLeft(path, ";")
		path = strings.TrimRight(path, ";")
		path = strings.Replace(path, ";;", "", -1)

		lua.LuaPathDefault = path
	}

	addPathPrefix()

	return nil
}

// ResetDefaultPath restore to LuaPathDefault.
func ResetDefaultPath() error {
	if err := checkError(); err != nil {
		return err
	}

	lua.LuaPathDefault = OldLuaPathDefault

	return nil
}
