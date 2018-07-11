// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package loglib

import (
	"fmt"
	"github.com/jefurry/logrus"
	"github.com/jefurry/logrus/hooks/rotatelog"
	"github.com/yuin/gopher-lua"
	"time"
)

func checkLogger(L *lua.LState, n int) *logrus.Logger {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*logrus.Logger); ok {
		return v
	}

	L.ArgError(n, fmt.Sprintf("%s expected, got %s", logLoggerTypeName, ud.Type()))

	return nil
}

func loggerWrap(L *lua.LState) (*logrus.Logger, []interface{}) {
	logger := checkLogger(L, 1)
	vals := wrapVals(L, 2)

	return logger, vals
}

func loggerWrapformat(L *lua.LState) (*logrus.Logger, string, []interface{}) {
	logger := checkLogger(L, 1)
	format := L.CheckString(2)
	vals := wrapVals(L, 3)

	return logger, format, vals
}

func checkEntry(L *lua.LState, n int) *logrus.Entry {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*logrus.Entry); ok {
		return v
	}

	L.ArgError(n, fmt.Sprintf("%s expected, got %s", logEntryTypeName, ud.Type()))

	return nil
}

func entryWrap(L *lua.LState) (*logrus.Entry, []interface{}) {
	entry := checkEntry(L, 1)
	vals := wrapVals(L, 2)

	return entry, vals
}

func entryWrapformat(L *lua.LState) (*logrus.Entry, string, []interface{}) {
	entry := checkEntry(L, 1)
	format := L.CheckString(2)
	vals := wrapVals(L, 3)

	return entry, format, vals
}

func wrapVals(L *lua.LState, start int) []interface{} {
	vals := make([]interface{}, 0, 3)
	top := L.GetTop()
	if top >= start {
		for i := start; i <= top; i++ {
			vals = append(vals, L.Get(i))
		}
	}

	return vals
}

func checkTextFormatter(L *lua.LState, n int) *logrus.TextFormatter {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*logrus.TextFormatter); ok {
		return v
	}

	L.ArgError(n, fmt.Sprintf("%s expected, got %s", logTextFormatterTypeName, ud.Type()))

	return nil
}

func checkJSONFormatter(L *lua.LState, n int) *logrus.JSONFormatter {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*logrus.JSONFormatter); ok {
		return v
	}

	L.ArgError(n, fmt.Sprintf("%s expected, got %s", logJSONFormatterTypeName, ud.Type()))

	return nil
}

func checkRotatelog(L *lua.LState, n int) *rotatelog.RotatelogHook {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*rotatelog.RotatelogHook); ok {
		return v
	}

	L.ArgError(n, fmt.Sprintf("%s expected, got %s", logRotatelogTypeName, ud.Type()))

	return nil
}

func checkRotatelogHookOptions(L *lua.LState, val *lua.LTable) []rotatelog.Option {
	opts := make([]rotatelog.Option, 0, 6)
	if val != nil {
		withClock := val.RawGetString("withClock")
		if withClock != lua.LNil {
			switch withClock {
			case lua.LString("UTC"):
				opts = append(opts, rotatelog.WithClock(rotatelog.UTC))
			case lua.LString("Local"):
				opts = append(opts, rotatelog.WithClock(rotatelog.Local))
			}
		}

		withLocation := val.RawGetString("withLocation")
		if withLocation != lua.LNil {
			zone, ok := withLocation.(lua.LString)
			if !ok {
				L.ArgError(2, fmt.Sprintf("%s expected, got %s", lua.LTString, withLocation.Type()))
			}

			loc, err := time.LoadLocation(string(zone))
			if err != nil {
				L.ArgError(2, fmt.Sprintf("cannot find time zone %v", zone))
			}

			opts = append(opts, rotatelog.WithLocation(loc))
		}

		withLinkName := val.RawGetString("withLinkName")
		if withLinkName != lua.LNil {
			linkName, ok := withLinkName.(lua.LString)
			if !ok {
				L.ArgError(2, fmt.Sprintf("%s expected, got %s", lua.LTString, withLinkName.Type()))
			}

			opts = append(opts, rotatelog.WithLinkName(string(linkName)))
		}

		withMaxAge := val.RawGetString("withMaxAge")
		if withMaxAge != lua.LNil {
			maxAge, ok := withMaxAge.(lua.LNumber)
			if !ok {
				L.ArgError(2, fmt.Sprintf("%s expected, got %s", lua.LTNumber, withMaxAge.Type()))
			}

			opts = append(opts, rotatelog.WithMaxAge(time.Duration(maxAge)))
		}

		withRotationTime := val.RawGetString("withRotationTime")
		if withRotationTime != lua.LNil {
			rotationTime, ok := withRotationTime.(lua.LNumber)
			if !ok {
				L.ArgError(2, fmt.Sprintf("%s expected, got %s", lua.LTNumber, withRotationTime.Type()))
			}

			opts = append(opts, rotatelog.WithRotationTime(time.Duration(rotationTime)))
		}

		withRotationCount := val.RawGetString("withRotationCount")
		if withRotationCount != lua.LNil {
			rotationCount, ok := withRotationCount.(lua.LNumber)
			if !ok {
				L.ArgError(2, fmt.Sprintf("%s expected, got %s", lua.LTNumber, withRotationCount.Type()))
			}

			opts = append(opts, rotatelog.WithRotationCount(uint(rotationCount)))
		}
	}

	return opts
}
