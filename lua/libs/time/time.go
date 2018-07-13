// (c) 2018, Jeff Chen <jefurry@qq.com>
//
// This file is part of Gola
//
// Copyright 2018 The Gola Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package time implements time for Lua.
package time

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	ttime "time"
)

const (
	TimeLibName = "time"
)

func Open(L *lua.LState) {
	L.PreloadModule(TimeLibName, Loader)
}

func Loader(L *lua.LState) int {
	dimod := L.SetFuncs(L.NewTable(), timeFuncs)
	L.Push(dimod)

	timeRegisterTimeMetatype(L)
	timeRegisterLocationMetatype(L)

	dimod.RawSetString("UTC", newLocation(L, ttime.UTC))
	dimod.RawSetString("Local", newLocation(L, ttime.Local))

	for k, v := range timeFields {
		dimod.RawSetString(k, v)
	}

	return 1
}

var timeFuncs = map[string]lua.LGFunction{
	"parse":                  timeParse,
	"parseInLocation":        timeParseInLocation,
	"parseDuration":          timeParseDuration,
	"fixedZone":              timeFixedZone,
	"loadLocation":           timeLoadLocation,
	"loadLocationFromTZData": timeLoadLocationFromTZData,
	"now":    timeNow,
	"date":   timeDate,
	"since":  timeSince,
	"until":  timeUntil,
	"unix":   timeUnix,
	"sleep":  timeSleep,
	"isLeap": timeIsLeap,
}

var timeFields = map[string]lua.LValue{
	"ANSIC":        lua.LString(ttime.ANSIC),
	"UNIX_DATE":    lua.LString(ttime.UnixDate),
	"RUBY_DATE":    lua.LString(ttime.RubyDate),
	"RFC822":       lua.LString(ttime.RFC822),
	"RFC822Z":      lua.LString(ttime.RFC822Z),
	"RFC850":       lua.LString(ttime.RFC850),
	"RFC1123":      lua.LString(ttime.RFC1123),
	"RFC1123Z":     lua.LString(ttime.RFC1123Z),
	"RFC3339":      lua.LString(ttime.RFC3339),
	"RFC3339_NANO": lua.LString(ttime.RFC3339Nano),
	"KIT_CHEN":     lua.LString(ttime.Kitchen),
	// Handy time stamps.
	"STAMP":       lua.LString(ttime.Stamp),
	"STAMP_MILLI": lua.LString(ttime.StampMilli),
	"STAMP_MICRO": lua.LString(ttime.StampMicro),
	"STAMP_NANO":  lua.LString(ttime.StampNano),

	// time.Duration
	"NANO_SECOND":  lua.LNumber(ttime.Nanosecond),
	"MICRO_SECOND": lua.LNumber(ttime.Microsecond),
	"MILLI_SECOND": lua.LNumber(ttime.Millisecond),
	"SECOND":       lua.LNumber(ttime.Second),
	"MINUTE":       lua.LNumber(ttime.Minute),
	"HOUR":         lua.LNumber(ttime.Hour),
}

func timeNow(L *lua.LState) int {
	t := ttime.Now()

	L.Push(newTime(L, t))

	return 1
}

func timeDate(L *lua.LState) int {
	tb := L.OptTable(1, nil)
	if tb == nil {
		tb = L.CreateTable(0, 0)
	}

	year := tb.RawGetString("year")
	if year == lua.LNil {
		year = lua.LNumber(1970)
	}

	month := tb.RawGetString("month")
	if month == lua.LNil {
		month = lua.LNumber(1)
	}

	day := tb.RawGetString("day")
	if day == lua.LNil {
		day = lua.LNumber(1)
	}

	hour := tb.RawGetString("hour")
	if hour == lua.LNil {
		hour = lua.LNumber(0)
	}

	min := tb.RawGetString("min")
	if min == lua.LNil {
		min = lua.LNumber(0)
	}

	sec := tb.RawGetString("sec")
	if sec == lua.LNil {
		sec = lua.LNumber(0)
	}

	nsec := tb.RawGetString("nsec")
	if nsec == lua.LNil {
		nsec = lua.LNumber(0)
	}

	loc := tb.RawGetString("loc")
	if loc == lua.LNil {
		loc = newLocation(L, ttime.UTC)
	}

	if year.Type() != lua.LTNumber {
		L.ArgError(1, fmt.Sprintf("%s expected, got %s", lua.LTNumber, year.Type()))
	}
	y, _ := year.(lua.LNumber)

	if month.Type() != lua.LTNumber {
		L.ArgError(1, fmt.Sprintf("%s expected, got %s", lua.LTNumber, month.Type()))
	}
	mon, _ := month.(lua.LNumber)

	if day.Type() != lua.LTNumber {
		L.ArgError(1, fmt.Sprintf("%s expected, got %s", lua.LTNumber, day.Type()))
	}
	d, _ := day.(lua.LNumber)

	if hour.Type() != lua.LTNumber {
		L.ArgError(1, fmt.Sprintf("%s expected, got %s", lua.LTNumber, hour.Type()))
	}
	h, _ := hour.(lua.LNumber)

	if min.Type() != lua.LTNumber {
		L.ArgError(1, fmt.Sprintf("%s expected, got %s", lua.LTNumber, min.Type()))
	}
	m, _ := min.(lua.LNumber)

	if sec.Type() != lua.LTNumber {
		L.ArgError(1, fmt.Sprintf("%s expected, got %s", lua.LTNumber, sec.Type()))
	}
	s, _ := sec.(lua.LNumber)

	if nsec.Type() != lua.LTNumber {
		L.ArgError(1, fmt.Sprintf("%s expected, got %s", lua.LTNumber, nsec.Type()))
	}
	ns, _ := nsec.(lua.LNumber)

	if loc.Type() != lua.LTUserData {
		L.ArgError(1, fmt.Sprintf("%s expected, got %s", lua.LTNumber, loc.Type()))
	}

	ud, _ := loc.(*lua.LUserData)
	l, ok := ud.Value.(*ttime.Location)
	if !ok {
		L.ArgError(1, "invalid time zone")
	}

	t := ttime.Date(int(y), ttime.Month(int(mon)), int(d), int(h), int(m), int(s), int(ns), l)

	L.Push(newTime(L, t))

	return 1
}

func timeSince(L *lua.LState) int {
	t := checkTime(L, 1)

	d := ttime.Since(t)

	L.Push(lua.LNumber(d))

	return 1
}

func timeUntil(L *lua.LState) int {
	t := checkTime(L, 1)

	d := ttime.Until(t)

	L.Push(lua.LNumber(d))

	return 1
}

func timeUnix(L *lua.LState) int {
	sec := L.CheckInt(1)
	nsec := L.CheckInt(2)

	u := ttime.Unix(int64(sec), int64(nsec))

	L.Push(newTime(L, u))

	return 1
}

func timeSleep(L *lua.LState) int {
	d := L.CheckInt(1)

	ttime.Sleep(ttime.Duration(d))

	return 0
}

func timeIsLeap(L *lua.LState) int {
	year := L.CheckInt(1)

	L.Push(lua.LBool(isLeap(year)))

	return 1
}
